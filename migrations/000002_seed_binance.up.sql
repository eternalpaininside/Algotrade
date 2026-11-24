BEGIN;

INSERT INTO algotrade.exchanges (code, name)
VALUES ('binance', 'Binance')
ON CONFLICT (code) DO NOTHING;

INSERT INTO algotrade.instruments (exchange_id, symbol, name, quote_currency)
SELECT exchange_id, 'BTCUSDT', 'BTC/USDT', 'USDT'
FROM algotrade.exchanges
WHERE code = 'binance'
ON CONFLICT (exchange_id, symbol) DO NOTHING;

COMMIT;