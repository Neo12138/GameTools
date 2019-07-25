package tools

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"sort"
	"strings"
	"time"
)

func RunGetConfigDef(args []string) {
	setting = Setting{
		SourceDir:        "./builds/,./builds2/",
		DefOutDir:        "./def/",
		ImportSuffix:     ".json",
		ConfigNamespace:  "ConfigData",
		ConfigNameOutDir: "./def/",
	}
	getDefReadStartupParams(args)
	readJSONFiles()
}

func getDefReadStartupParams(args []string) {
	if len(args) <= 1 {
		return
	}
	props := []string{"SourceDir", "DefOutDir", "ConfigNamespace", "ConfigNameOutDir"}
	args = args[1:]
	fmt.Println("启动参数：", args)
	ReflectAssign(&setting, props, args)
}

var dtsConfigDeclare bytes.Buffer
var dtsBuffer2 *bufio.Writer
var cn *os.File
var cnBuffer *bufio.Writer
var failCount2 = 0

func readJSONFiles() {
	pattern := fmt.Sprintf(`^[^~].*(%s)$`, setting.ImportSuffix)

	dirs := strings.Split(setting.SourceDir, ",")

	paths := make([]string, 0)
	names := make([]string, 0)
	length := 0
	for _, dir := range dirs {
		fmt.Println("read dir", dir)
		efcFiles := GetEffectiveFiles(dir, pattern)
		sort.Strings(efcFiles)
		p2 := make([]string, len(efcFiles))
		n2 := make([]string, len(efcFiles))
		paths = append(paths, p2...)
		names = append(names, n2...)
		//处理每个json
		for _, file := range efcFiles {
			paths[length] = dir + file
			names[length] = strings.TrimSuffix(file, setting.ImportSuffix)
			length++
		}
	}

	fmt.Printf("检测到有 %d 个json文件待处理\n", length)
	if length == 0 {
		return
	}

	//正式开始处理
	begin := time.Now()

	EnsureDir(setting.DefOutDir)

	dts, dtsBuffer2 = GetBufferWriter(setting.DefOutDir + "config.d.ts")
	defer dts.Close()
	_, _ = dtsBuffer2.WriteString("//由工具自动生成，请勿手动修改\n")
	_, _ = dtsBuffer2.WriteString("declare namespace " + setting.ConfigNamespace + " {")
	indent(1)

	cn, cnBuffer = GetBufferWriter(setting.ConfigNameOutDir + "config-names.ts")
	defer cn.Close()
	_, _ = cnBuffer.WriteString("//由工具自动生成，请勿手动修改\n")
	_, _ = cnBuffer.WriteString("namespace ConfigName {")

	for i, path := range paths {
		parseJSON(path, names[i])
	}

	_, _ = dtsBuffer2.WriteString(dtsConfigDeclare.String())
	_, _ = dtsBuffer2.WriteString("\n}")
	_ = dtsBuffer2.Flush()

	_, _ = cnBuffer.WriteString("\n}")
	_ = cnBuffer.Flush()

	CopyTo("config.d.ts", setting.DefOutDir, setting.DefMoveTo)

	fmt.Printf("所有json文件处理完成。成功%d个,失败%d个，总任务%d个，耗时：%s\n", length-failCount, failCount, length, time.Since(begin))
	if failCount > 0 {
		fmt.Println("错误项：\n", failDesc)
	}
}

func parseJSON(filename string, name string) {
	fmt.Println("\t start parse:", name, filename)
	byteValue, err := ioutil.ReadFile(filename)
	//去除bom标志
	byteValue = bytes.TrimPrefix(byteValue, []byte("\xef\xbb\xbf"))
	if err != nil {
		fmt.Println("read file error", err)
		failCount++
		failDesc = name + ": " + err.Error()
		return
	}

	var result []interface{}
	err = json.Unmarshal(byteValue, &result)

	if err != nil {
		fmt.Println("json parse error", err)
		failCount++
		failDesc += name + ": " + err.Error() + "\n"
		return
	}

	for _, v := range result {
		o := v.(map[string]interface{})

		//写入到dts文件
		if dtsBuffer2 != nil {
			_, _ = dtsBuffer2.WriteString(tab + "interface I" + name + " {")
			indent(1)
			for k, p := range o {
				_, _ = dtsBuffer2.WriteString(tab + "readonly " + k + ": " + getTargetLangType(reflect.TypeOf(p).String()) + ";")
			}
			_, _ = dtsBuffer2.WriteString(indent(-1) + "}\n")
		}
		dtsConfigDeclare.WriteString(tab + "const " + name + ": I" + name + "[];")
		_, _ = cnBuffer.WriteString("\n\texport const " + name + ": string = \"" + name + "\";")

		break
	}

}
