DROP TRIGGER IF EXISTS tg_tasks_after_insert ON tasks;

CREATE TRIGGER tg_tasks_after_insert
AFTER INSERT
ON tasks
FOR EACH ROW
EXECUTE PROCEDURE tg_tasks_move_to_first_position();
