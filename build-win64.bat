chcp 65001

@echo off


echo 设置win64编译环境

set GOOS=windows
set GOARCH=amd64

echo 开始执行编译

go build -o ym-img-caption.exe .

echo 编译结束

pause