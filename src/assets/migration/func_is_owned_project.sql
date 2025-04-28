CREATE OR REPLACE FUNCTION public.is_owned_project(p_user_id integer, p_project_id integer)
RETURNS boolean
AS $function$
BEGIN
	RETURN EXISTS (
        SELECT 1 FROM projects 
        WHERE id = p_project_id AND user_id = p_user_id 
        LIMIT 1
    );
END;
$function$ LANGUAGE plpgsql
