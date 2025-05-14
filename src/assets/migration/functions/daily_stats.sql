CREATE FUNCTION daily_stats(p_user_id bigint, p_current_date date)
    RETURNS TABLE(hour integer, count integer)
    LANGUAGE plpgsql
AS
$$
DECLARE
    v_start timestamptz;
    v_end timestamptz;
BEGIN
    v_start := p_current_date::timestamp;
    v_end := (p_current_date + INTERVAL '1 day')::timestamp - INTERVAL '1 second';

    RETURN QUERY
    SELECT
        EXTRACT(HOUR FROM completed_at)::integer AS hour,
        COUNT(*)::integer AS count
    FROM tasks
    INNER JOIN columns ON columns.id = tasks.column_id
    INNER JOIN projects ON projects.id = columns.project_id
    WHERE projects.user_id = p_user_id
        AND completed_at BETWEEN v_start AND v_end
    GROUP BY hour
    ORDER BY hour;
END;
$$;
