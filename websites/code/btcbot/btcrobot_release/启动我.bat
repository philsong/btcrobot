@echo off

setlocal

if exist 启动我.bat goto ok
echo 启动我.bat must be run from its folder
goto end

:ok

start /b bin\btcrobot &
Sleep 3
start http://localhost:9090
echo start successfully

:end