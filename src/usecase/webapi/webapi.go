package webapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

var uri = "https://pure-api3.soulplatform.com"

var (
	ErrTokenExpired         = errors.New("token expired")
	ErrFailedToRefreshToken = errors.New("failed to refresh token")
)

type Service struct {
	uri          string
	refreshToken string
	accesToken   string
}

func New(refreshToken string, accessToken string) *Service {
	return &Service{
		uri:          uri,
		refreshToken: refreshToken,
		accesToken:   accessToken,
	}
}

func (s Service) Like(userId string) error {
	payloadBytes, err := json.Marshal(map[string]any{
		"value":           "liked",
		"createdTime":     time.Now().Unix(),
		"screen":          "feed",
		"consumptionId":   nil,
		"match_condition": nil,
	})

	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/users/%s/reactions/sent/likes", s.uri, userId),
		bytes.NewReader(payloadBytes))

	if err != nil {
		return err
	}

	s.addHeaders(req)

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusUnauthorized {
		return ErrTokenExpired
	}

	if res.StatusCode != http.StatusCreated {
		return errors.New(fmt.Sprintf("%d: %s, failed to like", res.StatusCode, body))
	}

	return nil
}

func (s Service) SetLocation(lat, lng float64) error {
	payloadBytes, err := json.Marshal(map[string]any{
		"parameters": map[string]any{
			"filterable": map[string]any{
				"location": map[string]float64{
					"lat": lat,
					"lng": lng,
				},
			},
		},
	})

	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPatch, fmt.Sprintf("%s/me", s.uri), bytes.NewReader(payloadBytes))
	if err != nil {
		return err
	}

	s.addHeaders(req)

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		return ErrTokenExpired
	}

	return nil
}

func (s *Service) RefreshToken() (*RefreshTokenResp, error) {
	payloadBytes, err := json.Marshal(map[string]any{
		"api_key":       "48d95c045bfbf6544448fe07744e558b",
		"refresh_token": s.refreshToken,
	})

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/auth-companion-service/v1/refresh-token/", s.uri),
		bytes.NewReader(payloadBytes))

	if err != nil {
		return nil, err
	}

	s.addHeaders(req)

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, ErrFailedToRefreshToken
	}

	resp := RefreshTokenResp{}

	if err = json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return nil, err
	}

	s.accesToken = resp.AccessToken
	s.refreshToken = resp.RefreshToken

	return &resp, nil
}

func (s Service) GetWebsocketToken() (string, error) {
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/centrifugo/get_token/", s.uri), nil)
	if err != nil {
		return "", err
	}

	s.addHeaders(req)

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		return "", ErrTokenExpired
	}

	resp := GetWebsocketTokenResp{}

	if err = json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return "", err
	}

	return resp.Token, nil
}

func (s Service) addHeaders(req *http.Request) {
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/117.0")
	req.Header.Add("x-js-user-agent", "PureFTP/4.2.2 (JS 1.0; OS X 10.15 MacIntel Firefox 117.0; en-US) SoulSDK/0.26.0 (JS)")
	req.Header.Add("Authorization", s.accesToken)
	req.Header.Add("Content-Type", "application/json")
}
