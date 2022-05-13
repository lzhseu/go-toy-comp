// Package schema completes the mapping of go structure to SQL table
package schema

import (
	"fmt"
	"geeorm/dialect"
	"geeorm/log"
	"go/ast"
	"reflect"
	"strconv"
	"strings"
)

const (
	Tag     = "geeorm"
	TagSize = "size"
)

// Schema represents a table of database
type Schema struct {
	Model      interface{} // origin go value
	Name       string      // table name
	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field
}

// Field represents a column of data table
type Field struct {
	Name string
	Type string
	Tag  string
}

func (s *Schema) GetField(name string) *Field {
	return s.fieldMap[name]
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup(Tag); ok {
				tags := strings.Split(v, " ")
				var attr []string
				for _, tag := range tags {
					if strings.HasPrefix(tag, TagSize) {
						size, err := strconv.Atoi(tag[len(TagSize)+1 : len(tag)-1])
						if err != nil {
							log.Errorf("column(%s) size must be integer", field.Name)
							panic(err)
						}
						field.Type = fmt.Sprintf(field.Type, size)
					} else {
						attr = append(attr, tag)
					}
					field.Tag = strings.Join(attr, " ")
				}
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}

	return schema
}

// RecordValue converts value from object to database column in order
func (s *Schema) RecordValue(dest interface{}) (fieldValues []interface{}) {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	for _, field := range s.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return
}
