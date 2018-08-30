import { sessionTokenCookieName } from './config';

export function toPrettyNumber(num, countOfFixedNumbers = 8) {
  return num ? num.toFixed(countOfFixedNumbers) : 0;
}

export function getCountOfDecimals(num) {
  const pieces = num.toString().split('.');
  return pieces[1] && pieces[1].length;
}
/** Solution for avoiding rounding problem https://stackoverflow.com/a/12830454/3445316 */
export function round(num, scale) {
  if(!('' + num).includes('e')) {
    return +(Math.round(num + 'e+' + scale)  + 'e-' + scale);
  } else {
    const arr = ('' + num).split('e');
    let sig = '';
    if(+arr[1] + scale > 0) {
      sig = '+';
    }
    return +(Math.round(+arr[0] + 'e' + sig + (+arr[1] + scale)) + 'e-' + scale);
  }
}

export function roundToPrecision(num, precision) {
  const countOfDecimals = getCountOfDecimals(precision);
  return round(num, countOfDecimals) || precision;
}

export function readCookie(name) {
  const nameEQ = name + '=';
  const ca = document.cookie.split(';');
  for(let i=0;i < ca.length;i++) {
    let c = ca[i];
    while (c.charAt(0)==' ') c = c.substring(1,c.length);
    if (c.indexOf(nameEQ) == 0) return c.substring(nameEQ.length,c.length);
  }
  return null;
}

export function setCookie(name,value,days) {
  let expires = '';
  if (days) {
    let date = new Date();
    date.setTime(date.getTime() + (days*24*60*60*1000));
    expires = '; expires=' + date.toUTCString();
  }
  document.cookie = name + '=' + (value || '')  + expires + '; path=/';
}

/** Removes cookie 'session_token' */
export function removeSessionToken() {
  document.cookie = sessionTokenCookieName + '=; Path=/; Expires=Thu, 01 Jan 1970 00:00:01 GMT;';
}


