@echo off
title Export Excel T
:: 不存在的文件夹只能有一层，比如想创建 D://a/b/,如果a不存在，a,b都不会创建成功

:: excel目录
set SourceDir=./tables/
:: 配置表导出目录
set OutDir=./configs/
:: 导出标志，只导出包含标志的列
set ExportFlag=C
:: 导出的配置表后缀
set ExportSuffix=.txt
:: .d.ts文件导出目录
set DefOutDir=./def/
:: 配置表数据所在的名字空间
set DefNamespace=ConfigData
:: 配置表名称定义文件导出目录(如果不赋值，则不导出)
set ConfigNameOutDir=./def/
:: 配置表名称定义文件名字空间
set ConfigNameNamespace=ConfigName

export.exe %SourceDir% %OutDir% %ExportFlag% %ExportSuffix% %DefOutDir% %DefNamespace% %ConfigNameOutDir% %ConfigNameNamespace%