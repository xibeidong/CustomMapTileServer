package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/Comdex/imgo"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)
type mysqlInfo struct {
	MysqlDataSourceName string
	MapTableName        string
}
type confInfo struct {
	MysqlConf    *mysqlInfo
	SqliteResourceName string
	ResourcePath string
	SqlType string
}

type roadMapInfo struct {
	IdLevel uint8
	IdDir   uint64
	IdPng   uint64
	ID      string
	ImgData *[]byte
}

var db *sql.DB
var chDone chan int
var chWriteControl chan int
var chSqlWriteControl chan int
var failMap map[string]*roadMapInfo
var maplock sync.Mutex
var mConf *confInfo


func (roadMap *roadMapInfo) toSql(tableName string) bool {
	chSqlWriteControl<-1
	defer func() {
		<-chSqlWriteControl
	}()
	// sqlite3 不支持 on duplicate key update 这样的语法
	//smt, err := db.Prepare("insert into " + tableName + " (img,level_id,dir_id,png_id,id) values (?,?,?,?,?) on duplicate key update img = ?")
	smt, err := db.Prepare("insert into " + tableName + " (img,level_id,dir_id,png_id,id) values (?,?,?,?,?)")

	if err != nil {
		fmt.Println(err)
		return false
	}
	defer smt.Close()

	_, err = smt.Exec(roadMap.ImgData, roadMap.IdLevel, roadMap.IdDir, roadMap.IdPng, roadMap.ID)

	if err != nil {
		fmt.Println(err)
		return false
	} else {
		//fmt.Println(ret.RowsAffected())
		return true
	}

}

func init() {
	dir, e := os.Getwd()
	fmt.Println(dir)
	data, e := ioutil.ReadFile("conf.json")
	if e != nil {
		fmt.Println(e)
	}

	mConf = &confInfo{}
	mConf.MysqlConf = &mysqlInfo{}
	e = json.Unmarshal(data, mConf)
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Println(mConf.MysqlConf, mConf.ResourcePath)
	}
}
func main() {

	signalExit := make(chan int)

	defer db.Close()
	run()
	for {
		select {
		case a := <-signalExit:
			chDone <- -1
			os.Exit(a)
		default:
			time.Sleep(time.Second * 2)
		}
	}
}
func run() {
	failMap = make(map[string]*roadMapInfo)
	chDone = make(chan int, 1000)
	chWriteControl = make(chan int,15)
	chSqlWriteControl = make(chan int , 1)
	go monitor(chDone)
	go retry()
	var e error
	if mConf.SqlType == "mysql"{
		db, e = sql.Open("mysql", mConf.MysqlConf.MysqlDataSourceName)
	}else if mConf.SqlType == "sqlite"{
		db, e = sql.Open("sqlite3", mConf.SqliteResourceName)
	}

	if e != nil {
		fmt.Println(e.Error())
		//os.Exit(0)
	} else {
		fmt.Println("sql is ready ...")
	}

	if mConf.SqlType == "mysql"{
		db.SetMaxOpenConns(300)
	}
	creatTable(mConf.MysqlConf.MapTableName)
	praseImgDir(mConf.ResourcePath)
}
func creatTable(tableName string) bool {
	sqlStr :=fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id BIGINT(21),level_id INT(11),dir_id INT(11),png_id INT(11),img LONGBLOB)",tableName)
	result, err := db.Exec(sqlStr)
	if err!=nil{
		fmt.Println(err)
		return false
	}
	fmt.Println(result)
	return true
}
func retry() {
	for {
		time.Sleep(time.Millisecond * 10)
		maplock.Lock()
		for k, v := range failMap {
			chWriteControl<-1
			b :=  v.toSql(mConf.MysqlConf.MapTableName)
			if b {
				delete(failMap, k)
				chDone <- 2
				break
			} else {
				chDone <- 0
			}
			<-chWriteControl
		}
		maplock.Unlock()
	}

}
func monitor(ch <-chan int) {
	successNum := 0
	failNum := 0

	for {
		c := <-ch
		if c < 0 {
			return
		} else if c == 0 {
			failNum++
			fmt.Println("all fail Num = ", failNum)
			fmt.Println("now fail num = ", len(failMap))
		} else if c == 1 {
			successNum++
			if successNum%100 == 0 {
				fmt.Println(successNum, "has Done,协程数量 = ", runtime.NumGoroutine())
			}
		} else if c == 2 {
			fmt.Println("now fail num = ", len(failMap))
		}
	}
}
func praseImgDir(rootPath string) {

	levelDirs, e := ioutil.ReadDir(rootPath)
	if e != nil {
		fmt.Println(e)
		return
	}

	for _, levelDir := range levelDirs {
		lastDirs, e := ioutil.ReadDir(rootPath + "\\" + levelDir.Name())
		if e != nil {
			fmt.Println(e)
			continue
		}

		for _, lastdir := range lastDirs {

			if lastdir.IsDir() {
				path1 := rootPath + "\\" + levelDir.Name() + "\\" + lastdir.Name()
				files, e1 := ioutil.ReadDir(path1)
				if e1 != nil {
					fmt.Println(e1)
					continue
				}
				//fmt.Println(len(files))

				for _, file := range files {
					//fmt.Println(file.Name())

					//限制最多15个协程，超过会阻塞
					chWriteControl <-1


					path2 := rootPath + "\\" + levelDir.Name() + "\\" + lastdir.Name() + "\\" + file.Name()
					i := strings.Index(file.Name(), ".png")
					if i==-1{

						fmt.Println("删除=》 "+path2)
						err3:=os.Remove(path2)//删除非png格式的文件
						if err3!=nil{
							fmt.Println(err3)
						}
						continue
					}

					str := strings.TrimRight(file.Name(), ".png")

					go imgFile2Sql(
						path2,
						levelDir.Name()+":"+lastdir.Name()+":"+str)
				}
			}
		}

	}

}

func imgFile2Sql(path string, pathKey string) {
	defer func() {
		<- chWriteControl
	}()

	mapinfo := &roadMapInfo{}

	data, err := png2jpg(path) //png压缩成jpg

	if err!=nil{
		data, err = ioutil.ReadFile(path)

		if err != nil {
			fmt.Println(err)
			return
		}

	}


	mapinfo.ImgData = &data
	strs := strings.Split(pathKey, ":")
	if len(strs) == 3 {
		mapinfo.ID = strs[0] + strs[1] + strs[2]
		ret, e := strconv.ParseInt(strs[0], 10, 8)
		if e != nil {
			fmt.Println(e)
			return
		}
		mapinfo.IdLevel = uint8(ret)

		ret, e = strconv.ParseInt(strs[1], 10, 64)
		if e != nil {
			fmt.Println(e)
			return
		}
		mapinfo.IdDir = uint64(ret)

		ret, e = strconv.ParseInt(strs[2], 10, 64)
		if e != nil {
			fmt.Println(e)
			return
		}
		mapinfo.IdPng = uint64(ret)

		b := mapinfo.toSql(mConf.MysqlConf.MapTableName)
		if b {
			chDone <- 1
		} else {
			failMap[mapinfo.ID] = mapinfo
			chDone <- 0

		}
	}
}

func png2jpg(path string) ([]byte,error)  {
	name := strings.TrimRight(path, ".png") + ".jpg"

	imgMatrix := imgo.MustRead(path)
	err := imgo.SaveAsJPEG(name, imgMatrix, 50)
	if err!=nil{
		fmt.Println(err)
		return nil, err
	}
	defer func() {
		err2 := os.Remove(name)
		if err2!=nil{
			fmt.Println(err2)
		}
	}()
	bytes, err := ioutil.ReadFile(name)
	if err!=nil{
		fmt.Println(err)
		return nil,err
	}
	return bytes,nil

}