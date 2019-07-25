/**
 * 将excel配置表导出为纯文本格式配置表
 */
package main

import (
	"export"
	"os"
	"utils"
)

func main() {
	export.Run(os.Args)
	utils.Pause()
}
