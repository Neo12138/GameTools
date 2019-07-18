package practise

import (
	"bufio"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"os"
	"path/filepath"
	"strings"
)

func RunFile() {
	testFilePath()
}

func testFilePath() {
	//文件夹也会读取到
	//files, _ := ioutil.ReadDir("./")
	//for _, f := range files {
	//	fmt.Println(f.Name())
	//}

	files, _ := filepath.Glob("*.xlsx")
	for _, f := range files {
		fmt.Println(f)
		testExcel(f)
	}

}

func testExcel(filename string) {
	xlsx, err := excelize.OpenFile(filename)
	if err != nil {
		fmt.Println(err)
		return
	}

	//获取工作表及索引
	sheets := xlsx.GetSheetMap()
	fmt.Println(sheets)

	for index, name := range sheets {
		fmt.Println(index, name)
		parseSheet(xlsx, name)
	}
}

func parseSheet(xlsx *excelize.File, sheet string) {
	fmt.Print(xlsx.Path)
	// 获取 Sheet1 上所有单元格
	rows, err := xlsx.GetRows(sheet)
	if err != nil {
		fmt.Printf("读取Sheet(%s[%s])时发生了错误。 %s \n",xlsx.Path, sheet, err)
		return
	}
	if len(rows) < 4 {
		fmt.Printf("Sheet(%s[%s])行数不合法，少于4行。已经忽略\n",xlsx.Path, sheet)
		return
	}

	path := sheet + ".data"
	newFile, err := os.Create(path)
	checkErr(err)
	defer newFile.Close()

	buffer := bufio.NewWriter(newFile)

	effectiveColumns := getEffectiveColumns(rows[1], rows[2], rows[3])
	fmt.Println(effectiveColumns)
	if len(effectiveColumns) <= 0 {
		fmt.Printf("解析Sheet(%s[%s])失败，有效列数为0\n", xlsx.Path, sheet)
		return
	}

	//正式处理
	//处理头部
	var header string
	for _, i := range effectiveColumns {
		header += " " + rows[2][i] + "|" + getJSType(rows[1][i])
	}
	header = header[1:]
	_, _ = buffer.WriteString(header)

	//处理数据
	for _, row := range rows[4:] {
		var line string
		for _, i := range effectiveColumns {
			line += " " + encode(row[i])
		}

		line = "\n" + line[1:]
		_, _ = buffer.WriteString(line)
	}
	checkErr(buffer.Flush())
}

func getEffectiveColumns(names []string, types []string, kinds []string) []int {
	columns := make([]int, 0, len(kinds))
	for i, v := range kinds {
		if strings.Contains(v, C) && names[i] != "" && types[i] != "" {
			columns = append(columns, i)
		}
	}
	return columns
}
func encode(s string) string {
	s = strings.Replace(s, " ", "%20", -1)
	s = strings.Replace(s, "\\n", "%0A", -1)
	s = strings.Replace(s, "\\t", "%09", -1)
	//s = strings.Replace(s, "|", "%7C", -1)
	return s
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

const C = "C"
var TypeMap = map[string]string{
	"int":      "number",
	"float":    "number",
	"string":   "string",
	"bool":     "boolean",
	"int[]":    "number[]",
	"float[]":  "number[]",
	"string[]": "string[]",
	"bool[]":   "boolean[]",
}


func getJSType(t string) string {
	return TypeMap[t]
}
