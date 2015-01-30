@echo off

setlocal

if exist getpkg.bat goto ok
echo getpkg.bat must be run from its folder
goto end

:ok

set OLDGOPATH=%GOPATH%
set GOPATH=%~dp0

go get -u github.com/philsong/martini
go get -u github.com/codegangsta/martini-contrib/auth
go get -u github.com/philsong/goleveldb/leveldb
go get -u github.com/bitly/go-simplejson
go get -u github.com/philsong/go-bittrex
go get -u github.com/gorilla/websocket

set GOPATH=%OLDGOPATH%

:end
echo finished
