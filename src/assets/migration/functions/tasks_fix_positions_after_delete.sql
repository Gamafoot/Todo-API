CREATE OR REPLACE FUNCTION public.tasks_fix_positions_after_delete(
    p_old_column_id integer,
    p_old_position integer
)
RETURNS VOID AS $$
BEGIN
    UPDATE tasks
    SET position = position - 1
    WHERE column_id = p_old_column_id
    AND position > p_old_position;
END;
$$ LANGUAGE plpgsql;
