INSERT INTO customer (
    customer_id,
    nickname,
    birthday,
    location,
    registered_at
) VALUES (
    '019f0866-5d9e-7c53-a3ee-f39fbd967ac4',
    'たろう',
    '1998-01-23',
    'Asia/Tokyo',
    NOW()
), (
    '019f0866-5d9e-7c53-a3ee-f3ad58ebe5ff',
    'はなこ',
    '1989-11-11',
    'Asia/Tokyo',
    NOW()
), (
    '019f0866-5d9e-7c53-a3ee-f3b055baa122',
    'けんじ',
    '1973-09-23',
    'Asia/Tokyo',
    NOW()
), (
    '019f0866-5d9e-7c53-a3ee-f3cca52e0bc1',
    '',
    '1-01-01',
    'Asia/Tokyo',
    NULL
), (
    '019f0866-5d9e-7c53-a3ee-f3d9eef4cae7',
    'みさき',
    '1999-12-31',
    'Asia/Tokyo',
    NOW()
) ON DUPLICATE KEY UPDATE
nickname = VALUES (nickname),
birthday = VALUES (birthday),
location = VALUES (location),
registered_at = VALUES (registered_at);
