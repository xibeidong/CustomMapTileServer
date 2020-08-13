package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type mysqlInfo struct {
	MysqlDataSourceName string `json:"mysqlDataSourceName"`
	MapTableName        string `json:"mapTableName"`
}
type gpsIdInfo struct {
	Id int64 `json:"id"`
}

type infos struct {
	MysqlConf             mysqlInfo   `json:"mysqlConf"`
	MapHttpServerListen   string      `json:"mapHttpServerListen"`
	QdPortUdpServerListen string      `json:"qdPortUdpServerListen"`
	QdPortTcpServerListen string      `json:"qdPortTcpServerListen"`
	GpsIds                []gpsIdInfo `json:"gpsIds"`
}

var MyConfs *infos

func Init() error{

	MyConfs = &infos{}
	data, e := ioutil.ReadFile("../conf/conf.json")
	if e != nil {
		fmt.Println(e)
		return e
	}
	e = json.Unmarshal(data, MyConfs)
	if e != nil {
		fmt.Println(e)
		return e
	} else {
		fmt.Println("读取配置 ==》 ",MyConfs.MysqlConf)
		fmt.Println("读取配置 ==》 ",MyConfs.MapHttpServerListen)
		fmt.Println("读取配置 ==》 ",MyConfs.QdPortUdpServerListen)
		fmt.Println("读取配置 ==》 ",MyConfs.GpsIds)
	}

	return nil
}