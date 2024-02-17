CREATE OR REPLACE FUNCTION generate_barcode()
RETURNS TRIGGER AS $$
BEGIN
    NEW.barcode = LPAD(FLOOR(RANDOM() * 100000000)::TEXT, 8, '0');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER generate_barcode_trigger
BEFORE INSERT ON product
FOR EACH ROW
WHEN (NEW.barcode IS NULL)
EXECUTE FUNCTION generate_barcode();
