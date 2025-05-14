CREATE OR REPLACE FUNCTION is_owned_task(
    p_user_id integer, 
    p_task_id integer
)
RETURNS BOOLEAN AS $$
BEGIN
	RETURN EXISTS (
        SELECT * FROM tasks
        INNER JOIN columns ON columns.id = tasks.column_id 
        INNER JOIN projects ON projects.id = columns.project_id 
        WHERE projects.user_id = p_user_id AND tasks.id = p_task_id LIMIT 1
    );
END;
$$ LANGUAGE plpgsql;
