DROP TRIGGER IF EXISTS tg_subtasks_after_insert ON subtasks;

CREATE TRIGGER tg_subtasks_after_insert
AFTER INSERT
ON subtasks
FOR EACH ROW
EXECUTE PROCEDURE public.tg_subtasks_move_to_first_position();
