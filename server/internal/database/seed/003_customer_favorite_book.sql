INSERT INTO customer_favorite_book (
    customer_favorite_book_id,
    customer_id,
    book_id
) VALUES (
    '01a00132-2d61-767b-8305-b678932bc209',
    '019f0866-5d9e-7c53-a3ee-f39fbd967ac4',
    '019fea73-27ca-77b2-9f80-e90d5d5aaa09'
), (
    '01a00132-2d62-7f7a-b093-2b2b4c10686e',
    '019f0866-5d9e-7c53-a3ee-f39fbd967ac4',
    '01a0012c-a774-7f0c-b04e-c8a551ccfe94'
), (
    '01a00132-2d63-75aa-af0b-6532998716eb',
    '019f0866-5d9e-7c53-a3ee-f3ad58ebe5ff',
    '01a0012c-a772-70bd-9856-a094176e29a0'
), (
    '01a00132-2d64-73ce-9b8e-0885c80e4a31',
    '019f0866-5d9e-7c53-a3ee-f3ad58ebe5ff',
    '01a0012c-a77c-79d6-b80f-a448cda12866'
), (
    '01a00132-2d65-7219-b4fc-6dac2035feed',
    '019f0866-5d9e-7c53-a3ee-f3b055baa122',
    '01a0012c-a777-7979-9544-941a23ddc0b7'
), (
    '01a00132-2d66-71fe-a49a-74b48666a524',
    '019f0866-5d9e-7c53-a3ee-f3b055baa122',
    '01a0012c-a781-752c-9d9f-4b4ad40c6fe8'
), (
    '01a00132-2d67-7796-97f8-fadcc927135b',
    '019f0866-5d9e-7c53-a3ee-f3cca52e0bc1',
    '01a0012c-a779-79f3-abbc-101f3daacdc4'
), (
    '01a00132-2d68-78fa-ab4d-cd4f0209d9e7',
    '019f0866-5d9e-7c53-a3ee-f3d9eef4cae7',
    '01a0012c-a77e-7302-9932-7012a4d549bc'
), (
    '01a00132-2d69-7924-9d72-3d6a8a8d7f2f',
    '019f0866-5d9e-7c53-a3ee-f3d9eef4cae7',
    '01a0012c-a783-7905-8196-aeedbd9ac29b'
) ON DUPLICATE KEY UPDATE
customer_id = VALUES (customer_id),
book_id = VALUES (book_id);
