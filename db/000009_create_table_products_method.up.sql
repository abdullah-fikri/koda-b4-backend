CREATE TABLE product_method (
    id SERIAL PRIMARY KEY,
    product_id BIGINT NOT NULL,
    method_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),

    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (method_id) REFERENCES method(id) ON DELETE RESTRICT
);

