DROP TRIGGER IF EXISTS tg_tasks_before_update ON tasks;

CREATE TRIGGER tg_tasks_before_update
BEFORE UPDATE ON tasks
FOR EACH ROW
EXECUTE FUNCTION public.tg_tasks_update_completed_at();
