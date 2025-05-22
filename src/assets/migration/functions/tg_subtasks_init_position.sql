CREATE OR REPLACE FUNCTION tg_subtasks_init_position()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM subtasks_set_default_position(NEW.task_id, NEW.id);
    PERFORM subtasks_move_to_position(NEW.task_id, NEW.id, 1);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
