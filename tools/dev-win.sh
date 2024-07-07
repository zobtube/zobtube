#!/bin/bash

export ZT_DB_DRIVER="sqlite"
export ZT_DB_CONNSTRING="zt.db"
export ZT_MEDIA="test_data"
exec air --build.cmd 'go build -o tmp/zt.exe cli/main.go' --build.bin 'tmp\zt.exe' --build.log 'tmp\logs.txt'
