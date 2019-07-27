// main.go
package main

import (
	"fmt"
	"os"
	"utils"
)

func main() {
	//tools.Run(os.Args)
	//tools.RunMerge(os.Args)
	//test()
	//filename.Run(os.Args)
	//tools.RunGetConfigDef(os.Args)

	a := utils.StringToCamel(os.Args[1], false)
	b := utils.StringToCamel(os.Args[1], true)
	fmt.Printf("length: %d %d\n", len(a), len(b))
	fmt.Println(a)
	fmt.Println(b)
	//utils.Pause()

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
