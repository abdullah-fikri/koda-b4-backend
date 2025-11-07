
CREATE TABLE discount (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    start_discount TIMESTAMP,
    end_discount TIMESTAMP,
     is_active BOOLEAN DEFAULT TRUE,
    percent_discount FLOAT,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);


CREATE TABLE promo_card (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),  
    start TIMESTAMP,
    end_promo TIMESTAMP,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);


CREATE TABLE product_discount (
    id SERIAL PRIMARY KEY,
    product_id BIGINT,
    discount_id BIGINT,
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (discount_id) REFERENCES discount(id),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);


CREATE TABLE product_promo_card (
    id SERIAL PRIMARY KEY,
    promo_card_id BIGINT,
    product_id BIGINT,
    FOREIGN KEY (promo_card_id) REFERENCES promo_card(id),
    FOREIGN KEY (product_id) REFERENCES products(id),
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
);
