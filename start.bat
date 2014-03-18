@echo off

setlocal

if exist start.bat goto ok
echo start.bat must be run from its folder
goto end

:ok

start /b bin\btcrobot
Sleep 3
start http://localhost:9090
echo start successfully

:end