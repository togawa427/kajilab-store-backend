#!/bin/bash

# 現在の日付を取得
DATE=$(date +"%Y-%m-%d_%H-%M")

# SQLiteデータベースのバックアップ
sqlite3 kajilabstore.db ".backup backup/db_backup_$DATE.bk"