package main

import (
	"My/CustomTileMapServer/common"
	"My/Learn/zapLog"
	_ "My/Learn/zapLog"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type httpConf struct {
	MysqlConf        *common.MysqlInfo
	HttpServerListen string
}

var mConf *httpConf

func init() {
	mConf = &httpConf{&common.MysqlInfo{}, ""}
	data, e := ioutil.ReadFile("../conf/conf.json")
	if e != nil {
		fmt.Println(e)
	}
	e = json.Unmarshal(data, mConf)
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Println(mConf.MysqlConf, mConf.HttpServerListen)
	}
}
func logTest()  {
	for{

		zapLog.Logger.Error("error....")
		zapLog.Logger.Info("info....")
		zapLog.Logger.Warn("warn....")
		time.Sleep(time.Second*5)
	}
}

func main() {
	logTest()
	zapLog.Logger.Info("Run")
	signalExit := make(chan int)
	run()

	for {
		select {
		case a := <-signalExit:
			os.Exit(a)
		default:
			time.Sleep(time.Second * 2)
		}
	}
}

func run() {
	db, e := sql.Open("mysql", mConf.MysqlConf.MysqlDataSourceName)
	if e != nil {
		fmt.Println(e.Error())
	} else {
		fmt.Println("mysql open")
	}
	db.SetMaxOpenConns(300)
	db.SetConnMaxLifetime(time.Minute * 1)
	db.SetMaxIdleConns(300)

	common.DbMySql = db

	defer db.Close()

	startHttpServer()
}

func startHttpServer() {
	http.HandleFunc("/roadmap", roadMapTile)
	http.HandleFunc("/roadmap/",roadMapTile2)
	http.ListenAndServe(mConf.HttpServerListen, nil)

}
func roadMapTile2(w http.ResponseWriter, r *http.Request){
	strs:= strings.Split(r.URL.String(),"/")
	data := getTileData(strs[2])
	if data == nil {
		fmt.Println("data is nil :" + strs[2])

		w.WriteHeader(400)
	} else {

		w.WriteHeader(200)
		w.Write(data)
		fmt.Println(r.URL.RawQuery)
	}

}
func roadMapTile(w http.ResponseWriter, r *http.Request) {

	data := getTileData(r.URL.RawQuery)
	if data == nil {
		fmt.Println("data is nil :" + r.URL.RawQuery)
		w.WriteHeader(400)
	} else {

		w.WriteHeader(200)
		w.Write(data)
		fmt.Println(r.URL.RawQuery)
	}

}

func getTileData(id string) []byte {
	smt, e := common.DbMySql.Prepare("select img from " + mConf.MysqlConf.MapTableName + " where id=?")
	if e != nil {
		fmt.Println(e)
		return nil
	}
	defer smt.Close()
	rows, e := smt.Query(id)
	if e != nil {
		fmt.Println(e)
		return nil
	}
	defer rows.Close()
	var data []byte
	if rows.Next() {
		e = rows.Scan(&data)
		if e != nil {
			fmt.Println(e)
		}
		return data
	}

	return nil
}
