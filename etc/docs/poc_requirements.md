# Crypto Bot PoC Requirements

This documents refers to the initial crypto bot development and defines the scope of the work, acceptance criteria and timeline.

The crypto bot requests the current exchange rates of different cryptocurrencies on separate crypto exchanges, simulates fake exchange rates and provides an API to retrieve that data as well as a UI for the end user.

## Scope

The bot consists of 2 parts: back-end and front-end described below.

### Back-End

This section describes back-end functionality.

#### Keep exchange rates

The back-end keeps the current exchange rates of different asset pairs in the following way:
1. Request asset pairs (see the list of asset pairs below) current exchange rate at a single source (find the cryptocompare below).
2. Save the current exchange rate of each asset pair to the database.
3. Wait 10 minutes and repeat. New values should replace the old.

Use this source to retrieve asset pairs prices: https://www.cryptocompare.com/api 
It's as easy, for example, as to get BTC to ETH and USD: https://min-api.cryptocompare.com/data/price?fsym=BTC&tsyms=ETH,USD

The application should only care about the following asset pairs (they should also be saved in the database for further extension):
- BTC/Bitcoin Cash
- BTC/Ethereum
- BTC/Dash
- BTC/BitConnect
- BTC/Monero
- BTC/Litecoin

#### Store trades

The bot (the front-end) creates trades and can also retrieve them.

#### Provide API

###### GET /v1/user-info

Retrieves user-related information.

- balance
- volume per trade

###### GET /v1/exchange-rates

Provides a random simulated exchange rate for buying and selling separately for all the existing asset pairs. The simulated exchange rate should be modified randomly and unioformly from -25% to 25% for buy and sell prices separately.
The administrator can increase or decrease the number of positive and negative trades by tuning these parameters.

Then this exchange rate should be fixed for the following 2-5 seconds (randomly). This price also holds for creating trades (see `/v1/trades` below).

Data returned for each asset pair:
- buy price
- sell price
- crypto exchange name

Crypto exchange name is assigned randomly just as an example from the list taken from https://cryptocoincharts.info/markets/info

###### PUT /v1/trade-volume

Sets the default trade volume.

###### POST /v1/deposit

Deposits the given amount of money to the balance.

###### POST /v1/trades

Creates a trade. The trade is:
- type: buy|sell
- currency1
- currency2
- amount
- crypto exchange

In response it gets the price `currency2` was traded for, an indication if the trade was successful (it won't be in case of insufficient funds, for example) and the new balance.

###### GET /v1/trades

Retrieves a list of all trades from the database. The higher, the more recent.

###### GET /v1/trade-stats

- trades number per day
- % of trades with profit
- % of trades with lost

Trades with lost are also possible as the rate could change when after the bot's requested the actual exchange rates and before it made a trade.

#### Misc

1. Each asset pair always contains BTC and the balance is kept in BTC as well.
2. Every significant parameters should be kept in the config.

Technology:
1. Golang for the application.
2. Percona as the database.
3. Docker Compose for containers orchestration.

Requirements to the logic and concerns from the customer:
1. There are postiv results and negativ ones. postive results, are results, which are profitable for the customer. i would like to determine how often a positiv result is displayed in the front-end
2. Provide a PoC first and then work on its stability so that he doesn't need a 24/7 developer to look for the app.
3. i want 99,9% reliable product at the end with a great documentation

### Front-End

The front-end represents the actual crypto bot, its functions are described in following sections.

#### Show exchange rates

The front-end keep updating the table with the current exchange rates using endpoint `/v1/exchange-rates` every 1 second (should be configurable in the config).
The difference is computed on the front-end side.

#### Run trade simulation

The flow for starting simulating:
1. Deposit money (defaut prompt with the amount) adding them to the total amount by requesting endpoint `/v1/deposit`.
2. Set BTC amount per trade (using +/- buttons) and change it using `/v1/trade-volume`.
3. Start simulation by hitting the start button. After starting the simulation the bot makes trades each the second the table is updated. It may or may not create a new trade using `/v1/trades`.
4. See how the bot buys and sells BTCs and logs the trades to the trades tab. Balance should be changing correspondingly.
5. The simulation stops when all its funds are exhausted or when the pause button is hit.

The front-end is responsible to:
- keeps balance up-to-date constantly requesting it using `/v1/user-info`
- 

Notes:
- the user can deposit more money while the simulation is working
- deposit amount, and BTC amount per trade and simulation status (started/stopped) should survive page reloads

#### Misc

The job includes converting provided design from PSD to HTML and it consists of two pages:
1. Signals.
2. Last Trades.

It has no specific technical requirements in a sense of browsers and screen sizes as this is a PoC, we only support modern browsers.

### Deployment

After the bot has been verified by the customer, it should be deployed to the server specified by the customer.

## Dealine

Fri, Oct 27 evening on our servers.

## Acceptance Criteria

This section describes formal criteria which the bot should be compliant to to be accepted by the customer.

1. The functionality listed in this document is implemented.
2. There's documentation on:
    - how to run the application
    - public back-end API
    - deployment