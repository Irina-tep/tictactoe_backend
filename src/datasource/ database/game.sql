-- --таблица игроков

-- -- Создание расширения для UUID
-- CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
-- --таблица для игр
-- CREATE TABLE IF NOT EXISTS games (
--     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--     field JSONB NOT NULL, -- храним поле в JSON формате
--     game_state TEXT DEFAULT 'waiting'   --?
--     players JSONB NOT NULL, -- храним поле в JSON формате
   
--     --next_turn TEXT DEFAULT 'O', --?
--     --score_x INTEGER, --?
--     --score_y INTEGER, --?
--     --finished BOOLEAN, --?
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     game_type TEXT DEFAULT 'pvp'
-- );

-- -- Индекс для быстрого поиска активных игр
-- --CREATE INDEX idx_games_status ON games(status);

-- --таблица игроков
-- CREATE TABLE IF NOT EXISTS player (
--     id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
--     login TEXT UNIQUE NOT NULL,
--     password TEXT UNIQUE NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );

-- -- Индекс для быстрого поиска по логину
-- -- CREATE INDEX IF NOT EXISTS idx_player_login ON player(login);
-- -- DELETE FROM games;
-- -- DELETE FROM player;

Select * FROM games;