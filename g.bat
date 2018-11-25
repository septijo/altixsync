:start
go build -o altixSync.exe
@if %errorlevel% equ 1 goto rebuild
@if %errorlevel% equ 2 goto rebuild
@echo %errorlevel% ... build success ! press any key to run altixSync.exe
@pause > nul
cls
@altixSync.exe
@pause > nul
goto start
:rebuild
@echo .. error on building .exe .. press any key to rebuild
@pause > nul
goto start
