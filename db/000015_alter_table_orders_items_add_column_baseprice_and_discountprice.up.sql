ALTER TABLE order_items
ADD COLUMN base_price NUMERIC(10,2),
ADD COLUMN discount_price NUMERIC(10,2),
ADD COLUMN discount_percent NUMERIC(5,2);
