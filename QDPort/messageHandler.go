package QDPort

import (
	"My/CustomTileMapServer/common"
	"My/CustomTileMapServer/common/dbMysql"
	"My/Learn/zapLog"

	"encoding/json"
	"fmt"
	"net"
	"time"
)

//type gpsinfo struct {
//	GpsId int64
//	Lng float64
//	Lat float64
//	T string
//}

type playBackPositions struct {
	GpsId     int64
	TBegin    string
	TEnd      string
	Positions []gpsInfo
}
type gpsPositions struct {
	Positions []gpsInfo
}

func playbackHandler(conn *net.TCPConn,data []byte,messageId uint16)  {
	pb:=&playBackPositions{}
	err := json.Unmarshal(data, pb)
	if err!=nil{
		zapLog.Logger.Error(err)
	}
	//ToDo
	getPositionsByDb(pb)

	replyDataBody, err := json.Marshal(pb)
	if err!=nil{
		zapLog.Logger.Error(err)
		return
	}

	reply(conn,replyDataBody,messageId)

}
func heartHandler(conn *net.TCPConn,data []byte,messageId uint16)  {

}



func getPositionsByDb(pb *playBackPositions)  {
	begin, err := time.Parse(common.TimeFormat, pb.TBegin)
	if err!=nil{
		zapLog.Logger.Error(err)
		return
	}
	tableName1:=fmt.Sprintf("gpsinfo%02d%02d", begin.Year(), begin.Month())
	end, err := time.Parse(common.TimeFormat, pb.TEnd)
	if err!=nil{
		zapLog.Logger.Error(err)
		return
	}
	tableName2:=fmt.Sprintf("gpsinfo%02d%02d", end.Year(), end.Month())

	var sqlStr string
	if end.Month() == begin.Month() {
		sqlStr = fmt.Sprintf("select * from %s where gpsid = %v and t > '%s' and t <'%s'",tableName1,pb.GpsId,
		pb.TBegin,pb.TEnd)
	}else{
		sqlStr = fmt.Sprintf("(select * from %s where gpsid = %v and t > '%s' and t <'%s') union (select * from %s where gpsid = %v and t> '%s' and t<'%s')",tableName1,pb.GpsId,
			pb.TBegin,pb.TEnd,tableName2,pb.GpsId,pb.TBegin,pb.TEnd)
	}

	rows, err := dbMysql.Db.Query(sqlStr)
	if err!=nil{
		zapLog.Logger.Error(err)
		return
	}
	defer rows.Close()
	for rows.Next(){
		info:=&gpsInfo{}
		rows.Scan(&info.Id,&info.Lng,&info.Lat,&info.T)
		pb.Positions = append(pb.Positions, *info)
	}

	//fmt.Println("读取到定位点数量 = " ,len(pb.Positions))
}