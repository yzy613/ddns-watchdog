@echo off
mode con cols=65 lines=30
title=watchdog-ddns-client-startup-script
color 3f

SET conf=.\conf
if not exist %conf% (call watchdog-ddns-client -init) else goto start
echo watchdog-ddns-client has initialized
echo Please change .\conf\client.json
echo.
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
echo watchdog-ddns-client has exited.
echo Press any key to restart watchdog-ddns ...
pause>nul
cls
goto start
