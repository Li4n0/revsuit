package notice

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/li4n0/revsuit/internal/record"
)

func formatRecordField(r record.Record,fieldFormat string) (content string) {
	structType := reflect.ValueOf(r)
	for i := 0; i < structType.NumField(); i++ {
		structField := structType.Type().Field(i)
		fieldName := structField.Name
		tag := structField.Tag
		label := tag.Get("notice")
		if label == "-" {
			continue
		} else if label != "" {
			fieldName = label
		}
		value := structType.Field(i).Interface()
		if value, ok := value.(record.Record); ok {
			content += formatRecordField(value, fieldFormat)
			continue
		}
		content += fmt.Sprintf(fieldFormat+"\n", strings.ToUpper(fieldName), value)
	}
	strings.TrimSuffix(content, "\n")
	return
}
