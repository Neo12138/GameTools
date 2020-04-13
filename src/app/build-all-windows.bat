@echo off
go generate -ldflags="-s -w"
go build -ldflags="-s -w -H windowsgui" -o app.exe
pause
