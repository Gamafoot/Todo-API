CREATE FUNCTION project_progress(p_project_id bigint)
    RETURNS TABLE(day date, count bigint)
    LANGUAGE plpgsql
AS
$$
DECLARE
    project_deadline date;
BEGIN
    -- Получаем дедлайн проекта или текущую дату, если дедлайна нет
    SELECT COALESCE(deadline, NOW() AT TIME ZONE 'UTC')::date 
    INTO project_deadline
    FROM projects 
    WHERE id = p_project_id;

    RETURN QUERY
    WITH project_dates AS (
        SELECT 
            generate_series(
                (SELECT created_at::date FROM projects WHERE id = p_project_id),
                project_deadline,
                '1 day'::interval
            )::date AS day_date
    ),
    completed_tasks AS (
        SELECT 
            t.updated_at::date AS completion_date
        FROM tasks t
        JOIN columns c ON t.column_id = c.id
        WHERE c.project_id = p_project_id
        AND t.status = true
    )
    SELECT 
        pd.day_date,
        COUNT(ct.completion_date)::bigint AS completed_tasks
    FROM project_dates pd
    LEFT JOIN completed_tasks ct ON ct.completion_date <= pd.day_date
    GROUP BY pd.day_date
    ORDER BY pd.day_date;
END;
$$;
