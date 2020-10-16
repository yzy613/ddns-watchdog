@echo off
mode con cols=65 lines=30
title=ddns-watchdog-client-startup-script
color 3f

SET conf=.\conf
if not exist %conf% (call ddns-watchdog-client.exe -init 0123) else goto start
echo ddns-watchdog-client has initialized
echo Please change .\conf\client.json
echo.
echo Press any key to contiune ...
pause>nul
cls

:start
echo ddns-watchdog-client is starting...
echo.
echo =================================================================
echo               Welcome to use ddns-watchdog-client
echo                  ddns-watchdog-client has run
echo =================================================================
echo.
call ddns-watchdog-client.exe
echo.
echo ddns-watchdog-client has exited.
echo Press any key to restart ddns-watchdog ...
pause>nul
cls
goto start
