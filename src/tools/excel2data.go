package tools

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"log"
	"os"
	"strings"
	"time"
)

//可以是绝对路径
type Setting struct {
	SourceDir        string
	OutDir           string
	DefOutDir        string
	ImportSuffix     string
	ExportFlag       string
	ExportSuffix     string
	ConfigMoveTo     string
	DefMoveTo        string
	ConfigNamespace  string
	ConfigNameOutDir string
}

var setting Setting

var jsTypeMap = map[string]string{
	"int":      "number",
	"float":    "number",
	"string":   "string",
	"bool":     "boolean",
	"int[]":    "number[]",
	"float[]":  "number[]",
	"string[]": "string[]",
	"bool[]":   "boolean[]",

	"float64":        "number",
	"interface {}":   "any",
	"[]interface {}": "any[]",
}
var encodeMap = map[string]string{
	" ":   "%20",
	"\\n": "%0A",
	"\\t": "%09",
}

var failCount = 0
var failDesc string
var dtsConfigMapDeclare string
var dts *os.File
var dtsBuffer *bufio.Writer

func Run(args []string) {
	setting = Setting{
		SourceDir:       "./tables/",
		OutDir:          "./builds/",
		DefOutDir:       "./def/",
		ImportSuffix:    ".xlsx|.xls",
		ExportFlag:      "C",
		ExportSuffix:    ".txt",
		ConfigMoveTo:    "",
		DefMoveTo:       "",
		ConfigNamespace: "ConfigData",
		ConfigNameOutDir: "",
	}

	readStartupParams(args)

	readFiles()
}
func readStartupParams(args []string) {
	if len(args) <= 1 {
		return
	}
	props := []string{"SourceDir", "OutDir", "ExportFlag", "ExportSuffix", "DefOutDir", "ConfigNamespace", "ConfigNameMoveTo"}
	begin := time.Now()
	args = args[1:]
	fmt.Println("启动参数：", args)
	ReflectAssign(&setting, props, args)
	fmt.Println("反射赋值耗时：", time.Since(begin))

}

func readFiles() {
	pattern := fmt.Sprintf(`^[^~].*(%s)$`, setting.ImportSuffix)
	efcFiles := GetEffectiveFiles(setting.SourceDir, pattern)
	numExcels := len(efcFiles)
	fmt.Printf("检测到有 %d 个excel文件待处理\n", numExcels)

	if numExcels > 0 {
		EnsureDir(setting.OutDir)
		EnsureDir(setting.DefOutDir)
		EnsureDir(setting.ConfigMoveTo)
		EnsureDir(setting.DefMoveTo)

		begin := time.Now()

		dts, dtsBuffer = GetBufferWriter(setting.DefOutDir + "config.d.ts")
		defer dts.Close()
		_, _ = dtsBuffer.WriteString("//由工具自动生成，请勿手动修改\n")
		_, _ = dtsBuffer.WriteString("declare namespace "+setting.ConfigNamespace+" {")
		indent(1)

		//处理每个excel表
		for i, file := range efcFiles {
			parseExcel(file, i+1)
		}

		_, _ = dtsBuffer.WriteString(dtsConfigMapDeclare)
		_, _ = dtsBuffer.WriteString("\n}")
		_ = dtsBuffer.Flush()
		CopyTo("config.d.ts", setting.DefOutDir, setting.DefMoveTo)

		fmt.Printf("所有excel处理完成。成功%d个,失败%d个，总任务%d个，耗时：%s\n", numExcels-failCount, failCount, numExcels, time.Since(begin))
		if failCount > 0 {
			fmt.Println("错误项：\n", failDesc)
		}
	}
}

//解析Excel文件
func parseExcel(filename string, index int) {
	log.Println("任务", index)
	log.Printf("正在处理 %s...\n", filename)
	excel, err := excelize.OpenFile(setting.SourceDir + filename)
	if err != nil {
		printExcelError(filename, "", err)
		return
	}

	sheets := excel.GetSheetMap()
	for _, name := range sheets {
		parseSheet(excel, name)
	}
	log.Printf("%s 处理完成\n", excel.Path)
}

//解析每张工作簿
func parseSheet(excel *excelize.File, sheet string) {
	//读取表格中的所有行
	rows, err := excel.GetRows(sheet)
	if err != nil {
		printExcelError(excel.Path, sheet, err)
		return
	}

	numRows := len(rows)
	//检查行数
	if numRows < 4 {
		printExcelError(excel.Path, sheet, errors.New("表格行数至少为4行，实际为："+string(numRows)))
		return
	}

	//检查有效列
	types := rows[1]
	names := rows[2]
	effectiveColumns, err := getEffectiveColumns(types, names, rows[3])
	if err != nil {
		printExcelError(excel.Path, sheet, err)
		return
	}

	//正式开始处理
	outPath := setting.OutDir + sheet + setting.ExportSuffix
	file, buffer := GetBufferWriter(outPath)
	defer file.Close()

	//文件名
	_, _ = buffer.WriteString(sheet + "\n")

	//头部
	var header string
	for _, i := range effectiveColumns {
		header += " " + names[i] + "|" + getTargetLangType(types[i])
	}
	header = header[1:]
	_, _ = buffer.WriteString(header)

	//数据
	for _, row := range rows[4:] {
		var line string
		for _, i := range effectiveColumns {
			line += " " + encode(row[i])
		}
		line = "\n" + line[1:]

		_, _ = buffer.WriteString(line)
	}

	CheckError(buffer.Flush())

	CopyTo(sheet+setting.ExportSuffix, setting.OutDir, setting.ConfigMoveTo)
	SetReadOnly(outPath)
	SetReadOnly(setting.ConfigMoveTo + sheet + setting.ExportSuffix)

	log.Printf("导出 %s[%s] 到 %s%s 成功！\n", excel.Path, sheet, sheet, setting.ExportSuffix)

	//写入到dts文件
	if dtsBuffer != nil {
		_, _ = dtsBuffer.WriteString(tab + "interface I" + sheet + " {")
		indent(1)
		for _, i := range effectiveColumns {
			_, _ = dtsBuffer.WriteString(tab + "/** " + rows[0][i] + " */")
			_, _ = dtsBuffer.WriteString(tab + "readonly " + names[i] + ": " + getTargetLangType(types[i]) + ";")
		}
		_, _ = dtsBuffer.WriteString(indent(-1) + "}\n")
	}
	dtsConfigMapDeclare += tab + "let " + sheet + ": { [key: number]: I" + sheet + " };"
}

//获取表格中的有效列的索引
func getEffectiveColumns(types []string, names []string, exportFlags []string) ([]int, error) {
	columns := make([]int, 0, len(exportFlags))
	for i, v := range exportFlags {
		if !strings.Contains(v, setting.ExportFlag) {
			continue
		}
		if names[i] == "" || types[i] == "" {
			return nil, errors.New("存在无效列" + string(i))
		}
		columns = append(columns, i)
	}
	if len(columns) <= 0 {
		return nil, errors.New("有效列数为0")
	}

	return columns, nil
}

// 转义特殊字符
func encode(s string) string {
	for k, v := range encodeMap {
		s = strings.Replace(s, k, v, -1)
	}
	//s = strings.Replace(s, " ", "%20", -1)
	//s = strings.Replace(s, "\\n", "%0A", -1)
	//s = strings.Replace(s, "\\t", "%09", -1)
	//s = strings.Replace(s, "|", "%7C", -1)
	return s
}

//转成目标语言的类型
func getTargetLangType(t string) string {
	return jsTypeMap[t]
}

func printExcelError(path string, sheet string, err error) {
	out := fmt.Sprintf("读取(%s[%s])出错！%s。已经忽略导出\n", path, sheet, err)
	log.Print(out)
	failCount++
	failDesc += out
}

var tab = "\n"

func indent(delta int) string {
	if delta > 0 {
		tab += "\t"
	} else {
		numTab := len(tab)
		if len(tab) >= 2 {
			tab = tab[:numTab-1]
		} else {
			tab = "\n"
		}
	}
	return tab
}
