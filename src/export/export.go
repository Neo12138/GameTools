/**
 * 将excel配置表导出为纯文本格式配置表
 */
package export

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"log"
	"os"
	"strings"
	"time"
	"utils"
)

//可以是绝对路径
type Setting struct {
	//excel表的所在的目录，相对路径或绝对路径
	SourceDir string
	//excel表的后缀名
	SourceSuffix string
	//导出标志，excel中某些列是不需要导出的
	ExportFlag string
	//导出目录
	OutDir string
	//导出文件后缀
	ExportSuffix string
	//配置表数据的类型定义文件(.d.ts)文件导出目录
	DefOutDir string
	//配置表数据所在的名字空间
	DefNamespace string
	//配置表 表名文件导出目录
	ConfigNameOutDir string
	//配置表 表名文件所在的名字空间
	ConfigNameNamespace string
}

var setting Setting
var tsTypeMap map[string]string
var encodeMap map[string]string
var failCount = 0
var failDesc string
var dtsConfigMapDeclare string
var dts *os.File
var dtsBuffer *bufio.Writer
var fName *os.File
var bufferName *bufio.Writer

func Run(args []string) {
	initData()
	readStartupParams(args)
	readFiles()
}

/*初始化配置数据*/
func initData() {
	setting = Setting{
		SourceDir:           "./tables/",
		SourceSuffix:        ".xlsx|.xls",
		ExportFlag:          "C",
		OutDir:              "./configs/",
		ExportSuffix:        ".txt",
		DefOutDir:           "./def/",
		DefNamespace:        "ConfigData",
		//ConfigNameOutDir:    "./def/",
		//ConfigNameNamespace: "ConfigName",
	}

	tsTypeMap = map[string]string{
		"int":          "number",
		"float":        "number",
		"float64":      "number",
		"bool":         "boolean",
		"string":       "string",
		"interface {}": "any",

		"int[]":          "number[]",
		"float[]":        "number[]",
		"float64[]":      "number[]",
		"bool[]":         "boolean[]",
		"string[]":       "string[]",
		"[]interface {}": "any[]",
	}

	encodeMap = map[string]string{
		" ":   "%20",
		"\\n": "%0A",
		"\\t": "%09",
	}
}

/*读取启动参数*/
func readStartupParams(args []string) {
	if len(args) <= 1 {
		return
	}
	props := []string{"SourceDir",
		"OutDir", "ExportFlag", "ExportSuffix",
		"DefOutDir", "DefNamespace",
		"ConfigNameOutDir", "ConfigNameNamespace"}

	args = args[1:]
	fmt.Println("启动参数：", args)
	utils.ReflectAssign(&setting, props, args)
}

/*读取excel文件*/
func readFiles() {
	pattern := fmt.Sprintf(`^[^~].*(%s)$`, setting.SourceSuffix)
	efcFiles := utils.GetEffectiveFiles(setting.SourceDir, pattern)
	numExcels := len(efcFiles)
	fmt.Printf("检测到有 %d 个excel文件待处理\n", numExcels)

	if numExcels > 0 {
		utils.EnsureDir(setting.OutDir)
		utils.EnsureDir(setting.DefOutDir)
		utils.EnsureDir(setting.ConfigNameOutDir)

		begin := time.Now()

		//创建配置表定义文件 bufferWriter
		dts, dtsBuffer = utils.GetBufferWriter(setting.DefOutDir + "config.d.ts")
		defer dts.Close()
		_, _ = dtsBuffer.WriteString("//由工具自动生成，请勿手动修改\n")
		_, _ = dtsBuffer.WriteString("declare namespace " + setting.DefNamespace + " {")
		indent(1)

		//创建配置表名称定义文件 bufferWriter
		if setting.ConfigNameNamespace != "" {
			fName, bufferName = utils.GetBufferWriter(setting.ConfigNameOutDir + "config-names.ts")
			defer fName.Close()
			_, _ = bufferName.WriteString("//由工具自动生成，请勿手动修改\n")
			_, _ = bufferName.WriteString("namespace " + setting.ConfigNameNamespace + " {")
		}

		//处理每个excel表
		for i, file := range efcFiles {
			parseExcel(file, i+1)
		}

		_, _ = dtsBuffer.WriteString(dtsConfigMapDeclare)
		_, _ = dtsBuffer.WriteString("\n}")
		_ = dtsBuffer.Flush()

		if bufferName != nil {
			_, _ = bufferName.WriteString("\n}")
			_ = bufferName.Flush()
		}

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
	file, buffer := utils.GetBufferWriter(outPath)
	defer file.Close()

	interfaceConfig := utils.StringToCamel(sheet, true)
	propConfig := utils.StringToCamel(sheet, false)
	//文件名
	_, _ = buffer.WriteString(propConfig + "\n")

	//头部
	//var header string
	//header += " " + names[i] + "|" + getTargetLangType(types[i])
	//header = header[1:]
	//_, _ = buffer.WriteString(header)
	var strNames string
	var strTypes string
	for _, i := range effectiveColumns {
		strNames += " " + utils.StringToCamel(names[i], false)
		strTypes += " " + getTargetLangType(types[i])
	}
	_, _ = buffer.WriteString(strNames[1:])
	_, _ = buffer.WriteString("\n" + strTypes[1:])

	//数据
	for _, row := range rows[4:] {
		var line string
		for _, i := range effectiveColumns {
			line += " " + encode(row[i])
		}
		line = "\n" + line[1:]

		_, _ = buffer.WriteString(line)
	}

	utils.CheckError(buffer.Flush())
	utils.SetReadOnly(outPath)

	log.Printf("导出 %s[%s] 到 %s%s 成功！\n", excel.Path, sheet, sheet, setting.ExportSuffix)


	//写入到dts文件
	if dtsBuffer != nil {
		_, _ = dtsBuffer.WriteString(tab + "interface I" + interfaceConfig + " {")
		indent(1)
		for _, i := range effectiveColumns {
			_, _ = dtsBuffer.WriteString(tab + "/** " + rows[0][i] + " */")
			convertedName := utils.StringToCamel(names[i], false)
			_, _ = dtsBuffer.WriteString(tab + "readonly " + convertedName + ": " + getTargetLangType(types[i]) + ";")
		}
		_, _ = dtsBuffer.WriteString(indent(-1) + "}\n")
	}
	dtsConfigMapDeclare += tab + "let " + propConfig + ": { [key: number]: I" + interfaceConfig + " };"

	if bufferName != nil {
		_, _ = bufferName.WriteString("\n\texport const " + propConfig + ": string = \"" + sheet + "\";")
	}
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
	return s
}

//转成目标语言的类型
func getTargetLangType(t string) string {
	return tsTypeMap[t]
}

func printExcelError(path string, sheet string, err error) {
	out := fmt.Sprintf("读取(%s[%s])出错！%s。已经忽略导出\n", path, sheet, err)
	log.Print(out)
	failCount++
	failDesc += out + "\n"
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
