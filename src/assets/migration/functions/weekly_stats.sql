CREATE FUNCTION weekly_stats(p_user_id bigint, week_start_date date)
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
    WHERE projects.user_id = p_user_id AND completed_at >= week_start_date
      AND completed_at < (week_start_date + INTERVAL '7 days')
    GROUP BY day
    ORDER BY day;
END;
$$;
