@echo off

setlocal

if exist getpkg.bat goto ok
echo getpkg.bat must be run from its folder
goto end

:ok

set OLDGOPATH=%GOPATH%
set GOPATH=%~dp0

go get -u github.com/go-martini/martini
go get -u github.com/codegangsta/martini-contrib/auth

set GOPATH=%OLDGOPATH%

:end
echo finished