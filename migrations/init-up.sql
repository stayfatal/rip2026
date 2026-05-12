CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    login VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(100) NOT NULL,
    is_moderator BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS sharding_strategies (
    strategy_id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    short_description_en VARCHAR(100) NOT NULL DEFAULT '',
    description VARCHAR(1000) NOT NULL,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE,
    photo_url VARCHAR(255),
    video VARCHAR(255),
    latency_coefficient NUMERIC(5,2) NOT NULL DEFAULT 1.0,
    throughput_coefficient NUMERIC(5,2) NOT NULL DEFAULT 1.0,
    reliability_coefficient NUMERIC(5,2) NOT NULL DEFAULT 0.9
);

CREATE TABLE IF NOT EXISTS system_loads (
    system_load_id SERIAL PRIMARY KEY,
    status VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    creator_id INTEGER NOT NULL REFERENCES users(user_id),
    forming_date TIMESTAMP,
    finish_date TIMESTAMP,
    moderator_id INTEGER REFERENCES users(user_id),
    description VARCHAR(2000)
);

CREATE TABLE IF NOT EXISTS system_load_strategies (
    system_load_id INTEGER NOT NULL REFERENCES system_loads(system_load_id),
    strategy_id INTEGER NOT NULL REFERENCES sharding_strategies(strategy_id),
    data_volume INTEGER NOT NULL DEFAULT 0,
    query_count INTEGER NOT NULL DEFAULT 0,
    response_time NUMERIC(12,2),
    PRIMARY KEY (system_load_id, strategy_id)
);

INSERT INTO users (login, password, is_moderator) VALUES
    ('user1', 'pass1', false),
    ('moderator', 'modpass', true)
ON CONFLICT (login) DO NOTHING;

INSERT INTO sharding_strategies (title, short_description_en, description, is_deleted, photo_url, video, latency_coefficient, throughput_coefficient, reliability_coefficient) VALUES
    ('Range Sharding', 'Range shards group neighboring keys and speed up ordered scans for time-series data.', 'Разбивка по диапазону ключа. Для временных рядов, логов и быстрых range-запросов.', false, 'range_sharding.jpg', 'range_sharding.mp4', 1.8, 0.6, 0.70),
    ('Hash Sharding', 'Hash shards spread keys uniformly, balancing write load across partitions in OLTP systems.', 'Распределение по хэшу ключа. Равномерная нагрузка, стандарт для OLTP.', false, 'hash_sharding.jpg', 'hash_sharding.mp4', 1.0, 1.2, 0.90),
    ('Geo Sharding', 'Geo shards place data near users, reducing latency and supporting regional compliance needs.', 'Данные рядом с пользователем по географии. Низкая латентность и локализация (GDPR).', false, 'geo_sharding.jpg', 'geo_sharding.mp4', 1.4, 0.8, 0.85),
    ('Directory-Based Sharding', 'Directory shards map each tenant key to a shard for flexible routing and controlled movement.', 'Каталог «ключ → шард». Максимальная гибкость маршрутизации, мультитенантность.', false, 'directory_sharding.jpg', 'directory_sharding.mp4', 2.2, 0.5, 0.95),
    ('Composite Sharding', 'Composite shards combine rules like geo plus hash to scale globally with balanced traffic.', 'Гибрид: гео + хэш (или иные комбинации). Глобальные системы с равномерной нагрузкой.', false, 'composite_sharding.jpg', 'composite_sharding.mp4', 1.2, 1.0, 0.92),
    ('Dynamic Sharding', 'Dynamic shards split and merge automatically as traffic changes to keep cluster usage stable.', 'Авто split/merge шардов по нагрузке. Облачные БД, эластичное масштабирование.', false, 'dynamic_sharding.jpg', 'dynamic_sharding.mp4', 0.8, 1.5, 0.88);
