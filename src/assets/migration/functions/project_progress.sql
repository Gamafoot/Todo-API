CREATE FUNCTION project_progress(p_project_id bigint)
    RETURNS TABLE(day date, count bigint)
    LANGUAGE plpgsql
AS
$$
BEGIN
    RETURN QUERY
    SELECT 
        DATE(tasks.completed_at) AS day,
        COUNT(*)::bigint AS count
    FROM tasks
    JOIN columns ON tasks.column_id = columns.id
    JOIN projects ON columns.project_id = projects.id
    WHERE 
        projects.id = p_project_id
        AND tasks.completed_at IS NOT NULL
        AND tasks.status = true
    GROUP BY DATE(tasks.completed_at)
    ORDER BY DATE(tasks.completed_at);
END;
$$;
