import React from 'react';
import PropTypes from 'prop-types';

import gb from '../img/gb.svg';
import de from '../img/de.svg';
import fr from '../img/fr.png';
import es from '../img/es.png';

export default class TourLanguageModal extends React.Component {
  render() {
    const { onLangSelect } = this.props;

    return (
      <div className="tour-language-modal">
        <div className="tour-language-modal__overlay" />
        <div className="tour-language-modal__content">
          <p className="tour-language-modal__title">Please choose your language</p>
          <div>
            <button className="btn btn_language" onClick={() => onLangSelect('EN')}>
              <img src={gb} width="20" alt="UK" /> English
            </button>
            <button className="btn btn_language" onClick={() => onLangSelect('DE')}>
              <img src={de} width="20" alt="DE" /> Deutsch
            </button>
            <button className="btn btn_language" disabled>
              <img src={fr} width="20" alt="FR" /> Fran&#231;ais
            </button>
            <button className="btn btn_language" disabled>
              <img src={es} width="20" alt="ES" /> Espa&#241;ol
            </button>
          </div>
        </div>
      </div>
    );
  }
}

TourLanguageModal.propTypes = {
  onLangSelect: PropTypes.func
};
