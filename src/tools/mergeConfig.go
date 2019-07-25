package tools

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

var mSetting Setting

func RunMerge(args []string) {
	mSetting = Setting{
		SourceDir:    "./builds/",
		OutDir:       "./merged/",
		ImportSuffix: ".txt",
		ConfigMoveTo: "",
	}
	mergeReadStartupParams(args)
	mergeReadFiles()
}

func mergeReadStartupParams(args []string) {
	if len(args) <= 1 {
		return
	}
	props := []string{"SourceDir", "OutDir", "ImportSuffix", "ConfigMoveTo"}
	begin := time.Now()
	args = args[1:]
	fmt.Println("启动参数：", args)
	ReflectAssign(&mSetting, props, args)
	fmt.Println("设置参数：", mSetting)
	fmt.Println("反射赋值耗时：", time.Since(begin))
}

func mergeReadFiles() {
	pattern := fmt.Sprintf(`^[^~].*(%s)$`, mSetting.ImportSuffix)
	efcFiles := GetEffectiveFiles(mSetting.SourceDir, pattern)
	numConfigs := len(efcFiles)
	fmt.Printf("检测到有 %d 个文件待处理\n", numConfigs)

	if numConfigs > 0 {
		EnsureDir(mSetting.OutDir)
		EnsureDir(mSetting.ConfigMoveTo)

		begin := time.Now()
		filename := "merged"
		mergedFilePath := mSetting.OutDir + filename + mSetting.ImportSuffix
		mergedFile, bWriter := GetBufferWriter(mergedFilePath)
		defer mergedFile.Close()

		var source *os.File
		var err error
		for i, n := range efcFiles {
			srcPath := mSetting.SourceDir + n
			SetReadWrite(srcPath)
			source, err = os.Open(srcPath)

			if err != nil {
				fmt.Printf("无法打开文件 %v\n", err)
				return
			}

			////写入原文件名
			//index := strings.LastIndex(n, ".")
			//_, _ = bWriter.WriteString(n[:index] + "\n")

			//写入文件内容
			bReader := bufio.NewReader(source)
			for {
				buffer := make([]byte, 2048)
				readCount, err := bReader.Read(buffer)
				if err == io.EOF {
					if i < numConfigs-1 {
						_, _ = bWriter.Write([]byte("\n\n"))
					}
					break
				} else {
					_, _ = bWriter.Write(buffer[:readCount])
				}
			}
			SetReadOnly(srcPath)
		}
		defer source.Close()
		_ = bWriter.Flush()

		CopyTo(filename+mSetting.ImportSuffix, mSetting.OutDir, mSetting.ConfigMoveTo)
		SetReadOnly(mergedFilePath)
		SetReadOnly(mSetting.ConfigMoveTo + filename + mSetting.ImportSuffix)
		fmt.Println("合并完成，耗时:", time.Since(begin))
	}
}
