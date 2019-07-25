@echo off
title Export Excel To Text
:: 不存在的文件夹只能有一层，比如想创建 D://a/b/,如果a不存在，a,b都不会创建成功
:: bat 代码创建时编码必须得是GBK,不能转换编码，否的会有各种问题

:: excel目录
set SOURCE_DIR=./tables/
:: 配置表导出目录
set OUT_DIR=./configs/
:: 导出标志，只导出包含标志的列
set EXPORT_FLAG=C
:: 导出的配置表后缀
set EXPORT_SUFFIX=.txt
:: .d.ts文件导出目录
set DEF_OUT_DIR=./def/
:: 配置表数据所在的名字空间
set DEF_NAMESPACE=ConfigData
:: 配置表名称定义文件导出目录【如果ConfigNameOutDir和ConfigNameNamespace不赋值，则不导出】
:: set CONFIG_NAME_OUT_DIR=./def/
:: 配置表名称定义文件名字空间
:: set CONFIG_NAME_NAMESPACE=ConfigName

export.exe %SOURCE_DIR% %OUT_DIR% %EXPORT_FLAG% %EXPORT_SUFFIX% %DEF_OUT_DIR% %DEF_NAMESPACE% %CONFIG_NAME_OUT_DIR% %CONFIG_NAME_NAMESPACE%