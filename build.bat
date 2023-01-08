@echo off

if "%1%"=="clean" goto CLEAN
if "%1%"=="build" goto BUILD
if "%1%"=="prod"  goto PROD
goto BUILD

:CLEAN
ninja -C driver/build-debug clean
ninja -C driver/build-release clean
del static\dist\*.js
cd server
go clean
cd ..
goto EXIT

:BUILD
ninja -C driver/build-debug driver install
call npm run --prefix static build

cd server
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=1
echo building server...
go clean
go build -tags debug
cd ..
goto EXIT

:PROD
ninja -C driver/build-release driver install
call npm run --prefix static build-prod

cp -r static/dist server/server
cd server
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=1
echo building production server...
go clean
go build -tags release -ldflags "-s -w"
cd ..
rmdir /Q /S .\server\server\dist\
goto EXIT

:EXIT
echo done!