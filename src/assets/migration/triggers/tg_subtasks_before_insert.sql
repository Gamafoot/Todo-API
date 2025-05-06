DROP TRIGGER IF EXISTS tg_subtasks_before_insert ON subtasks;

CREATE TRIGGER tg_subtasks_before_insert
BEFORE INSERT
ON subtasks
FOR EACH ROW
EXECUTE PROCEDURE public.tg_subtasks_set_default_position();
