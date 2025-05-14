CREATE OR REPLACE FUNCTION tg_tasks_fix_positions_after_delete()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM public.tasks_fix_positions_after_delete(OLD.column_id, OLD.position);
END;
$$ LANGUAGE plpgsql;
