CREATE OR REPLACE FUNCTION tg_tasks_init_position()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM tasks_set_default_position(NEW.column_id, NEW.id);
    PERFORM tasks_move_to_position(NEW.column_id, NEW.id, 1);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
