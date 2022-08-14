#!/bin/bash

tick="$(date -u +%s)"

# install dependencies
sudo apt install libevdev-dev libgtk-3-dev
git submodule update --init --recursive --remote

# bootstrap driver
cd driver
meson build-debug --buildtype debug
meson build-release --buildtype release
cd ..

# bootstrap frontend
cd frontend
meson build-debug --buildtype debug
meson build-release --buildtype release
cd ..

# bootstrap site
cd static
npm install
cd ..

tock="$(date -u +%s)"
elapsed="$(($tock-$tick))"
echo "finished in $elapsed seconds"

