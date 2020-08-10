package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"time"
)

type postionXYZ struct {
	lng float64
	lat float64
}
const rootPath = "D:/MapDownCustom"
var urls  = [4]string {
	"https://mt0.googleapis.com/vt/lyrs=y&hl=zh-cn&x=%v&y=%v&z=%v",
	"https://mt1.googleapis.com/vt/lyrs=y&hl=zh-cn&x=%v&y=%v&z=%v",
	"https://mt2.googleapis.com/vt/lyrs=y&hl=zh-cn&x=%v&y=%v&z=%v",
	"https://mt3.googleapis.com/vt/lyrs=y&hl=zh-cn&x=%v&y=%v&z=%v",
}

func main() {
	downByZoom(
		postionXYZ{120.1042556762695312, 35.9743881225585938},
		postionXYZ{120.2487945556640625, 36.0701751708984375},
		15)
	ch:=make(chan int)
	<-ch
}
/**
 * 谷歌下转换经纬度对应的层行列
 *
 * @param lon  经度
 * @param lat  维度
 * @param zoom 在第zoom层进行转换
 * @return
 */
func GoogleLngLatToTileXY(lng ,lat ,zoom float64) (x,y float64)  {
	n:= math.Pow(2,zoom)
	tileX:=(lng+180)/360 * n
	tileY:= (1-(math.Log(math.Tan(Augular2Radain(lat))+(1/math.Cos(Augular2Radain(lat))))/math.Pi))/2*n

	return math.Floor(tileX),math.Floor(tileY);
}
/**
 * 层行列转经纬度
 *
 * @param x
 * @param y
 * @param z
 * @return
 */
func XYZtoLngLat(z,x,y float64) (lng,lat float64)  {
	n:=math.Pow(2,z)
	lng =x/n*360-180
	lat = math.Atan(math.Sinh(math.Pi*(1-2*y/n)))
	lat = lat*180/math.Pi
	return
}
//角度转弧度
func Augular2Radain(augular float64) float64{
	return augular*math.Pi/180
}

func downByZoom(minPos,maxPos postionXYZ,z float64)  {
	minX,minY:=GoogleLngLatToTileXY(minPos.lng,minPos.lat,z)
	maxX,maxY:=GoogleLngLatToTileXY(maxPos.lng,maxPos.lat,z);

	fmt.Println(minX,minY)
	fmt.Println(maxX,maxY)
	x1:=minX
	x2:=maxX
	y1:=maxY
	y2:=minY

	url1 := fmt.Sprintf(urls[0],x1,y1,z)
	fmt.Println(url1)
	url2 :=fmt.Sprintf(urls[1],x2,y2,z)
	fmt.Println(url2)
	downTile(url1,x1,y1,z)

}
//http://119.167.141.201/
func downTile(url string,x,y,z float64)  {
	Second("https://www.google.com/","http://119.167.141.201:14390")

	//resp,e := http.Get("https://www.google.com")
	//if e!=nil {
	//	fmt.Println(e.Error())
	//	return
	//}
	//defer resp.Body.Close()
	//body,_:=ioutil.ReadAll(resp.Body)
	//f,e :=os.Create(rootPath+"/123.png")
	//if e!=nil {
	//	fmt.Println(e)
	//}else {
	//	f.Write(body)
	//	f.Close()
	//}
}

func Second(webUrl, proxyUrl string) {
	/*
		1. 代理请求
		2. 跳过https不安全验证
	*/
	// webUrl := "http://ip.gs/"
	// proxyUrl := "http://115.215.71.12:808"

	proxy, _ := url.Parse(proxyUrl)
	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 5, //超时时间
	}

	resp, err := client.Get(webUrl)
	if err != nil {
		fmt.Println("出错了", err)
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

}
