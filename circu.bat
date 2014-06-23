echo off

setlocal

if exist circu.bat goto ok
echo circu.bat must be run from its folder
goto end

:ok

:: stop
taskkill /im btcrobot.exe /f
del /q /f /a pid

Sleep 3

:: start
start /b bin\btcrobot

Sleep 3600

goto ok

echo restart successfully

:end