CREATE OR REPLACE FUNCTION columns_set_default_position(
    p_project_id bigint,
    p_column_id bigint
)
RETURNS VOID AS $$
DECLARE
    v_new_position integer;
BEGIN
    SELECT COUNT(*) INTO v_new_position FROM columns 
    WHERE project_id = p_project_id;
    
    UPDATE columns SET position = v_new_position 
    WHERE id = p_column_id;
END;
$$ LANGUAGE plpgsql;
