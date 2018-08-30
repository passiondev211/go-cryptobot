/* eslint-disable indent */

import Axios from 'axios';
import {toast} from 'react-toastify';
import {apiTimeout, apiUrl} from './config';

/** Returns true if status is correct */
function statusIsCorrect(status) {
    return status >= 200 && status < 300;
}

const axios = Axios.create({
    baseURL: apiUrl,
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
        toasts[text] = {id: toast.error(text), timerId: setTimeout(() => toasts[text].timerId = null, 100)};
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
        addToast('Bad internet connection');
        return Promise.reject({serverUnavailable: true});
    } else {
        // Something happened in setting up the request that triggered an Error
        toast.error(`Server error: ${error.message}`);
        addToast('Bad internet connection');
        return Promise.reject({serverUnavailable: true});
    }
}

export default class Service {
    static options = {
        headers: {
            'Content-Type': 'application/json'
        },
        hash: '9c262e49102a982dc2e033544dd3973d'
    };
    static timeout = 5000;

    /** Sets header for all requests */
    static setHeader(key, val) {
        this.options.headers[key] = val;
    }

    static setApiTimeout(timeout) {
        axios.defaults.timeout = timeout;
    }

    static getConfig() {
        return axios.get('/app-config', this.options)
            .catch(handleError);
    }

    static getTrades(query) {
        return axios.get(`/trades?limit=${query.limit}`, this.options)
            .catch(handleError);
    }

    static getDashboardInfo() {
        return axios.get('/dashboard-info', this.options)
            .catch(handleError);
    }

    static getFreshChatId() {
        return axios.get('/chat-id', this.options)
            .catch(handleError);
    }

    static setFreshChatId(chatId) {
        return axios.post('/chat-id', {chat_id:chatId}, this.options)
            .catch(handleError);
    }

    static toggleBot(state) {
        return axios.post(`/bot-toggle?state=${state ? 'on' : 'off'}`, null, this.options)
            .catch(handleError);
    }

    static getTourState() {
        return axios.get('/tour-visited', this.options)
            .catch(handleError);
    }

    static setTourState() {
        return axios.post('/tour-visited', null, this.options)
            .catch(handleError);
    }

    static logOut() {
        return axios.post('/logout', null, this.options)
            .catch(handleError);
    }

    static getUserInfo() {
        return axios.get('/user-info', this.options)
            .catch(handleError);
    }

    static getNotification() {
        return axios.get('/notification', this.options)
            .catch(handleError);
    }

    static getGTMEvents() {
        return axios.get('/gtm-events', this.options)
            .catch(handleError);
    }
}
