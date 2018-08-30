import Axios from 'axios';
import { toast } from 'react-toastify';
import { referApiUrl, apiTimeout } from './config';

function statusIsCorrect(status) {
  return status >= 200 && status < 300;
}

const axios = Axios.create({
  baseURL: referApiUrl,
  timeout: apiTimeout,
  validateStatus: function (status) {
    return statusIsCorrect(status);
  },
});

const axiosUser = Axios.create({
  baseURL: 'https://profitcoins.io/wp-json/pbc/v1',
  timeout: apiTimeout,
  validateStatus: function (status) {
    return statusIsCorrect(status);
  },
});

/** Validator for server response */
function validateServerResponse(response) {
  if (!response) {
    throw Promise.reject('Server returned incorrect response');
  }
  if (response.data.error) {
    throw Promise.reject('Server Error: ' + response.data.error);
  }
  if (statusIsCorrect(response.status)) {
    return;
  }
  throw Promise.reject('Server Error: ' + response.statusText);
}

/** interceptor for server response. These callbacks fire before any 'than' or 'catch' */
axios.interceptors.response.use(function (response) {
  validateServerResponse(response);
  return (response.data && response.data.response) || response;
}, function (error) {
  return Promise.reject(error);
});

const toasts = {};
function addToast(text) {
  if (!toasts[text] || (!toast.isActive(toasts[text].id) && !toasts[text].timerId)) {
    toasts[text] = { id: toast.error(text), timerId: setTimeout(() => toasts[text].timerId = null, 100)};
  }
}

/**
 * Handler for server errors
 * */
function handleError(error) {
  if (error.response) {
    // The request was made and the server responded with a status code
    // that falls out of the range of 2xx
    if (error.response.status >= 400 && error.response.status < 500) {
      addToast(error.response.data && error.response.data.error);
      return Promise.reject(error);
    }
    return Promise.reject(error);
  } else if (error.request) {
    // The request was made but no response was received
    // `error.request` is an instance of XMLHttpRequest in the browser and an instance of
    // http.ClientRequest in node.js
    addToast('Server is not responding. Please try later');
    return Promise.reject({ serverUnavailable: true });
  } else {
    // Something happened in setting up the request that triggered an Error
    toast.error(`Server error: ${error.message}`);
    addToast('Server is not responding. Please try later');
    return Promise.reject({ serverUnavailable: true });
  }
}


export default class ReferService {
  static options = { headers: {} };
  static userOptions = {
    headers: {
      'Content-Type': 'application/json'
    }
  };
  static timeout = 5000;
  static hash = '9c262e49102a982dc2e033544dd3973d'

  /** Sets header for all requests */
  static setHeader(key, val) {
    this.options.headers[key] = val;
  }

  static setApiTimeout(timeout) {
    axios.defaults.timeout = timeout;
  }


  static getUserReferalLink(query){
    return axios.get('/'+ query.userId + '/', this.options)
      .catch(handleError);
  }

  static getUserVerificationState(email) {
    return axiosUser.get(`/user_is_verified?email=${email}&hash=${this.hash}`, this.userOptions).catch(handleError);
  }

  static resendVerificationEmail(email) {
    return axiosUser.get(`/send_email?email=${email}&hash=${this.hash}`, this.userOptions).catch(handleError);
  }

  static getReferredFriends(query){
    return axios.get('/referred/'+ query.userId + '/', this.options)
      .catch(handleError);
  }

  static getInvestedFriends(query){
    return axios.get('/invested/'+ query.userId + '/', this.options)
      .catch(handleError);
  }

  static getCommission(query){
    return axios.get('/commission/'+ query.userId + '/', this.options)
      .catch(handleError);
  }
}
