CREATE OR REPLACE FUNCTION public.is_owned_subtask(p_user_id integer, p_subtask_id integer)
RETURNS boolean
AS $function$
BEGIN
	RETURN EXISTS (
        SELECT 1 FROM subtasks 
        INNER JOIN tasks ON tasks.id = subtasks.task_id 
        INNER JOIN columns ON columns.id = tasks.column_id 
        INNER JOIN projects ON projects.id = columns.project_id 
        WHERE projects.user_id = p_user_id AND subtasks.id = p_subtask_id LIMIT 1
    );
END;
$function$ LANGUAGE plpgsql
