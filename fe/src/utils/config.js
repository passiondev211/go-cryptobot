/* eslint-disable no-undef */
// export const apiUrl =  __DEV__ ? 'http://localhost:9512/api/v1' : 'https://api.profitcoins.io/api/v1'; //eslint-disable-line
//export const apiUrl = 'https://app.profitcoins.io/api/v1'; //eslint-disable-line
//export const referApiUrl = 'https://profitcoins.io/wp-json/pcr/v1/user';
const domain = __DOMAIN__ || 'profitcoins.io'; //eslint-disable-line

var _apiUrl;
if(__DEV__) {
  _apiUrl = 'http://localhost:9512/api/v1';
} else {
  _apiUrl = 'https://app.' + domain + '/api/v1';
}

export const apiUrl = _apiUrl;
export const referApiUrl = 'https://' + domain + '/wp-json/pcr/v1/user';
export const sessionTokenCookieName = 'session_token';
export const requestsIntervalForUnavailableServer = 30000; //ms
export const apiTimeout = 10000;

export const tradeStatuses = {
  get ANALYZE() { return 'ANALYZE'; },
  get ANALYZE_CANCELED() { return 'ANALYZE_CANCELED'; },
  get BUYING() { return 'BUYING'; },
  get SELLING() { return 'SELLING'; },
  get SUCCESS_TRADE() { return 'SUCCESS_TRADE'; },
  get UNSUCCESS_TRADE() { return 'UNSUCCESS_TRADE'; },
  get UPDATE_RATES() { return 'UPDATE_RATES'; }
};

export const tradeStatusesTimes = {
  get [tradeStatuses.ANALYZE]() { return 0; }, //ms
  get [tradeStatuses.ANALYZE_CANCELED]() { return 300; }, //ms
  get [tradeStatuses.BUYING]() { return 3000; }, //ms
  get [tradeStatuses.SELLING]() { return 3000; }, //ms
  get [tradeStatuses.SUCCESS_TRADE]() { return 1000; }, //ms
  get [tradeStatuses.UNSUCCESS_TRADE]() { return 1000; }, //ms
};

export const tradeStatusesAfterSuccess = {
  get [tradeStatuses.ANALYZE]() { return tradeStatuses.BUYING; },
  get [tradeStatuses.ANALYZE_CANCELED]() { return tradeStatuses.ANALYZE; },
  get [tradeStatuses.BUYING]() { return tradeStatuses.SELLING; },
  get [tradeStatuses.SELLING]() { return tradeStatuses.SUCCESS_TRADE; },
  get [tradeStatuses.SUCCESS_TRADE]() { return tradeStatuses.UPDATE_RATES; },
  get [tradeStatuses.UNSUCCESS_TRADE]() { return tradeStatuses.ANALYZE; },
};

export const tradeStatusesAfterFail = {
  get [tradeStatuses.SELLING]() { return tradeStatuses.UNSUCCESS_TRADE; },
};

export const statusesSuccessChecks = {
  [tradeStatuses.SELLING]: (exchangeRate) => exchangeRate.trade.actualProfit >= 0,
};
