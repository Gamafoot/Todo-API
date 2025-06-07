CREATE OR REPLACE FUNCTION project_metrics(p_project_id bigint)
    RETURNS TABLE(total_tasks integer, done_tasks integer, rem_tasks integer, days_elapsed integer, days_left integer)
    LANGUAGE plpgsql
AS
$$
DECLARE
  v_start_date date;
  v_deadline date;
BEGIN
  -- Получаем даты начала и дедлайна проекта
  SELECT created_at, deadline INTO v_start_date, v_deadline
  FROM projects
  WHERE id = p_project_id;

  -- Вычисляем метрики задач
  SELECT
    COUNT(*)::int,
    SUM(CASE WHEN status = true THEN 1 ELSE 0 END)::int,
    SUM(CASE WHEN status = false THEN 1 ELSE 0 END)::int,
    (CURRENT_DATE - v_start_date)::int,
    (v_deadline - CURRENT_DATE)::int
  INTO total_tasks, done_tasks, rem_tasks, days_elapsed, days_left
  FROM tasks t
  JOIN columns c ON t.column_id = c.id
  WHERE c.project_id = p_project_id;

  RETURN NEXT;
  RETURN;
END;
$$;

