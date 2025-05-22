CREATE OR REPLACE FUNCTION tasks_move_to_column(
    p_new_column_id bigint,
    p_task_id bigint,
    p_new_position integer
)
RETURNS VOID AS $$
DECLARE
    v_old_column_id bigint;
    v_old_position integer;
BEGIN
    SELECT column_id, position INTO v_old_column_id, v_old_position FROM tasks
    WHERE id = p_task_id;

    IF v_old_column_id = p_new_column_id THEN
        RETURN;
    END IF;

    UPDATE tasks SET column_id = p_new_column_id
    WHERE id = p_task_id;
    
    PERFORM tasks_set_default_position(p_new_column_id, p_task_id);
    PERFORM tasks_move_to_position(p_new_column_id, p_task_id, p_new_position);
    PERFORM tasks_fix_positions_after_delete(v_old_column_id, v_old_position);
END;
$$ LANGUAGE plpgsql;
