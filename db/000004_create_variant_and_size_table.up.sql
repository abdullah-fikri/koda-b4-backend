
CREATE TABLE variant (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);


CREATE TABLE size (
    id SERIAL PRIMARY KEY,
    name varchar(50),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);


CREATE TABLE product_variant (
    id SERIAL PRIMARY KEY,
    variant_id BIGINT,
    product_id BIGINT,
    FOREIGN KEY (variant_id) REFERENCES variant(id),
    FOREIGN KEY (product_id) REFERENCES products(id),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);


CREATE TABLE product_size (
    id SERIAL PRIMARY KEY,
    product_id BIGINT,
    size_id BIGINT,
    price NUMERIC(10,2) NOT NULL,        
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (size_id) REFERENCES size(id),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);


