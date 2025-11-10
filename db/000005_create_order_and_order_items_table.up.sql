
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    users_id BIGINT,
    payment_id BIGINT,
    shipping_id BIGINT,
    order_date TIMESTAMP,
    customer_name VARCHAR(150),
    customer_phone VARCHAR(20),
    customer_address TEXT,
    status VARCHAR(50),
    FOREIGN KEY (users_id) REFERENCES users(id),
    FOREIGN KEY (payment_id) REFERENCES payment(id),
    FOREIGN KEY (shipping_id) REFERENCES shippings(id),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);


CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id BIGINT,
    product_id BIGINT,
    variant_id BIGINT,
    size_id BIGINT,
    qty INT,
    subtotal NUMERIC(10,2),
    status VARCHAR(50),
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (variant_id) REFERENCES variant(id),
    FOREIGN KEY (size_id) REFERENCES size(id),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);