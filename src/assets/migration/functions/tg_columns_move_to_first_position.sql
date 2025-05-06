CREATE OR REPLACE FUNCTION public.tg_columns_move_to_first_position()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM public.columns_move_to_position(NEW.project_id::INT, NEW.id::INT, 1);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
