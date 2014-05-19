@echo off

setlocal

if exist install.bat goto ok
echo install.bat must be run from its folder
goto end

:ok

set OLDGOPATH=%GOPATH%
set GOPATH=%~dp0

go install github.com/go-martini/martini
go install github.com/codegangsta/martini-contrib/auth
go install github.com/syndtr/goleveldb/leveldb

set GOPATH=%OLDGOPATH%

:end
echo finished