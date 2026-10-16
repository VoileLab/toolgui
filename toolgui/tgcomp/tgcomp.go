package tgcomp

//go:generate go run ./internal/gen

// The rest of this package is generated from the internal tc* declarations
// marked //tgcomp:export. These two are not plain forwards.

import (
	"fmt"
	"go/parser"
	"runtime"
	"strings"
	"sync"

	"github.com/voilelab/toolgui/toolgui/tgcomp/internal/tccontent"
	"github.com/voilelab/toolgui/toolgui/tgframe"
)

// Base is embedded in every component conf, and is where a conf's ID comes
// from. A third-party component's conf embeds it the same way the built-in
// ones do.
type Base = tgframe.Base

// Echo will execute lambda and show the code in the lambda.
// To use Echo, we need to store the code in advance (usually by embedded).
//
//	//go:embed main.go
//	var code string
//	// ...
//	// ok, echo will execute and show `tgcomp.Text(c, "hello echo")`
//	tgcomp.Echo(c, code, func() {
//		tgcomp.Text(c, "hello echo")
//	})
//
//	// panic, since Echo only parse code line by line
//	tgcomp.Echo(c, code, func() {tgcomp.Text(c, "hello echo")})
//
//	// panic, since Echo only parse code that start from caller
//	myFunc := func() {
//		tgcomp.Text(c, "hello echo")
//	}
//	tgcomp.Echo(c, code, myFunc)
func Echo(c *tgframe.Container, code string, lambda func()) {
	_, filename, line, ok := runtime.Caller(1)
	if !ok {
		panic("Unable to get caller")
	}

	curCode, ok := codeCache.get(filename, line)
	if !ok {
		codeLines := strings.Split(code, "\n")
		left := line - 1
		right := line
		for i := line; i < len(codeLines); i++ {
			_, err := parser.ParseExpr(strings.Join(codeLines[left:i], "\n"))
			if err == nil {
				right = i
				break
			}
		}

		if left+1 > right-1 {
			panic("Fail to obtain code")
		}

		curCode = strings.Join(removeIndent(codeLines[left+1:right-1]), "\n")
		codeCache.set(filename, line, curCode)
	}

	lambda()
	tccontent.Code(c, curCode)
}

type echoCodeCache struct {
	data sync.Map
}

func (cc *echoCodeCache) get(filename string, line int) (string, bool) {
	key := fmt.Sprintf("%s\t%d", filename, line)
	val, ok := cc.data.Load(key)
	if ok {
		return val.(string), true
	}
	return "", false
}

func (cc *echoCodeCache) set(filename string, line int, code string) {
	key := fmt.Sprintf("%s\t%d", filename, line)
	cc.data.Store(key, code)
}

var codeCache echoCodeCache

// countIndent returns how many tabs line opens with.
func countIndent(line string) int {
	return len(line) - len(strings.TrimLeft(line, "\t"))
}

func removeIndent(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}

	minIndent := -1
	for _, l := range lines {
		if strings.TrimSpace(l) == "" {
			continue
		}

		indent := countIndent(l)
		if minIndent == -1 || minIndent > indent {
			minIndent = indent
		}
	}

	if minIndent == -1 {
		minIndent = 0
	}

	ret := make([]string, len(lines))
	for i, l := range lines {
		if strings.TrimSpace(l) == "" {
			ret[i] = ""
		} else {
			ret[i] = l[minIndent:]
		}
	}
	return ret
}
