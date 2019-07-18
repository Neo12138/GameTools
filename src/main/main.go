// main.go
package main

import (
	"fmt"
	"os"
	"tools"
)

func main() {
	//tools.Run(os.Args)
	//tools.RunMerge(os.Args)
	//test()
	//filename.Run(os.Args)
	tools.RunGetConfigDef(os.Args)

	fmt.Print("输入任意字符结束：")
	var enter string
	fmt.Scanln(&enter)
	fmt.Printf("%s\n", enter)

}

func test() {
	//practise.Run()
	//practise.RunTime()
	//
	//practise.RunFile()
	//path := "./builds/activity.txt"
	//
	//stat, _ := os.Stat(path)
	//fmt.Println("打开前: ", stat.Name(), stat.ModTime(), stat.Mode())
	//tools.SetReadWrite(path)
	//
	//file, _ := os.Open(path)
	//defer file.Close()
	//
	//stat, _ = file.Stat()
	//fmt.Println("打开后: ", stat.Name(), stat.ModTime(), stat.Mode())
	//
	//tools.SetReadOnly(path)
	//stat, _ = file.Stat()
	//fmt.Println("修改后: ", stat.Name(), stat.ModTime(), stat.Mode())
}
