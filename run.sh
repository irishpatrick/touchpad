#!/bin/bash

./frontend/build/frontend > /dev/null & disown
LD_LIBRARY_PATH=`pwd`/driver/build-debug ./touchpad $@

