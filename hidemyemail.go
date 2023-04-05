package hidemyemail

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type GenerateResponse struct {
	Success   bool   `json:"success"`
	Timestamp int    `json:"timestamp"`
	Result    Result `json:"result"`
}

type Result struct {
	Hme string `json:"hme"`
}

type ReserveResponse struct {
	Success   bool          `json:"success,omitempty"`
	Timestamp int           `json:"timestamp,omitempty"`
	Result    ReserveResult `json:"result,omitempty"`
	Error     ReserveError  `json:"error,omitempty"`
}

type Hme struct {
	Origin          string `json:"origin,omitempty"`
	AnonymousID     string `json:"anonymousId,omitempty"`
	Domain          string `json:"domain,omitempty"`
	Hme             string `json:"hme,omitempty"`
	Label           string `json:"label,omitempty"`
	Note            string `json:"note,omitempty"`
	CreateTimestamp int    `json:"createTimestamp,omitempty"`
	IsActive        bool   `json:"isActive,omitempty"`
	RecipientMailID string `json:"recipientMailId,omitempty"`
}

type ReserveResult struct {
	Hme Hme `json:"hme,omitempty"`
}

type ReserveError struct {
	ErrorCode    string `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	RetryAfter   int    `json:"retryAfter,omitempty"`
}

type HideMyEmail struct {
	Label   string
	Cookies string
}

func (hme *HideMyEmail) Generate() (*GenerateResponse, error) {
	req, err := http.NewRequest(http.MethodPost, "https://p68-maildomainws.icloud.com/v1/hme/generate", nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %v", err)
	}

	q := req.URL.Query()
	q.Add("clientBuildNumber", "2206Hotfix11")
	q.Add("clientMasteringNumber", "2206Hotfix11")
	q.Add("clientId", "")
	q.Add("dsid", "")
	req.URL.RawQuery = q.Encode()

	req.Header = http.Header{
		"Connection":      {"keep-alive"},
		"Pragma":          {"no-cache"},
		"Cache-Control":   {"no-cache"},
		"User-Agent":      {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.109 Safari/537.36"},
		"Content-Type":    {"text/plain"},
		"Accept":          {"*/*"},
		"Sec-GPC":         {"1"},
		"Origin":          {"https://www.icloud.com"},
		"Sec-Fetch-Site":  {"same-site"},
		"Sec-Fetch-Mode":  {"cors"},
		"Sec-Fetch-Dest":  {"empty"},
		"Referer":         {"https://www.icloud.com/"},
		"Accept-Language": {"en-US,en-GB;q=0.9,en;q=0.8,cs;q=0.7"},
		"Cookie":          {hme.Cookies},
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network err: %v", err)
	}
	defer res.Body.Close()

	var resp GenerateResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		b, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("unable to read res body: %s", err)
		}
		return nil, fmt.Errorf("unable to decode body (%s) into json: %v", string(b), err)
	}

	return &resp, nil
}

type ReservePayload struct {
	HME   string `json:"hme"`
	Label string `json:"label"`
	Note  string `json:"note"`
}

var ErrTimeLimit = errors.New("you have reached the limit of addresses you can create right now")

func (hme *HideMyEmail) Reserve(email string) (*ReserveResponse, error) {
	payload := ReservePayload{
		HME:   email,
		Label: hme.Label,
		Note:  "hey",
	}

	payloadBuf := new(bytes.Buffer)
	if err := json.NewEncoder(payloadBuf).Encode(payload); err != nil {
		return nil, fmt.Errorf("unable to encode JSON: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://p68-maildomainws.icloud.com/v1/hme/reserve", payloadBuf)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %v", err)
	}

	q := req.URL.Query()
	q.Add("clientBuildNumber", "2206Hotfix11")
	q.Add("clientMasteringNumber", "2206Hotfix11")
	q.Add("clientId", "")
	q.Add("dsid", "")
	req.URL.RawQuery = q.Encode()

	req.Header = http.Header{
		"Connection":      {"keep-alive"},
		"Pragma":          {"no-cache"},
		"Cache-Control":   {"no-cache"},
		"User-Agent":      {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.109 Safari/537.36"},
		"Content-Type":    {"text/plain"},
		"Accept":          {"*/*"},
		"Sec-GPC":         {"1"},
		"Origin":          {"https://www.icloud.com"},
		"Sec-Fetch-Site":  {"same-site"},
		"Sec-Fetch-Mode":  {"cors"},
		"Sec-Fetch-Dest":  {"empty"},
		"Referer":         {"https://www.icloud.com/"},
		"Accept-Language": {"en-US,en-GB;q=0.9,en;q=0.8,cs;q=0.7"},
		"Cookie":          {hme.Cookies},
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network err: %v", err)
	}
	defer res.Body.Close()

	var resp ReserveResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, fmt.Errorf("unable to decode body into json: %v", err)
	}

	if !resp.Success {
		return nil, ErrTimeLimit
	}

	return &resp, nil
}

type ListEmailsResp struct {
	Success   bool             `json:"success,omitempty"`
	Timestamp int              `json:"timestamp,omitempty"`
	Result    ListEmailsResult `json:"result,omitempty"`
}

type HmeEmails struct {
	Origin          string `json:"origin,omitempty"`
	AnonymousID     string `json:"anonymousId,omitempty"`
	Domain          string `json:"domain,omitempty"`
	ForwardToEmail  string `json:"forwardToEmail,omitempty"`
	Hme             string `json:"hme,omitempty"`
	Label           string `json:"label,omitempty"`
	Note            string `json:"note,omitempty"`
	CreateTimestamp int64  `json:"createTimestamp,omitempty"`
	IsActive        bool   `json:"isActive,omitempty"`
	RecipientMailID string `json:"recipientMailId,omitempty"`
}

type ListEmailsResult struct {
	ForwardToEmails   []interface{} `json:"forwardToEmails,omitempty"`
	HmeEmails         []HmeEmails   `json:"hmeEmails,omitempty"`
	SelectedForwardTo string        `json:"selectedForwardTo,omitempty"`
}

func (hme *HideMyEmail) List() ([]string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://p68-maildomainws.icloud.com/v1/hme/list", nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create request: %v", err)
	}

	q := req.URL.Query()
	q.Add("clientBuildNumber", "2206Hotfix11")
	q.Add("clientMasteringNumber", "2206Hotfix11")
	q.Add("clientId", "")
	q.Add("dsid", "")
	req.URL.RawQuery = q.Encode()

	req.Header = http.Header{
		"Connection":      {"keep-alive"},
		"Pragma":          {"no-cache"},
		"Cache-Control":   {"no-cache"},
		"User-Agent":      {"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.109 Safari/537.36"},
		"Content-Type":    {"text/plain"},
		"Accept":          {"*/*"},
		"Sec-GPC":         {"1"},
		"Origin":          {"https://www.icloud.com"},
		"Sec-Fetch-Site":  {"same-site"},
		"Sec-Fetch-Mode":  {"cors"},
		"Sec-Fetch-Dest":  {"empty"},
		"Referer":         {"https://www.icloud.com/"},
		"Accept-Language": {"en-US,en-GB;q=0.9,en;q=0.8,cs;q=0.7"},
		"Cookie":          {hme.Cookies},
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("network err: %v", err)
	}
	defer res.Body.Close()

	var listEmailsResp ListEmailsResp
	if err := json.NewDecoder(res.Body).Decode(&listEmailsResp); err != nil {
		return nil, fmt.Errorf("unable to decode http response into json: %v", err)
	}

	var emails []string

	for _, item := range listEmailsResp.Result.HmeEmails {
		emails = append(emails, item.Hme)
	}
	return emails, nil
}
