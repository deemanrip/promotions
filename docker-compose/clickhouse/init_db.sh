#!/bin/bash
set -e

clickhouse client -n <<-EOSQL
    CREATE DATABASE IF NOT EXISTS promotions ENGINE = Atomic;
    CREATE TABLE IF NOT EXISTS promotions.promotions
    (
      id UUID,
      price Decimal64(6),
      expiration_date String
    ) ENGINE = MergeTree()
    ORDER BY id;
    CREATE TABLE IF NOT EXISTS promotions.promotions_tmp
    (
      id UUID,
      price Decimal64(6),
      expiration_date String
    ) ENGINE = MergeTree()
    ORDER BY id;
    CREATE USER IF NOT EXISTS app IDENTIFIED WITH plaintext_password BY 'test12345';
    GRANT SELECT, INSERT, DROP TABLE, CREATE TABLE ON promotions.* TO app WITH GRANT OPTION;
    GRANT CREATE TEMPORARY TABLE, S3 ON *.* TO app;
    GRANT TRUNCATE ON promotions.promotions_tmp TO app
EOSQL