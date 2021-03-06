#!/bin/bash

if [ "$1" = "clean" ]; then
    # clean driver
    ninja -C driver/build clean
    
    # clean site
    rm static/dist/*.js

    # clean server
    go clean
else
    # build driver
    ninja -C driver/build

    # build site
    cd static
    npm run build
    cd ..

    # build server
    go build
fi

