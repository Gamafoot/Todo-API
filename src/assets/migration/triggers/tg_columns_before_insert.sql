DROP TRIGGER IF EXISTS tg_columns_before_insert ON columns;

CREATE TRIGGER tg_columns_before_insert
BEFORE INSERT
ON columns
FOR EACH ROW
EXECUTE PROCEDURE public.tg_columns_set_default_position();
