import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { CurrencyIcon } from './currency-icon';
import AnimateStatus from './animate-status';
import { toPrettyNumber } from '../utils/helpers';

export class ExchangeRatePair extends Component {
  render() {
    const { t, row, botStarted, setBuyingStatus } = this.props;
    const { currency, buyPrice, sellPrice, buyExchange, sellExchange } = row;
    return (
      <div className="card dashboard-card mt-3">
        <div className="card-body">
          <div className="row align-items-center justify-content-center">
            <div className="col-sm-6 text-center d-md-none trade-data-wrap">
              <div className="rate-row__currency">
                <CurrencyIcon iconName={currency} />
                <div className="rate-row__currency-name">
                  {currency}
                </div>
              </div>

              <div className="row currency-title">
                <div className="col buyat">{buyExchange}</div>
                <div className="col sellto bdl1">{sellExchange}</div>
              </div>

              <div className="row value">
                <div className="col">{toPrettyNumber(buyPrice)}<b> {currency}</b></div>
                <div className="col bdl1">{toPrettyNumber(sellPrice)}<b> {currency}</b></div>
              </div>
            </div>

            <br />
            <br />
            <br />
            <br />

            <div className="col-sm-4 text-center text-lg-left d-none d-md-block">
              <div className="rate-row__currency">
                <CurrencyIcon iconName={currency} />
                <div className="rate-row__currency-name">
                  {currency}
                </div>
              </div>

              <h5 className="card-title">{t('signalsTable.buyAt')}</h5>
              <p className="buyat">{buyExchange}</p>
              <p className="value"><b>{t('signalsTable.for')} {toPrettyNumber(buyPrice)}</b> {currency}</p>
              <div className="card-arrow" />
            </div>
            <div className="col-sm-4 mt-5 text-center text-lg-left d-none d-md-block">
              <h5 className="card-title">{t('signalsTable.sellTo')}</h5>
              <p className="sellto">{sellExchange}</p>
              <p className="value"><b>{t('signalsTable.for')} {toPrettyNumber(sellPrice)}</b> {currency}</p>
            </div>

            <div className="animation-wrapper col-sm-4 text-center text-lg-left">
              { botStarted ? ( <AnimateStatus setBuyingStatus={setBuyingStatus} exchangeRate={row} />) : null }
            </div>
          </div>
        </div>
      </div>
    );
  }
}
ExchangeRatePair.propTypes = {
  row: PropTypes.shape({
    currency: PropTypes.string,
    buyPrice: PropTypes.number,
    sellPrice: PropTypes.number,
    buyExchange: PropTypes.string,
    sellExchange: PropTypes.string,
    trade: PropTypes.shape({
      expectedProfit: PropTypes.number,
      actualProfit: PropTypes.number,
    }),
  }),
  botStarted: PropTypes.bool,
  t: PropTypes.func,
  setBuyingStatus: PropTypes.func,
};
