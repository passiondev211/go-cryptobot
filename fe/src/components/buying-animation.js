import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { toPrettyNumber } from '../utils/helpers';

export default class BuyingAnimation extends Component {
  constructor(props) {
    super(props);
    this.state = {
      counterValue: 6,
    };
  }

  componentDidMount() {
    this.setCounterTimer(5);
  }

  componentWillUnmount() {
    if (this.timer) {
      clearTimeout(this.timer);
      this.timer = null;
    }
  }

  setCounterTimer(newVal) {
    if (this.state.counterValue <= 3) return;
    this.timer = setTimeout(() => {
      this.setState(() => {
        return {
          counterValue: newVal
        };
      });
      this.setCounterTimer(newVal - 1);
    }, 1000);
  }

  render() {
    const { currency } = this.props;

    return (
      <div className="buying-container counter-container">
        <div className="counter-wrap">
          <div className="coin">
            <svg xmlns="http://www.w3.org/2000/svg" width={42} height={22} xmlnsXlink="http://www.w3.org/1999/xlink" className="coin__svg">
              <use fill={`url(#${currency}_b)`} fillOpacity=".5" xlinkHref={`#${currency}_a`} transform="rotate(-90 13.32 6.25)" />
              <use fill={`url(#${currency}_d)`} fillOpacity=".5" xlinkHref={`#${currency}_c`} transform="rotate(-90 7.64 7.64)" />
              <g transform="translate(-1195 1284)">
                <use xlinkHref={`#${currency}_e`} transform="rotate(-90 -23.39 -1239.61)" className="main-color" />
                <mask id={`${currency}_f`}>
                  <path fill="#fff" d="M1215.21-1262v-22h22v22z" />
                  <use xlinkHref={`#${currency}_e`} transform="rotate(-90 -23.39 -1239.61)" />
                </mask>
                <use fill="#FFF" xlinkHref={`#${currency}_g`} transform="rotate(-90 -23.39 -1239.61)" mask={`url(#${currency}_f)`} />
              </g>
              <g transform="translate(-1195 1284)">
                <use xlinkHref={`#${currency}_h`} transform="rotate(-90 -23.42 -1243.87)" className="second-color" />
                <mask id={`${currency}_i`}>
                  <path fill="#fff" d="M1219.45-1266.29v-14h14v14z" />
                  <use xlinkHref={`#${currency}_h`} transform="rotate(-90 -23.42 -1243.87)" />
                </mask>
                <use fill="#FFF" fillOpacity=".6" xlinkHref={`#${currency}_j`} transform="rotate(-90 -23.42 -1243.87)" mask={`url(#${currency}_i)`} />
              </g>
              <defs>
                <linearGradient id={`${currency}_b`} x2={1} gradientUnits="userSpaceOnUse" gradientTransform="matrix(0 -33.93 17.14 0 0 33.93)">
                  <stop offset={0} className="coin-tail" />
                  <stop offset={1} className="coin-tail" stopOpacity={0} />
                </linearGradient>
                <linearGradient  id={`${currency}_d`} x2={1} gradientUnits="userSpaceOnUse" gradientTransform="matrix(0 -41 8.57 0 0 41)">
                  <stop offset={0} className="coin-tail" />
                  <stop offset={1} className="coin-tail" stopOpacity={0} />
                </linearGradient>
                <path id={`${currency}_a`} d="M0 0h17.143v25.36c0 4.733-3.838 8.57-8.572 8.57C3.84 33.93 0 30.094 0 25.36V0z" />
                <path id={`${currency}_c`} d="M0 0h8.57v36.714C8.57 39.08 6.654 41 4.287 41S0 39.08 0 36.714V0z" />
                <path id={`${currency}_e`} d="M20 9.897c0 5.465-4.477 9.896-10 9.896S0 15.363 0 9.897C0 4.43 4.477 0 10 0s10 4.43 10 9.897z" />
                <path id={`${currency}_g`} d="M19 9.897c0 4.903-4.02 8.896-9 8.896v2c6.065 0 11-4.87 11-10.896h-2zm-9 8.896c-4.98 0-9-3.993-9-8.896h-2c0 6.027 4.935 10.896 11 10.896v-2zM1 9.897C1 4.993 5.02 1 10 1v-2C3.935-1-1 3.87-1 9.897h2zM10 1c4.98 0 9 3.993 9 8.897h2C21 3.87 16.065-1 10-1v2z" />
                <path id={`${currency}_h`} d="M11.43 5.655c0 3.123-2.56 5.655-5.716 5.655S0 8.778 0 5.655 2.558 0 5.714 0C8.87 0 11.43 2.532 11.43 5.655z" />
                <path id={`${currency}_j`} d="M10.43 5.655c0 2.56-2.102 4.655-4.716 4.655v2c3.7 0 6.715-2.97 6.715-6.655h-2zM5.713 10.31C3.1 10.31 1 8.216 1 5.655h-2c0 3.686 3.016 6.655 6.714 6.655v-2zM1 5.655C1 3.095 3.1 1 5.714 1v-2C2.016-1-1 1.97-1 5.655h2zM5.714 1c2.614 0 4.715 2.094 4.715 4.655h2C12.43 1.97 9.412-1 5.713-1v2z" />
              </defs>
            </svg>
          </div>
          <div className="counter">
            <div className="counter__number"><span>{this.state.counterValue}</span></div>
            <div className="counter__border" />
            <svg xmlns="http://www.w3.org/2000/svg" width={60} height={60} xmlnsXlink="http://www.w3.org/1999/xlink" className="counter-svg">
              <use xlinkHref={`#${currency}_a1`} className="inner-fill" />
              <g transform="translate(-1267 1306)">
                <mask id={`${currency}_c1`}>
                  <path fill="#fff" d="M1271-1302h52v52h-52z" />
                  <use xlinkHref={`#${currency}_b1`} transform="translate(1272 -1301)" />
                </mask>
                <use fill="#FFF" fillOpacity=".2" xlinkHref={`#${currency}_d`} transform="translate(1272 -1301)" mask={`url(#${currency}_c)`} />
              </g>
              <defs>
                <path id={`${currency}_a1`} d="M60 30C60 46.57 46.57 60 30 60 13.43 60 0 46.57 0 30 0 13.43 13.43 0 30 0 46.57 0 60 13.43 60 30z" />
                <path id={`${currency}_b1`} d="M50 25c0 13.807-11.193 25-25 25S0 38.807 0 25 11.193 0 25 0s25 11.193 25 25z" />
                <path id={`${currency}_d1`} d="M49 25c0 13.255-10.745 24-24 24v2c14.36 0 26-11.64 26-26h-2zM25 49C11.745 49 1 38.255 1 25h-2c0 14.36 11.64 26 26 26v-2zM1 25C1 11.745 11.745 1 25 1v-2C10.64-1-1 10.64-1 25h2zM25 1c13.255 0 24 10.745 24 24h2C51 10.64 39.36-1 25-1v2z" />
              </defs>
            </svg>
          </div>
        </div>
        <div className="counter__text">
          <div className="counter__title element-animation">Buying {currency} at {this.props.buyExchange}</div>
          <div className="counter__course">Expected Profit: {toPrettyNumber(this.props.profit)}</div>
        </div>
      </div>
    );
  }
}

BuyingAnimation.propTypes = {
  currency: PropTypes.string,
  buyExchange: PropTypes.string,
  profit: PropTypes.number,
};