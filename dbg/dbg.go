/*
	Debug tools
*/
package dbg

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"sync"
)

// Print functions generated by FlagPrintf.
type PrintFunc func(fmts string, arg ...face{}) (int, error)

// To be embedded in structures to add a debug flag and Dprintf function.
struct Flag {
	Tag   string // prefixing debug prints
	Debug bool   // enable debug prints
}

var (
	ExitDumpsStacks bool // If set Fatal/Exit will dump the stacks
	lk              sync.Mutex
)

func xexit(sts int) {
	if ExitDumpsStacks {
		var buf [64 * 1024]byte
		n := runtime.Stack(buf[0:], true)
		os.Stderr.Write(buf[:n])
	}
	os.Exit(sts)
}

func fatal(warn bool, args ...face{}) {
	if len(args) == 0 || args[0] == nil {
		xexit(0)
	}
	if s, ok := args[0].(string); ok {
		if s == "" {
			xexit(0)
		}
		if warn {
			Warn(s, args[1:]...)
		}
	} else if e, ok := args[0].(error); ok {
		if e == nil {
			xexit(0)
		}
		if warn {
			Warn("%s", e)
		}
	} else if warn {
		Warn("fatal")
	}
	xexit(1)
}

// Printf with d.Tag if d.Debug is set
func (d Flag) Dprintf(str string, args ...face{}) (n int, err error) {
	if d.Debug {
		return Printf("%s: %s", d.Tag, fmt.Sprintf(str, args...))
	}
	return 0, nil
}

// Atomic print to stderr
func Printf(str string, args ...face{}) (n int, err error) {
	lk.Lock()
	defer lk.Unlock()
	return fmt.Fprintf(os.Stderr, str, args...)
}

// Return a function that calls Printf but only when flag is set.
func FlagPrintf(flag *bool) PrintFunc {
	return func(fmts string, arg ...face{}) (int, error) {
		if *flag {
			return Printf(fmts, arg...)
		}
		return 0, nil
	}
}

// Printf to stderr, prefixed with program name and terminating with \n.
// Each warn is atomic.
func Warn(str string, args ...face{}) (n int, err error) {
	return Printf("%s: %s\n", os.Args[0], fmt.Sprintf(str, args...))
}

// Warn and exit
func Fatal(args ...face{}) {
	if len(args) == 0 {
		args = append(args, "fatal errors")
	}
	fatal(true, args...)
}

// Exit with 0 or 1 depending on the args (like Fatal), but do not print.
func Exit(args ...face{}) {
	fatal(false, args...)
}

// Return a string with the hex dump of at most n bytes from data
// (all if n is 0) for debugging.
func HexStr(data []byte, n int) string {
	var b bytes.Buffer
	if len(data)%2 != 0 {
		fmt.Fprintf(&b, " %02x", data[0])
		data = data[1:]
	}
	if n <= 0 || n > len(data) {
		n = len(data)
	}
	for i := 0; i < n/2; i++ {
		fmt.Fprintf(&b, " %02x", data[i])
	}
	if n < len(data) {
		fmt.Fprintf(&b, "...")
	}
	for i := len(data) - n/2; i < len(data); i++ {
		fmt.Fprintf(&b, " %02x", data[i])
	}
	return b.String()
}

// Return a string with at most n runes for debugging
func Str(s string, n int) string {
	if n <= 0 || n >= len(s) {
		return s
	}
	rs := []rune(s)
	if n <= 0 || n > len(rs) {
		n = len(rs)
	}
	if n == len(rs) {
		return s
	}
	return string(rs[:n/2]) + "..." + string(rs[len(rs)-n/2:])
}
