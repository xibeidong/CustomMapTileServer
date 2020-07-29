package main

import (
	"My/CustomTileMapServer/common"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type confInfo struct {
	MysqlConf    *common.MysqlInfo
	ResourcePath string
}

var chDone chan int
var failMap map[string]*common.RoadMapInfo
var maplock sync.Mutex
var mConf *confInfo

func init() {
	dir, e := os.Getwd()
	fmt.Println(dir)
	data, e := ioutil.ReadFile("conf.json")
	if e != nil {
		fmt.Println(e)
	}

	mConf = &confInfo{&common.MysqlInfo{}, ""}
	e = json.Unmarshal(data, mConf)
	if e != nil {
		fmt.Println(e)
	} else {
		fmt.Println(mConf.MysqlConf, mConf.ResourcePath)
	}
}
func main() {

	signalExit := make(chan int)

	defer func() {
		if common.DbMySql != nil {
			common.DbMySql.Close()
		}
	}()
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
	failMap = make(map[string]*common.RoadMapInfo)
	chDone = make(chan int, 1000)
	go monitor(chDone)
	go retry()
	db, e := sql.Open("mysql", mConf.MysqlConf.MysqlDataSourceName)
	if e != nil {
		fmt.Println(e.Error())
		//os.Exit(0)
	} else {
		fmt.Println("mysql open")
	}
	db.SetMaxOpenConns(300)
	common.DbMySql = db
	praseImgDir(mConf.ResourcePath)
}
func retry() {
	for {
		time.Sleep(time.Millisecond * 10)
		maplock.Lock()
		for k, v := range failMap {
			b := v.ToMySql(mConf.MysqlConf.MapTableName)
			if b {
				delete(failMap, k)
				chDone <- 2
				break
			} else {
				chDone <- 0
			}
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
			if successNum%10 == 0 {
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

		//doImgFile2Mysql(rootPath,lastDirs,levelDir.Name())
		for _, lastdir := range lastDirs {
			if lastdir.IsDir() {
				path1 := rootPath + "\\" + levelDir.Name() + "\\" + lastdir.Name()
				files, e1 := ioutil.ReadDir(path1)
				if e1 != nil {
					fmt.Println(e1)
					continue
				}
				fmt.Println(len(files))

				for _, file := range files {
					name1 := file.Name()
					fmt.Println(name1)
					str := strings.TrimRight(file.Name(), ".png")
					path2 := rootPath + "\\" + levelDir.Name() + "\\" + lastdir.Name() + "\\" + file.Name()
					imgFile2mysql(
						path2,
						levelDir.Name()+":"+lastdir.Name()+":"+str)
				}
			}
		}

	}
}

//遍历level文件夹,弃用
func doImgFile2Mysql(rootPath string, dirs []os.FileInfo, level string) {

	for _, dir := range dirs {
		if dir.IsDir() {
			path1 := rootPath + "\\" + level + "\\" + dir.Name()
			files, e1 := ioutil.ReadDir(path1)
			if e1 != nil {
				fmt.Println(e1)
				continue
			}
			fmt.Println(len(files))

			for _, file := range files {
				name1 := file.Name()
				fmt.Println(name1)
				str := strings.TrimRight(file.Name(), ".jpg")
				path2 := rootPath + "\\" + level + "\\" + dir.Name() + "\\" + file.Name()
				imgFile2mysql(
					path2,
					level+":"+dir.Name()+":"+str)
			}
		}
	}
}
func imgFile2mysql(path string, pathKey string) {
	mapinfo := &common.RoadMapInfo{}

	data, e := ioutil.ReadFile(path)
	if e != nil {
		fmt.Println(e)
		return
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
		mapinfo.ID_level = uint8(ret)

		ret, e = strconv.ParseInt(strs[1], 10, 64)
		if e != nil {
			fmt.Println(e)
			return
		}
		mapinfo.ID_dir = uint64(ret)

		ret, e = strconv.ParseInt(strs[2], 10, 64)
		if e != nil {
			fmt.Println(e)
			return
		}
		mapinfo.ID_png = uint64(ret)

		b := mapinfo.ToMySql(mConf.MysqlConf.MapTableName)
		if b {
			chDone <- 1
		} else {
			failMap[mapinfo.ID] = mapinfo
			chDone <- 0

		}
	}
}
