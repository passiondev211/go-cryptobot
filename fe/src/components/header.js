import React, { Component } from 'react';
import PropTypes from 'prop-types';
import { removeSessionToken } from '../utils/helpers';
import ClickOutside from 'react-click-outside';

import Service from '../utils/service';

import gb from '../img/gb.svg';
import de from '../img/de.svg';
import fr from '../img/fr.png';
import es from '../img/es.png';

export class Header extends Component {
  constructor(props) {
    super(props);

    this.state = {
      menuHidden: true,
      menuMobileHidden: true,
      langHidden: true
    };

    this.logOut = this.logOut.bind(this);
  }

  toggleUserMenuBlock() {
    this.setState(() => {
      return {
        menuHidden: !this.state.menuHidden
      };
    });
  }

  closeUserMenuBlock() {
    this.setState(() => {
      return {
        menuHidden: true
      };
    });
  }

  toggleUserMenuMobileBlock() {
    this.setState(() => {
      document.documentElement.classList.toggle('no-scroll');
      return {
        menuMobileHidden: !this.state.menuMobileHidden
      };
    });
  }

  closeUserMenuMobileBlock() {
    this.setState(() => {
      return {
        menuMobileHidden: true
      };
    });
  }

  closeLangMenuBlock() {
    this.setState(() => {
      return {
        langHidden: true
      };
    });
  }

  toggleLangMenuBlock() {
    this.setState(() => {
      return {
        langHidden: !this.state.langHidden
      };
    });
  }

  logOut(event) {
    event.preventDefault();
    Service.setHeader('Authorization', null);
    removeSessionToken();
    window.location.href = this.props.config.logoutRedirectPath;
  }

  render() {
    let unathorized;
    let email;

    if (this.props.authData && this.props.authData.decodedData.email) {
      email = this.props.authData.decodedData.email;
      unathorized = false;
    } else {
      email = '';
      unathorized = true;
    }

    return (
      <header>
        <nav className="header navbar navbar-expand-lg navbar-light bg-white">
          <div className="container">
            <a className="navbar-brand" href="/" />

            <ul className="faq-link-wrapper navbar-nav ml-auto">
              <li className="nav-item dropdown">
                <a className="dropdown-item" href="https://profitcoins.io/faq/" target="_blank">FAQ</a>
              </li>
            </ul>

            <button className="navbar-toggler" type="button" onClick={() => this.toggleUserMenuMobileBlock()}>
              <span className="navbar-toggler-icon" />
            </button>

            <div className="collapse navbar-collapse" id="navbarSupportedContent">
              <ul className="navbar-nav ml-auto">
                {!unathorized ?
                  (
                    <li className="nav-item dropdown">
                      <ClickOutside onClickOutside={() => this.closeUserMenuBlock()}>
                        <a className="nav-link dropdown-toggle" href="#" onClick={() => this.toggleUserMenuBlock()}>{email}</a>

                        <div className={`dropdown-menu ${!this.state.menuHidden ? 'show': ''}`} aria-labelledby="navbarDropdown">
                          <a className="dropdown-item" href={this.props.config.firstLinkRedirectPath}>{this.props.config.firstLinkName || 'Settings'}</a>
                          <a className="dropdown-item" href={this.props.config.secondLinkRedirectPath}>{this.props.config.secondLinkName || 'Invest'}</a>
                          <div className="dropdown-divider" />
                          <a className="dropdown-item" href="https://profitcoins.io/logout" onClick={this.logOut}>Logout</a>
                        </div>
                      </ClickOutside>
                    </li>
                  ) : ''
                }
                <li className="nav-item dropdown">
                  <ClickOutside onClickOutside={() => this.closeLangMenuBlock()}>
                    <a className="nav-link dropdown-toggle" href="#" onClick={() => this.toggleLangMenuBlock()}>
                      <img src={gb} width="20" alt="UK" />
                    </a>

                    <div className={`dropdown-menu ${!this.state.langHidden ? 'show': ''}`} style={{ minWidth: 50 }}>
                      <a className="dropdown-item button-disabled" href="#"><img src={de} width="20" alt="UK" /></a>
                      <a className="dropdown-item button-disabled" href="#"><img src={fr} width="20" alt="UK" /></a>
                      <a className="dropdown-item button-disabled" href="#"><img src={es} width="20" alt="UK" /></a>
                    </div>
                  </ClickOutside>
                </li>
              </ul>
            </div>
          </div>
        </nav>

        <div className={`mobile-menu-container ${!this.state.menuMobileHidden ? 'show-menu' : ''}`}>
          <ul className="mobile-menu-links">
            <li>
              <ul className="mobile-sub-menu">
                <li><a href={this.props.config.firstLinkRedirectPath}>{this.props.config.firstLinkName || 'Settings'}</a></li>
                <li><a href={this.props.config.secondLinkRedirectPath}>{this.props.config.secondLinkName || 'Invest'}</a></li>
                <li><a href="#" onClick={this.logOut}>Logout</a></li>
              </ul>
            </li>
            <li>
              <a href="/"><img src={gb} width="20" alt="UK" /> English</a>
              <ul className="mobile-sub-menu">
                <li><a href="javascript:void(0)" className="button-disabled"><img src={de} className="image-mask" width="20" alt="DE" /> Deutsch - (Coming Soon)</a></li>
                <li><a href="javascript:void(0)" className="button-disabled"><img src={fr} className="image-mask" width="20" alt="FR" /> Fran&#231;ais - (Coming Soon)</a></li>
                <li><a href="javascript:void(0)" className="button-disabled"><img src={es} className="image-mask" width="20" alt="ES" /> Espanol - (Coming Soon)</a></li>
              </ul>
            </li>
          </ul>
        </div>
      </header>
    );
  }
}

//onClick={() => this.props.toggleLanguage()}

Header.propTypes = {
  //toggleLanguage: PropTypes.func
  authData: PropTypes.shape( {
    token: PropTypes.string,
    decodedData: PropTypes.shape( {
      email: PropTypes.string,
      id: PropTypes.number
    })
  }),
  config: PropTypes.shape( {
    logoutRedirectPath: PropTypes.string,
    firstLinkRedirectPath: PropTypes.string,
    secondLinkRedirectPath: PropTypes.string,
    firstLinkName: PropTypes.string,
    secondLinkName: PropTypes.string
  })
};
