CREATE TABLE IF NOT EXISTS income (
    id UUID PRIMARY KEY,
    branch_id UUID REFERENCES branch(id),
    price    numeric(100,4)
);
CREATE TABLE IF NOT EXISTS income_product (
    id UUID PRIMARY KEY,
    income_id UUID REFERENCES income(id),
    product_id UUID REFERENCES product(id),
    price numeric(100,4),
    count int
);