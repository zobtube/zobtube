#!/bin/sh

export ZT_SERVER_BIND='127.0.0.1:8080'
export ZT_DB_DRIVER='sqlite'
export ZT_DB_CONNSTRING='/tmp/zt-db.sqlite3'
export ZT_MEDIA_PATH='/tmp/zt-data'

echo 'delete existing database'
rm -f $ZT_DB_CONNSTRING

echo 'run zobtube to generate database'
timeout 5 /tmp/zt

echo 'insert fake user'
sqlite3 $ZT_DB_CONNSTRING "insert into users values ('b23f4f4a-1c5c-11f0-8822-305a3a05e04d', date('now'), date('now'), null, 'test', 'invalid', 1);"

echo 'done'
