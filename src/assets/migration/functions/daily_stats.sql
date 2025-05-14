CREATE OR REPLACE FUNCTION daily_stats(p_current_date date)
RETURNS TABLE(
    hour integer, 
    count integer
)
AS $$
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
    WHERE status = TRUE
        AND completed_at BETWEEN v_start AND v_end
    GROUP BY hour
    ORDER BY hour;
END;
$$ LANGUAGE plpgsql;

