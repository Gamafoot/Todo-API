CREATE OR REPLACE FUNCTION tasks_set_default_position(
    p_column_id bigint,
    p_task_id bigint
)
RETURNS VOID AS $$
DECLARE
    v_new_position integer;
BEGIN
    SELECT COUNT(*) + 1 INTO v_new_position FROM tasks 
    WHERE column_id = p_column_id AND archived = false;
    
    UPDATE tasks SET position = v_new_position 
    WHERE id = p_task_id;
END;
$$ LANGUAGE plpgsql;
