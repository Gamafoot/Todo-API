DROP TRIGGER IF EXISTS tg_subtasks_after_delete ON subtasks;

CREATE TRIGGER tg_subtasks_after_delete
AFTER DELETE
ON subtasks
FOR EACH ROW
EXECUTE PROCEDURE public.tg_subtasks_fix_positions_after_delete();
