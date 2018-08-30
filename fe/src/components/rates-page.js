/* eslint-disable indent */

import React, {Component} from 'react';
import PropTypes from 'prop-types';
import {Trans} from 'react-i18next';
import {ExchangeRatePair} from './exchange-rate-pair';

export class RatesPage extends Component {

    render() {
        let {t, rows, botStarted, setBuyingStatus} = this.props;
        const pairs = rows && rows.map((row, index) => (
            <ExchangeRatePair setBuyingStatus={setBuyingStatus} botStarted={botStarted} t={t} key={index} row={row}/>));
        return (
            <div className="container page">
                <div className="page-text">
                    <Trans i18nKey="descriptions.signals">
                        <p>With the Bitcoins you transfer to your account, the Profit Coins software buys
                            cryptocurrencies from one</p>
                        <p>marketplace and sells them for a higher price to another marketplace.</p>
                        <p>Automatically.</p>
                    </Trans>
                </div>
                {pairs}
            </div>
        );
    }
}

RatesPage.propTypes = {
    rows: PropTypes.array,
    t: PropTypes.func,
    botStarted: PropTypes.bool,
    setBuyingStatus: PropTypes.func,
};