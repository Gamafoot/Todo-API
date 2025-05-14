CREATE OR REPLACE FUNCTION tg_tasks_set_default_position()
RETURNS TRIGGER AS $$
BEGIN
    SELECT COUNT(*) + 1 INTO NEW.position FROM tasks 
    WHERE column_id = NEW.column_id AND archived = false;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
