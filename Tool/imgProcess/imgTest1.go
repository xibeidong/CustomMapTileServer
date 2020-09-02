package main

import (
	"bytes"
	"fmt"
	"github.com/Comdex/imgo"
	"image"
	"io/ioutil"
	"os"
)

func main() {

	imgMatrix := imgo.MustRead("D:\\SGDownload\\png\\1.png")
	err := imgo.SaveAsJPEG("D:\\SGDownload\\png\\1_1.jpg",imgMatrix, 50)
	if err!=nil{
		fmt.Println(err)
	}
	select {

	}

}
func testClarlityPNG(){
	path:="D:\\SGDownload\\png\\3225.jpg"
	data, err := ioutil.ReadFile(path)

	if err!=nil{
		panic(err)
	}
	img, s, _ := bytesToImage(data)
	fmt.Println(s)

	imgMatrix := imgo.ImageToImgMatrix(img)

	height := len(imgMatrix)
	width := len(imgMatrix[0])
	for i:=0;i<height;i++{
		for j:=0;j<width;j++{
			if imgMatrix[i][j][0] >= 225 &&imgMatrix[i][j][1] >= 225 && imgMatrix[i][j][2] >=225{
				imgMatrix[i][j][3] = 0
			}
		}
	}

	data2, err := getBytesAfterSaveAsJPED(imgMatrix, 50)
	if err!=nil{
		panic(err)
	}
	err = ioutil.WriteFile("D:\\SGDownload\\png\\3225_2.jpg", data2, 0666)
	if err!=nil{
		panic(err)
	}
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

	_, err2 := os.Stat("./temp.jpeg")
	if err2 == nil{
		err3 := os.Remove("./temp.jpeg")
		if err3!=nil{
			return nil,err3
		}
	}

	err := imgo.SaveAsJPEG("./temp.jpeg", imgMatrix, quality)
	if err!=nil{
		return nil,err
	}

	data, err := ioutil.ReadFile("./temp.jpeg")
	if err!=nil{
		return nil,err
	}
	return data,nil
}

//这个方法需要添加到包 imgo 下面
//func ImageToImgMatrix(img image.Image) (imgMatrix [][][]uint8) {
//
//	bounds := img.Bounds()
//	width := bounds.Max.X
//	height := bounds.Max.Y
//
//	src := convertToNRGBA(img)
//	imgMatrix = NewRGBAMatrix(height, width)
//
//	for i := 0; i < height; i++ {
//		for j := 0; j < width; j++ {
//			c := src.At(j, i)
//			r, g, b, a := c.RGBA()
//			imgMatrix[i][j][0] = uint8(r)
//			imgMatrix[i][j][1] = uint8(g)
//			imgMatrix[i][j][2] = uint8(b)
//			imgMatrix[i][j][3] = uint8(a)
//
//		}
//	}
//	return
//}