@echo off
mode con cols=65 lines=30
title=watchdog-ddns-client-startup-script
color 3f

SET conf=.\conf
if not exist %conf% (call watchdog-ddns-client -ini)else goto start
echo.
echo Program Init.
echo Please change .\conf\client.json
echo Press any key to contiune ...
pause>nul
cls

:start
echo watchdog-ddns-client is starting...
echo.
echo =================================================================
echo               Welcome to use watchdog-ddns-client
echo                  watchdog-ddns-client has run
echo =================================================================
echo.
call watchdog-ddns-client.exe
echo.
echo Program exit. Press any key to restart watchdog-ddns ...
pause>nul
cls
goto start
