:: 不存在的文件夹只能有一层，比如想创建 D://a/b/,如果a不存在，a,b都不会创建成功
:: "SourceDir", "OutDir", "ExportFlag", "ExportSuffix", "DefOutDir", "ConfigMoveTo", "DefMoveTo"
:: 								     E:/Configs/     E:/Def/

:: start
export.exe ./tables/ ./builds/ C .csv ./def/  