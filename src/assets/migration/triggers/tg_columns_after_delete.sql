DROP TRIGGER IF EXISTS tg_columns_after_delete ON columns;

CREATE TRIGGER tg_columns_after_delete
AFTER DELETE
ON columns
FOR EACH ROW
EXECUTE PROCEDURE public.tg_columns_fix_positions_after_delete();
