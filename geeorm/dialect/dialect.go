package dialect

import "reflect"

var dialectMap = make(map[string]Dialect)

type Dialect interface {
	DataTypeOf(typ reflect.Value) string                    // mapping from go type to sql type
	TableExistSQL(tableName string) (string, []interface{}) // returns the SQL statement to determine whether a table exists
}

func RegisterDialect(name string, dialect Dialect) {
	dialectMap[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectMap[name]
	return
}
