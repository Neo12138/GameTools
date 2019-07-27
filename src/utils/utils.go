/**
 * created by wangcheng at 2019/7/25 15:00
 */
package utils

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"
)

//通过反射赋值
func ReflectAssign(i interface{}, props []string, values []string) {
	refType := reflect.ValueOf(i).Elem()
	for i, v := range values {
		refType.FieldByName(props[i]).SetString(v)
	}
}

/*获取有效文件名*/
func GetEffectiveFiles(path string, pattern string) []string {
	files, _ := ioutil.ReadDir(path)
	efcFiles := make([]string, 0, len(files))

	reg, _ := regexp.Compile(pattern)

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		filename := f.Name()
		match := reg.Match([]byte(filename))
		if match {
			efcFiles = append(efcFiles, filename)
		}
	}

	return efcFiles
}

//复制指定路径的文件到另一个路径
func CopyTo(filename, srcDir string, destDir string) {
	if destDir != "" {
		SetReadWrite(srcDir + filename)
		source, _ := os.Open(srcDir + filename)
		defer source.Close()

		SetReadWrite(destDir + filename)
		dest, _ := os.Create(destDir + filename)
		defer dest.Close()

		_, err := io.Copy(dest, source)
		if err == nil {
			log.Printf("拷贝 %s 到 %s 成功！", filename, destDir)
		} else {
			fmt.Println("copy", err)
		}
	}
}

//判断文件夹是否存在
func IsDirExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return false
}

//确保文件夹存在，如果不存在则创建
func EnsureDir(path string) {
	if path != "" {
		//自动生成文件夹
		if !IsDirExist(path) {
			err := os.Mkdir(path, os.ModePerm)
			fmt.Println("create dir", path)
			CheckError(err)
		}
	}
}

func GetBufferWriter(filePath string) (*os.File, *bufio.Writer) {
	SetReadWrite(filePath)
	file, err := os.Create(filePath)
	fmt.Println("create file", filePath)
	CheckError(err)

	buffer := bufio.NewWriter(file)
	return file, buffer
}

func CheckError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func SetReadOnly(path string) {
	_ = os.Chmod(path, os.FileMode(0400))
}
func SetReadWrite(path string) {
	_ = os.Chmod(path, os.FileMode(0600))
}

func Pause() {
	fmt.Print("输入任意字符结束：")
	var enter string
	_, _ = fmt.Scanln(&enter)
}

//将字符串转换成下划线的格式
func StringToUnderline(str string, upper bool) string {
	data := make([]byte, 0)
	strLength := len(str)
	data = append(data, str[0])
	for i := 1; i < strLength; i++ {
		char := str[i]
		if char >= 'A' && char <= 'Z' && str[i-1] != '_' {
			data = append(data, '_')
		}
		data = append(data, char)
	}

	if upper {
		return strings.ToUpper(string(data))
	} else {
		return strings.ToLower(string(data))
	}
}

// 将字符串转换成驼峰的格式
func StringToCamel(str string, firstUpper bool) string {
	data := make([]byte, 0)
	strLength := len(str)
	var offsetAa uint8 = 32
	if strings.IndexByte(str, '_') >= 0 {
		str = strings.ToLower(str)
	}
	//处理首字母
	firstChar := str[0]
	if firstUpper && firstChar >= 'a' && firstChar <= 'z' {
		firstChar -= offsetAa
	}
	data = append(data, firstChar)

	for i := 1; i < strLength; i++ {
		char := str[i]
		if char == '_' {
			i++
			char = str[i] - offsetAa
		}
		data = append(data, char)
	}

	return string(data)
}
