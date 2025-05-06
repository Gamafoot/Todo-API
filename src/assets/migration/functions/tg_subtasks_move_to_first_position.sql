CREATE OR REPLACE FUNCTION public.tg_subtasks_move_to_first_position()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM public.subtasks_move_to_position(NEW.column_id::INT, NEW.id::INT, 1);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
