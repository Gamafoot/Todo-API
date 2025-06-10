CREATE OR REPLACE FUNCTION heatmap_add_activity(p_user_id bigint) 
RETURNS void AS $$
DECLARE
    v_current_date date := (CURRENT_TIMESTAMP AT TIME ZONE 'UTC')::date;
BEGIN
    WITH existing AS (
        SELECT count FROM heatmaps 
        WHERE user_id = p_user_id AND date = v_current_date
        FOR UPDATE
    )
    INSERT INTO heatmaps (user_id, date, count)
    SELECT p_user_id, v_current_date, 
        COALESCE((SELECT count FROM existing), 0) + 1
    WHERE NOT EXISTS (SELECT 1 FROM existing);
    
    IF NOT FOUND THEN
        UPDATE heatmaps 
        SET count = count + 1
        WHERE user_id = p_user_id AND date = v_current_date;
    END IF;
END;
$$ LANGUAGE plpgsql;
