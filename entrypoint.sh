#!/bin/bash -e

cat > ./conf/config.json << EOF
{
  "port": "80",
  "maxTradesPerDay": 0,
  "maxMarginPerDay": 0,
  "Auth": {
    "accessToken": "5G8AXoY8ASASm943ZQb9iDmp8EWEVuvB",
    "authTokenTtlSeconds": 2419200,
    "jwtSecret": "N7QCg3fz89YmtsbMzSkMuMDsmuJxZD32",
    "jwtTokenTtlMinutes": 10080
  },
  "rates": {
    "buyingMinPercents": 1,
    "buyingMaxPercents": 1.001,
    "sellingMinPercents": 1,
    "sellingMaxPercents": 1.001,
    "tableUpdateIntervalSeconds": 600,
    "demoRatesUpdateIntervalMinSeconds": 10,
    "demoRatesUpdateIntervalMaxSeconds": 20,
    "realRatesUpdateIntervalMinSeconds": 1800,
    "realRatesUpdateIntervalMaxSeconds": 1900
  },
  "database": {
    "name": "${CRYPTOBOT_DB_NAME}",
    "user": "${CRYPTOBOT_DB_USER}",
    "host": "${CRYPTOBOT_DB_HOST}",
    "port": "3306",
    "password": "${CRYPTOBOT_DB_PASSWORD}"
  },
  "frontEnd": {
    "apiTimeout": 5000,
    "logoutRedirectPath": "https://${DOMAIN}/logout/",
    "depositRedirectPath": "https://${DOMAIN}/my-account/invest/",
    "firstLinkRedirectPath": "https://${DOMAIN}/my-account/",
    "secondLinkRedirectPath": "https://${DOMAIN}/my-account/invest/",
    "firstLinkName": "Account",
    "secondLinkName": "Invest",
    "invalidJwtRedirectPath": "https://${DOMAIN}"
  },
  "demoBalance": 1,
  "minBalance": 0.00009,
  "leverages": [
    {
      "minBalance": 0.0,
      "leverage": 10,
      "profitPerDay": 0.25,
      "dailyTradeBoundLower": 0.0005,
      "dailyTradeBoundUpper": 0.002,
      "resultTradeBoundLower": 0.0014,
      "resultTradeBoundUpper": 0.0016
    },
    {
      "minBalance": 0.5,
      "leverage": 20,
      "profitPerDay": 0.5,
      "dailyTradeBoundLower": 0.001,
      "dailyTradeBoundUpper": 0.004,
      "resultTradeBoundLower": 0.0028,
      "resultTradeBoundUpper": 0.0032
    },
    {
      "minBalance": 1,
      "leverage": 30,
      "profitPerDay": 1,
      "dailyTradeBoundLower": 0.003,
      "dailyTradeBoundUpper": 0.007,
      "resultTradeBoundLower": 0.004,
      "resultTradeBoundUpper": 0.006
    },
    {
      "minBalance": 2,
      "leverage": 40,
      "profitPerDay": 1.5,
      "dailyTradeBoundLower": 0.007,
      "dailyTradeBoundUpper": 0.013,
      "resultTradeBoundLower": 0.009,
      "resultTradeBoundUpper": 0.011
    }
  ],
  "debug": {
    "log": true,
    "db": false
  },
  "robot": {
    "timeoutAfterNewTradeS": 10,
    "timeToExecutingTradeMS": 1600
  },
  "fixFactor": {
    "min": -100,
    "max": 100
  },
  "security": {
    "allowed_admin_ip": [
      "172.18.0.0/24",
      "172.17.0.8/16",
      "46.28.204.34", 
      "46.28.207.150", 
      "95.183.55.80",
      "10.42.153.45",
      "10.42.40.164",
      "10.42.29.34",
      "10.42.14.97",
      "10.42.1.210",
      "10.0.0.0/8",
      "185.35.139.32",
      "185.35.137.91",
      "103.100.134.104"
    ]
  }
}
EOF

go run main.go ./conf/config.json
