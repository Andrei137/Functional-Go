@echo off

cd ..

if not exist bin (
    echo Creating bin directory...
    mkdir bin
)

echo Compiling...
go build -o bin/main.exe src/main.go
echo Compiled successfully!

cd scripts
endlocal
