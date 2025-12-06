-- DOWN: Kembalikan ke default (tanpa ON DELETE CASCADE)
ALTER TABLE cart_items
DROP CONSTRAINT cart_items_product_id_fkey;

ALTER TABLE cart_items
ADD CONSTRAINT cart_items_product_id_fkey
FOREIGN KEY (product_id)
REFERENCES products(id);
