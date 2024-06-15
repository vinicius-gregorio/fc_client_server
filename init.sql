CREATE DATABASE IF NOT EXISTS challenge01;

USE challenge01;

-- Create the quotes table
CREATE TABLE IF NOT EXISTS quotes (
    id SERIAL PRIMARY KEY,
    code varchar(255),
    codein varchar(255),
    name varchar(255),
    high float,
    low float,
    varBid float,
    pctChange float,
    bid float,
    ask float,
    timestamp varchar(255),
    create_date date
);
