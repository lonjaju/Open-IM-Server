#!/usr/bin/env bash
minio server /Users/loong/data/minio/shared   --address ":39005" > logs/mino-server.log 2>&1 &