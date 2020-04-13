package main

import (
	"fmt"
	"utils"
)

/**
1.读取所有文件，并对每个进行错误检查，存在错误汇总输出，无错误则下一步
2.判断文件的结构：对象，对象数组
3.针对不同文件结构进行不同的解析。解析结果包括数据定义和数据
4.将所有文件中的数据定义合并输出，将所有数据合并输出(如有必要加密)

rsrc -manifest test.manifest -o rsrc.syso
gui程序打包命令 go build -ldflags="-H windowsgui"
*/

type Object map[string]interface{}
type ObjectArray []Object
func main() {
	//file_parser.Run()
	//utils.Pause()

	//f,_ := ioutil.ReadFile("新建项目.gt")
	//m1 := utils.ParseMap(f)
	//fmt.Println(m1["name"], m1["source"])
	//
	//fmt.Println()
	//f2,_ := ioutil.ReadFile(".projects.list")
	//m2 := utils.ParseMap(f2)
	//for s, i := range m2 {
	//	fmt.Println(s, i)
	//}
	a := make([]Object, 0, 3)
	a = append(a, map[string]interface{}{"id": 1})
	a = append(a, map[string]interface{}{"id": 2})
	a = append(a, map[string]interface{}{"id": 3})
	a = append(a, map[string]interface{}{"id": 4})

	fmt.Println(utils.StringifyMap(a))

	var v Object
	var index = -1
	var id = 3
	for i, object := range a {
		if object["id"] == id {
			v = object
			index = i
			break
		}
	}
	fmt.Println("find: ", index, v)

}
