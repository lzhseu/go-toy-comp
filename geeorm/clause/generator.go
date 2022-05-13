package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

func genBindVars(num int) string {
	var builder strings.Builder
	for i := 0; i < num; i++ {
		builder.WriteString("?, ")
	}
	return builder.String()[:builder.Len()-2]
}

// values[0]: table name
// values[1]: columns
func _insert(values ...interface{}) (stat string, vars []interface{}) {
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	stat = fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields)
	return
}

// Supported batch insert. Each value in values represents one record
func _values(values ...interface{}) (stat string, vars []interface{}) {
	var bindStr string
	var builder strings.Builder
	builder.WriteString("VALUES ")
	for _, value := range values {
		v := value.([]interface{})
		if bindStr == "" {
			bindStr = genBindVars(len(v))
		}
		builder.WriteString(fmt.Sprintf("(%v), ", bindStr))
		vars = append(vars, v...)
	}
	stat = builder.String()[:builder.Len()-2]
	return
}

// values[0]: table name
// values[1]: select columns
func _select(values ...interface{}) (stat string, vars []interface{}) {
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	stat = fmt.Sprintf("SELECT %v FROM %s", fields, tableName)
	return
}

func _limit(values ...interface{}) (stat string, vars []interface{}) {
	return "LIMIT ?", values
}

func _where(values ...interface{}) (stat string, vars []interface{}) {
	desc := values[0]
	return fmt.Sprintf("WHERE %s", desc), values[1:]
}

func _orderBy(values ...interface{}) (stat string, vars []interface{}) {
	stat = fmt.Sprintf("ORDER BY %s", values[0])
	return
}

// values[0]: table name
// values[1]: a map stores key-value pairs to be updated
func _update(values ...interface{}) (stat string, vars []interface{}) {
	tableName := values[0]
	m := values[1].(map[string]interface{})
	var keys []string
	for k, v := range m {
		keys = append(keys, k+" = ?")
		vars = append(vars, v)
	}
	stat = fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(keys, ", "))
	return
}

func _delete(values ...interface{}) (stat string, vars []interface{}) {
	stat = fmt.Sprintf("DELETE FROM %s", values[0])
	return
}

func _count(values ...interface{}) (stat string, vars []interface{}) {
	return _select(values[0], []string{"count(*)"})
}

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[VALUES] = _values
	generators[SELECT] = _select
	generators[LIMIT] = _limit
	generators[WHERE] = _where
	generators[ORDERBY] = _orderBy
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
}
