#!/bin/sh

export ZT_SERVER_BIND='127.0.0.1:8069'
export ZT_DB_DRIVER='sqlite'
export ZT_DB_CONNSTRING='/tmp/zt-db.sqlite3'
export ZT_MEDIA_PATH='/tmp/zt-data'

echo 'delete existing database'
rm -f $ZT_DB_CONNSTRING

echo 'run zobtube to generate database'
timeout 5 /tmp/zt

set -e

echo 'insert default configuration with authentication'
sqlite3 $ZT_DB_CONNSTRING "update configurations set user_authentication = 1;"

echo 'insert fake users'
sqlite3 $ZT_DB_CONNSTRING "insert into users values ('b23f4f4a-1c5c-11f0-8822-305a3a05e04d', date('now'), date('now'), null, 'validation', '98c41dcd20b86b86830ec0794559835614458ceaae0f0ec77a3ed1cd3a1f7d55', 1);"
sqlite3 $ZT_DB_CONNSTRING "insert into users values ('21e55ff0-1dc1-11f0-9c1f-305a3a05e04d', date('now'), date('now'), null, 'non-admin', '030d96618d48820fd7ad11e5fd465972b013aed8fdd2bdfc7b02d979a2d4be98', 0);"

echo 'insert fake actor'
sqlite3 $ZT_DB_CONNSTRING "insert into actors values ('045e1b0e-1dc4-11f0-a04a-305a3a05e04d', date('now'), date('now'), null, 'test', 0, 'f');"

echo 'insert fake channel'
sqlite3 $ZT_DB_CONNSTRING "insert into channels values ('8c50735e-1dc4-11f0-b1fc-305a3a05e04d', date('now'), date('now'), null, 'test', 0);"

echo 'insert fake video'
sqlite3 $ZT_DB_CONNSTRING "insert into videos (id, created_at, updated_at, deleted_at, name, filename, thumbnail, thumbnail_mini, duration, type, imported, status, channel_id) values ('d8045d56-1dc4-11f0-9970-305a3a05e04d', date('now'), date('now'), null, 'test', 'test-filename', 0, 0, 0, 'v', 1, 'ready', null);"

echo 'copy fake video'
mkdir -p $ZT_MEDIA_PATH/videos/d8045d56-1dc4-11f0-9970-305a3a05e04d
cp test/video/Big_Buck_Bunny_360_10s_1MB.mp4 $ZT_MEDIA_PATH/videos/d8045d56-1dc4-11f0-9970-305a3a05e04d/video.mp4

echo 'done'
