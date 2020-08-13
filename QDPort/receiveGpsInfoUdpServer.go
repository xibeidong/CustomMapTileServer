package QDPort

import (
	"My/CustomTileMapServer/common/dbMysql"
	"My/CustomTileMapServer/conf"
	"My/Learn/zapLog"
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"strconv"
	"sync"
	"time"
)

type gpsInfo struct {
	Id int64 	`json:"id`
	Lng float64
	Lat float64
	T string
}
var GpsInfoMap sync.Map
func StartUdpServer()  {

	initMap()

	addr,e:= net.ResolveUDPAddr("udp",conf.MyConfs.QdPortUdpServerListen)
	if e!=nil {
		zapLog.Logger.Error(e)
		return
	}
	conn,e:= net.ListenUDP("udp",addr)
	if e!=nil{
		zapLog.Logger.Error(e)
		return
	}
	defer conn.Close()

	for{
		data:=make([]byte,512)
		udplen, udpAddr, e := conn.ReadFromUDP(data)
		if e!=nil{
			zapLog.Logger.Error(e)
			return
		}
		fmt.Println("len=",udplen,"UdpServer receive message from ",udpAddr)
		go resolve(data[:udplen])
	}

}
func initMap()  {
	//GpsInfoMap = make(sync.Map[int64] *gpsinfo)
	for _,v:= range conf.MyConfs.GpsIds{
		GpsInfoMap.Store(v.Id,&gpsInfo{})
	}
}

func resolve(data []byte)  {
	//info:=&gpsInfo{}
	if len(data)>=44{
		if data[0]==0xaa && data[1]==0x00{ //标志头
			if data[2]==0xcc && data[3]==0x00{//0xcc 表示GPS信息
				//6个字节 GPSID
				idStr:=fmt.Sprintf("%d%02d%02d%02d%02d%02d",data[4],data[5],data[6],data[7],data[8],data[9])
				id,e:= strconv.ParseInt(idStr,10,64)
				if e!=nil{
					zapLog.Logger.Error(e)
					return
				}
				v,ok:=GpsInfoMap.Load(id)
				if !ok {
					return
				}
				info  := v.(*gpsInfo)
				info.Id = id
				bits:=binary.LittleEndian.Uint64(data[10:18])
				info.Lng = math.Float64frombits(bits)

				bits = binary.LittleEndian.Uint64(data[18:26])
				info.Lat = math.Float64frombits(bits)

				info.T = time.Now().Format("2006-01-02 15:04:05")

				go updateToDb(info)
			}
		}

	}
}

func updateToDb(info *gpsInfo)  {
	fmt.Println(info)
	strsql:=fmt.Sprintf("insert into gpsinfo%02d%02d value (%d,%v,%v,'%s')",time.Now().Year(),time.Now().Month(),
		info.Id,info.Lng,info.Lat,info.T)

	_, err := dbMysql.Db.Exec(strsql)
	if err!=nil{
		zapLog.Logger.Error(err)
		return
	}
	go BoastToAllClients(info,99)
}