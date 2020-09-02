package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("请输入需要移除的目录：")
	if scanner.Scan(){
		path:=scanner.Text()
		fmt.Println("\n正在执行删除任务=》",path)
		err := os.RemoveAll(path)
		if err!=nil{
			fmt.Println(err)
		}
	}
	fmt.Println("\nOK，请关闭 Ctrl+C")
	for{}
}
