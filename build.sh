#!/bin/bash

tick="$(date -u +%s)"

if [ "$1" = "clean" ]; then
    # clean driver
    ninja -C driver/build-debug clean
    ninja -C driver/build-release clean

    # clean frontend
    ninja -C frontend/build-debug clean
    ninja -C frontend/build-release clean
    
    # clean site
    rm static/dist/*.js

    # clean server
    go clean
elif [ "$1" = "prod" ]; then
    # build driver
    ninja -C driver/build-release driver install

    # build frontend
    ninja -C frontend/build-release

    # build site
    cd static
    npm run build-prod
    cd ..

    # build server
    cp -r static/dist server/server
    cd server
    echo 'building server...'
    go clean
    go build -tags release -ldflags "-s -w"
    cd ..
    rm -rf server/server/dist
else
    # build driver
    ninja -C driver/build-debug driver install

    # build frontend
    ninja -C frontend/build-debug

    # build site
    cd static
    npm run build
    cd ..

    # build server
    echo 'building server...'
    go clean
    go build -tags debug
fi

tock="$(date -u +%s)"
elapsed="$(($tock-$tick))"
echo "finished in $elapsed seconds"
