CREATE OR REPLACE FUNCTION tg_columns_init_position()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM columns_set_default_position(NEW.project_id, NEW.id);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
