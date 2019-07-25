package tools

import (
	"fmt"
	"os"
	"strings"
)

func RunReanme(args []string) {
	fmt.Println(args, len(args))
	if len(args) > 2 {
		rename(args[1], args[2])
	}
}

func rename(dir string, suffix string) {
	pattern := fmt.Sprintf(`^[^~].*(%s)$`, suffix)
	filenames := GetEffectiveFiles(dir, pattern)
	for _, fn := range filenames {
		oldName := fn
		newName := strings.ToLower(oldName)

		err := os.Rename(dir+oldName, dir+newName)

		if err != nil {
			fmt.Println("rename error", err)
			fmt.Printf("rename %s => %s fail!\n", oldName, newName)
			continue
		} else {
			fmt.Printf("rename %s => %s success!\n", oldName, newName)
		}
	}
}
