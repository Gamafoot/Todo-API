CREATE OR REPLACE FUNCTION public.tg_tasks_set_default_position()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM public.tasks_set_default_position(NEW.column_id, NEW.id);
END;
$$ LANGUAGE plpgsql;
