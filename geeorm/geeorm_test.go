package geeorm

import (
	"errors"
	"geeorm/session"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"testing"
)

func OpenDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("mysql", "root:sy091314@/go")
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	return engine
}

func TestNewEngine(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
}

type User struct {
	Name string `geeorm:"size(20) PRIMARY KEY"`
	Age  int
}

func transactionRollback(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_ = s.Model(&User{}).CreateTable()
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		_, err = s.Insert(&User{Name: "Tom", Age: 18})
		_, err = s.Insert(&User{Name: "Sam", Age: 20})
		return nil, errors.New("error test")
	})
	if err == nil {
		t.Fatal("failed to rollback")
	}
}

func transactionCommit(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_ = s.Model(&User{}).CreateTable()
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		_, err = s.Insert(&User{Name: "Tom", Age: 18})
		_, err = s.Insert(&User{Name: "Sam", Age: 20})
		return
	})
	u := &User{}
	_ = s.First(u)
	if err != nil || u.Name != "Sam" {
		t.Fatal("failed to commit")
	}
}

func TestEngine_Transaction(t *testing.T) {
	t.Run("rollback", func(t *testing.T) {
		transactionRollback(t)
	})
	t.Run("commit", func(t *testing.T) {
		transactionCommit(t)
	})
}

func TestEngine_Migrate(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS User;").Exec()
	_, _ = s.Raw("CREATE TABLE User(Name varchar(20) PRIMARY KEY, XXX integer);").Exec()
	_, _ = s.Raw("INSERT INTO User(`Name`) values (?), (?)", "Tom", "Sam").Exec()
	_ = engine.Migrate(&User{})

	rows, _ := s.Raw("SELECT * FROM User").QueryRows()
	columns, _ := rows.Columns()
	if !reflect.DeepEqual(columns, []string{"Name", "Age"}) {
		t.Fatal("Failed to migrate table User, got columns", columns)
	}
}
