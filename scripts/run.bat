@echo off
setlocal enabledelayedexpansion

set build=False
set clear=False

:parse_flags
if "%1"=="" goto end_parse
set flag=%1
if "%flag:~0,1%"=="-" (
    set flag=%flag:~1%
) else (
    goto help_message
)
for /l %%i in (0,1,2) do (
    set "char=!flag!"
    set "char=!char:~%%i,1!"
    if "!char!"=="" goto end_parse
    if "!char!"=="b" (
        set build=True
    ) else if "!char!"=="c" (
        set clear=True
    ) else (
        goto help_message
    )
)
shift
goto parse_flags
:end_parse

if %build%==True (
    if %clear%==True (
        cls && call build.bat %type% && cls
    ) else (
        call build.bat %type%
    )
) else (
    if %clear%==True (
        cls
    )
)

cd ../bin
echo Running...
main.exe
endlocal
exit /b 0

:help_message
echo Usage: %0 [PARAMS]
echo Param:
echo  -h     Display this help message
echo  -b     Build the project
echo  -c     Clear the console
exit /b
