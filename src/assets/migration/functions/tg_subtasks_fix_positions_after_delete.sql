CREATE OR REPLACE FUNCTION tg_subtasks_fix_positions_after_delete()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE subtasks
    SET position = position - 1
    WHERE task_id = OLD.task_id
    AND position > OLD.position;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;
