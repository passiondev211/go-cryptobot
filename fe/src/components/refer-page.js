/* eslint-disable no-unused-vars */
/* eslint-disable indent */

import React, {Component} from 'react';
import PropTypes from 'prop-types';
import {Trans} from 'react-i18next';
import moment from 'moment';
import {CopyToClipboard} from 'react-copy-to-clipboard';
import ReferService from '../utils/refer-service';
import Service from '../utils/service';

const customStyle = {
    blueStyle: {
        color: '#1175ba'
    }
};

export class ReferPage extends Component {


    constructor(props) {
        super(props);

        this.state = {
            userData: {},
            userId: 5,
            referLinkCopied: false,
            referLink: '',
            dateValue: moment().format('DD/MM/YYYY'),
            referFriends: [],
            referFriendsCount: 0,
            investedFriends: 0,
            commission: '0,-',
            chartData: {
                labels: ['', '01.01', '02.01', '03.01', '04.01', '05.01', '06.01', '07.01', '08.01', '09.01', '10.01', ''],
                datasets: [{
                    fill: false,
                    borderColor: '#d4d9e0',
                    pointRadius: 5,
                    pointBorderColor: '#1175bb',
                    pointBackgroundColor: '#1175bb',
                    data: [null, 125, 100, 250, 325, 175, 100, 400, 350, 500, 500, null]
                }]
            },


        };

        this.getUserInfo = this.getUserInfo.bind(this);
        this.getUserReferalLink = this.getUserReferalLink.bind(this);
        this.getReferredFriends = this.getReferredFriends.bind(this);
        this.getInvestedFriends = this.getInvestedFriends.bind(this);
        this.getCommission = this.getCommission.bind(this);

        this.getUserInfo();


    }

    getUserInfo() {
        Service.getUserInfo()
            .then((response) => {
                    let data = response || {id: 0};
                    this.getUserReferalLink(data);
                },
                (error) => {
                });

    }

    getUserReferalLink(data) {
        ReferService.getUserReferalLink({userId: data.id})
            .then((response) => {
                    let link = response.data || '';
                    this.setState({referLink: link, referLinkCopied: false, userData: data, userId: data.id});
                    this.getReferredFriends();
                    this.getInvestedFriends();
                    this.getCommission();
                },
                (error) => {
                });
    }

    getReferredFriends() {
        ReferService.getReferredFriends({userId: this.state.userId})
            .then((response) => {
                    let friends = parseInt(response.data);
                    this.setState({referFriends: [], referFriendsCount: friends || 0});
                },
                (error) => {
                });
    }

    getInvestedFriends() {
        ReferService.getInvestedFriends({userId: this.state.userId})
            .then((response) => {
                    let invested = parseInt(response.data);
                    this.setState({investedFriends: isNaN(invested) ? 0 : invested});
                },
                (error) => {
                });
    }

    getCommission() {
        ReferService.getCommission({userId: this.state.userId})
            .then((response) => {
                    let comm = '0,-';
                    if (!isNaN(parseFloat(response.data)) && parseFloat(response.data) > 0) {
                        comm = response.data.toString().replace('.', ',');
                    }
                    this.setState({commission: comm});
                },
                (error) => {
                });
    }

    onCopy() {
        //this.state.referLinkCopied = true;
        this.setState({referLinkCopied: true});
        setTimeout(function () {
            this.setState({referLinkCopied: false});
        }.bind(this), 1000);
    }

    onDateChange(moment) {
        //console.log(moment);
        this.setState({dateValue: moment.format('DD/MM/YYYY')});
    }

    render() {
        const chartOptions = {
            maintainAspectRatio: false,
            layout: {
                padding: {
                    left: 0,
                    right: 0,
                    top: 0,
                    bottom: 0
                }
            },
            legend: {
                display: false
            },
            elements: {
                line: {
                    tension: 0, // disables bezier curves
                }
            },
            scales: {
                xAxes: [{
                    gridLines: {
                        display: false
                    }
                }],
                yAxes: [{
                    gridLines: {
                        display: false
                    },
                    fontFamily: 'Roboto',
                    fontSize: 13,
                    fontColor: '#a0acb9',
                    ticks: {
                        min: 50,
                        max: 550,
                        step: 50,
                        callback: function (value, index, values) {
                            if (value == 550) return '';
                            return '$' + value + ',-';
                        }
                    }
                }]
            },
            tooltips: {
                backgroundColor: '#2976d1',
                position: 'nearest',
                xPadding: 15,
                yPadding: 15,
                caretSize: 5,
                xAlign: 'center',
                yAlign: 'bottom',
                bodySpacing: 10,
                titleMarginBottom: 10,
                displayColors: false,
                callbacks: {
                    title: function (tooltipItem, data) {
                        return '2 Friends Refferred';
                    },
                    label: function (tooltipItem, data) {
                        return '3 Friends Invested';
                    },
                    afterLabel: function (tooltipItem, data) {
                        // var val = data['datasets'][0]['data'][tooltipItem['index']];

                        return '$12,- Earned';
                    }
                }
            }
        };

        const copyText = this.state.referLinkCopied ? 'Copied' : 'Copy Link';
        let {referLink, dateValue, referFriendsCount, investedFriends, commission, chartData} = this.state;
        return (
            <div className="page refer-page container">
                <div className="page-text">
                    <Trans i18nKey="descriptions.referFriends">
                        <p>Help us growing and get great rewards! Get 0.005 BTC for every
                            friend you invite to Profitcoins.io.</p>
                    </Trans>
                </div>
                <div className="page-bg">
                    <div className="row justify-content-center mb-4">
                        <div className="col-lg-5 col-md-8">
                            <div className="form">
                                <div className="form-group mb-4">
                                    <label htmlFor="copy-link" className="sr-only">Refer Link</label>
                                    <div className="input-group">
                                        <input ref="copyLink" id="copy-link" type="text" className="form-control"
                                               value={referLink}/>
                                        <div className="input-group-append">
                                            <CopyToClipboard text={referLink} onCopy={this.onCopy.bind(this)}>
                                                <button type="button" className="btn btn-primary">{copyText}</button>
                                            </CopyToClipboard>
                                        </div>
                                    </div>
                                </div>
                                <div className="form-group mb-5">
                                    <label htmlFor="email" className="sr-only">Refer Link</label>
                                    <div className="input-group">
                                        <input id="email" name="email" type="text" className="form-control" value=""
                                               placeholder="Friend’s email address"/>
                                        <div className="input-group-append">
                                            <button type="button" className="btn btn-orange">Invite</button>
                                        </div>
                                    </div>
                                </div>
                                <div className="form-group d-flex justify-content-center">
                                    <a href={'https://www.facebook.com/sharer/sharer.php?u=' + referLink}
                                       target="_blank" className="btn btn-share-round btn-facebook"><i
                                        className="fa fa-facebook"/> Share</a>
                                    <a href={'https://twitter.com/home?status=' + referLink} target="_blank"
                                       className="btn btn-share-round btn-twitter"><i className="fa fa-twitter"/> Share</a>
                                    <a href={'https://plus.google.com/share?url=' + referLink} target="_blank"
                                       className="btn btn-share-round btn-google-plus"><i
                                        className="fa fa-google-plus"/> Share</a>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div className="row justify-content-center">
                        <div className="col-lg-11 col-md-9 refer-list-center">
                            <ul className="row">
                                <li className="col-lg-4"><i className="icon icon-link"/> Send your unique link to
                                    friends so that they can sign up
                                </li>
                                <li className="col-lg-4"><i className="icon icon-bitcoin"/> Your friends register and
                                    make a minimum deposit of 0.07 BTC
                                </li>
                                <li className="col-lg-4"><i className="icon icon-gift"/> As a reward, we will credit you
                                    0.005 BTC
                                </li>
                            </ul>
                        </div>
                    </div>
                    {/*<div className="row justify-content-center">
                   <div className="col-lg-3">
                   <div className="form-group form-group ref-calendar">
                   <label htmlFor="ref-calendar-input" className="sr-only">Refer Link</label>
                   <Datetime
                   defaultValue={dateValue}
                   inputProps={{ placeholder: 'dd/mm/yyyy', className: 'form-control datepicker'}}
                   locale={'en'}
                   onChange={moment => this.onDateChange(moment)}
                   dateFormat={'DD/MM/YYYY'}
                   timeFormat={''} />
                   </div>
                   </div>
                   </div>*/}
                    <div className="row justify-content-center" style={{marginTop: '50px'}}>
                        <div className="col-lg-9 col-md-12">
                            <div className="circle-block">
                                <div className="circle-block-value">{referFriendsCount}</div>
                                <div className="circle-block-descr">Friends Referred</div>
                            </div>
                            <div className="circle-block">
                                <div className="circle-block-value">{investedFriends}</div>
                                <div className="circle-block-descr">Friends Invested</div>
                            </div>
                            <div className="circle-block">
                                <div className="circle-block-value">{commission}</div>
                                <div className="circle-block-descr">BTC Earned</div>
                            </div>
                        </div>

                    </div>
                    {/*<div className='row justify-content-center'>
                   <div className='col refer-chart-col'>
                   <div className="chart-container">
                   <Line data={chartData} options={chartOptions} height={400} />
                   </div>
                   </div>
                   </div>*/}
                </div>
            </div>
        );
    }

}

ReferPage.propTypes = {
    onDateChange: PropTypes.func,
    t: PropTypes.func
};
