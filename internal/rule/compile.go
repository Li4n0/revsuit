package rule

import (
	"strings"
)

// CompileTpl receive []byte or string type tpl and variables map.
// Return the template after variable substitution
func CompileTpl(tpl interface{}, vars map[string]string) (compiled string) {
	switch v := tpl.(type) {
	case string:
		compiled = v
	case []byte:
		compiled = string(v)
	}
	for n, v := range vars {
		compiled = strings.ReplaceAll(compiled, "${"+n+"}", v)
	}
	return compiled
}
