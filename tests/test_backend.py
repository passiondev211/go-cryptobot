import requests
import json
from pprint import pprint
from time import sleep
s = requests.Session()
s.headers.update({
        'Content-Type': 'application/json',
        'Authorization': 'Bearer 5G8AXoY8ASASm943ZQb9iDmp8EWEVuvB'
    })

def makeRequst(reqType, endpoint, dataObj=None, params=None, auth_token=None):
    url = "http://localhost:9512/api/v1/" + endpoint
    payload = json.dumps(dataObj) if dataObj else None
    headers = {'Authorization': 'Bearer ' + auth_token} if auth_token else None
    response = s.request(reqType, url, data=payload, params=params, headers=headers)
    print(response.status_code)
    try:
        return json.loads(response.text)
    except:
        return response.text

for outerUid in range(2000, 2001):
    makeRequst("POST", "signup", dataObj={"data": {"userId": outerUid, "language": "en", "email": "1@2.com"}})
    token = makeRequst("POST", "signin", {"data": {"userId": outerUid}})['response']
    print(token)
    # makeRequst("GET", "login", params={'auth_token': token})
    # user_token = s.cookies.get_dict()['session_token']
    # makeRequst("GET", "dashboard-info", auth_token=user_token)
    # print(makeRequst("POST", "bot-toggle", auth_token=user_token))
    # sleep(1)
    # print(makeRequst("POST", "bot-toggle", auth_token=user_token))

