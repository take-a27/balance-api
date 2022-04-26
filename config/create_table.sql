CREATE DATABASE balance;
USE balance;

CREATE TABLE IF NOT EXISTS user_balance (
      id varchar(50) PRIMARY KEY NOT NULL,
      amount integer NOT NULL
);

CREATE TABLE IF NOT EXISTS balance_operation (
    id varchar(50) PRIMARY KEY NOT NULL,
    user_id varchar(50) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
    timestamp datetime NOT NULL,
    difference integer NOT NULL,
    is_transaction boolean NOT NULL
);
