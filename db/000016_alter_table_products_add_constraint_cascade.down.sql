ALTER TABLE product_img
DROP CONSTRAINT IF EXISTS product_img_product_id_fkey,
ADD CONSTRAINT product_img_product_id_fkey
FOREIGN KEY (product_id) REFERENCES products(id);

ALTER TABLE product_variant
DROP CONSTRAINT IF EXISTS product_variant_product_id_fkey,
ADD CONSTRAINT product_variant_product_id_fkey
FOREIGN KEY (product_id) REFERENCES products(id);

ALTER TABLE product_size
DROP CONSTRAINT IF EXISTS product_size_product_id_fkey,
ADD CONSTRAINT product_size_product_id_fkey
FOREIGN KEY (product_id) REFERENCES products(id);


ALTER TABLE order_items
DROP CONSTRAINT IF EXISTS order_items_product_id_fkey,
ADD CONSTRAINT order_items_product_id_fkey
FOREIGN KEY (product_id) REFERENCES products(id);
