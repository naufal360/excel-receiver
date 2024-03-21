CREATE DATABASE IF NOT EXISTS excel_rec_db;
USE excel_rec_db;

CREATE TABLE IF NOT EXISTS token (
    id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    token VARCHAR(255) NOT NULL,
    expired_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE IF NOT EXISTS request (
    id int NOT NULL AUTO_INCREMENT PRIMARY KEY,
    request_id VARCHAR(255) NOT NULL,
    status VARCHAR(255) NOT NULL,
    file_path VARCHAR(255) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO token (token, expired_at, created_at) VALUES(
    "token_test",
    "2024-03-30 12:59:56",
    "2024-03-17 12:59:56"
),(
    "expired_token",
    "2024-03-17 12:59:56",
    "2024-03-17 12:59:56"
);