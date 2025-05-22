DROP TRIGGER IF EXISTS tg_columns_after_insert ON columns;

CREATE TRIGGER tg_columns_after_insert
AFTER INSERT
ON columns
FOR EACH ROW
EXECUTE PROCEDURE tg_columns_init_position();
