@echo off
mode con cols=65 lines=30
title=watchdog-ddns-client-startup-script
color 3f
echo watchdog-ddns-client is starting...
echo.
echo =================================================================
echo               Welcome to use watchdog-ddns-client
echo                  watchdog-ddns-client has run
echo =================================================================
echo.
call watchdog-ddns-client.exe
echo.
pause
