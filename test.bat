@echo off

setlocal

if exist test.bat goto ok
echo test.bat must be run from its folder
goto end

:ok

start /b bin\backtest

:end