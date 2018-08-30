import React, { Component } from 'react';
import { Trans } from 'react-i18next';
import PropTypes from 'prop-types';
import moment from 'moment';

import { TradeRecord } from './trade-record';
import Service from '../utils/service';

const customStyles = {
  tradeBlock: {
    minHeight: '550px'
  }
};

export class TradesPage extends Component {
  constructor(props) {
    super(props);

    this.state = {
      trades: [],
      tradesLimit: 10,
      showLoadButton: true
    };

    this.loadMore = this.loadMore.bind(this);
    this.getTradesRecords = this.getTradesRecords.bind(this);

    this.getTradesRecords();
  }

  loadMore() {
    let nextTradeLimit = this.state.tradesLimit + 10;
    this.setState({ tradesLimit: nextTradeLimit }, () => {
      this.getTradesRecords();
    });
  }

  getTradesRecords() {
    Service
      .getTrades({
        limit: this.state.tradesLimit,
        since: moment().subtract(6, 'months').startOf('month').format('YYYY-MM-DD HH:mm:ss')
      })
      .then((data) => {
        this.setState({
          trades: data,
          showLoadButton: data.length >= 10 && this.state.tradesLimit <= data.length
        });
      });
  }

  render() {
    const { t } = this.props;
    const records = this.state.trades && this.state.trades.map((trade, i) => <TradeRecord key={i} item={trade} />);

    return (
      <div className="page trades-page container">
        <div className="page-text">
          <Trans i18nKey="descriptions.trades">
            <p>On this easy to navigate dashboard you can see all your past trades,</p>
            <p>your profit and your total return.</p>
          </Trans>
        </div>
        <div className="page-bg trade-block" style={customStyles.tradeBlock}>
          <div className="page-title">{t('headers.tradesHistory')}</div>
          <table className="trades-table">
            <thead>
              <tr>
                <th>{t('tableHeaders.trades')}</th>
                <th>{t('tableHeaders.date')}</th>
                <th>{t('tableHeaders.profit')}</th>
              </tr>
            </thead>
            <tbody>
              {records}
            </tbody>
          </table>

          {this.state.showLoadButton &&
            <div className="text-center py-4">
              <button className="btn btn_blue btn-lg" onClick={this.loadMore}>Load more</button>
            </div>
          }
        </div>
      </div>
    );
  }
}

TradesPage.propTypes = {
  trades: PropTypes.array,
  t: PropTypes.func
};
