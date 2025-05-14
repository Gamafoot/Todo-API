CREATE OR REPLACE FUNCTION is_owned_column(
    p_user_id integer,
    p_column_id integer
)
RETURNS BOOLEAN AS $$
BEGIN
	RETURN EXISTS (
        SELECT * FROM columns
        INNER JOIN projects ON projects.id = columns.project_id
        WHERE projects.user_id = p_user_id AND columns.id = p_column_id
        LIMIT 1
    );
END;
$$ LANGUAGE plpgsql;
