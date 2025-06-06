DROP TRIGGER IF EXISTS tg_subtasks_after_insert ON subtasks;

CREATE TRIGGER tg_subtasks_after_insert
AFTER INSERT
ON subtasks
FOR EACH ROW
EXECUTE PROCEDURE tg_subtasks_init_position();
