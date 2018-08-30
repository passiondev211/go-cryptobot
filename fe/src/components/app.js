/* eslint-disable indent */

import React, {Component} from 'react';
import {HashRouter, NavLink, Route, Switch} from 'react-router-dom';
import {translate} from 'react-i18next';
import PropTypes from 'prop-types';
import {toast, ToastContainer} from 'react-toastify';
import Joyride from 'react-joyride';

import {RatesPage} from './rates-page';
import {ReferPage} from './refer-page';
import {TradesPage} from './trades-page';
import {UserPanel} from './user-panel';
import {Header} from './header';
import TourLanguageModal from './tour-language-modal';

import {sessionTokenCookieName} from '../utils/config';
import Service from '../utils/service';
import ReferService from '../utils/refer-service';
import getSteps from '../utils/tourSteps';
import {readCookie, removeSessionToken} from '../utils/helpers';
import FreshChat from 'react-freshchat';

class Root extends Component {
    constructor(props) {
        super(props);

        this.bindMethods();

        let authData = this.checkSessionToken();
        if (authData) {
            Service.setHeader('Authorization', 'Bearer ' + authData.token);
            ReferService.setHeader('Authorization', 'Bearer ' + authData.token);
        }

        this.state = {
            userInfo: {
                balance: 0,
                profit: 0,
                margin: 0,
                currentLeverageValue: 1,
                botStarted: false,
                nextLeverage: {
                    minBalance: 0,
                    value: 1,
                },
            },
            badInternet: false,
            rows: [],
            authData,
            config: {
                apiTimeout: 0,
            },
            showTourLanguage: false,
            isRunning: false,
            isVerified: true,
            stepIndex: 0,
            steps: [],

            loadFreshChat:false,
            freshChatId:null
        };

        this.unsuccessRequestsCount = 0;

        this.getConfig(() => {
            Service.getTourState().then((tourStateData) => {
                const data = tourStateData.data ? tourStateData.data.response : tourStateData;
                if (data) {
                    if (this.userAuthorized) this.setIntervalUpdateDashboardInfo();
                } else {
                    this.setState({showTourLanguage: true});
                }
            });
        });

        if (authData) {
            Service.getFreshChatId().then((data)=>{
                if (!data) {
                    //console.log('User does not have freshchat id')
                    this.setState({loadFreshChat: true});
                    return;
                }
                //console.log('Successfully retrieved  freshchat id:'+data.chat_id)
                this.setState({freshChatId:data.chat_id}, ()=>{
                    this.setState({loadFreshChat: true});
                });
            });
        }
    }

    componentWillMount() {
        Service.getNotification().then(response => {
            if (response && response.length > 0) {
                toast.success(response, {autoClose: 10000});
            }
        });
        Service.getGTMEvents().then(response => {
            let event = {
                'event': response.Event,
                'userID': response.User,
            };
            if (response.Amount) {
                event.amount = response.Amount;
            }
            window.dataLayer = window.dataLayer || [];
            window.dataLayer.push(event);
        });
    }

    componentDidMount() {
        const { authData } = this.state;
        if (authData) {
            this.checkIfVerified(authData.decodedData.email);
        }
    }

    setIntervalUpdateDashboardInfo() {
        return new Promise((resolve) => {
            this.setIntervalOfMethod(this.updateDashboardInfo.bind(this, resolve), 1000);
        });
    }

    setBuyingStatus(buyingInProgress) {
        this.buyingInProgress = buyingInProgress;
    }

    setIntervalOfMethod(method, timeout) {
        return method().then(() => {
            this.unsuccessRequestsCount = 0;
            setTimeout(() => {
                this.setIntervalOfMethod(method, timeout);
            }, timeout);
        }).catch((err) => {
            this.handleServerError(err);
            if (this.unsuccessRequestsCount > 2) return;
            setTimeout(() => {
                this.setIntervalOfMethod(method, timeout);
            }, timeout);
        });
    }

    updateDashboardInfo(resolve = null) {
        return Service.getDashboardInfo().then((data) => {
            if (!data || this.buyingInProgress) return;

            this.setState(() => {
                const newTradeRow = data.rows.find(t => t.trade);
                const newState = {rows: data.rows};
                if (!newTradeRow) newState.userInfo = data.newUserInfo;
                newState.badInternet = false;
                return newState;
            });

            resolve();
        }).catch(er => this.handleServerError(er));
    }

    /** Checks for cookie 'session_token', takes jwt from it, parses and returns */
    checkSessionToken() {
        const cookie = readCookie(sessionTokenCookieName);
        if (!cookie) return;
        const jwtDecode = require('jwt-decode');
        return {token: cookie, decodedData: jwtDecode(cookie)};
    }

    /** Binds context for callbacks */
    bindMethods() {
        this.checkAuthError = this.checkAuthError.bind(this);
        this.setBuyingStatus = this.setBuyingStatus.bind(this);
        this.joyrideNextStep = this.joyrideNextStep.bind(this);
        this.setIntervalUpdateDashboardInfo = this.setIntervalUpdateDashboardInfo.bind(this);
        this.handleServerError = this.handleServerError.bind(this);
        //this.toggleLanguage = this.toggleLanguage.bind(this);
    }

    get userAuthorized() {
        return !!this.state.authData;
    }

    handleServerError(data) {
        if (data.response && data.response.status >= 400 && data.response.status < 500) {
            this.checkAuthError(data);
        } else if (data.serverUnavailable) {
            this.setState({
                badInternet: true
            });
            // this.unsuccessRequestsCount += 1;
        }
    }

    /** Checks server error, if it has code 401,
     * clears user info, auth cookie and token
     * */
    checkAuthError(data) {
        if (!data.response || data.response.status !== 401) {
            return;
        }
        

        //console.log('Got auth eror')
        let authData = this.checkSessionToken();
        if (authData) {
            Service.setHeader('Authorization', 'Bearer ' + authData.token);
            ReferService.setHeader('Authorization', 'Bearer ' + authData.token);

            this.setState({
                authData:authData
            });
        }
        else {
            //console.log('Session token was removed')
            removeSessionToken();
            Service.setHeader('Authorization', null);
            this.setState(() => {
                return {
                    authData: null,
                    botStarted: false
                };
            });
        }
        /*removeSessionToken();
        Service.setHeader('Authorization', null);
        this.setState(() => {
            return {
                authData: null,
                botStarted: false
            };
        });*/
    }

    checkIfVerified(email) {
        const { authData } = this.state;
        if (authData) {
            ReferService.getUserVerificationState(email).then((response) => {
                this.setState({
                    isVerified: response.data.response
                });
            });
        }
    }

    resendVerificationEmail() {
        const { authData } = this.state;
        if (authData) {
            ReferService.resendVerificationEmail(authData.decodedData.email)
            .then(() => {
                this.setState({
                    emailSentNotified: true
                });
                setTimeout(() => {
                    this.setState({
                        emailSentNotified: false
                    });
                }, 5000);
            });
        }
    }

    /** Gets config from server */
    getConfig(cb) {
        Service.getConfig().then((data) => {
            Service.setApiTimeout(data.apiTimeout);
            this.setState({config: data, badInternet: false});
            cb();
        }).catch(er => this.handleServerError(er));
    }

    onLangSelect = (lang) => {
        this.setState({
            isRunning: true,
            language: lang,
            steps: getSteps(lang),
            showTourLanguage: false
        });
    };

    tourCallback = ({action, type}) => {
        if (action === 'skip' || type === 'finished') {
            Service.setTourState().then(() => {
                this.setState({isRunning: false});
            });
        }
    };

    joyrideNextStep() {
        this.joyride.next();
    }

    // toggleLanguage() {
    //   if (this.props.i18n.language === 'en') {
    //     this.props.i18n.changeLanguage('ru');
    //   } else {
    //     this.props.i18n.changeLanguage('en');
    //   }
    // }

    _renderFreshchat() {
        const { authData, loadFreshChat, freshChatId} = this.state;

        if (authData&&authData.decodedData.email&&loadFreshChat) {
            if (freshChatId&&freshChatId.length>0) {
                return <div>
                <FreshChat
                    token = "62991200-f1e8-4cf7-89e9-1ce0a8e1279f"
                    externalId={authData.decodedData.email}
                    restoreId={freshChatId}   
                    onInit={widget => {
                        widget.user.get(resp=>{
                            var status = resp && resp.status;
                            //data = resp && resp.data;
                            //console.log('freshchat-subscribe-auth', resp)
                            if (status !== 200) {
                                widget.user.setProperties({
                                    email: authData.decodedData.email,    // user's email address
                                });
                            }
                        });
    
                        widget.on('user:created', function(resp) {
                            var status = resp && resp.status;
                            var data = resp && resp.data;
                            //console.log('freshchat-user-create', resp)
                            if (status === 200) {
                                if (data.restoreId) {
                                    Service.setFreshChatId(data.restoreId)
                                    .then(()=>{
                                        //console.log('FreshChatId was stored for this user')
                                    });                          // Update restoreId in your database
                                }
                            }
                        });
                    }}
                />
                </div>;
            }
            else {
                return <div>
                <FreshChat
                    token = "62991200-f1e8-4cf7-89e9-1ce0a8e1279f"
                    externalId={authData.decodedData.email}    
                    onInit={widget => {
                        widget.user.get(resp=>{
                            var status = resp && resp.status;
                            //data = resp && resp.data;
                            //console.log('freshchat-subscribe-auth', resp)
                            if (status !== 200) {
                                widget.user.setProperties({
                                    email: authData.decodedData.email,    // user's email address
                                });
                            }
                        });
    
                        widget.on('user:created', function(resp) {
                            var status = resp && resp.status,
                                data = resp && resp.data;
                            //console.log('freshchat-user-create', resp)
                            if (status === 200) {
                                if (data.restoreId) {
                                    Service.setFreshChatId(data.restoreId)
                                    .then(()=>{
                                        //console.log('FreshChatId was stored for this user')
                                    });                          // Update restoreId in your database
                                }
                            }
                        });
                    }}
                />
                </div>;
            }
            
        }
    }

    render() {
        const {t} = this.props;
        const {
            showTourLanguage,
            isRunning,
            stepIndex,
            steps,
            isVerified,
            emailSentNotified
        } = this.state;

        return (
            <HashRouter>
                <div>
                    <Joyride
                        ref={c => (this.joyride = c)}
                        debug={false}
                        locale={{
                            back: (this.state.language === 'DE' ? <span>ZURÜCK</span> : <span>BACK</span>),
                            close: (this.state.language === 'DE' ? <span>ÜBERSPRINGEN</span> : <span>CLOSE</span>),
                            last: (this.state.language === 'DE' ? <span>FERTIG</span> : <span>DONE</span>),
                            next: (this.state.language === 'DE' ? <span>WEITER</span> : <span>NEXT</span>),
                        }}
                        run={isRunning}
                        showOverlay={true}
                        autoStart
                        showSkipButton={true}
                        showStepsProgress={true}
                        disableOverlay={true}
                        stepIndex={stepIndex}
                        steps={steps}
                        type="continuous"
                        keyboardNavigation={false}
                        callback={this.tourCallback}
                    />
                    {!isVerified &&
                        <div className="topbar">
                            {!emailSentNotified &&
                                <div>
                                    Please verify your account.&nbsp;
                                    <a className="underlined" onClick={() => this.resendVerificationEmail()}>
                                        Click here
                                    </a>
                                    &nbsp;to resend the verification email.
                                </div>
                            }
                            {emailSentNotified &&
                                <div>
                                    We just sent you a new verification email.
                                </div>
                            }
                        </div>
                    }
                    <Header
                        //toggleLanguage={this.toggleLanguage}
                        authData={this.state.authData}
                        config={this.state.config}
                    />
                    <div style={{background: '#fff'}}>
                        <UserPanel
                            t={t}
                            joyrideIsRunning={isRunning}
                            joyrideNextStep={this.joyrideNextStep}
                            authData={this.state.authData}
                            userInfo={this.state.userInfo}
                            badInternet={this.state.badInternet}
                            config={this.state.config}
                            setIntervalUpdateDashboardInfo={this.setIntervalUpdateDashboardInfo}
                            handleServerError={this.handleServerError}
                        />
                    </div>
                    <nav className="navigation">
                        <NavLink exact to="/" className="navigation__btn"
                                 activeClassName="navigation__btn_selected">{t('nav.signals')}</NavLink>
                        <NavLink to="/refer" className="navigation__btn"
                                 activeClassName="navigation__btn_selected">{t('nav.refer')}</NavLink>
                        <NavLink to="/trades" className="navigation__btn"
                                 activeClassName="navigation__btn_selected">{t('nav.reports&stats')}</NavLink>
                    </nav>
                    <div className="page-wrapper">
                        <Switch>
                            <Route path="/trades" render={props => <TradesPage t={t} {...props} />}/>
                            <Route path="/refer" render={props => <ReferPage t={t} {...props} />}/>
                            <Route exact path="/"
                                   render={props => <RatesPage t={t} setBuyingStatus={this.setBuyingStatus}
                                                               rows={this.state.rows}
                                                               botStarted={this.state.userInfo.botStarted && !this.state.badInternet} {...props} />}/>
                        </Switch>
                    </div>
                    <ToastContainer
                        position="top-center"
                        autoClose={50000}
                        hideProgressBar={true}
                        newestOnTop={false}
                        closeOnClick={true}
                        pauseOnHover={true}
                        closeButton={false}
                        className={'profit-toast'}
                    />
                    {showTourLanguage && <TourLanguageModal onLangSelect={this.onLangSelect} t={t}/>}
                    {this._renderFreshchat()}
                </div>
            </HashRouter>
        );
    }
}

Root.propTypes = {
    t: PropTypes.func,
    i18n: PropTypes.object
};

export default translate('translations')(Root);
