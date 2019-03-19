package table

import (
	"bytes"
	"github.com/yinheli/dev-tools/pkg/database"
	"golang.org/x/tools/imports"
	"text/template"
)

type (
	Table struct {
		Name          string
		TitleCaseName string
		Comment       string
		Columns       []*Column
	}

	Column struct {
		Name          string
		Comment       string
		ColumnType    string
		Nullable      bool
		TitleCaseName string
		CamelCaseName string
		GoType        string
		Tag           string
	}
)

var (
	tpl = template.Must(template.New("struct").Parse(`
package model

// {{.TitleCaseName}} table: {{.Name}} {{.Comment}}
type {{.TitleCaseName}} struct {
    {{range .Columns -}}
        {{- .TitleCaseName}} {{.GoType}} {{.Tag}}
    {{end -}}
}
`))
)

func ToGo(table string) (string, error) {
	var t Table
	err := database.DB.Raw(`select t.table_name name, t.table_comment comment from information_schema.tables t where t.table_schema = database() and t.table_name=?`, table).Scan(&t).Error
	if err != nil {
		return "", err
	}
	var cols []*Column
	err = database.DB.Raw(`select column_name name, column_type, column_comment comment, if(lower(is_nullable)='yes', 1, 0) is_nullable from information_schema.columns t where t.table_schema=database() and t.table_name=?`, table).Find(&cols).Error
	if err != nil {
		return "", err
	}

	t.Columns = cols

	t.TitleCaseName = TitleCase(t.Name)
	for _, c := range t.Columns {
		c.TitleCaseName = TitleCase(c.Name)
		c.CamelCaseName = CamelCase(c.Name)
		c.GoType = DataType(c.ColumnType, c.Nullable, map[string]string{})
		c.Tag = Tag(c)
	}

	var b bytes.Buffer
	err = tpl.Execute(&b, &t)
	if err != nil {
		return "", err
	}

	data, err := imports.Process("", b.Bytes(), nil)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
