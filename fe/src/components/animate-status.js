import React, {
  Component
} from 'react';
import PropTypes from 'prop-types';
import {
  tradeStatuses,
  tradeStatusesAfterSuccess,
  tradeStatusesAfterFail,
  tradeStatusesTimes,
  statusesSuccessChecks,
} from '../utils/config';
import RadarAnimation from './radar-animation';
import WinAnimation from './win-animation';
import LoseAnimation from './lose-animation';
import BuyingAnimation from './buying-animation';
import SellingAnimation from './selling-animation';

export default class AnimateStatus extends Component {
  constructor(props) {
    super(props);
    this.state = {
      status: null,
    };
  }

  componentDidMount() {
    this.updateCurrentStatus(this.props.exchangeRate.trade ? tradeStatuses.BUYING : tradeStatuses.ANALYZE);
  }

  componentWillReceiveProps(nextProps) {
    const oldRates = this.props.exchangeRate;
    const newRates = nextProps.exchangeRate;
    const tradeCreated = newRates.trade && this.state.status === tradeStatuses.ANALYZE;
    const ratesWasUpdated = oldRates.buyPrice !== newRates.buyPrice || oldRates.sellPrice !== oldRates.sellPrice;
    if (tradeCreated) {
      this.updateCurrentStatus(tradeStatuses.BUYING);
    } else if (ratesWasUpdated) {
      this.updateCurrentStatus(tradeStatuses.ANALYZE_CANCELED);
    }
  }

  componentWillUnmount() {
    this.statusTimer && clearTimeout(this.statusTimer);
  }

  setNextStatusTimeout(currentStatus) {
    const success = !statusesSuccessChecks[currentStatus] || statusesSuccessChecks[currentStatus](this.props.exchangeRate);
    const nextStatus = success ? tradeStatusesAfterSuccess[currentStatus] : tradeStatusesAfterFail[currentStatus];
    const currentStatusTime = tradeStatusesTimes[currentStatus];
    if (currentStatusTime) {
      this.statusTimer && clearTimeout(this.statusTimer);
      this.statusTimer = setTimeout(() => this.updateCurrentStatus(nextStatus), currentStatusTime);
    }
  }

  onStatusChanged(oldStatus, newStatus) {
    const { ANALYZE, BUYING, SELLING } = tradeStatuses;
    if (oldStatus === ANALYZE && newStatus === BUYING) {
      this.props.setBuyingStatus(true);
    } else if (oldStatus === SELLING) {
      this.props.setBuyingStatus(false);
    }
  }

  updateCurrentStatus(newStatus) {
    const oldStatus = this.state.status;
    this.setState(() => {
      return {
        status: newStatus,
      };
    });
    this.onStatusChanged(oldStatus, newStatus);
    newStatus && this.setNextStatusTimeout(newStatus);
  }

  render() {
    const { currency, buyExchange, sellExchange, trade } = this.props.exchangeRate;
    switch (this.state.status) {
    case tradeStatuses.ANALYZE:
      return (<RadarAnimation currency={currency} className="status"/>);
    case tradeStatuses.ANALYZE_CANCELED:
      return (<RadarAnimation currency={currency} red={true} className="status"/>);
    case tradeStatuses.BUYING:
      return (trade && <BuyingAnimation currency={currency} buyExchange={buyExchange} profit={trade.expectedProfit?trade.expectedProfit:0} className="status"/>);
    case tradeStatuses.SELLING:
      return (trade && <SellingAnimation currency={currency} sellExchange={sellExchange} profit={trade.expectedProfit?trade.expectedProfit:0} className="status"/>);
    case tradeStatuses.SUCCESS_TRADE:
      return (trade && <WinAnimation profit={trade.actualProfit} className="status"/>);
    case tradeStatuses.UNSUCCESS_TRADE:
      return (trade && <LoseAnimation profit={trade.actualProfit} sellExchange={sellExchange} className="status"/>);
    }
    return (<div className="status status_hidden"/>);
  }
}

AnimateStatus.propTypes = {
  exchangeRate: PropTypes.shape({
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
  setBuyingStatus: PropTypes.func,
};
