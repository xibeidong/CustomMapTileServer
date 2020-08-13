package main

import (
	"My/CustomTileMapServer/QDPort"
	"My/CustomTileMapServer/common/dbMysql"
	"My/CustomTileMapServer/conf"
	"My/CustomTileMapServer/tileMapServer"
	"My/Learn/zapLog"
	_ "My/Learn/zapLog"
	"os"
	"time"
)

var signalExit chan int

func logTest()  {
	for{

		zapLog.Logger.Error("error....")
		zapLog.Logger.Info("info....")
		zapLog.Logger.Warn("warn....")
		time.Sleep(time.Second*5)
	}
}

func main() {
	//logTest()
	zapLog.Logger.Info("Run")
	signalExit = make(chan int)

	 e:= conf.Init()
	 if e!=nil{
	 	stop()
	 }
	dbMysql.Init(conf.MyConfs.MysqlConf.MysqlDataSourceName)
	go QDPort.StartUdpServer()
	go QDPort.StartTcpServer()
	go tileMapServer.Run()

	stop()
}

func stop()  {
	for {
		select {
		case a := <-signalExit:
			os.Exit(a)
		default:
			time.Sleep(time.Second * 2)
		}
	}
}
