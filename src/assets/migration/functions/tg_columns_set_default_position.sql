CREATE OR REPLACE FUNCTION public.tg_columns_set_default_position() 
RETURNS TRIGGER AS $$
BEGIN
    SELECT COUNT(*) + 1 INTO NEW.position FROM columns 
    WHERE project_id = NEW.project_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
