CREATE OR REPLACE FUNCTION public.tg_subtasks_set_default_position() 
RETURNS TRIGGER AS $$
BEGIN
    SELECT COUNT(*) + 1 INTO NEW.position FROM subtasks 
    WHERE task_id = NEW.task_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
