@echo off
mode con cols=80 lines=30
title=watchdog-ddns
color 3f
echo 》watchdog-ddns 正在启动，请稍后...
echo 》                           
echo 》=================================================================
echo 》                   欢迎使用 watchdog-ddns
echo 》                        程序已运行
echo 》=================================================================
echo 》                                     
call ddns-client.exe
echo.
pause
