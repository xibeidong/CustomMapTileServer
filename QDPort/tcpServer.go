package QDPort

import (
	"My/CustomTileMapServer/conf"
	"My/Learn/zapLog"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"
)
var tcpConnMap sync.Map
func StartTcpServer()  {
	addr, err := net.ResolveTCPAddr("tcp", conf.MyConfs.QdPortTcpServerListen)
	if err!=nil{
		zapLog.Logger.Error(err)
		return
	}
	tcpListener, err := net.ListenTCP("tcp", addr)
	if err!=nil{
		zapLog.Logger.Error(err)
		return
	}
	defer tcpListener.Close()

	for{
		tcpConn, err := tcpListener.AcceptTCP()
		if err!=nil{
			zapLog.Logger.Error(err)
			continue
		}
		zapLog.Logger.Info("新Tcp连接 ==》 ",tcpConn.RemoteAddr().String())

		tcpConnMap.Store(tcpConn.RemoteAddr(),tcpConn)
		go tcpPipe(tcpConn)
	}
}

func tcpPipe(conn *net.TCPConn)  {
	defer func() {
		tcpConnMap.Delete(conn.RemoteAddr())
		zapLog.Logger.Warn("close tcp connect :",conn.RemoteAddr())
		conn.Close()

	}()
	for  {
		bufHead:=make([]byte,6,6)
		len, err := conn.Read(bufHead)
		if err!=nil{
			zapLog.Logger.Error(err)

			return
		}
		fmt.Println("Head len = ",len)
		messageId:= binary.BigEndian.Uint16(bufHead[:2])
		bodyLen:=binary.BigEndian.Uint32(bufHead[2:6])

		bufBody:=make([]byte,bodyLen,bodyLen)
		len, err = conn.Read(bufBody)
		if err!=nil{
			zapLog.Logger.Error(err)
			continue
		}
		fmt.Println("Body len = ",len," ; messageID = ",messageId)
		switch messageId {
		case 101:
			heartHandler(conn,bufBody,101)
		case 103:
			playbackHandler(conn,bufBody,103)
		}
	}

}

func BoastToAllClients(info *gpsInfo,messageId uint16)  {

	messageBody, err := json.Marshal(info)
	if err!=nil{
		zapLog.Logger.Error(err)
	}

	idBytes:=make([]byte ,2,2)
	binary.LittleEndian.PutUint16(idBytes,messageId)

	bodyLenByte:=make([]byte,4,4)
	binary.LittleEndian.PutUint32(bodyLenByte,uint32(len(messageBody)))
	fmt.Println("messageLen = ",len(messageBody))

	data := bytesCombine(idBytes, bodyLenByte, messageBody)

	tcpConnMap.Range(func(key , value interface{}) bool {
		tcpConn:=value.(*net.TCPConn)
		tcpConn.SetWriteDeadline(time.Now().Add(time.Second*5))
		_, err2 := tcpConn.Write(data)
		if err2!=nil{
			zapLog.Logger.Error(err2)
		}
		return true
	})
}
func reply(conn *net.TCPConn,messageBody []byte,messageId uint16)  {

	idBytes:=make([]byte ,2,2)
	binary.BigEndian.PutUint16(idBytes,messageId)

	bodyLenByte:=make([]byte,4,4)
	binary.BigEndian.PutUint32(bodyLenByte,uint32(len(messageBody)))

	data := bytesCombine(idBytes, bodyLenByte, messageBody)

	conn.Write(data)

}
func uint16ToBytes(num uint16) []byte {
	b:=make([]byte,2)
	binary.BigEndian.PutUint16(b,num)
	return b
}

func bytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}