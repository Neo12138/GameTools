@echo off
title Export Excel T
:: �����ڵ��ļ���ֻ����һ�㣬�����봴�� D://a/b/,���a�����ڣ�a,b�����ᴴ���ɹ�

:: excelĿ¼
set SourceDir=./tables/
:: ���ñ���Ŀ¼
set OutDir=./configs/
:: ������־��ֻ����������־����
set ExportFlag=C
:: ���������ñ��׺
set ExportSuffix=.txt
:: .d.ts�ļ�����Ŀ¼
set DefOutDir=./def/
:: ���ñ��������ڵ����ֿռ�
set DefNamespace=ConfigData
:: ���ñ����ƶ����ļ�����Ŀ¼(�������ֵ���򲻵���)
set ConfigNameOutDir=./def/
:: ���ñ����ƶ����ļ����ֿռ�
set ConfigNameNamespace=ConfigName

export.exe %SourceDir% %OutDir% %ExportFlag% %ExportSuffix% %DefOutDir% %DefNamespace% %ConfigNameOutDir% %ConfigNameNamespace%