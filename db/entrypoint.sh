#!/bin/sh

cd /db/
./test_data.py
sleep 15
psql -h grants-database -p 5432 -U grants -d grants -f test_data.sql