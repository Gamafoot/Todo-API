CREATE OR REPLACE FUNCTION public.tg_tasks_set_default_position()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM public.tasks_set_default_position(NEW.column_id::integer, NEW.id::integer);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
