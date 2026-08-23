INSERT INTO customer_favorite_product (
    customer_favorite_product_id,
    customer_id,
    product_id
) VALUES (
    '019f0e8b-cdab-7d43-b0a9-5105344b1614',
    '019f0866-5d9e-7c53-a3ee-f39fbd967ac4',
    '019ef097-4931-7c91-ac71-8b8862b3e9d9'
), (
    '019f0e8b-cdb3-7be2-9dfa-25de4cc85901',
    '019f0866-5d9e-7c53-a3ee-f39fbd967ac4',
    '019ef097-4931-7c91-ac71-8bc21eb34f8d'
), (
    '019f0e8b-cdb9-7710-b228-9ef487368cc5',
    '019f0866-5d9e-7c53-a3ee-f3ad58ebe5ff',
    '019ef097-4931-7c91-ac71-8b9c2ac94c3a'
), (
    '019f0e8b-cdbe-79e2-9a1c-ab8dc3dad8c8',
    '019f0866-5d9e-7c53-a3ee-f3ad58ebe5ff',
    '019ef097-4931-7c91-ac71-8bb7ce91e63b'
), (
    '019f0e8b-cdc3-72c1-976c-5abbd6855f2e',
    '019f0866-5d9e-7c53-a3ee-f3ad58ebe5ff',
    '019ef097-4931-7c91-ac71-8bd1d0f424b2'
), (
    '019f0e8b-cdc8-7af3-94f9-34303be09808',
    '019f0866-5d9e-7c53-a3ee-f3b055baa122',
    '019ef097-4931-7c91-ac71-8ba21535fe2d'
), (
    '019f0e8b-cdcc-7eb0-913b-b3ab400deb72',
    '019f0866-5d9e-7c53-a3ee-f3b055baa122',
    '019ef097-4931-7c91-ac71-8bef02eba5db'
), (
    '019f0e8b-cdd2-77c3-acf4-83b6702086bb',
    '019f0866-5d9e-7c53-a3ee-f3cca52e0bc1',
    '019ef097-4931-7c91-ac71-8bf7418601e8'
), (
    '019f0e8b-cdd8-7de2-9039-e665c00702d5',
    '019f0866-5d9e-7c53-a3ee-f3cca52e0bc1',
    '019ef097-4931-7c91-ac71-8c07c3dba6d6'
), (
    '019f0e8b-cddd-7d43-958e-2c5fb67b8c5d',
    '019f0866-5d9e-7c53-a3ee-f3d9eef4cae7',
    '019ef097-4931-7c91-ac71-8c1d049c254f'
), (
    '019f0e8b-cde2-70f0-b31f-be35c9aefeb3',
    '019f0866-5d9e-7c53-a3ee-f3d9eef4cae7',
    '019ef097-4931-7c91-ac71-8b8862b3e9d9'
) ON DUPLICATE KEY UPDATE
customer_id = VALUES (customer_id),
product_id = VALUES (product_id);
