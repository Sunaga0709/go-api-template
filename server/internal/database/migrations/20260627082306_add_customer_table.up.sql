CREATE TABLE IF NOT EXISTS customer (
    customer_id CHAR(36) PRIMARY KEY COMMENT '顧客ID',
    nickname VARCHAR(255) NOT NULL COMMENT 'ニックネーム',
    birthday DATE NOT NULL COMMENT '生年月日',
    location VARCHAR(32) NOT NULL COMMENT 'ロケーション',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '作成日時'
);

CREATE INDEX customer_idx_birthday ON customer (birthday);
