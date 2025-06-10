CREATE OR REPLACE FUNCTION heatmap_remove_activity(p_user_id bigint)
RETURNS VOID AS $$
BEGIN
    UPDATE heatmaps SET count = heatmaps.count - 1 WHERE user_id = p_user_id;
END;
$$ LANGUAGE plpgsql;
