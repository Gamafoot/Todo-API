CREATE OR REPLACE FUNCTION tg_tasks_fix_positions_after_delete()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE tasks
    SET position = position - 1
    WHERE column_id = OLD.column_id
    AND position > OLD.position;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;
