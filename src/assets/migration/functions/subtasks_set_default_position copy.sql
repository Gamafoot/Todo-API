CREATE OR REPLACE FUNCTION subtasks_set_default_position(
    p_task_id bigint,
    p_subtask_id bigint
)
RETURNS VOID AS $$
DECLARE
    v_new_position integer;
BEGIN
    SELECT COUNT(*) INTO v_new_position FROM subtasks 
    WHERE task_id = p_task_id;
    
    UPDATE subtasks SET position = v_new_position 
    WHERE id = p_subtask_id;
END;
$$ LANGUAGE plpgsql;
