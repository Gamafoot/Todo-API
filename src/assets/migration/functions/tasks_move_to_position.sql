CREATE OR REPLACE FUNCTION tasks_move_to_position(
    p_column_id bigint, 
    p_task_id bigint, 
    p_new_position integer
)
RETURNS VOID AS $$
DECLARE
    v_current_position INT;
    v_max_position INT;
BEGIN
    -- Получаем текущее состояние задачи
    SELECT position INTO v_current_position
    FROM tasks 
    WHERE id = p_task_id AND column_id = p_column_id;
    
    IF NOT FOUND THEN
        RAISE EXCEPTION SQLSTATE 'P0002';
    END IF;
    
    -- Если позиция не изменилась - ничего не делаем
    IF v_current_position = p_new_position THEN
        RETURN;
    END IF;
    
    SELECT COALESCE(MAX(position), 1) INTO v_max_position
    FROM tasks 
    WHERE column_id = p_column_id;
    
    -- Корректируем позицию если она выходит за пределы
    p_new_position := LEAST(p_new_position, v_max_position);
    
    -- Обновляем позиции только среди задач с тем же архивным статусом
    IF p_new_position < v_current_position THEN
        -- Перемещение вверх (позиция уменьшается)
        UPDATE tasks
        SET position = position + 1
        WHERE column_id = p_column_id
          AND position >= p_new_position
          AND position < v_current_position;
    ELSE
        -- Перемещение вниз (позиция увеличивается)
        UPDATE tasks
        SET position = position - 1
        WHERE column_id = p_column_id
          AND position > v_current_position
          AND position <= p_new_position;
    END IF;
    
    -- Устанавливаем новую позицию для задачи
    UPDATE tasks
    SET position = p_new_position
    WHERE id = p_task_id;
END;
$$ LANGUAGE plpgsql;
