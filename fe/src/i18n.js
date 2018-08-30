import i18n from 'i18next';
//import LanguageDetector from 'i18next-browser-languagedetector';
import { reactI18nextModule } from 'react-i18next';
import { en } from './locales/';

i18n // eslint-disable-line
  //.use(LanguageDetector)
  .use(reactI18nextModule)
  .init({
    fallbackLng: 'en',
    // have a common namespace used around the full app
    ns: ['translations'],
    defaultNS: 'translations',
    lng: 'en',
    debug: true,
    resources: {
      en: {
        translations: en
      }
    },
    react: {
      wait: true
    }
  });


export default i18n;