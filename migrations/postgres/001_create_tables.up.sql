CREATE TABLE IF NOT EXISTS branch (
   id UUID PRIMARY KEY,
   name VARCHAR(75) NOT NULL,
   address VARCHAR(75) NOT NULL,
   created_at TIMESTAMP DEFAULT NOW(),
   updated_at TIMESTAMP,
   deleted_at TIMESTAMP
);
CREATE TABLE IF NOT EXISTS category (
    id UUID PRIMARY KEY,
    name VARCHAR(75) NOT NULL,
    parent_id UUID REFERENCES category(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE TABLE IF NOT EXISTS product (
    id UUID PRIMARY KEY,
    name VARCHAR(75) NOT NULL,
    price numeric(75,4) NOT NULL,
    barcode VARCHAR(10) UNIQUE NOT NULL,
    category_id UUID REFERENCES category(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE TABLE IF NOT EXISTS  storage (
    id UUID PRIMARY KEY,
    product_id UUID REFERENCES product(id),
    branch_id UUID REFERENCES branch(id),
    count INT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE TABLE IF NOT EXISTS sale (
    id UUID PRIMARY KEY,
    branch_id UUID REFERENCES branch(id),
    shop_assistent_id VARCHAR(10),
    chashier_id  VARCHAR(10) NOT NULL,
    payment_type VARCHAR(20) CHECK (payment_type IN('card', 'cash')),
    price numeric(75,4) NOT NULL,
    status  VARCHAR(20) CHECK (status IN('in_procces', 'succes', 'cancel')),
    client_name VARCHAR(75) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE TABLE IF NOT EXISTS basket (
    id UUID PRIMARY KEY,
    sale_id UUID REFERENCES sale(id),
    product_id UUID REFERENCES product(id),
    quantity INT NOT NULL,
    price numeric(75,4) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE TABLE IF NOT EXISTS tarif (
    id UUID PRIMARY KEY,
    name VARCHAR(75) NOT NULL,
    tarif_type VARCHAR(20) CHECK (tarif_type IN('percent', 'fixed')),
    amount_for_cash numeric(75,4) NOT NULL,
    amount_for_card numeric(75,4) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE TABLE IF NOT EXISTS staff (
    id UUID PRIMARY KEY,
    branch_id UUID REFERENCES branch(id) NOT NULL,
    tarif_id UUID REFERENCES tarif(id) NOT NULL,
    type_staff VARCHAR(20) CHECK (type_staff IN('shop_assistant', 'chashier')) NOT NULL,
    name VARCHAR(75) NOT NULL,
    balance numeric(75,4) NOT NULL,
    birth_date date NOT NULL,
    age INT,
    gender VARCHAR(10) CHECK (gender IN('male', 'female')),
    logIN VARCHAR(75) NOT NULL,
    password VARCHAR(128) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE TABLE IF NOT EXISTS transaction (
    id UUID PRIMARY KEY,
    sale_id UUID REFERENCES sale(id),
    staff_id UUID REFERENCES staff(id),
    transaction_type VARCHAR(20) CHECK (transaction_type IN ('withdraw', 'topup')),
    source_type VARCHAR(20) CHECK (source_type IN ('bonus', 'sales')),
    amount numeric(75,4) NOT NULL,
    description text NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE TABLE IF NOT EXISTS storage_transaction (
    id UUID PRIMARY KEY,
    staff_id UUID REFERENCES staff(id),
    product_id UUID REFERENCES product(id),
    storage_transaction_type VARCHAR(20) CHECK (storage_transaction_type IN ('minus', 'plus')),
    price numeric(75,4) NOT NULL,
    quantity numeric(75,4) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);


