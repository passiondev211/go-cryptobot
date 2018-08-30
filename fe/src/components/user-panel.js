/* eslint-disable indent */

import React, {Component} from 'react';
import PropTypes from 'prop-types';
import Joyride from 'react-joyride';
import Modal from 'react-modal';
import Service from '../utils/service';
import {toPrettyNumber} from '../utils/helpers';
import {MarginChart} from './margin-chart';

const customStyles = {
    overlay: {
        backgroundColor: 'rgba(0,0,0,0.4)',
        zIndex: 1600
    },
    content: {
        top: '50%',
        left: '50%',
        right: 'auto',
        bottom: 'auto',
        border: 'none',
        borderRadius: '15px',
        boxShadow: '0 0 20px 0 rgba(0, 0, 0, 0.05)',
        marginRight: '-50%',
        transform: 'translate(-50%, -50%)',
        maxHeight: '90vh',
        overflow: 'auto',
    }
};

export class UserPanel extends Component {
    constructor(props) {
        super(props);

        this.state = {
            disabledToggleButton: false,
            modalIsOpen: false,
            tooltip: false,
            toggleBotState: false
        };

        this.botToggle = this.botToggle.bind(this);
        this.openModal = this.openModal.bind(this);
        this.closeModal = this.closeModal.bind(this);
    }

    componentDidMount() {
        this.joyride.addTooltip({
            text: 'Leverage refers to the leverage effect of our capital on your investment. The leverage can significantly increase the returns on a given investment amount. Depending on your investment, you will get a leverage of 1:10, 1:20, 1:30 or 1:40. Leverage 1:10 means that if you invest 0,1 Bitcoins, you still can trade with 1 Bitcoin. We increase your investment by the factor ten.',
            selector: '.account-data__question',
            trigger: '.account-data__question',
            type: window.innerWidth >= 768 ? 'hover' : 'click',
            event: window.innerWidth >= 768 ? 'hover' : 'click',
            position: 'bottom',
            style: {
                backgroundColor: 'white',
                borderRadius: '0.2rem',
                color: '#000',
                textAlign: 'center',
                width: '30rem',
                zIndex: 150
            }
        });
    }

    componentWillReceiveProps(nextProps) {
        const userAuthorized = !!nextProps.authData;
        this.changeToggleButtonState(!userAuthorized);
    }

    /** Toggles bot state  */
    botToggle() {
        const {authData, userInfo, joyrideIsRunning, joyrideNextStep, setIntervalUpdateDashboardInfo} = this.props;

        const userAuthorized = !!authData;
        if (!userAuthorized) return;


        // this.changeToggleButtonState(false);

        Service.toggleBot(!userInfo.botStarted).then(() => {
            Service.getTourState().then((tourStateData) => {
                const data = tourStateData.data ? tourStateData.data.response : tourStateData;
                if (!data) {
                    setIntervalUpdateDashboardInfo().then(() => {
                        if (joyrideIsRunning) joyrideNextStep();
                    });
                }

                // this.changeToggleButtonState(true);
            }).catch(er => this.props.handleServerError(er));
        }).catch(er => this.props.handleServerError(er));
    }

    openModal(e) {
        e.preventDefault();
        this.setState({modalIsOpen: true});
    }

    closeModal(e) {
        e.preventDefault();
        this.setState({modalIsOpen: false});
    }

    /**
     * Changes Toggle button Disable attribute state(if true button will be disabled)
     * @param state New state
     */
    changeToggleButtonState(state) {
        this.setState(() => {
            return {
                disabledToggleButton: state,
            };
        });
    }

    addInvest() {
        window.location.href = this.props.config.depositRedirectPath;
    }

    getModeLabel() {
        const {t, userInfo} = this.props;
        const realMode = userInfo.mode === 'real';
        if (realMode) return null;
        return (<div className="text account-data__demo">{t('demo')}</div>);
    }

    getStartButtonText() {
        const {t, userInfo} = this.props;
        const realMode = userInfo.mode === 'real';
        if (realMode) {
            return userInfo.botStarted ? t('buttons.stopReal') : t('buttons.startReal');
        }
        return userInfo.botStarted ? t('buttons.stopDemo') : t('buttons.startDemo');
    }

    getLeverageBlock() {
        const {currentLeverageValue} = this.props.userInfo;
        const leverageBlock = [(<span key={0} className="bold">1:{currentLeverageValue}</span>)];
        /*    if (nextLeverage) {
                      leverageBlock[1] = (
                        <span key={1}>
                          &nbsp;(<span>{t('headers.invest')}</span>&nbsp;
                          <span className="bold">{toPrettyNumber(nextLeverage.minBalance - balance)}</span>
                          <span className="text text_currency text_lite">BTC</span>&nbsp;
                          <span>{t('headers.for')}</span>&nbsp;
                          <span className="bold">1:{nextLeverage.value}</span>)
                        </span>
                      );
                    }*/
        return leverageBlock;
    }

    render() {
        const {t, userInfo} = this.props;
        const {balance, profit, margin} = userInfo;

        return (
            <div className="container" style={{maxWidth: 'initial'}}>
                <div className="row align-items-center profit-block py-5">
                    <div className="col-md col-sm text-center order-md-2 order-sm-2 profit-center">
                        <h5>{t('headers.yourProfit')}</h5>
                        <p className="profit" id="tour-profit">
              <span className="tour-profit-highlighter">
                <span className="tour-profit-value">{toPrettyNumber(profit)}</span> <span
                  className="profit-currency">BTC</span>
              </span>
                        </p>
                        <a className="earn-more" href="#" onClick={this.openModal}>Want to earn more?</a>
                        <MarginChart t={t} value={margin}/>
                    </div>
                    <div className="col-md col-sm text-center text-sm-left order-md-3 order-sm-3 mt-5 mt-lg-1">
                        <h5 className="d-inline profit-subheader">{t('headers.leverage')}</h5>
                        <div className="account-data__question"/>
                        <p id="tour-leverage" className="leverage mt-2">
              <span className="tour-leverage-highlighter">
                {this.getLeverageBlock()}
              </span>
                        </p>
                        <button disabled={this.state.disabledToggleButton}
                                onClick={this.botToggle} className="btn btn_blue btn-lg demo-btn">
                            {this.getStartButtonText()}
                        </button>
                    </div>
                    <div className="col-md col-sm text-center text-sm-right order-md-1 order-sm-1 mt-5 mt-lg-1">
                        {this.getModeLabel()}
                        <h5 className="d-inline profit-subheader">{t('headers.accountBalance')}</h5>
                        <p id="tour-balance" className="balance mt-2">
              <span className="tour-balance-highlighter">
                {toPrettyNumber(balance)} <span className="profit-currency">BTC</span>
              </span>
                        </p>
                        <button disabled={this.props.badInternet} onClick={() => this.addInvest()}
                                className="btn btn_orange btn-lg invest-btn">{t('buttons.invest')}</button>
                    </div>

                    <div/>

                    <Joyride
                        ref={c => this.joyride = c}
                        steps={[]}
                    />
                    <Modal
                        isOpen={this.state.modalIsOpen}
                        style={customStyles}
                        onRequestClose={this.closeModal}
                        contentLabel="Want to earn more?"
                    >
                        <div className="increase_close">
                            <a href="#" onClick={this.closeModal}/>
                        </div>
                        <div className="increase_block">
                            <div className="wrapper">
                                <h2>Invest more. Earn more.</h2>
                                <p>
                                    You want to make higher profit? No problem!
                                    Depending on your investment you are able to trade with a certain leverage.
                                    The higher your leverage, the higher your daily profit.
                                </p>
                            </div>
                        </div>
                        <div className="increase_list">
                            <div className="wrapper">
                                <ul>
                                    <li>
                                        <header>
                                            <h6>daily profit</h6>
                                            <div>0.15%</div>
                                        </header>
                                        <footer>
                                            Leverage 1:10 <p>Investment: <strong>0.0 BTC</strong></p></footer>
                                    </li>
                                    <li>
                                        <header>
                                            <h6>daily profit</h6>
                                            <div>0.30%</div>
                                        </header>
                                        <footer>
                                            Leverage 1:20 <p>Investment: <strong>0.5 BTC</strong></p></footer>
                                    </li>
                                    <li>
                                        <header>
                                            <h6>daily profit</h6>
                                            <div>0.55%</div>
                                        </header>
                                        <footer>
                                            Leverage 1:30 <p>Investment: <strong>1.0 BTC</strong></p></footer>
                                    </li>
                                    <li>
                                        <header>
                                            <h6>daily profit</h6>
                                            <div>1.00%</div>
                                        </header>
                                        <footer>
                                            Leverage 1:40 <p>Investment: <strong>2.0 BTC</strong></p></footer>
                                    </li>
                                </ul>
                            </div>
                            <div className="invest-btn">
                                <button onClick={() => {
                                    this.addInvest();
                                    return false;
                                }} className="btn btn_orange btn-lg">{t('buttons.invest')}</button>
                            </div>
                        </div>
                    </Modal>
                </div>
            </div>
        );
    }
}

UserPanel.propTypes = {
    authData: PropTypes.any,
    t: PropTypes.any,
    config: PropTypes.any,
    joyrideIsRunning: PropTypes.bool,
    joyrideNextStep: PropTypes.func,
    botToggle: PropTypes.func,
    setIntervalUpdateDashboardInfo: PropTypes.func,
    handleServerError: PropTypes.func,
    badInternet: PropTypes.bool,
    userInfo: PropTypes.shape({
        balance: PropTypes.number,
        profit: PropTypes.number,
        margin: PropTypes.number,
        currentLeverageValue: PropTypes.number,
        nextLeverage: PropTypes.shape({
            minBalance: PropTypes.number,
            value: PropTypes.number,
        }),
        config: PropTypes.shape({
            depositRedirectPath: PropTypes.string
        })
    })
};
