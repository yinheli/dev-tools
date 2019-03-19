package table

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	numberSequence    = regexp.MustCompile(`([a-zA-Z])(\d+)([a-zA-Z]?)`)
	numberReplacement = []byte(`$1 $2 $3`)
)

func TitleCase(str string) string {
	return toCamelCase(str, true)
}

func CamelCase(str string) string {
	return toCamelCase(str, false)
}

func DataType(dataType string, nullable bool, typeMapping map[string]string) string {
	dataType = strings.ToLower(strings.TrimSpace(dataType))

	goType := "string"

	dataType = strings.ToLower(strings.TrimSpace(dataType))

	newType := dataType
	bracketIndex := strings.Index(newType, "(")
	if bracketIndex > 0 {
		newType = newType[0:bracketIndex]
	}

	if strings.Contains(dataType, "unsigned") {
		newType = "u" + newType
	}

	switch newType {
	case "int", "tinyint":
		goType = "int32"
	case "uint", "utinyint":
		goType = "uint32"
	case "bigint":
		goType = "int64"
	case "ubigint":
		goType = "uint64"
	case "date", "datetime", "timestamp":
		goType = "time.Time"
	case "float", "decimal", "double":
		goType = "float64"
	}

	if v, ok := typeMapping[goType]; ok {
		goType = v
	}

	if nullable {
		return fmt.Sprintf("*%s", goType)
	}

	return goType
}

func addWordBoundariesToNumbers(s string) string {
	b := []byte(s)
	b = numberSequence.ReplaceAll(b, numberReplacement)
	return string(b)
}

// Converts a string to CamelCase
func toCamelCase(s string, first bool) string {
	s = addWordBoundariesToNumbers(s)
	s = strings.Trim(s, " ")
	n := ""
	capNext := first
	for _, v := range s {
		if v >= 'A' && v <= 'Z' {
			n += string(v)
		}
		if v >= '0' && v <= '9' {
			n += string(v)
		}
		if v >= 'a' && v <= 'z' {
			if capNext {
				n += strings.ToUpper(string(v))
			} else {
				n += string(v)
			}
		}
		if v == '_' || v == ' ' || v == '-' {
			capNext = true
		} else {
			capNext = false
		}
	}
	return n
}

func Tag(column *Column) string {
	jsonTag := fmt.Sprintf(`json:"%s"`, JsonTag(column.CamelCaseName, column.GoType))
	validateTag := ""
	if !column.Nullable {
		switch column.TitleCaseName {
		case "Id", "CreatedAt", "UpdatedAt":
		default:
			validateTag = ` validate:"required"`
		}
	}

	return fmt.Sprintf("`%s%s`", jsonTag, validateTag)
}

func JsonTag(colunm string, goType string) string {
	switch goType {
	case "uint32", "int64", "uint64",
		"*uint32", "*int64", "*uint64":
		colunm = fmt.Sprintf("%s,string", colunm)
	}
	return colunm
}
