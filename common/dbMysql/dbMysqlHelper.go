package dbMysql

import (
	"My/Learn/zapLog"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/cron/v3"
	"time"
)

var Db *sql.DB

func Init(url string)  {
	db, e := sql.Open("mysql", url)
	if e != nil {
		zapLog.Logger.Error(e)
	} else {
		zapLog.Logger.Info("mysql ready connect...")
	}

	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Minute * 1)
	db.SetMaxIdleConns(25)
	Db = db

	//开启定时建表任务
	creatTableTimerTask("gpsinfo")
	//fmt.Println(time.Now().Add(time.Hour*24).Format("2006-01-02 15:04:05"))

}

func Close()  {
	Db.Close()
}

//定时建表任务
func creatTableTimerTask(templatename string)  {
	creatTable(templatename,time.Now())
	c := cron.New(cron.WithSeconds())
	//每月28号执行
	c.AddFunc("0 0 0 28 * *", func() {
		creatTable(templatename,time.Now().Add(time.Hour*24*7))
	})
}

func creatTable(templatename string ,t time.Time)  {

	tablename:=fmt.Sprintf("%s%02d%02d",templatename,t.Year(),t.Month())
	strSql:= fmt.Sprintf("create table IF NOT EXISTS  %v like %v",tablename,templatename)
	result, err := Db.Exec(strSql)

	if err!=nil{
		zapLog.Logger.Error(err)
		return
	}
	fmt.Print(strSql," ==> ")
	fmt.Println(result.RowsAffected())
}