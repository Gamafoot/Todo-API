CREATE OR REPLACE FUNCTION public.tg_columns_fix_positions_after_delete()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE columns
    SET position = position - 1
    WHERE project_id = OLD.project_id
    AND position > OLD.position;
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;
