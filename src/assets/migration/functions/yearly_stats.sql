CREATE FUNCTION yearly_stats(p_user_id bigint, year_num integer)
    RETURNS TABLE(month integer, count bigint)
    LANGUAGE plpgsql
AS
$$
BEGIN
    RETURN QUERY
    SELECT EXTRACT(MONTH FROM completed_at)::INT AS month, 
           COUNT(*) AS count
    FROM tasks
    INNER JOIN columns ON columns.id = tasks.column_id
    INNER JOIN projects ON projects.id = columns.project_id
    WHERE projects.user_id = p_user_id 
      AND EXTRACT(YEAR FROM completed_at) = year_num
    GROUP BY month
    ORDER BY month;
END;
$$;
