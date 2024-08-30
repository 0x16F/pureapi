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
	ErrTokenExpired = errors.New("token expired")
)

type Service struct {
	uri   string
	token string
}

func New(apiToken string) *Service {
	return &Service{
		uri:   uri,
		token: apiToken,
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

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/users/%s/reactions/sent/likes", s.uri, userId), bytes.NewReader(payloadBytes))
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

func (s Service) addHeaders(req *http.Request) {
	req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:109.0) Gecko/20100101 Firefox/117.0")
	req.Header.Add("x-js-user-agent", "PureFTP/4.2.2 (JS 1.0; OS X 10.15 MacIntel Firefox 117.0; en-US) SoulSDK/0.26.0 (JS)")
	req.Header.Add("Authorization", s.token)
	req.Header.Add("Content-Type", "application/json")
}
