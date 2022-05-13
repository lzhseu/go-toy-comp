package schema

import (
	"fmt"
	"geeorm/dialect"
	"reflect"
	"testing"
)

func TestReflect(t *testing.T) {
	var x = &Field{
		Name: "name",
		Type: "type",
		Tag:  "tag",
	}
	typ1 := reflect.ValueOf(x).Type()
	typ2 := reflect.Indirect(reflect.ValueOf(x)).Type()
	typ3 := reflect.TypeOf(x)
	fmt.Println(typ1, typ2, typ3)

	f1 := typ1.Kind()
	f2 := typ2.Kind()
	f3 := typ3.Kind()
	fmt.Println(f1, f2, f3)

	typ1.Field(0)
}

type User struct {
	Name string `geeorm:"size(20) PRIMARY KEY"`
	Age  int
}

var TestDial, _ = dialect.GetDialect("mysql")

func TestParse(t *testing.T) {
	schema := Parse(&User{}, TestDial)
	if schema.Name != "User" || len(schema.Fields) != 2 {
		t.Fatal("failed to parse User struct")
	}
	if schema.GetField("Name").Tag != "PRIMARY KEY" {
		t.Fatal("failed to parse primary key")
	}
	fmt.Printf("Name type: %s, Age type: %s\n", schema.GetField("Name").Type, schema.GetField("Age").Type)
	fmt.Printf("Name Tag: %s, Age Tag: %s\n", schema.GetField("Name").Tag, schema.GetField("Age").Tag)
}
