/*
	Ink exec.
	A shell for clive using ink.
*/
package main

import (
	"clive/cmd"
	"clive/cmd/opt"
	"clive/net/ink"
	"strings"
	"clive/ch"
	"clive/cmd/run"
	"strconv"
	"time"
	"net/url"
)

struct Cmds {
	win *ink.Txt
}

func pgdrop(pg *ink.Pg, ev *ink.Ev) {
	args := ev.Args
	what := args[1]
	uri, err := url.Parse(what)
	if err != nil {
		cmd.Dprintf("pgdrop: not a url: %s\n", err)
		return
	}
	if !uri.IsAbs() {
		cmd.Dprintf("pgdrop: not an abs url\n")
		return
	}
	pg.Add(ink.Url(args[1]))
	
}

func pgevs(pg *ink.Pg) {
	for ev := range pg.Events() {
		ev := ev
		if ev == nil || len(ev.Args) == 0 {
			continue
		}
		cmd.Dprintf("ix: pg %s ev: %v\n", ev.Src, ev.Args)
		switch ev.Args[0] {
		case "click4":
			if len(ev.Args) < 2 {
				cmd.Dprintf("ix: short drop ev\n")
				continue
			}
			pgdrop(pg, ev)
		}
	}
}

func newCmds() *Cmds {
	win := ink.NewTxt("Online and ready");
	win.SetTag("cmds")
	cin := win.Events()
	c := &Cmds{win: win}
	go func() {
		for ev := range cin {
			ev := ev
			if ev == nil || len(ev.Args) == 0 {
				continue
			}
			cmd.Dprintf("cmds ev: %v\n", ev.Args)
			switch ev.Args[0] {
			case "click2":
				go c.exec(ev)
			}
		}
	}()
	return c
}

func (c *Cmds) exec(ev *ink.Ev) {
	if len(ev.Args) < 4 {
		cmd.Warn("bad args in click2")
		return
	}
	pos, err := strconv.Atoi(ev.Args[3])
	if err != nil {
		cmd.Warn("bad p1 in click2")
		return
	}
	line := ev.Args[1]
	hasnl := len(line) > 0 && line[len(line)-1] == '\n'
	ln := strings.TrimSpace(line)
	if len(ln) == 0 {
		return
	}
	c.win.SetMark("cmd", pos)
	cmd.Dprintf("exec %s\n", ln)
	x, err := run.UnixCmd(strings.Fields(ln)...)
	if err != nil {
		cmd.Warn("run: %s", err)
		c.win.MarkIns("cmd", []rune("error: " + err.Error() + "\n"));
		return
	}
	go func() {
		for m := range ch.GroupBytes(ch.Merge(x.Out, x.Err), time.Second, 4096) {
			switch m := m.(type) {
			case []byte:
				cmd.Dprintf("-> [%d] bytes\n", len(m))
				s := string(m)
				if !hasnl {
					hasnl = true
					s = "\n" + s
				}
				c.win.MarkIns("cmd", []rune(s))
			default:
				cmd.Dprintf("got type %T\n", m)
			}
		}
		if err := cerror(x.Err); err != nil {
			c.win.MarkIns("cmd", []rune("error: " + err.Error() + "\n"));
		}
		c.win.DelMark("cmd")
	}()
}


func main() {
	opts := opt.New("")
	c := cmd.AppCtx()
	opts.NewFlag("D", "debug", &c.Debug)
	cmd.UnixIO()
	args := opts.Parse()
	if len(args) != 0 {
		opts.Usage()
	}
	cmds1 := newCmds();
	cmds2 := newCmds();
	col1 := []face{}{"XX", "YY", cmds1.win}
	col2 := []face{}{"XX", "ZZ", cmds2.win}
	col3 := []face{}{"XX", ink.Url("http://lsub.org")}
	pg := ink.NewColsPg("/", col1, col2, col3)
	pg.Tag = "Ink cmds"
	go pgevs(pg)
	go ink.Serve(":8181")
	go func() {
		time.Sleep(5*time.Second)
		pg.Del(cmds2.win.Id)
		cmd.Warn("pg views: %v", pg.Views())
		cmds1.win.Dirty()
		cmds1.win.SetTag("ql")
		time.Sleep(15*time.Second)
		cmds1.win.Clean()
		time.Sleep(5*time.Second)
		cmds1.win.Dirty()
	}()
	cmds1.win.Wait()
	cmds2.win.Wait()
}
