CREATE TABLE branch (
   id UUID PRIMARY KEY,
   name VARCHAR(75) NOT NULL,
   address VARCHAR(75) NOT NULL,
   created_at TIMESTAMP DEFAULT NOW(),
   updated_at TIMESTAMP,
   deleted_at TIMESTAMP
);
CREATE TABLE category (
    id UUID PRIMARY KEY,
    name VARCHAR(75) NOT NULL,
    parent_id uuid references category(id),
    created_at timestamp DEFAULT now(),
    updated_at timestamp,
    deleted_at timestamp
);
CREATE TABLE product (
    id uuid PRIMARY key,
    name VARCHAR(75) not null,
    price numeric(75,4) not null,
    barcode varchar(10) unique not null,
    category_id uuid references category(id),
    created_at timestamp DEFAULT now(),
    updated_at timestamp,
    deleted_at timestamp
);
create table  storage (
    id uuid PRIMARY key,
    product_id uuid references product(id),
    branch_id uuid references branch(id),
    count numeric(75,4) not null,
    created_at timestamp DEFAULT now(),
    updated_at timestamp,
    deleted_at timestamp
);
create table sale (
    id uuid primary key,
    branch_id uuid references branch(id),
    shop_assistent_id varchar(10),
    chashier_id  serial not null,
    payment_type varchar(20) check (payment_type in('card', 'cash')),
    price numeric(75,4) not null,
    status  varchar(20) check (status in('in_procces', 'succes', 'cancel')),
    client_name varchar(75) not null,
    created_at timestamp DEFAULT now(),
    updated_at timestamp,
    deleted_at timestamp
);
create table basket (
    id uuid primary key,
    sale_id uuid references sale(id),
    product_id uuid references product(id),
    quantity numeric(75,4) not null,
    price numeric(75,4) not null,
    created_at timestamp DEFAULT now(),
    updated_at timestamp,
    deleted_at timestamp
);
create table tarif (
    id uuid primary key,
    name varchar(75) not null,
    tarif_type varchar(20) check (tarif_type in('percent', 'field')),
    amount_for_cash numeric(75,4) not null,
    amount_for_card numeric(75,4) not null,
    created_at timestamp DEFAULT now(),
    updated_at timestamp,
    deleted_at timestamp
);
create table staff (
   id uuid primary key,
   branch_id uuid references branch(id) not null,
   tarif_id uuid references tarif(id) not null,
   type_staff varchar(20) check (type_staff in('shop_assistant', 'chashier')) not NULL,
   name varchar(75) not null,
   balance numeric(75,4) not null,
   birth_date date not null,
   age int,
   gender varchar(10) check (gender in('male', 'female')),
   login varchar(75) not null,
   password varchar(128) not null,
   created_at timestamp DEFAULT now(),
   updated_at timestamp,
   deleted_at timestamp
);
create table transaction (
  id uuid primary key,
  sale_id uuid references sale(id),
  staff_id uuid references staff(id),
  transaction_type varchar(20) check (transaction_type in ('withdraw', 'topup')),
  source_type varchar(20) check (source_type in ('bonus', 'sales')),
  amount numeric(75,4) not null,
  description text not null,
  created_at timestamp DEFAULT now(),
  updated_at timestamp,
  deleted_at timestamp
);
create table storage_transaction (
    id uuid primary key,
    staff_id uuid references staff(id),
    product_id uuid references product(id),
    storage_transaction_type varchar(20) check (storage_transaction_type in ('minus', 'plus')),
    price numeric(75,4) not null,
    quantity numeric(75,4) not null,
    created_at timestamp DEFAULT now(),
    updated_at timestamp,
    deleted_at timestamp
);


