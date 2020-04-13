/*常用的字符串处理方法*/
package utils

import (
	"strings"
)

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
