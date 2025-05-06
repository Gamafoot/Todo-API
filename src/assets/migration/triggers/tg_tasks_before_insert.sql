DROP TRIGGER IF EXISTS tg_tasks_before_insert ON tasks;

CREATE TRIGGER tg_tasks_before_insert
BEFORE INSERT
ON tasks
FOR EACH ROW
EXECUTE PROCEDURE public.tg_tasks_set_default_position();
