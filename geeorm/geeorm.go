package geeorm

import (
	"database/sql"
	"fmt"
	"geeorm/dialect"
	"geeorm/log"
	"geeorm/session"
	"strings"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver, dataSource string) (e *Engine, err error) {
	db, err := sql.Open(driver, dataSource)
	if err != nil {
		log.Error(err)
		return
	}

	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}

	dial, ok := dialect.GetDialect(driver)
	if !ok {
		err = fmt.Errorf("dial %s Not Found", driver)
		log.Errorf(err.Error())
		return
	}
	e = &Engine{db: db, dialect: dial}
	log.Info("Connect database success")
	return
}

func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

// TxFunc will be called between tx.Begin() and tx.Commit()
type TxFunc func(s *session.Session) (interface{}, error)

// Transaction todo: mysql 中表操作（如创建表）都是默认自动提交的，所以在事务中并没有用。后续看怎么优化这一点
func (e *Engine) Transaction(fn TxFunc) (result interface{}, err error) {
	s := e.NewSession()
	if err = s.Begin(); err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil { // TxFun may panic
			_ = s.Rollback()
			panic(p)
		} else if err != nil { // TxFun may error
			_ = s.Rollback()
		} else {
			defer func() {
				if err != nil {
					_ = s.Rollback()
				}
			}()
			err = s.Commit()
		}
	}()
	return fn(s)
}

func (e *Engine) Migrate(value interface{}) error {
	_, err := e.Transaction(func(s *session.Session) (result interface{}, err error) {
		if !s.Model(value).HasTable() {
			log.Infof("table %s doesn't exist", s.RefTable().Name)
			return nil, s.CreateTable()
		}
		table := s.RefTable()
		rows, _ := s.Raw(fmt.Sprintf("SELECT * FROM %s LIMIT 1", table.Name)).QueryRows()
		columns, _ := rows.Columns()
		addCols := difference(table.FieldNames, columns) // new_table - old_table = add_cols
		delCols := difference(columns, table.FieldNames) // old_table - new_table = delete_cols
		log.Infof("added cols %v, deleted cols %v\n", addCols, delCols)

		for _, col := range addCols {
			field := table.GetField(col)
			sqlStr := fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table.Name, field.Name, field.Type)
			if _, err = s.Raw(sqlStr).Exec(); err != nil {
				return
			}
		}

		if len(delCols) == 0 {
			return
		}

		tmp := "tmp_" + table.Name
		fieldStr := strings.Join(table.FieldNames, ", ")

		defer func() {
			if p := recover(); p != nil {
				return
			}
		}()

		_, err = s.Raw(fmt.Sprintf("CREATE TABLE %s AS SELECT %s FROM %s;", tmp, fieldStr, table.Name)).Exec()
		checkErr(err)
		_, err = s.Raw(fmt.Sprintf("DROP TABLE %s;", table.Name)).Exec()
		checkErr(err)
		_, err = s.Raw(fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tmp, table.Name)).Exec()
		checkErr(err)
		return
	})
	return err
}

// difference returns a - b
func difference(a, b []string) (diff []string) {
	mapB := make(map[string]bool)
	for _, v := range b {
		mapB[v] = true
	}

	for _, v := range a {
		if _, ok := mapB[v]; !ok {
			diff = append(diff, v)
		}
	}
	return
}

func checkErr(err error) {
	if err != nil {
		log.Error(err)
		panic(err)
	}
}
