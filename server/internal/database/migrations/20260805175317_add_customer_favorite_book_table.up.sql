CREATE TABLE IF NOT EXISTS customer_favorite_book (
    customer_favorite_book_id CHAR(36) PRIMARY KEY COMMENT '顧客お気に入り書籍ID',
    customer_id CHAR(36) NOT NULL COMMENT '顧客ID',
    book_id CHAR(36) NOT NULL COMMENT '書籍ID',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '作成日時',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新日時',

    CONSTRAINT cfb_fk_customer_id FOREIGN KEY (customer_id) REFERENCES customer (customer_id) ON DELETE CASCADE,
    CONSTRAINT cfb_fk_book_id FOREIGN KEY (book_id) REFERENCES book (book_id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX cfb_idx_customer_id_book_id ON customer_favorite_book (customer_id, book_id);
