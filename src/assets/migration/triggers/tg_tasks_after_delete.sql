DROP TRIGGER IF EXISTS tg_tasks_after_delete ON tasks;

CREATE TRIGGER tg_tasks_after_delete
AFTER DELETE
ON tasks
FOR EACH ROW
EXECUTE PROCEDURE tg_tasks_fix_positions_after_delete();
