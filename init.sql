CREATE DATABASE IF NOT EXISTS challenge01;

USE challenge01;


-- Create the quotes table
CREATE TABLE IF NOT EXISTS quotes (
    id SERIAL PRIMARY KEY,
    name varchar(255),
    bid float,
    ask float
);
