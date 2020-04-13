@echo off
go build -ldflags="-s -w -H windowsgui" -o app.exe
pause
