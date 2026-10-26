BEGIN;

CREATE SCHEMA algotrade;
SET LOCAL search_path TO algotrade;
-- 1. Пользователи

CREATE TABLE users (
    user_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_name VARCHAR(50) NOT NULL
        CHECK (length(trim(user_name)) > 0),
    email VARCHAR(254) NOT NULL
        CHECK (length(trim(email)) > 0 AND email = trim(email)),
    password TEXT NOT NULL
        CHECK (length(trim(password)) > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX users_email_unique
ON users (lower(email));


-- 2. Стратегии пользователей

CREATE TABLE strategies (
    strategy_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(user_id),
    name VARCHAR(100) NOT NULL
        CHECK (length(trim(name)) > 0),
    description TEXT,
    code TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX strategies_user_idx
ON strategies (user_id);

-- Автоматическое обновление даты изменения стратегии

CREATE FUNCTION update_strategy_timestamp()
RETURNS TRIGGER
LANGUAGE plpgsql
AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$;

CREATE TRIGGER strategies_updated_at_trigger
BEFORE UPDATE ON strategies
FOR EACH ROW
EXECUTE FUNCTION update_strategy_timestamp();


-- 3. Биржи

CREATE TABLE exchanges (
    exchange_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    code VARCHAR(30) NOT NULL UNIQUE
        CHECK (length(trim(code)) > 0),
    name VARCHAR(100) NOT NULL
        CHECK (length(trim(name)) > 0)
);


-- 4. Торговые инструменты
-- quote_currency — валюта, в которой выражена цена

CREATE TABLE instruments (
    instrument_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    exchange_id BIGINT NOT NULL REFERENCES exchanges(exchange_id),
    symbol VARCHAR(50) NOT NULL
        CHECK (length(trim(symbol)) > 0),
    name VARCHAR(150) NOT NULL
        CHECK (length(trim(name)) > 0),
    quote_currency VARCHAR(20) NOT NULL
        CHECK (length(trim(quote_currency)) > 0),

    UNIQUE (exchange_id, symbol)
);


-- 5. Исторические котировки в виде свечей
-- interval_seconds — длительность свечи в секундах

CREATE TABLE candles (
    instrument_id BIGINT NOT NULL REFERENCES instruments(instrument_id),
    interval_seconds INTEGER NOT NULL
        CHECK (interval_seconds > 0),
    open_time TIMESTAMPTZ NOT NULL,

    open_price NUMERIC(28, 10) NOT NULL,
    high_price NUMERIC(28, 10) NOT NULL,
    low_price NUMERIC(28, 10) NOT NULL,
    close_price NUMERIC(28, 10) NOT NULL,
    volume NUMERIC(38, 10) NOT NULL
        CHECK (volume >= 0),

    PRIMARY KEY (instrument_id, interval_seconds, open_time),

    CHECK (high_price >= low_price),
    CHECK (open_price BETWEEN low_price AND high_price),
    CHECK (close_price BETWEEN low_price AND high_price)
);


-- 6. Запуски тестирования
-- code_snapshot — копия кода на момент запуска
-- parameters — параметры стратегии на момент запуска
-- Период: period_start включительно, period_end не включительно

CREATE TABLE backtests (
    backtest_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    strategy_id BIGINT NOT NULL REFERENCES strategies(strategy_id),
    code_snapshot TEXT NOT NULL,

    parameters JSONB NOT NULL DEFAULT '{}'
        CHECK (jsonb_typeof(parameters) = 'object'),

    period_start TIMESTAMPTZ NOT NULL,
    period_end TIMESTAMPTZ NOT NULL,
    interval_seconds INTEGER NOT NULL
        CHECK (interval_seconds > 0),

    status VARCHAR(20) NOT NULL DEFAULT 'queued'
        CHECK (
            status IN (
                'queued',
                'running',
                'completed',
                'failed',
                'cancelled'
            )
        ),

    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMPTZ,
    finished_at TIMESTAMPTZ,

    CHECK (period_end > period_start),
    CHECK (
        started_at IS NULL
        OR started_at >= created_at
    ),
    CHECK (
        finished_at IS NULL
        OR finished_at >= created_at
    ),
    CHECK (
        started_at IS NULL
        OR finished_at IS NULL
        OR finished_at >= started_at
    )
);

CREATE INDEX backtests_strategy_created_idx
ON backtests (strategy_id, created_at DESC);

CREATE INDEX backtests_status_created_idx
ON backtests (status, created_at);


-- 7. Инструменты, выбранные для тестирования

CREATE TABLE backtest_instruments (
    backtest_id BIGINT NOT NULL
        REFERENCES backtests(backtest_id) ON DELETE CASCADE,
    instrument_id BIGINT NOT NULL
        REFERENCES instruments(instrument_id),

    PRIMARY KEY (backtest_id, instrument_id)
);


-- 8. Прогнозы стратегии по времени
-- prediction_time — момент, на который сделан прогноз
--                      с использованием доступной тогда истории
-- target_time — будущий момент, для которого прогнозируется цена
-- actual_price — фактическая цена на target_time

CREATE TABLE backtest_results (
    backtest_id BIGINT NOT NULL,
    instrument_id BIGINT NOT NULL,
    prediction_time TIMESTAMPTZ NOT NULL,
    target_time TIMESTAMPTZ NOT NULL,
    predicted_price NUMERIC(28, 10) NOT NULL,
    actual_price NUMERIC(28, 10),

    PRIMARY KEY (
        backtest_id,
        instrument_id,
        prediction_time,
        target_time
    ),

    FOREIGN KEY (backtest_id, instrument_id)
        REFERENCES backtest_instruments(backtest_id, instrument_id)
        ON DELETE CASCADE,

    CHECK (target_time > prediction_time)
);


-- 9. Итоговые показатели теста целиком
-- metric_unit — единица измерения: percent, seconds и т. д.
-- Показатели в разных валютах нельзя складывать без пересчёта

CREATE TABLE backtest_metrics (
    backtest_id BIGINT NOT NULL
        REFERENCES backtests(backtest_id) ON DELETE CASCADE,
    metric_name VARCHAR(100) NOT NULL
        CHECK (length(trim(metric_name)) > 0),
    metric_value NUMERIC NOT NULL,
    metric_unit VARCHAR(30) NOT NULL
        CHECK (length(trim(metric_unit)) > 0),

    PRIMARY KEY (backtest_id, metric_name)
);


-- 10. Итоговые показатели по каждому инструменту

CREATE TABLE backtest_instrument_metrics (
    backtest_id BIGINT NOT NULL,
    instrument_id BIGINT NOT NULL,
    metric_name VARCHAR(100) NOT NULL
        CHECK (length(trim(metric_name)) > 0),
    metric_value NUMERIC NOT NULL,
    metric_unit VARCHAR(30) NOT NULL
        CHECK (length(trim(metric_unit)) > 0),

    PRIMARY KEY (backtest_id, instrument_id, metric_name),

    FOREIGN KEY (backtest_id, instrument_id)
        REFERENCES backtest_instruments(backtest_id, instrument_id)
        ON DELETE CASCADE
);


-- 11. Журнал выполнения тестов

CREATE TABLE backtest_logs (
    log_id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    backtest_id BIGINT NOT NULL
        REFERENCES backtests(backtest_id) ON DELETE CASCADE,
    level VARCHAR(10) NOT NULL DEFAULT 'info'
        CHECK (level IN ('info', 'warning', 'error')),
    message TEXT NOT NULL
        CHECK (length(trim(message)) > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX backtest_logs_test_time_idx
ON backtest_logs (backtest_id, created_at, log_id);

COMMIT;
