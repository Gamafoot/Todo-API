CREATE OR REPLACE FUNCTION tg_tasks_update_completed_at()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.status = TRUE AND (OLD.status IS DISTINCT FROM NEW.status) THEN
        NEW.completed_at = NOW();
    ELSIF NEW.status = FALSE AND (OLD.status IS DISTINCT FROM NEW.status) THEN
        NEW.completed_at = NULL;
    END IF;
    
    NEW.updated_at = NOW();
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
