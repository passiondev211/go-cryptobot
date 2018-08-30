import React, { Component } from 'react';
import PropTypes from 'prop-types';

export default class RadarAnimation extends Component {
  render() {
    const { currency } = this.props;
    return (
      <div className="radar">
        <div className={'target target-1' + (this.props.red ? ' target_red' : '')} />
        <div className={'target target-2' + (this.props.red ? ' target_red' : '')} />
        <div className={'target target-3 big-dot' + (this.props.red ? ' target_red' : '')} />
        <div className={'target target-4' + (this.props.red ? ' target_red' : '')} />
        <div className={'target target-5 big-dot' + (this.props.red ? ' target_red' : '')} />
        <div className={'target target-6' + (this.props.red ? ' target_red' : '')} />
        <div className={'target target-7' + (this.props.red ? ' target_red' : '')} />
        <div className={'target target-8 big-dot' + (this.props.red ? ' target_red' : '')} />
        <div className={'target target-9' + (this.props.red ? ' target_red' : '')} />
        <div className={'target target-10 big-dot' + (this.props.red ? ' target_red' : '')} />
        <div className="line line-1" />
        <div className="line line-2" />
        <div className="circle circle-1" />
        <div className="circle circle-2" />
        <div className="circle circle-3" />
        <div className="circle circle-4" />
        <div className={'title' + (this.props.red ? ' title_red' : '')} />
        <svg xmlns="http://www.w3.org/2000/svg" className={'arrow' + (this.props.red ? ' arrow_red' : '')} width="90" height="90" xmlnsXlink="http://www.w3.org/1999/xlink">
          <use fill="#FFF" fillOpacity=".01" xlinkHref={`#${currency}_a5`} transform="matrix(1 0 0 1 0 1)" />
          <use fill={`url(#${currency}_c5)`} xlinkHref={`#${currency}_b5`} transform="rotate(-61 61 -6) translate(0 -18.99)" opacity="1" />
          <defs>
            <linearGradient id={`${currency}_c5`} x2={1} gradientUnits="userSpaceOnUse" gradientTransform="matrix(-58.427 -24.1543 23.986 -58.0057 43.063 99.78)">
              <stop offset="0" stopColor="#22A7FE" />
              <stop offset="1" stopColor="#6DC0F7" stopOpacity="0" />
            </linearGradient>
            <linearGradient id={`${currency}_e5`} x2={1} gradientUnits="userSpaceOnUse" gradientTransform="matrix(-66.0826 -.2576 30.835 -7909.83 64.482 3955.68)">
              <stop offset="0" stopColor="#6DC0F7" />
              <stop offset="1" stopColor="#6DC0F7" stopOpacity="0" />
            </linearGradient>
            <path id={`${currency}_a5`} d="M162 81c0 44.735-36.265 81-81 81S0 125.735 0 81 36.265 0 81 0s81 36.265 81 81z" />
            <path id={`${currency}_b5`} fillRule="evenodd" d="M42.317 59.982l30.527-14.32C59.747 18.634 32.05 0 0 0v33.704c18.554 0 34.608 10.706 42.317 26.278z" />
            <path id={`${currency}_d5`} d="M-.004.163C-.28.165-.502.39-.5.667c.002.276.228.498.504.496l-.008-1zM79.342.5c.276-.002.498-.228.496-.504-.002-.276-.228-.498-.504-.496l.008 1zM.004 1.163L79.342.5l-.008-1L-.004.163l.008 1z" />
          </defs>
        </svg>
      </div>
    );
  }
}

RadarAnimation.propTypes = {
  red: PropTypes.bool,
  currency: PropTypes.string,
};