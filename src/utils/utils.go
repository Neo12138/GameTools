/**
 * created by wangcheng at 2019/7/25 15:00
 */
package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
)

//通过反射赋值
func ReflectAssign(i interface{}, props []string, values []string) {
	refType := reflect.ValueOf(i).Elem()
	for i, v := range values {
		refType.FieldByName(props[i]).SetString(v)
	}
}

func CheckError(err error) {
	if err != nil {
		log.Println(err)
	}
}

func Pause() {
	fmt.Print("输入任意字符结束：")
	var enter string
	_, _ = fmt.Scanln(&enter)
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
