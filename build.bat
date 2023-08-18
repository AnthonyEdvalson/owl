@echo off


del /s /q bin
mkdir bin
go build -o bin
xcopy lib bin\lib /s /e /y