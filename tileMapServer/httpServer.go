package tileMapServer

import (
	"My/CustomTileMapServer/common/dbMysql"
	"My/CustomTileMapServer/conf"
	"fmt"
	"net/http"
	"strings"
)

func Run() {
	startHttpServer()
}

func startHttpServer() {
	http.HandleFunc("/roadmap", roadMapTile)
	http.HandleFunc("/roadmap/",roadMapTile2)
	http.ListenAndServe(conf.MyConfs.MapHttpServerListen, nil)

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