import React, { Component } from 'react';
import PropTypes from 'prop-types';
import moment from 'moment';

import { CurrencyIcon } from './currency-icon';
import { toPrettyNumber } from '../utils/helpers';

export class TradeRecord extends Component {
  render() {
    const { item } = this.props;

    const profitIsPositive = item.profit > 0;
    const createdAt = moment(item.createdAt);

    const date = createdAt.format('DD.MM.YYYY');
    const hours = createdAt.format('HH:mm:ss');

    return (
      <tr>
        <td>
          <CurrencyIcon iconName={this.props.item.currency} />
          <span className="currency-name">{this.props.item.currency}</span>
        </td>
        <td className="currency-date">
          {date} {hours}
        </td>
        <td className={'currency-profit' + (profitIsPositive ? ' currency-profit_positive' : '')}>
          {(profitIsPositive ? '+' : '') + toPrettyNumber(this.props.item.profit)} <span className="text text_lite">BTC</span>
        </td>
      </tr>
    );
  }
}

TradeRecord.propTypes = {
  item: PropTypes.shape({
    currency: PropTypes.string,
    buyPrice: PropTypes.number,
    sellPrice: PropTypes.number,
    exchangeName1: PropTypes.string,
    exchangeName2: PropTypes.string,
    amount: PropTypes.number,
    createdAt: PropTypes.string,
    profit: PropTypes.number,
    profitIsPositive: PropTypes.bool
  })
};