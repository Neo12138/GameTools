/*常用的文件处理方法*/
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

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
			_ = os.Mkdir(path, os.ModePerm)
			fmt.Println("create dir", path)
		}
	}
}

func GetFillAndBufferWriter(filePath string) (*os.File, *bufio.Writer) {
	SetReadWrite(filePath)
	file, _ := os.Create(filePath)
	fmt.Println("create file", filePath)

	buffer := bufio.NewWriter(file)
	return file, buffer
}

//设置文件只读
func SetReadOnly(path string) {
	_ = os.Chmod(path, os.FileMode(0400))
}

//设置文件只写
func SetReadWrite(path string) {
	_ = os.Chmod(path, os.FileMode(0600))
}

//获取当前路径，比如：E:/abc/data/test
func GetCurrentPath() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
func StringifyMap(r interface{}) string {
	var b bytes.Buffer
	s, e := json.Marshal(r)
	if e == nil {
		_ = json.Indent(&b, s, "", "\t")
		return b.String()
	}
	return ""
}
func ParseMap(str []byte, v interface{}) error {
	return json.Unmarshal(str, &v)
}

/*常用的字符串处理方法*/

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

type Tab struct {
	Content string
}

func (s *Tab) Indent(delta int) string {
	if s.Content == "" {
		s.Content = "\n"
	}
	if delta > 0 {
		s.Content += "\t"
	} else {
		numTab := len(s.Content)
		if numTab >= 2 {
			s.Content = s.Content[:numTab-1]
		} else {
			s.Content = "\n"
		}
	}

	return s.Content
}
