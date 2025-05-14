CREATE OR REPLACE FUNCTION tg_subtasks_move_to_first_position()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM subtasks_move_to_position(NEW.column_id::INT, NEW.id::INT, 1);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
