CREATE OR REPLACE FUNCTION columns_move_to_position(
    p_project_id integer, 
    p_column_id integer, 
    p_new_position integer
) 
RETURNS VOID AS $$
DECLARE
    v_current_position INT;
    v_max_position INT;
BEGIN
    -- Получаем текущее состояние задачи
    SELECT position INTO v_current_position
    FROM columns
    WHERE id = p_column_id AND project_id = p_project_id;

    IF NOT FOUND THEN
        RAISE EXCEPTION SQLSTATE 'P0002';
    END IF;

    -- Если позиция не изменилась - ничего не делаем
    IF v_current_position = p_new_position THEN
        RETURN;
    END IF;

    -- Получаем максимальную позицию для соответствующего архивированного статуса
    SELECT COALESCE(MAX(position), 0) INTO v_max_position
    FROM columns
    WHERE project_id = p_project_id;

    -- Корректируем позицию если она выходит за пределы
    p_new_position := GREATEST(1, LEAST(p_new_position, v_max_position + 1));

    -- Обновляем позиции только среди задач с тем же архивным статусом
    IF p_new_position < v_current_position THEN
        -- Перемещение вверх (позиция уменьшается)
        UPDATE columns
        SET position = position + 1
        WHERE project_id = p_project_id
          AND position >= p_new_position
          AND position < v_current_position;
    ELSE
        -- Перемещение вниз (позиция увеличивается)
        UPDATE columns
        SET position = position - 1
        WHERE project_id = p_project_id
          AND position > v_current_position
          AND position <= p_new_position;
    END IF;

    -- Устанавливаем новую позицию для задачи
    UPDATE columns
    SET position = p_new_position
    WHERE id = p_column_id;
END;
$$ LANGUAGE plpgsql;
