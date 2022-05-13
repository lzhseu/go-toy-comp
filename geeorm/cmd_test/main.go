package main

import (
	"fmt"
	"geeorm"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	engine, _ := geeorm.NewEngine("mysql", "root:sy091314@/go")
	defer engine.Close()
	s := engine.NewSession()
	result, _ := s.Raw("INSERT INTO user(`Name`, `age`) values (?, ?), (?, ?)", "Sally", 22, "James", 19).Exec()
	count, _ := result.RowsAffected()
	fmt.Printf("Exec success, %d affected\n", count)
}
