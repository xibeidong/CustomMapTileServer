package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
)
//测试使用
func main() {

	dir, err2 := os.Getwd() //这样获取系统目录有时候出错误，遇到再说
	if err2!=nil {
		fmt.Println(err2)
	}
	fmt.Println(dir)

	_, err2 = os.Stat(dir+"\\main.sqlitedb")
	if err2!=nil{
		fmt.Println(err2)
	}

	//sqlite3 不能用相对路径，会出奇怪的错误
	db, err := sql.Open("sqlite3", dir+"\\main.sqlitedb")
	if err!=nil{
		fmt.Println(err)
	}
	defer db.Close()

	rows, err := db.Query("select x,y,z from tiles where x<10")
	//rows, err := db.Query("select * from info ")
	if err!=nil{
		fmt.Println(err)
		return
	}
	for rows.Next(){
		var x,y,z int
		rows.Scan(&x,&y,&z)
		fmt.Println(x,":",y,":",z)
	}
	select {

	}
}
