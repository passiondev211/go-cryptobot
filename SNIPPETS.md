# Snippets

## Reset user

```sql
UPDATE users SET balance=1, leverage=40, bot_started=NULL, current_profit=0, tour_visited=1; DELETE FROM trades;
```

## Reset user and it's trading bounds

```sql
UPDATE users SET balance=1, leverage=40, bot_started=NULL, current_profit=0, tour_visited=1, daily_trade_profit_bounds_lower=-0.002, daily_trade_profit_bounds_upper=0.002, result_trade_profit_bounds_lower=0.001, result_trade_profit_bounds_upper=0.002; DELETE FROM trades;
```

## Set config values via API

```shell
curl -H 'Authorization: Bearer 5G8AXoY8ASASm943ZQb9iDmp8EWEVuvB' -X POST -d '{"data":{"demoBalance":1,"leverages":[{"minBalance":0,"leverage":10,"profitPerDay":0.25,"dailyTradeBoundLower":-0.003,"dailyTradeBoundUpper":0.003,"resultTradeBoundLower":0.002,"resultTradeBoundUpper":0.003},{"minBalance":0.036624,"leverage":20,"profitPerDay":0.5,"dailyTradeBoundLower":-0.0055,"dailyTradeBoundUpper":0.0055,"resultTradeBoundLower":0.0045,"resultTradeBoundUpper":0.0055},{"minBalance":0.182174,"leverage":30,"profitPerDay":1,"dailyTradeBoundLower":-0.0105,"dailyTradeBoundUpper":0.0105,"resultTradeBoundLower":0.0095,"resultTradeBoundUpper":0.0105},{"minBalance":0.364348,"leverage":40,"profitPerDay":1.5,"dailyTradeBoundLower":-0.0155,"dailyTradeBoundUpper":0.0165,"resultTradeBoundLower":0.0145,"resultTradeBoundUpper":0.0165}]}}' http://127.0.0.1:9512/api/v1/config
```
