@echo off

git submodule update --init --recursive --remote

cd driver
mkdir build-debug
cd build-debug
cmake .. -GNinja -DCMAKE_BUILD_TYPE=Debug -DMINGW=ON
cd ..

mkdir build-release
cd build-release
cmake .. -GNinja  -DCMAKE_BUILD_TYPE=Debug -DMINGW=ON
cd ..

cd ..

cd static
npm install
cd ..
