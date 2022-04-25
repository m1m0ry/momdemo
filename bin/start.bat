@echo off
echo nsq
start /b bin\nsqlookupd
start /b bin\nsqd --lookupd-tcp-address=127.0.0.1:4160
start /b bin\nsqadmin --lookupd-http-address=127.0.0.1:4161
echo web
start /b random
start /b mean
start /b variance
start /b maxAndmin
start /b chart