@echo off
mode con cols=80 lines=30
title=看门狗DDNS
color 3f
echo 》DDNS服务正在启动，请稍后
echo 》                           
echo 》=================================================================
echo 》                   欢迎使用看门狗DDNS
echo 》                        程序已运行
echo 》=================================================================
echo 》                                     
call ddns-client.exe
