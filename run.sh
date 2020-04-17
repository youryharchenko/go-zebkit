#!/bin/bash

JWTSECRET="secret+secret" ./server/cmd/server \
    --dir=$PWD/static \
    --db=$PWD/server/db/test.db \
    --query=$PWD/server/query

