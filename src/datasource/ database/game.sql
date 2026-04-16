--таблица игроков

-- Создание расширения для UUID
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
--таблица для игр
CREATE TABLE IF NOT EXISTS games (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    field JSONB NOT NULL, -- храним поле в JSON формате
    --status TEXT DEFAULT 'in_progress'   --?
    --next_turn TEXT DEFAULT 'O', --?
    --score_x INTEGER, --?
    --score_y INTEGER, --?
    --finished BOOLEAN, --?
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индекс для быстрого поиска активных игр
--CREATE INDEX idx_games_status ON games(status);