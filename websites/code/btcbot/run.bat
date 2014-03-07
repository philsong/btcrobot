@echo off

setlocal

if exist run.bat goto ok
echo run.bat must be run from its folder
goto end

:ok

go run src\btcbot\btcrobot.go

echo run successfully

:end