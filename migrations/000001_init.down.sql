BEGIN;

DROP TABLE algotrade.backtest_logs;
DROP TABLE algotrade.backtest_instrument_metrics;
DROP TABLE algotrade.backtest_metrics;
DROP TABLE algotrade.backtest_results;
DROP TABLE algotrade.backtest_instruments;
DROP TABLE algotrade.backtests;
DROP TABLE algotrade.candles;
DROP TABLE algotrade.instruments;
DROP TABLE algotrade.exchanges;
DROP TABLE algotrade.strategies;
DROP TABLE algotrade.users;

DROP FUNCTION algotrade.update_strategy_timestamp();

DROP SCHEMA algotrade;

COMMIT;