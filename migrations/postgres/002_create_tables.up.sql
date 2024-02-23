CREATE TABLE IF NOT EXISTS income (
    id UUID PRIMARY KEY,
    branch_id UUID REFERENCES branch(id),
    price    numeric(100,4),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE TABLE IF NOT EXISTS income_products (
    id UUID PRIMARY KEY,
    income_id UUID REFERENCES income(id),
    product_id UUID REFERENCES product(id),
    price numeric(100,4),
    count int,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);