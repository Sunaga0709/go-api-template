ALTER TABLE customer_favorite_book
DROP FOREIGN KEY cfb_fk_customer_id,
DROP FOREIGN KEY cfb_fk_book_id;

DROP INDEX cfb_idx_customer_id_book_id ON customer_favorite_book;

DROP TABLE IF EXISTS customer_favorite_book;
