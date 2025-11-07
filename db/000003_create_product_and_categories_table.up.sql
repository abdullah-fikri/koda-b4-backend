
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name varchar(100), 
    description TEXT,
    price NUMERIC(10,2),
    stock INT,
    category_id BIGINT,
    FOREIGN KEY (category_id) REFERENCES categories(id),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);


CREATE TABLE product_img (
    id SERIAL PRIMARY KEY,
    image TEXT,
    product_id BIGINT,
    FOREIGN KEY(product_id) REFERENCES products(id),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
