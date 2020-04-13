package file_parser

import (
	"bytes"
	"app"
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"regexp"
	"strings"
	"time"
	"utils"
)

/**
文件解析结果定义
*/
type ParseResult struct {
	//数据定义
	Define string
	//数据
	Data string
	//警告
	Warning string
}

type ExcelReadError struct {
	msg      string
	filename string
	taskId   int
}
type SheetError struct {
	msg   string
	sheet string
}

type DataStruct int8

const (
	//未知
	UNKNOWN DataStruct = 0
	//对象
	OBJECT DataStruct = 1
	//对象数据
	OBJECT_ARRAY DataStruct = 2
)

var tsTypeMap map[string]string
var encodeMap map[string]string

/*初始化配置数据*/
func initData() {
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

func (e *ExcelReadError) Error() string {
	return fmt.Sprintf("任务%d Excel[%s]出错: %s", e.taskId, e.filename, e.msg)
}
func (e *SheetError) Error() string {
	return fmt.Sprintf("Sheet[%s]出错: %s", e.sheet, e.msg)
}

func Run() {
	initData()

	suffix := ".xlsx|.xls"
	sourceDir := "./assets/"
	outDir := "./configs/"
	outDefDir := "./def/"
	defNamespace := "ConfigData"

	pattern := fmt.Sprintf(`^[^~].*(%s)$`, suffix)
	files := main.GetEffectiveFiles(sourceDir, pattern)
	numFiles := len(files)
	fmt.Printf("检测到有 %d 个excel文件待处理\n", numFiles)

	if numFiles > 0 {
		begin:=time.Now()
		main.EnsureDir(outDir)
		main.EnsureDir(outDefDir)

		var bDefine, bData, bError bytes.Buffer
		bDefine.WriteString("//由工具自动生成，请勿手动修改\n")
		bDefine.WriteString("declare namespace " + defNamespace + " {")
		for i, file := range files {
			ret, err := parseOneExcel(sourceDir+file, file, i+1)
			if err == nil {
				bDefine.WriteString(ret.Define + "\n")
				bData.WriteString(ret.Data + "\n")
				fmt.Printf("%d. %s 导出完成\n", i+1, file)
			} else {
				bError.WriteString(err.Error() + "\n")
			}
		}
		bDefine.WriteString("}")

		//最终输出
		if bError.Len() > 0 {
			fmt.Printf("\n错误：\n%s", bError.String())
		} else {
			fDTS, bDTS := main.GetFillAndBufferWriter(outDefDir + "config.d.ts")
			defer fDTS.Close()
			_, _ = bDTS.Write(bDefine.Bytes())
			_ = bDTS.Flush()

			fConfig, bConfig := main.GetFillAndBufferWriter(outDir + "config.csv")
			defer fConfig.Close()
			_, _ = bConfig.Write(bData.Bytes())
			_ = bConfig.Flush()
			fmt.Printf("所有excel处理完成。总任务%d个，耗时：%s\n", numFiles, time.Since(begin))
		}
	}
}

func parseOneExcel(path string, filename string, taskId int) (*ParseResult, error) {
	excel, err := openExcel(path, filename, taskId)
	if err != nil {
		return nil, err
	}
	return parseExcel(excel)
}

func openExcel(path string, filename string, taskId int) (*excelize.File, error) {
	excel, err := excelize.OpenFile(path)
	if err != nil {
		return nil, &ExcelReadError{msg: err.Error(), filename: filename, taskId: taskId}
	}
	return excel, nil
}

/**
一张excel文件可能有多个sheet,需要逐个处理
一张excel文件中可能有多个错误，需要全部检查收集
*/
func parseExcel(excel *excelize.File) (*ParseResult, error) {
	var bDefine, bData, bError bytes.Buffer
	sheets := excel.GetSheetMap()

	//分别处理每张sheet
	for _, name := range sheets {
		ret, err := parseSheet(excel, name)
		if err != nil {
			bError.WriteString(err.Error())
		} else {
			bDefine.WriteString(ret.Define)
			bData.WriteString(ret.Data)
		}
	}
	// 如果没错误，返回此excel的数据定义和数据的合并结果
	if bError.Len() == 0 {
		return &ParseResult{Define: bDefine.String(), Data: bData.String()}, nil
	}
	return nil, errors.New(bError.String())
}

/**
解析一张Sheet
判断数据结构，分别解析
*/
func parseSheet(excel *excelize.File, sheet string) (*ParseResult, error) {
	//读取表格中的所有行
	rows, err := excel.GetRows(sheet)
	if err != nil {
		return nil, &SheetError{msg: err.Error(), sheet: sheet}
	}

	structType := getDataStruct(rows)
	var ret *ParseResult
	switch structType {
	case OBJECT:
		ret, err = parseSheetOfObject(rows, sheet)
		break
	case OBJECT_ARRAY:
		ret, err = parseSheetOfObjectArray(rows, sheet)
		break
	default:
		return nil, &SheetError{msg: "表格格式不正确", sheet: sheet}
	}
	if err != nil {
		return nil, &SheetError{msg: err.Error(), sheet: sheet}
	}
	return ret, nil
}

//解析对象结构的表格
//对于对象的结构，表中至少1行3列，其中第1列为注释，第2列为类型，第3列为字段名，第4列为数据
func parseSheetOfObject(rows [][]string, sheet string) (*ParseResult, error) {
	var bDefine, bData, bError bytes.Buffer
	var interfaceName = utils.StringToCamel(sheet, true)
	var attrName = utils.StringToCamel(sheet, false)
	var tab = &utils.Tab{Content: ""}
	tab.Indent(1)
	bDefine.WriteString(tab.Content + "interface I" + interfaceName + " {")
	tab.Indent(1)
	bData.WriteString(interfaceName + " object\n")

	for i, row := range rows {
		var hasError bool = false
		if row[1] == "" || row[2] == "" {
			bError.WriteString(fmt.Sprintf("第%d行1、2列存在空项。\n", i+1))
			hasError = true
		}
		if !isKeywordsOfType(row[1]) {
			bError.WriteString(fmt.Sprintf("第%d行2列，变量类型 '%s' 不合法。\n", i+1, row[1]))
			hasError = true
		}
		if !isValidVariable(row[2]) {
			bError.WriteString(fmt.Sprintf("第%d行3列，变量名 '%s' 不合法。\n", i+1, row[2]))
			hasError = true
		}
		if hasError {
			continue
		}

		/** 定义 */
		bDefine.WriteString(tab.Content + "/** " + row[0] + " */")
		fieldName := utils.StringToCamel(row[2], false)
		bDefine.WriteString(tab.Content + "readonly " + fieldName + ": " + getTSType(row[1]) + ";")

		//数据
		bData.WriteString(row[1] + " " + row[2] + " " + encode(row[3]) + "\n")
	}
	tab.Indent(-1)
	bDefine.WriteString(tab.Content + "}")
	bDefine.WriteString(tab.Content + "let " + attrName + ": I" + interfaceName + ";\n")
	if bError.Len() > 0 {
		return nil, errors.New(bError.String())
	}
	return &ParseResult{Define: bDefine.String(), Data: bData.String()}, nil
}

//解析对象数组结构的表格
func parseSheetOfObjectArray(rows [][]string, sheet string) (*ParseResult, error) {
	var bDefine, bData, bError bytes.Buffer
	var numCol = len(rows[1])

	/** 定义 */
	var interfaceName = utils.StringToCamel(sheet, true)
	var attrName = utils.StringToCamel(sheet, false)
	var tab = &utils.Tab{Content: ""}
	tab.Indent(1)
	bDefine.WriteString(tab.Content + "interface I" + interfaceName + " {")
	tab.Indent(1)
	colComment := len(rows[0])
	for i := 0; i < numCol; i++ {
		var hasError bool = false
		if !isKeywordsOfType(rows[1][i]) {
			bError.WriteString(fmt.Sprintf("第2行%d列，变量类型 '%s' 不合法。\n", i+1, rows[1][i]))
			hasError = true
		}
		if !isValidVariable(rows[2][i]) {
			bError.WriteString(fmt.Sprintf("第3行%d列，变量名 '%s' 不合法。\n", i+1, rows[2][i]))
			hasError = true
		}
		if hasError {
			continue
		}

		if i < colComment && rows[0][i] != "" {
			bDefine.WriteString(tab.Content + "/** " + rows[0][i] + " */")
		}
		fieldName := utils.StringToCamel(rows[2][i], false)
		bDefine.WriteString(tab.Content + "readonly " + fieldName + ": " + getTSType(rows[1][i]) + ";")
	}
	tab.Indent(-1)
	bDefine.WriteString(tab.Content + "}")
	bDefine.WriteString(tab.Content + "let " + attrName + ": { [key: number]: I" + interfaceName + "};")

	if bError.Len() > 0 {
		return nil, errors.New(bError.String())
	}
	/** 数据 */
	bData.WriteString(interfaceName + " object_array\n")
	//写入类型
	for i := 0; i < numCol; i++ {
		bData.WriteString(getTSType(rows[1][i]))
		if i != numCol-1 {
			bData.WriteString(" ")
		} else {
			bData.WriteString("\n")
		}

	}
	//写入字段名
	for i := 0; i < numCol; i++ {
		bData.WriteString(utils.StringToCamel(rows[2][i], false))
		if i != numCol-1 {
			bData.WriteString(" ")
		} else {
			bData.WriteString("\n")
		}
	}
	//写入数据
	for _, row := range rows[3:] {
		numCol := len(row)
		for i, s := range row {
			bData.WriteString(encode(s))
			if i != numCol-1 {
				bData.WriteString(" ")
			} else {
				bData.WriteString("\n")
			}
		}
	}
	return &ParseResult{Define: bDefine.String(), Data: bData.String()}, nil
}

/**
 * 判断数据结构，允许空表
	对于对象数组的结构，表中至少有1列3行，其中第1行为注释，第2行为字段名，第3行为类型，第4行开始为数据
	对于对象的结构，表中至少1行3列，其中第1列为注释，第2列为字段名，第3列为类型，第4列为数据
*/
func getDataStruct(rows [][]string) DataStruct {
	numRow := len(rows)
	if numRow < 1 {
		return UNKNOWN
	}
	numCol := len(rows[0])
	if numCol < 1 {
		return UNKNOWN
	}
	/*到这里，表格中至少1行1列*/

	//判断对象的结构 至少3列至多4列，并且第3列为数据类型关键字
	if numCol >= 3 && numCol <= 4 && isKeywordsOfType(rows[0][1]) {
		return OBJECT
	}
	//判断对象数组的结构 至少3行，并且第3行为数据类型关键字
	if numRow >= 3 && !hasEmptyCell(rows[1]) && !hasEmptyCell(rows[2]) && isKeywordsOfType(rows[1][0]) {
		return OBJECT_ARRAY
	}
	return UNKNOWN
}

func isKeywordsOfType(str string) bool {
	pattern := `^int|float|string|bool|int\[\]|float\[\]|string\[\]|bool\[\]$`
	ret, _ := regexp.MatchString(pattern, str)
	return ret
}
func isValidVariable(str string) bool {
	pattern := `^[a-zA-Z_$]\w*$`
	ret, _ := regexp.MatchString(pattern, str)
	return ret
}
func isEmptyRow(row []string) bool {
	return len(row) == 0
}
func hasEmptyCell(row []string) bool {
	if len(row) == 0 {
		return true
	}
	for _, s := range row {
		if s == "" {
			return true
		}
	}
	return false
}

func getTSType(t string) string {
	return tsTypeMap[t]
}

// 转义特殊字符
func encode(s string) string {
	for k, v := range encodeMap {
		s = strings.Replace(s, k, v, -1)
	}
	return s
}
