CREATE OR REPLACE FUNCTION monthly_stats(p_user_id bigint, year_num integer, month_num integer)
    RETURNS TABLE(day date, count bigint)
    LANGUAGE plpgsql
AS
$$
BEGIN
    RETURN QUERY
    SELECT completed_at::DATE AS day, 
           COUNT(*) AS count
    FROM tasks
    INNER JOIN columns ON columns.id = tasks.column_id
    INNER JOIN projects ON projects.id = columns.project_id
    WHERE projects.user_id = p_user_id 
      AND EXTRACT(YEAR FROM completed_at) = year_num
      AND EXTRACT(MONTH FROM completed_at) = month_num
    GROUP BY day
    ORDER BY day;
END;
$$;
