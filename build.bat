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
cd static
call npm run build
cd ..
cd server
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=1
echo building server...
go build
cd ..
goto EXIT

:PROD
ninja -C driver/build-release driver install
cd static
call npm run build-prod
cd ..
cd server
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=1
echo building production server...
go build -ldflags "-s -w"
cd ..
goto EXIT

:EXIT
echo done!