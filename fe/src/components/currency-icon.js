import React, { Component } from 'react';
import PropTypes from 'prop-types';

export class CurrencyIcon extends Component {
  getIconByName(name) {
    switch(name) {
    case 'LTC':
      return (<div className="currency-icon currency-icon_ltc" />);
    case 'BCH':
      return (<div className="currency-icon currency-icon_bch" />);
    case 'ETH':
      return (<div className="currency-icon currency-icon_eth" />);
    case 'DASH':
      return (<div className="currency-icon currency-icon_dash" />);
    case 'BCCOIN':
      return (<div className="currency-icon currency-icon_bccoin" />);
    case 'XMR':
      return (<div className="currency-icon currency-icon_xmr" />);
    case 'XRP':
      return (<div className="currency-icon currency-icon_xrp" />);
    case 'BCC':
      return (<div className="currency-icon currency-icon_bcc" />);
    case 'ZEC':
      return (<div className="currency-icon currency-icon_zec" />);
    case 'ETC':
      return (<div className="currency-icon currency-icon_etc" />);
    default:
      return null;
    }
  }
  render() {
    const icon = this.getIconByName(this.props.iconName);
    return icon;
  }
}

CurrencyIcon.propTypes = {
  iconName: PropTypes.string
};