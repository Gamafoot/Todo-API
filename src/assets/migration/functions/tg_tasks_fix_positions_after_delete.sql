CREATE OR REPLACE FUNCTION tg_tasks_fix_positions_after_delete()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM tasks_fix_positions_after_delete(OLD.column_id, OLD.position);
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;
