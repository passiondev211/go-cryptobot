import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { toPrettyNumber } from '../utils/helpers';

export class MarginChart extends Component {
  render() {
    const value = this.props.value;
    const marginIsNegative = value < 0;
    return (
      <div className="profit__chart">
        <div className="profit__circle">
          <div className="profit__inside-circle">
            <div className={`profit__margin-value ${marginIsNegative ? 'profit__margin-value_red' : ''}`}>
              {toPrettyNumber(value, 2)}%
            </div>
            <div className="profit__margin-label">
              {this.props.t('headers.margin')}
            </div>
          </div>
          {/*<div className={`profit__indicator ${marginIsNegative ? 'profit__indicator_red' : '' }`} style={{height: Math.abs(value) + '%'}} />*/}
        </div>
      </div>
    );
  }
}

MarginChart.propTypes = {
  value: PropTypes.number,
  t: PropTypes.func
};