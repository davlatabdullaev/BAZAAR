CREATE OR REPLACE FUNCTION generate_barcode()
RETURNS TRIGGER AS $$
BEGIN
    NEW.barcode = LPAD(FLOOR(RANDOM() * 10000000000)::VARCHAR, 10, '0');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_insert_product
BEFORE INSERT ON product
FOR EACH ROW
WHEN (NEW.barcode IS NULL)
EXECUTE FUNCTION generate_barcode();
