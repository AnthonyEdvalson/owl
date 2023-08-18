# Shell script for Linux
rm -rf bin
mkdir bin
go build -o bin
cp -r lib bin/lib
