@echo off
title Export Excel To Text
:: �����ڵ��ļ���ֻ����һ�㣬�����봴�� D://a/b/,���a�����ڣ�a,b�����ᴴ���ɹ�
:: bat ���봴��ʱ����������GBK,����ת�����룬��Ļ��и�������

:: excelĿ¼
set SOURCE_DIR=./tables/
:: ���ñ���Ŀ¼
set OUT_DIR=./configs/
:: ������־��ֻ����������־����
set EXPORT_FLAG=C
:: ���������ñ��׺
set EXPORT_SUFFIX=.txt
:: .d.ts�ļ�����Ŀ¼
set DEF_OUT_DIR=./def/
:: ���ñ��������ڵ����ֿռ�
set DEF_NAMESPACE=ConfigData
:: ���ñ����ƶ����ļ�����Ŀ¼�����ConfigNameOutDir��ConfigNameNamespace����ֵ���򲻵�����
:: set CONFIG_NAME_OUT_DIR=./def/
:: ���ñ����ƶ����ļ����ֿռ�
:: set CONFIG_NAME_NAMESPACE=ConfigName

export.exe %SOURCE_DIR% %OUT_DIR% %EXPORT_FLAG% %EXPORT_SUFFIX% %DEF_OUT_DIR% %DEF_NAMESPACE% %CONFIG_NAME_OUT_DIR% %CONFIG_NAME_NAMESPACE%