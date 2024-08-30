package wsconnect

import (
	"github.com/gorilla/websocket"
)

var addr = "wss://x.soulplatform.com/connection/websocket"

type Service struct {
	conn      *websocket.Conn
	token     string
	sessionID string
}

func New(wsToken string) (*Service, error) {
	c, _, err := websocket.DefaultDialer.Dial(addr, nil)
	if err != nil {
		return nil, err
	}

	conn := &Service{
		conn:  c,
		token: wsToken,
	}

	if err = conn.connect(); err != nil {
		return nil, err
	}

	return conn, nil
}

func (c *Service) Close() error {
	return c.conn.Close()
}

func (c *Service) connect() error {
	err := c.conn.WriteJSON(map[string]any{
		"params": map[string]any{
			"token": c.token,
			"name":  "js",
		},
		"id": 1,
	})

	if err != nil {
		return err
	}

	resp := connectResp{}

	if err = c.conn.ReadJSON(&resp); err != nil {
		return err
	}

	c.sessionID = resp.Result.SessionID

	return nil
}

func (c *Service) GetUsers(lastID int, filters GetUsersFilters) ([]User, error) {
	req := getUsersReq{
		Method: MethodIDGetUsers,
		Params: getUsersParams{
			Data: getUsersData{
				SessionID: c.sessionID,
				Filters:   filters,
				Ab: ab{
					SmartFeedLogic: smartFeedLogic1000kmUsersCountryRadiusNewUsers,
				},
			},
			Method: paramsMethodSmartFeedRead,
		},
		ID: lastID,
	}

	err := c.conn.WriteJSON(&req)
	if err != nil {
		return nil, err
	}

	resp := getUsersResp{}

	if err = c.conn.ReadJSON(&resp); err != nil {
		return nil, err
	}

	return resp.Result.Data.Data.Results, nil
}
