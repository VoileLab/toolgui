package tcmisc

import (
	"fmt"
	"go/parser"
	"runtime"
	"strings"
	"sync"

	"github.com/mudream4869/toolgui/toolgui/component"
	"github.com/mudream4869/toolgui/toolgui/framework"
)

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

func countIndent(line string) int {
	cnt := 0
	for _, c := range line {
		if c == '\t' {
			cnt++
		} else {
			return cnt
		}
	}
	return 0
}

func removeIndent(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}

	minIndent := countIndent(lines[0])
	for _, l := range lines {
		minIndent = min(minIndent, countIndent(l))
	}

	ret := make([]string, len(lines))
	for i, l := range lines {
		ret[i] = l[minIndent:]
	}
	return ret
}

func Echo(c *framework.Container, code string, lambda func()) {
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
	component.Code(c, curCode, "go")
}
