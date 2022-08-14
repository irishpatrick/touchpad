#!/bin/bash

# install dependencies
sudo apt install libevdev-dev libgtk-3-dev
git submodule update --init --recursive --remote

# bootstrap driver
cd driver
meson build
cd ..

# bootstrap frontend
cd frontend
meson build
cd ..

# bootstrap site
cd static
npm install
cd ..

