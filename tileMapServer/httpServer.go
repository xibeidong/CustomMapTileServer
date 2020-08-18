package tileMapServer

import (
	"My/CustomTileMapServer/common/dbMysql"
	"My/CustomTileMapServer/conf"
	"My/Learn/zapLog"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"os"
	"strings"
)

var sqliteDBMap map[string] *sql.DB
func Run() {
	sqliteDBMap = make(map[string] *sql.DB)

	var rootPath string
	if conf.MyConfs.MapResourceDBsRootPath !=""{
		rootPath = conf.MyConfs.MapResourceDBsRootPath
	}else {
		rootPath, _ = os.Getwd()
	}

	for _,v:=range conf.MyConfs.MapResourceDBs{
		sourceName:=rootPath+"/"+v
		db, err := sql.Open("sqlite3", sourceName)
		if err!=nil {
			zapLog.Logger.Error(err)
			return
		}
		sqliteDBMap[sourceName] = db
	}
	startHttpServer()
}


func startHttpServer() {
	http.HandleFunc("/roadmap", roadMapTile)
	http.HandleFunc("/roadmap/",roadMapTile2)
	http.HandleFunc("/sqlite3_roadmap",sqlite3RoadMapTile)
	http.ListenAndServe(conf.MyConfs.MapHttpServerListen, nil)

}
func sqlite3RoadMapTile(w http.ResponseWriter, r *http.Request)  {
	fmt.Println(r.URL.RawQuery)

	split := strings.Split(r.URL.RawQuery,"?")
	data := getTileFromSqlite3(split[0], split[1], split[2])
	if data == nil {
		//fmt.Println("data is nil :" + strs[2])

		w.WriteHeader(400)
	} else {
		w.WriteHeader(200)
		w.Write(data)
		//fmt.Println(r.URL)
	}
}
func roadMapTile2(w http.ResponseWriter, r *http.Request){
	strs:= strings.Split(r.URL.String(),"/")
	data := getTileData(strs[2])
	if data == nil {
		//fmt.Println("data is nil :" + strs[2])

		w.WriteHeader(400)
	} else {
		w.WriteHeader(200)
		w.Write(data)
		//fmt.Println(r.URL)
	}

}
func roadMapTile(w http.ResponseWriter, r *http.Request) {

	data := getTileData(r.URL.RawQuery)
	if data == nil {
		//fmt.Println("data is nil :" + r.URL.RawQuery)
		w.WriteHeader(400)
	} else {

		w.WriteHeader(200)
		w.Write(data)
		//fmt.Println(r.URL.RawQuery)
	}

}
func getTileFromSqlite3(z,x,y string) []byte  {

	sqlStr:= fmt.Sprintf( "select image from tiles where x=%s and y=%s and z=%s ",x,y,z)

	//查找每一个db文件
	for _,db := range sqliteDBMap{
		rows, err := db.Query(sqlStr)
		if err!=nil{
			zapLog.Logger.Error(err)
			return nil
		}
		for rows.Next() {
			var data []byte
			err = rows.Scan(&data)
			if err!=nil{
				zapLog.Logger.Error(err)
				return nil
			}
			return data
		}

	}
	return nil
}
func getTileData(id string) []byte {
	smt, e := dbMysql.Db.Prepare("select img from " + conf.MyConfs.MysqlConf.MapTableName + " where id=?")
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