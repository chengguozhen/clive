package main

import (
	"clive/ch"
	"clive/cmd"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

struct xFd {
	sync.Mutex
	ref  int
	fd   *os.File
	path string
	isIn bool
}

struct pFd {
	r, w *os.File
}

struct bgCmds {
	sync.Mutex
	cmds  map[*xEnv]bool
	waits map[string]chan bool
	wall  chan bool
}

// Execution environment for nodes.
// We use unix processes to run clive commands, so for now we
// use actual file descriptors for cmd IO.
// The IO environment is named in clive, 0, 1, 2 are "in", "out", "err",
// other names can be passed using environment variables that map
// the name to the unix file descriptor.
struct xEnv {
	fds   map[string]*xFd
	waits []chan bool
	bgtag string
	isbg  bool // this cmd is a child of a bg command
	xctx  *cmd.Ctx
}

var bgcmds = bgCmds{
	cmds:  map[*xEnv]bool{},
	waits: map[string]chan bool{},
	wall:  make(chan bool),
}

func (x *xEnv) Printf(fmts string, arg ...face{}) (int, error) {
	if ofd, ok := x.fds["out"]; ok {
		s := fmt.Sprintf(fmts, arg...)
		_, err := ch.WriteMsg(ofd.fd, 1, []byte(s))
		return len(s), err
	}
	return cmd.Printf(fmts, arg...)
}

func (x *xEnv) Eprintf(fmts string, arg ...face{}) (int, error) {
	if ofd, ok := x.fds["err"]; ok {
		s := fmt.Sprintf(fmts, arg...)
		_, err := ch.WriteMsg(ofd.fd, 1, []byte(s))
		return len(s), err
	}
	return cmd.Eprintf(fmts, arg...)
}

func (xfd *xFd) addref() {
	xfd.Lock()
	if xfd.ref > 0 {
		xfd.ref++
	}
	xfd.Unlock()
}

func (xfd *xFd) Close() error {
	xfd.Lock()
	if xfd.ref > 0 {
		xfd.ref--
		if xfd.ref == 0 {
			xfd.fd.Close()
		}
	}
	xfd.Unlock()
	return nil
}

func (x *xEnv) Close() error {
	for _, fd := range x.fds {
		if fd != nil {
			fd.Close()
		}
	}
	for _, wc := range x.waits {
		<-wc
	}
	x.waits = nil
	return nil
}

// Functions run...() run the command and wait for it to complete.
// The pipe creates a cmd ctx for each command in the pipe and perhaps waits for the last.
// Each run/start function receives a environment that can be changed
// by the function. If it creates children commands it should dup the environment
// for them, so they can change their own env.
// The errors returned by the functions indicate errors that lead to stop
// in the execution of commands
// The error returned by wait matches the exit status of the command.

func (b *bgCmds) add(x *xEnv) {
	b.Lock()
	defer b.Unlock()
	b.cmds[x] = true
	if x.bgtag != "&" && b.waits[x.bgtag] == nil {
		b.waits[x.bgtag] = make(chan bool)
	}
}

func (b *bgCmds) del(x *xEnv) {
	b.Lock()
	defer b.Unlock()
	delete(b.cmds, x)
	if x.bgtag != "&" {
		wc := b.waits[x.bgtag]
		delete(b.waits, x.bgtag)
		if wc != nil {
			close(wc)
		}
	}
	if len(b.cmds) == 0 {
		close(b.wall)
		b.wall = make(chan bool)
	}
}

func (b *bgCmds) wait(tag string) {
	b.Lock()
	if tag != "" {
		wc := b.waits[tag]
		b.Unlock()
		if wc != nil {
			<-wc
		}
		return
	}
	if len(b.cmds) == 0 {
		b.Unlock()
		return
	}
	wc := b.wall
	b.Unlock()
	<-wc
}

var env = os.Environ()

func (x *xEnv) addUXio() {
	for _, v := range env {
		if !strings.HasPrefix(v, "cliveio#") {
			continue
		}
		toks := strings.Split(v, "=")
		if len(toks) != 2 {
			continue
		}
		cname := strings.TrimPrefix(toks[0], "cliveio#")
		fdval := toks[1]
		if len(fdval) < 2 {
			continue
		}
		if cname == "in" || cname == "out" || cname == "err" {
			continue
		}
		dir := fdval[0]
		if dir != '<' && dir != '>' {
			continue
		}
		fdname := fdval[1:]
		n, err := strconv.Atoi(fdname)
		if err != nil {
			continue
		}
		fd := os.NewFile(uintptr(n), cname)
		x.fds[cname] = &xFd{fd: fd, path: cname, ref: -1, isIn: dir == '<'}
	}
}

func newEnv() *xEnv {
	x := &xEnv{
		fds: map[string]*xFd{
			"in":  &xFd{fd: os.Stdin, path: "in", ref: -1, isIn: true},
			"out": &xFd{fd: os.Stdout, path: "out", ref: -1, isIn: false},
			"err": &xFd{fd: os.Stderr, path: "err", ref: -1, isIn: false},
		},
	}
	x.addUXio()
	return x
}

func (x *xEnv) dup() *xEnv {
	ne := &xEnv{
		fds:  map[string]*xFd{},
		isbg: x.isbg,
	}
	for k, f := range x.fds {
		f.addref()
		ne.fds[k] = f
	}
	return ne
}

// run a top-level command
func (nd *Nd) run() error {
	if nd == nil {
		return nil
	}
	nprintf("cmd:\n%s\n", dnd{nd})
	defer nprintf("cmd done\n")
	if dry || yylex.nerrors > 0 {
		cmd.Dprintf("skip cmd run (dry|error)")
		yylex.nerrors = 0
		return nil
	}
	x := newEnv()
	// bgpipe or func
	var err error
	switch nd.typ {
	case Npipe:
		err = nd.runPipe(x)
	case Nsrc:
		err = nd.runSrc(x)
	case Nfunc:
		err = nd.runFunc(x)
	default:
		panic(fmt.Errorf("run: bad type %s", nd.typ))
	}
	if isExit(err) {
		cmd.Exit(strings.TrimPrefix(err.Error(), "qlexit"))
	}
	return err
}

func (nd *Nd) chk(k ...NdType) {
	if nd == nil {
		panic("chk: nil node")
	}
	for i := 0; i < len(k); i++ {
		if k[i] == nd.typ {
			cmd.Dprintf("chk %s...\n", nd.typ)
			return
		}
	}
	panic(fmt.Errorf("not %v; type %v", k, nd.typ))
}

func (nd *Nd) runSrc(x *xEnv) error {
	nd.chk(Nsrc)
	if len(nd.Child) != 1 {
		panic("runSrc: bad children")
	}
	nms, err := nd.Child[0].expand1(x)
	if err != nil {
		cmd.Warn("source: %s", err)
		cmd.SetEnv("sts", "source: expand")
		return nil
	}
	if len(nms) != 1 {
		cmd.Warn("source: more than one name given")
		cmd.SetEnv("sts", "source: more than one name given")
		return nil
	}
	cmd.VWarn("source %s\n", nms[0])
	if err := yylex.source(nms[0]); err != nil {
		cmd.SetEnv("sts", err.Error())
	} else {
		cmd.SetEnv("sts", "")
	}
	return nil
}

func (nd *Nd) runFunc(x *xEnv) error {
	nd.chk(Nfunc)
	if len(nd.Args) != 1 {
		panic("runFunc: bad args")
	}
	if len(nd.Child) != 1 {
		panic("runFunc: bad children")
	}
	newFunc(nd)
	return nil
}

// make xEnvs for pipe children
func (nd *Nd) mkChildEnvs(x *xEnv) ([]*xEnv, error) {
	var err error
	var pcloses []io.Closer
	nc := len(nd.Child)
	cxs := make([]*xEnv, nc)
	pipes := map[string]pFd{}
	defer func() {
		if err != nil {
			for _, x := range pcloses {
				x.Close()
			}
			return
		}
	}()
	for i, c := range nd.Child {
		cx := x.dup()
		pcloses = append(pcloses, cx)
		cxs[i] = cx
		if dry {
			continue
		}
		ccloses, err := c.applyRedirs(x, cx, pipes)
		pcloses = append(pcloses, ccloses...)
		if err != nil {
			return nil, err
		}
	}
	return cxs, nil
}

func (x *xEnv) wait() error {
	if x == nil || x.xctx == nil {
		return nil
	}
	wc := x.xctx.Waitc()
	<-wc
	err := cerror(wc)
	if err == nil {
		cmd.SetEnv("sts", "")
	} else {
		cmd.SetEnv("sts", err.Error())
	}
	return err
}

func (x *xEnv) bg(tag string) {
	x.bgtag = tag
	bgcmds.add(x)
	if x.xctx != nil {
		wc := x.xctx.Waitc()
		go func() {
			<-wc
			bgcmds.del(x)
		}()
	}
	cmd.SetEnv("sts", "")
}

// children may be cmd, block, for, while, cond, set
func (nd *Nd) runPipe(x *xEnv) error {
	nd.chk(Npipe)
	cxs, err := nd.mkChildEnvs(x)
	if err != nil {
		return err
	}
	bg := nd.Args[0]
	for i, c := range nd.Child {
		c := c
		cx := cxs[i]
		if bg != "" {
			cx.isbg = true
		}
		cx.xctx = cmd.New(func() {
			defer cx.Close()
			if bg != "" || i < len(nd.Child)-1 {
				cmd.ForkEnv()
				cmd.ForkNS()
				cmd.ForkDot()
			}
			switch c.typ {
			case Ncmd:
				err = c.runCmd(cx)
			case Nblock:
				err = c.runBlock(cx)
			case Nfor:
				err = c.runFor(cx)
			case Nwhile:
				err = c.runWhile(cx)
			case Ncond:
				err = c.runCond(cx)
			case Nset:
				err = c.runSet(cx)
			case Nsetmap:
				err = c.runSetMap(cx)
			default:
				panic(fmt.Errorf("run: bad pipe child type %s", c.typ))
			}
			if err != nil {
				cmd.Exit(err)
			}
			sts := cmd.GetEnv("sts")
			if sts != "" {
				cmd.Exit(sts)
			}
		})
	}
	if err != nil {
		return err
	}
	cx := cxs[len(nd.Child)-1]
	if bg != "" {
		cx.bg(bg)
	} else {
		if err := cx.wait(); isBreak(err) || isExit(err) {
			return err
		}
	}
	return nil
}

func (nd *Nd) varLen() int {
	nd.chk(Nlen)
	if len(nd.Args) != 1 {
		panic("bad Nlen arg list")
	}
	v := cmd.GetEnv(nd.Args[0])
	if v == "" {
		return 0
	}
	if cmd.IsEnvMap(v) {
		return len(cmd.EnvMap(v))
	}
	return len(cmd.EnvList(v))
}

func (nd *Nd) varValue(x *xEnv) (names []string) {
	var err error
	nd.chk(Nval, Nsingle)
	if len(nd.Args) != 1 {
		panic("bad var node args")
	}
	v := cmd.GetEnv(nd.Args[0])
	switch len(nd.Child) {
	case 0: // $a
		if cmd.IsEnvMap(v) {
			names = mapKeys(cmd.EnvMap(v))
		} else {
			names = cmd.EnvList(v)
		}
	case 1: // $a[b]
		c := nd.Child[0]
		names, err = c.expand1(x)
		if err != nil {
			cmd.Warn("expand: %s", err)
			break
		}
		if len(names) != 1 {
			cmd.Warn("$%s[...]: not a single index name", nd.Args[0])
			break
		}
		if cmd.IsEnvMap(v) {
			m := cmd.EnvMap(v)
			cmd.Dprintf("map value %v [ %s ] ...\n", m, names[0])
			names = m[names[0]]
		} else {
			lst := cmd.EnvList(v)
			cmd.Dprintf("lst value %v [ %s ] ...\n", lst, names[0])
			el := listEl(lst, names[0])
			if el == "" {
				names = []string{}
			} else {
				names = []string{el}
			}
		}
	default:
		panic("bad Nvar children list")
	}
	if nd.typ == Nsingle {
		name := strings.Join(names, " ")
		names = []string{name}
	}
	return names
}

func (nd *Nd) appNames(x *xEnv) (names []string) {
	nd.chk(Napp)
	if len(nd.Child) != 2 {
		panic("bad app node children")
	}
	var left, right []string
	var err error
	if nd.Child[0].typ == Nnames {
		left, err = nd.Child[0].expand(x)
	} else {
		left, err = nd.Child[0].expand1(x)
	}
	if err != nil {
		cmd.Warn("expand: append: %s", err)
		return nil
	}
	if nd.Child[1].typ == Nnames {
		right, err = nd.Child[1].expand(x)
	} else {
		right, err = nd.Child[1].expand1(x)
	}
	if err != nil {
		cmd.Warn("expand: append: %s", err)
		return nil
	}
	if len(left) == 0 {
		return right
	}
	if len(right) == 0 {
		return left
	}
	if len(left) == 1 {
		for i := 0; i < len(right); i++ {
			right[i] = left[0] + right[i]
		}
		return right
	}
	if len(right) == 1 {
		for i := 0; i < len(left); i++ {
			left[i] += right[0]
		}
		return left
	}
	if len(left) != len(right) {
		cmd.Warn("expand: different list lengths")
		return nil
	}
	for i := 0; i < len(left); i++ {
		left[i] += right[i]
	}
	return left
}

func (nd *Nd) pipeFrom(x *xEnv, cname string) (*xFd, error) {
	if cname == "" || cname == "in" {
		return nil, errors.New("can't pipe from command's input")
	}
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	cx := x.dup()
	if fd, ok := cx.fds[cname]; ok {
		fd.Close()
	}
	cx.fds[cname] = &xFd{fd: w, path: cname, ref: 1, isIn: false}
	cx.xctx = cmd.New(func() {
		defer cx.Close()
		if err := nd.runBlock(cx); err != nil {
			cmd.Exit(err)
		}
	})
	return &xFd{fd: r, path: "pipe", ref: 1, isIn: true}, nil
}

func (nd *Nd) pipeTo(x *xEnv, cname string) (*xFd, error) {
	if cname == "" || cname == "out" || cname == "err" {
		return nil, errors.New("can't pipe to command's output")
	}
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	cx := x.dup()
	if fd, ok := cx.fds[cname]; ok {
		fd.Close()
	}
	cx.fds[cname] = &xFd{fd: r, path: cname, ref: 1, isIn: true}
	cx.xctx = cmd.New(func() {
		defer cx.Close()
		if err := nd.runBlock(cx); err != nil {
			cmd.Exit(err)
		}
	})
	return &xFd{fd: w, path: "pipe", ref: 1, isIn: false}, nil
}

func collectNames(xfd *xFd) ([]string, error) {
	defer xfd.Close()
	names := []string{}
	for {
		_, _, m, err := ch.ReadMsg(xfd.fd)
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return names, nil
		}
		switch m := m.(type) {
		case []byte:
			s := string(m)
			if s != "" && s[len(s)-1] == '\n' {
				s = s[:len(s)-1]
			}
			cmd.Dprintf("expand io: msg '%s'\n", s)
			names = append(names, s)
		default:
			cmd.Dprintf("expand io: ignored %T\n", m)
		case error:
			cmd.Warn("expand: io: %s", m)
		}
	}
}

func (nd *Nd) expandIO(x *xEnv) ([]string, error) {
	nd.chk(Nioblk)
	// Either <{...} or <[names]{....} or >[name]{....}
	// The children is a block, including perhaps redirs.
	if len(nd.Args) == 0 || len(nd.Args) > 2 {
		panic("bad ioblk arg list")
	}
	if len(nd.Args) == 1 {
		fd, err := nd.pipeFrom(x, "out")
		if err != nil {
			cmd.Warn("expand: io: %s", err)
			return nil, nil
		}
		return collectNames(fd)
	}
	switch nd.Args[0] {
	case ">":
		tag := nd.Args[1]
		if strings.ContainsAny(tag, ";,") {
			return nil, fmt.Errorf("tag not implemented (only 'x' and 'x:y')")
		}
		tags := fields(tag, ":")
		if len(tags) == 1 {
			tags = []string{"in", tags[0]}
		}
		if len(tags) != 2 {
			return nil, fmt.Errorf("bad >{} tag '%s'", tag)
		}
		// eg: in: err	(block's in is a new err stream and we get |err as an arg)
		cname, nname := tags[0], tags[1]
		if cname == "out" || cname == "err" {
			return nil, errors.New("can't pipe to command's output")
		}
		if nname == "in" {
			return nil, errors.New("can't pipe from command's input")
		}
		pfd, err := nd.pipeTo(x, cname)
		if err != nil {
			cmd.Warn("expand: io: %s", err)
			return nil, nil
		}
		if fd, ok := x.fds[nname]; ok {
			fd.Close()
		}
		x.fds[nname] = pfd
		return []string{"|>" + nname}, nil
	case "<":
		tag := nd.Args[1]
		if strings.ContainsAny(tag, ";,") {
			return nil, fmt.Errorf("tag not implemented (only 'x' and 'x:y')")
		}
		tags := fields(tag, ":")
		if len(tags) == 1 {
			tags = append(tags, "out")
		}
		if len(tags) != 2 {
			return nil, fmt.Errorf("bad <{} tag '%s'", tag)
		}
		nname, cname := tags[0], tags[1]
		if nname == "out" || nname == "err" {
			return nil, errors.New("can't pipe to command's output")
		}
		pfd, err := nd.pipeFrom(x, cname)
		if err != nil {
			cmd.Warn("expand: io: %s", err)
			return nil, nil
		}
		if fd, ok := x.fds[nname]; ok {
			fd.Close()
		}
		x.fds[nname] = pfd
		return []string{"|<" + nname}, nil
	default:
		panic("bad ioblk arg")
	}

}

func (nd *Nd) expand1(x *xEnv) (nargs []string, err error) {
	nd.chk(Nname, Napp, Nlen, Nval, Nsingle, Nioblk)
	switch nd.typ {
	case Nname:
		nargs = nd.Args
	case Napp:
		nargs = nd.appNames(x)
	case Nlen:
		n := nd.varLen()
		nargs = []string{strconv.Itoa(n)}
	case Nval, Nsingle:
		nargs = nd.varValue(x)
	case Nioblk:
		nargs, err = nd.expandIO(x)
	default:
		panic(fmt.Errorf("expand1: bad names child type %s", nd.typ))
	}
	return nargs, err
}

// expand names: children can be name, app, len, single, val, ioblnk
func (nd *Nd) expand(x *xEnv) ([]string, error) {
	nd.chk(Nnames)
	xs := []string{}
	for _, c := range nd.Child {
		var nargs []string
		var err error
		if c.typ == Nnames {
			nargs, err = c.expand(x)
		} else {
			nargs, err = c.expand1(x)
		}
		if err != nil {
			return nil, err
		}
		xs = append(xs, nargs...)
	}
	nprintf("expanded: %v\n", xs)
	return xs, nil
}

func (nd *Nd) eval(x *xEnv, argv ...string) error {
	nd.chk(Nfunc)
	cmd.SetEnv("argv0", argv[0])
	e := cmd.ListEnv(argv[1:])
	cmd.SetEnv("argv", e)
	return nd.Child[0].runBlock(x)
}

func cleanenv(env []string) []string {
	for i := 0; i < len(env); {
		if strings.HasPrefix(env[i], "clivebg") {
			copy(env[i:], env[i+1:])
			env = env[:len(env)-1]
		} else {
			i++
		}
	}
	return env
}

func (nd *Nd) runCmd(x *xEnv) error {
	nd.chk(Ncmd)
	if len(nd.Child) != 1 {
		panic("bad Ncmd children")
	}
	args, err := nd.Child[0].expand(x)
	if err != nil {
		cmd.Warn("expand: %s", err)
		return err
	}
	cmd.VWarn("run: %s", dnames(args))
	if len(args) == 0 || len(args[0]) == 0 {
		err := errors.New("empty command name")
		cmd.Warn("run cmd: %s", err)
		return err
	}
	if dry {
		return nil
	}
	if args[0] == "builtin" {
		args = args[1:]
	} else {
		if fnd := getFunc(args[0]); fnd != nil {
			return fnd.eval(x, args...)
		}
	}
	if bfn := builtins[args[0]]; bfn != nil {
		return bfn(x, args...)
	}
	xc := exec.Command(args[0], args[1:]...)
	xc.Dir = cmd.Dot()
	xc.Env = cleanenv(cmd.OSEnv())
	if p := cmd.LookPath(args[0]); p != "" {
		xc.Path = p
	}
	for cname, xfd := range x.fds {
		switch cname {
		case "in":
			xc.Stdin = xfd.fd
		case "out":
			xc.Stdout = xfd.fd
		case "err":
			xc.Stderr = xfd.fd
		default:
			// XXX: TODO: set vars for In or out, not just for io
			dir := ">"
			if xfd.isIn {
				dir = "<"
			}
			no := 3 + len(xc.ExtraFiles)
			ev := fmt.Sprintf("cliveio#%s=%s%d", cname, dir, no)
			xc.Env = append(xc.Env, ev)
			xc.ExtraFiles = append(xc.ExtraFiles, xfd.fd)
		}
	}
	ev := fmt.Sprintf("dot=%s", cmd.Dot())
	xc.Env = append(xc.Env, ev)
	if x.isbg {
		xc.Env = append(xc.Env, "clivebg=y")
	}
	if err := xc.Start(); err != nil {
		cmd.Warn("%s", err)
		return nil
	}
	if err := xc.Wait(); err != nil {
		cmd.SetEnv("sts", err.Error())
		return nil
	} else {
		cmd.SetEnv("sts", "")
	}
	return nil
}

// block cmds are pipes or sources, there're also io blocks
func (nd *Nd) runBlock(x *xEnv) error {
	nd.chk(Nblock, Nioblk)
	if len(nd.Child) < 1 {
		panic("bad block children")
	}
	var err error
	for _, c := range nd.Child {
		cx := x.dup()
		defer cx.Close()
		switch c.typ {
		case Npipe:
			err = c.runPipe(cx)
		case Nsrc:
			err = c.runSrc(cx)
		default:
			panic(fmt.Errorf("runblock: bad child type %s", c.typ))
		}
		if err != nil {
			if isBreak(err) || isExit(err) {
				return err
			}
			break
		}
	}
	return nil
}

func (nd *Nd) runFor(x *xEnv) error {
	nd.chk(Nfor)
	if len(nd.Child) != 2 {
		panic("bad for children")
	}
	c0, blk := nd.Child[0], nd.Child[1]
	names, err := c0.expand(x)
	if err != nil {
		return err
	}
	if len(names) == 0 {
		cmd.Warn("missing for variable name")
		return fmt.Errorf("no variable name")
	}
	name, values := names[0], names[1:]
	if len(values) == 0 {
		fd := x.fds["in"]
		if fd != nil && fd.isIn {
			values, err = collectNames(fd)
			if err != nil {
				cmd.Warn("%s", err)
				cmd.SetEnv("sts", err.Error())
			}
		}
	}
	for _, v := range values {
		cmd.SetEnv(name, v)
		cx := x.dup()
		defer cx.Close()
		err = blk.runBlock(cx)
		if err != nil {
			if isExit(err) {
				return err
			}
			break
		}
	}
	cmd.SetEnv("sts", "")
	return nil
}

func (nd *Nd) runWhile(x *xEnv) error {
	nd.chk(Nwhile)
	if len(nd.Child) != 2 {
		panic("bad for children")
	}
	pipe, blk := nd.Child[0], nd.Child[1]
	var err error
	for {
		cx := x.dup()
		defer cx.Close()
		if err = pipe.runPipe(cx); err != nil {
			if isExit(err) {
				return err
			}
			break
		}
		if sts := cmd.GetEnv("sts"); sts != "" {
			break
		}
		cx2 := x.dup()
		defer cx2.Close()
		if err = blk.runBlock(cx2); err != nil {
			if isExit(err) {
				return err
			}
			break
		}
	}
	cmd.SetEnv("sts", "")
	return nil
}

var orSuccess = errors.New("or sucessful")

// like a block w/o redirs
// As soon as a child is not sucessful, we stop and return nil
// if the last child does run, we must return orSuccess
// so startCond() knows it has to stop
func (nd *Nd) runOr(x *xEnv) error {
	nd.chk(Nor)
	if len(nd.Child) == 0 {
		panic("bad or children")
	}
	for i, c := range nd.Child {
		var err error
		cx := x.dup()
		defer cx.Close()
		switch c.typ {
		case Npipe:
			err = c.runPipe(cx)
		case Nsrc:
			err = c.runSrc(cx)
		default:
			panic(fmt.Errorf("runor: bad child type %s", c.typ))
		}
		if err != nil {
			return err
		}
		if i < len(nd.Child)-1 {
			if sts := cmd.GetEnv("sts"); sts != "" {
				return nil
			}
		}
	}
	cmd.SetEnv("sts", "")
	return orSuccess
}

// children are or nodes that are like blocks (w/o redirs)
// and the nd has a final redir child
func (nd *Nd) runCond(x *xEnv) error {
	nd.chk(Ncond)
	if len(nd.Child) == 0 {
		// at least an or
		panic("bad cond children")
	}
	var err error
	for _, or1 := range nd.Child {
		cx := x.dup()
		defer cx.Close()
		if err = or1.runOr(cx); err != nil {
			if err == orSuccess {
				err = nil
			}
			if isExit(err) {
				return err
			}
			if isBreak(err) {
				err = nil
			}
			break
		}
	}
	cmd.SetEnv("sts", "")
	return err
}

func (nd *Nd) runSet(x *xEnv) error {
	nd.chk(Nset)
	if len(nd.Args) == 0 {
		panic("bad set args")
	}
	if len(nd.Child) == 0 || len(nd.Child) > 2 {
		panic("bad set children")
	}
	name := nd.Args[0]
	switch len(nd.Child) {
	case 1: // $name = ...
		c0 := nd.Child[0]
		vals, err := c0.expand(x)
		if err != nil {
			return err
		}
		cmd.VWarn("set %s = %s", name, dnames(vals))
		cmd.SetEnvList(name, vals)
	case 2: // $name[name] = ...
		c0, c1 := nd.Child[0], nd.Child[1]
		idxs, err := c0.expand1(x)
		if err != nil {
			return err
		}
		if len(idxs) == 0 {
			cmd.Warn("set %s: empty index", name)
			return nil
		}
		if len(idxs) > 1 {
			cmd.Warn("set %s: multiple index", name)
			return nil
		}
		idx := idxs[0]
		vals, err := c1.expand(x)
		if err != nil {
			return err
		}
		e := cmd.GetEnv(name)
		if cmd.IsEnvMap(e) {
			m := cmd.EnvMap(e)
			m[idx] = vals
			cmd.VWarn("set %s[%s] = %s", name, idx, dnames(vals))
			cmd.SetEnvMap(name, m)
		} else {
			lst := cmd.EnvList(e)
			lst = setListEl(lst, idx, strings.Join(vals, " "))
			cmd.VWarn("set %s[%s] = %s", name, idx, dnames(vals))
			cmd.SetEnvList(name, lst)
		}
	default:
	}
	return nil
}

func (nd *Nd) runSetMap(x *xEnv) error {
	nd.chk(Nsetmap)
	if len(nd.Args) == 0 {
		panic("bad set args")
	}
	name := nd.Args[0]
	m := map[string][]string{}
	for _, c := range nd.Child {
		nms, err := c.expand(x)
		if err != nil {
			return err
		}
		if len(nms) == 0 {
			continue
		}
		m[nms[0]] = nms[1:]
	}
	cmd.VWarn("set %s = %v", name, m)
	cmd.SetEnvMap(name, m)
	return nil
}
