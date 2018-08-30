/* eslint-disable import/default */
import React from 'react';
import { render } from 'react-dom';
import { AppContainer } from 'react-hot-loader';
import './styles/styles.scss';
import App from './components/app';
import './i18n';

// Favicons
require('./favicons/android-icon-36x36.png');
require('./favicons/android-icon-48x48.png');
require('./favicons/android-icon-72x72.png');
require('./favicons/android-icon-96x96.png');
require('./favicons/android-icon-144x144.png');
require('./favicons/android-icon-192x192.png');

require('./favicons/apple-icon.png');
require('./favicons/apple-icon-57x57.png');
require('./favicons/apple-icon-60x60.png');
require('./favicons/apple-icon-72x72.png');
require('./favicons/apple-icon-76x76.png');
require('./favicons/apple-icon-114x114.png');
require('./favicons/apple-icon-120x120.png');
require('./favicons/apple-icon-144x144.png');
require('./favicons/apple-icon-152x152.png');
require('./favicons/apple-icon-180x180.png');
require('./favicons/apple-icon-precomposed.png');

require('./favicons/favicon.ico');
require('./favicons/favicon-16x16.png');
require('./favicons/favicon-32x32.png');
require('./favicons/favicon-96x96.png');

require('./favicons/ms-icon-70x70.png');
require('./favicons/ms-icon-144x144.png');
require('./favicons/ms-icon-150x150.png');
require('./favicons/ms-icon-310x310.png');

/** If you are not able to run server locally then store `session_token` like this */
// import { setCookie } from './utils/helpers';
// setCookie('session_token', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MiwiZW1haWwiOiIxQDIuY29tIiwibGFuZ3VhZ2UiOiJlbiIsImV4cCI6MTUyMTk1NTQ3Nn0.EEuU9ET1EJ2K7LjTAxClhDQ4Ryy_RMjl6dEcJaKc5sQ', 7);

/** Root of react application */
render(
  <AppContainer>
    <App />
  </AppContainer>,
  document.getElementById('app')
);
