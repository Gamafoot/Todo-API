CREATE OR REPLACE FUNCTION public.is_owned_column(p_user_id integer, p_column_id integer)
RETURNS boolean
AS $function$
BEGIN
	RETURN EXISTS (
        SELECT 1 FROM tasks INNER JOIN columns ON columns.id = tasks.column_id 
        INNER JOIN projects ON projects.id = columns.id 
        WHERE projects.user_id = p_user_id 
        LIMIT 1
    );
END;
$function$ LANGUAGE plpgsql
