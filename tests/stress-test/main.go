package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"sync"
	"time"
)

type callStat struct {
	name     string
	duration time.Duration
	err      error
}
type user struct {
	baseURL   string
	client    *http.Client
	id        int64
	authToken string
	apiToken  string
	timings   map[string]time.Duration
	mx        *sync.RWMutex
	statsChan chan callStat
}

func newUser(id int64, stats chan callStat) *user {
	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: cookieJar}
	u := &user{
		baseURL:   "http://localhost:9512/api/v1/",
		client:    client,
		id:        id,
		mx:        &sync.RWMutex{},
		statsChan: stats,
	}
	return u
}

type allowedCodes map[int]struct{}

// errorString is a trivial implementation of error.
type skipError struct {
	s string
}

func (e *skipError) Error() string {
	return e.s
}

func (u *user) request(
	reqType string,
	endpoint string,
	jsonContent interface{},
	params map[string]string, codes allowedCodes) (string, error) {
	url := u.baseURL + endpoint
	var jsData []byte
	if jsonContent != nil {
		jsData, _ = json.Marshal(jsonContent)
	}
	req, err := http.NewRequest(reqType, url, bytes.NewBuffer(jsData))
	req.Header.Set("Authorization", "Bearer 5G8AXoY8ASASm943ZQb9iDmp8EWEVuvB")
	req.Header.Set("Content-Type", "application/json")
	u.mx.RLock()
	if u.apiToken != "" {
		req.Header.Set("Authorization", "Bearer "+u.apiToken)
	}
	u.mx.RUnlock()
	if params != nil && len(params) > 0 {
		values := req.URL.Query()
		for k, v := range params {
			values.Add(k, v)
		}
		req.URL.RawQuery = values.Encode()
	}
	resp, err := u.client.Do(req)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if _, ok := codes[resp.StatusCode]; !ok {
		return "", fmt.Errorf("Code %v", resp.StatusCode)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	return string(body), nil
}

func (u *user) signUp() {
	type SignUpStruct struct {
		UserID   int64  `json:"userId"`
		Language string `json:"language"`
		Email    string `json:"email"`
	}
	type SignUpPld struct {
		Data SignUpStruct `json:"data"`
	}
	payload := SignUpPld{
		SignUpStruct{u.id, "en", "1@2.com"},
	}
	if _, err := u.request("POST", "signup", payload, nil, allowedCodes{200: {}, 500: {}}); err != nil {
		log.Println(err)
	} else {
		// fmt.Println(resp)
	}
}

func (u *user) signIn() error {
	u.mx.Lock()
	u.authToken = ""
	u.apiToken = ""
	u.mx.Unlock()
	type SignUpStruct struct {
		UserID int64 `json:"userId"`
	}
	type SignUpPld struct {
		Data SignUpStruct `json:"data"`
	}
	payload := SignUpPld{
		SignUpStruct{u.id},
	}
	type R struct {
		Response string
	}
	if resp, err := u.request("POST", "signin", payload, nil, allowedCodes{201: {}}); err != nil {
		return err
	} else {
		r := &R{}
		json.Unmarshal([]byte(resp), r)
		u.mx.Lock()
		u.authToken = r.Response
		u.mx.Unlock()
		if r.Response == "" {
			return fmt.Errorf("signin err")
		}
	}
	return nil
}

func (u *user) login() error {
	var tk string
	u.mx.RLock()
	tk = u.authToken
	u.mx.RUnlock()
	if tk == "" {
		return &skipError{}
	}
	resp, err := u.request("GET", "login", nil, map[string]string{"auth_token": tk}, allowedCodes{200: {}})
	if err != nil {
		log.Println(resp, err)
		return err
	}
	url, _ := url.Parse(u.baseURL)
	u.mx.Lock()
	for _, c := range u.client.Jar.Cookies(url) {
		if c.Name == "session_token" {
			u.apiToken = c.Value
		}
	}
	if u.apiToken == "" {
		u.mx.Unlock()
		return fmt.Errorf("no api token")
	}
	u.mx.Unlock()
	return nil
}

func (u *user) dashboardInfo() error {
	u.mx.RLock()
	if u.apiToken == "" {
		u.mx.RUnlock()
		return &skipError{}
	}
	u.mx.RUnlock()
	_, err := u.request("GET", "dashboard-info", nil, nil, allowedCodes{200: {}})
	return err
}

func (u *user) trades() error {
	u.mx.RLock()
	if u.apiToken == "" {
		u.mx.RUnlock()
		return &skipError{}
	}
	u.mx.RUnlock()
	_, err := u.request("GET", "trades", nil, nil, allowedCodes{200: {}})
	return err
}

func (u *user) appConfig() error {
	u.mx.RLock()
	if u.apiToken == "" {
		u.mx.RUnlock()
		return &skipError{}
	}
	u.mx.RUnlock()
	_, err := u.request("GET", "app-config", nil, nil, allowedCodes{200: {}})
	return err
}

func (u *user) botToggle() error {
	u.mx.RLock()
	if u.apiToken == "" {
		u.mx.RUnlock()
		return &skipError{}
	}
	u.mx.RUnlock()
	_, err := u.request("POST", "bot-toggle", nil, nil, allowedCodes{200: {}})
	return err
}

func (u *user) runWithRandomInterval(
	stopChan chan struct{},
	name string,
	f func() error,
	interval time.Duration,
	wg *sync.WaitGroup) {
	defer wg.Done()
	timer := time.NewTicker(interval)
	for {
		select {
		case <-stopChan:
			return
		case <-timer.C:
			t := time.Now()
			err := f()
			if _, ok := err.(*skipError); !ok {
				u.statsChan <- callStat{name, time.Since(t), err}
				if err != nil {
					log.Println(err)
				}
			}

		}
	}

}
func randomDuration(seconds float64) time.Duration {
	val := 0.8*seconds + rand.Float64()*seconds*0.4
	return time.Microsecond * time.Duration(val*1000000)
}

func (u *user) run(duration time.Duration) {
	val := rand.Float64() * 120
	time.Sleep(time.Microsecond * time.Duration(val*1000000))
	// u.signUp()
	u.signIn()
	u.login()

	stop := make(chan struct{})

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go u.runWithRandomInterval(stop, "dashboard-info", func() error {
		return u.dashboardInfo()
	}, randomDuration(2), wg)

	wg.Add(1)
	go u.runWithRandomInterval(stop, "signin-login", func() error {
		if err := u.signIn(); err != nil {
			return err
		}
		if err := u.login(); err != nil {
			return err
		}
		return nil
	}, randomDuration(120), wg)

	wg.Add(1)
	go u.runWithRandomInterval(stop, "app-config", func() error {
		return u.appConfig()
	}, randomDuration(60), wg)

	wg.Add(1)
	go u.runWithRandomInterval(stop, "get-trades", func() error {
		return u.trades()
	}, randomDuration(60), wg)

	wg.Add(1)
	go u.runWithRandomInterval(stop, "bot-toggle", func() error {
		return u.botToggle()
	}, randomDuration(60), wg)

	time.AfterFunc(duration, func() {
		close(stop)
	})
	wg.Wait()
}
func main() {
	testDuration := time.Minute * 10
	// testDuration := time.Second * 20
	var userIDMin int64 = 0
	var userIDMax int64 = 5000

	log.SetFlags(log.Lshortfile | log.Ltime)
	wg := sync.WaitGroup{}
	users := []*user{}
	stats := make(chan callStat)
	var i int64
	for i = userIDMin; i < userIDMax; i++ {
		users = append(users, newUser(i, stats))
	}
	for _, u := range users {
		wg.Add(1)
		go func(usr *user) {
			usr.run(testDuration)
			// log.Printf("User %v done", usr.id)
			wg.Done()
		}(u)
	}
	type functStat struct {
		errors        int
		count         int
		totalDuration time.Duration
		maxDuration   time.Duration
		minDuration   time.Duration
	}
	statsMap := make(map[string]*functStat)

	statsFinished := make(chan struct{})
	go func() {
		defer func() { statsFinished <- struct{}{} }()
		ticker := time.NewTicker(time.Second * 10)
		for {
			select {
			case st, isOpened := <-stats:
				if !isOpened {
					return
				}
				if sf, ok := statsMap[st.name]; !ok {
					stat := &functStat{0, 0, st.duration, st.duration, st.duration}
					if st.err == nil {
						stat.count = 1
					} else {
						stat.errors = 1
					}
					statsMap[st.name] = stat
				} else {
					if st.err == nil {
						sf.count++
						sf.totalDuration += st.duration
						if sf.maxDuration < st.duration {
							sf.maxDuration = st.duration
						}
						if sf.minDuration > st.duration {
							sf.minDuration = st.duration
						}
					} else {
						sf.errors++
					}
				}
			case <-ticker.C:
				fmt.Println("#####")
				for k, v := range statsMap {
					fmt.Printf("%v called %v times with %v errors\n", k, v.count+v.errors, v.errors)
					fmt.Printf("    avg: %v ms\n", (v.totalDuration.Seconds()*1000.0)/float64(v.count))
					fmt.Printf("    min: %v ms\n", v.minDuration.Seconds()*1000.0)
					fmt.Printf("    max: %v ms\n", v.maxDuration.Seconds()*1000.0)
					fmt.Println()
				}

			}
		}

	}()
	wg.Wait()
	close(stats)
	<-statsFinished

}
