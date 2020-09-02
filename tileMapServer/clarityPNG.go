package tileMapServer

import (
	"My/Learn/zapLog"
	"bytes"
	"github.com/Comdex/imgo"
	"image"
	"io/ioutil"
	"os"
	"strconv"
)

//弃用的，地图瓦片的白色部分变成透明，效果不理想，有白色锯齿，白色标注也会受到影响
//quality must be  1-100
func clarityImg(data []byte,quality int) ([]byte,error)  {
	img, _, _ := bytesToImage(data)
	//fmt.Println(s)

	imgMatrix := imgo.ImageToImgMatrix(img)

	height := len(imgMatrix)
	width := len(imgMatrix[0])
	for i:=0;i<height;i++{
		for j:=0;j<width;j++{
			if imgMatrix[i][j][0] >= 245 &&imgMatrix[i][j][1] >= 245 && imgMatrix[i][j][2] >=245{
				imgMatrix[i][j][3] = 0
			}
		}
	}

	data2, err := getBytesAfterSaveAsJPED(imgMatrix, quality)
	if err!=nil{
		return nil,err
	}
	return data2,nil
}

func bytesToImage(data []byte) ( image.Image,string ,error) {
	reader := bytes.NewReader(data)
	img, s, err := image.Decode(reader)
	if err!=nil{
		panic(err)
	}
	return img,s,err
}

func getBytesAfterSaveAsJPED( imgMatrix [][][]uint8, quality int) ([]byte,error) {
	count:=<-ch
	unix:= strconv.FormatInt(count,10)
	count++
	ch<-count
	name:="./"+unix+"_temp.jpeg"
	//_, err2 := os.Stat(name)
	//if err2 == nil{
	//	err3 := os.Remove(name)
	//	if err3!=nil{
	//		return nil,err3
	//	}
	//}

	err := imgo.SaveAsJPEG(name, imgMatrix, quality)
	if err!=nil{
		return nil,err
	}

	defer func() {
		err:=os.Remove(name)
		if err!=nil{
			zapLog.Logger.Warn(err)
		}
	}()

	data, err := ioutil.ReadFile(name)
	if err!=nil{
		return nil,err
	}
	return data,nil
}
