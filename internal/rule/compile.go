package rule

import "strings"

func CompileTpl(tpl string, vars map[string]string) (compiled string) {
	compiled = tpl
	for n, v := range vars {
		compiled = strings.ReplaceAll(compiled, "${"+n+"}", v)
	}

	return compiled
}
