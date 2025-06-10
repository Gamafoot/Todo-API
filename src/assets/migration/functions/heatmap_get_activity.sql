CREATE OR REPLACE FUNCTION heatmap_get_activity(
    p_user_id INTEGER,
    p_year INTEGER DEFAULT EXTRACT(YEAR FROM CURRENT_DATE)
)
RETURNS TABLE (
    day DATE,
    count INTEGER
) AS $$
BEGIN
    RETURN QUERY
    WITH date_series AS (
        SELECT generate_series(
            make_date(p_year, 1, 1),
            make_date(p_year, 12, 31),
            interval '1 day'
        )::date AS day
    ),
    activity_data AS (
        SELECT 
            h.date AS activity_date,
            COALESCE(SUM(h.count), 0) AS activity
        FROM heatmaps h
        WHERE h.user_id = p_user_id
        AND EXTRACT(YEAR FROM h.date) = p_year
        GROUP BY h.date
    )
    SELECT 
        ds.day,
        COALESCE(ad.activity, 0)::INTEGER
    FROM date_series ds
    LEFT JOIN activity_data ad ON ds.day = ad.activity_date
    ORDER BY ds.day;
END;
$$ LANGUAGE plpgsql;
