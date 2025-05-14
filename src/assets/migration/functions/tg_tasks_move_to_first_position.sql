CREATE OR REPLACE FUNCTION tg_tasks_move_to_first_position()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM tasks_move_to_position(NEW.column_id::INT, NEW.id::INT, 1);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
