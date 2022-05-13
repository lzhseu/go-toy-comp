package clause

import (
	"reflect"
	"testing"
)

func testSelect(t *testing.T) {
	var clause Clause
	clause.Set(LIMIT, 3)
	clause.Set(SELECT, "user", []string{"*"})
	clause.Set(WHERE, "name = ?", "Tom")
	clause.Set(ORDERBY, "age ASC")
	stat, vars := clause.Build(SELECT, WHERE, ORDERBY, LIMIT)
	t.Log(stat, vars)
	if stat != "SELECT * FROM user WHERE name = ? ORDER BY age ASC LIMIT ?" {
		t.Fatal("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{"Tom", 3}) {
		t.Fatal("failed to build SQLVars")
	}
}

func TestClause_Build(t *testing.T) {
	t.Run("select", func(t *testing.T) {
		testSelect(t)
	})
}

func TestOther(t *testing.T) {
	var s []interface{}
	if s == nil {
		t.Log("nil===")
	}
	t.Logf("%v", s)
}
