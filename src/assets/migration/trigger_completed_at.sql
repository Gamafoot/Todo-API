-- Create a trigger for a completed_at field in tasks
CREATE OR REPLACE FUNCTION public.update_completed_at()
RETURNS TRIGGER AS $function$
BEGIN
    IF NEW.status = TRUE AND (OLD.status IS DISTINCT FROM NEW.status) THEN
        NEW.completed_at = NOW();
    ELSIF NEW.status = FALSE AND (OLD.status IS DISTINCT FROM NEW.status) THEN
        NEW.completed_at = NULL;
    END IF;
    
    NEW.updated_at = NOW();
    
    RETURN NEW;
END;
$function$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trigger_update_completed_at ON tasks;

CREATE TRIGGER trigger_update_completed_at
BEFORE UPDATE ON tasks
FOR EACH ROW
EXECUTE FUNCTION update_completed_at();
