#!/bin/bash

# 現在の日付を取得
# DATE=$(date +"%Y-%m-%d_%H-%M")
OUTPUT_FILE=$1

# SQLiteデータベースのバックアップ
# sqlite3 kajilabstore.db ".backup backup/db_backup_$DATE.bk"
sqlite3 kajilabstore.db ".dump" > $OUTPUT_FILE

# SQLiteのバックアップをリストア
#sqlite3 kajilabstore.db ".restore backup/db_backup_2024-09-05_16-56.bk"
# sqlite3 kajilabstore.db "alter table products drop safety_stock"