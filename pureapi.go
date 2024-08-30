package pureapi

import (
	"github.com/0x16f/pureapi/src/usecase/webapi"
	"github.com/0x16f/pureapi/src/usecase/wsconnect"
)

type GetUsersFilters wsconnect.GetUsersFilters

type Users []wsconnect.User

type PureAPI struct {
	wsapi  *wsconnect.Service
	webapi *webapi.Service
}

func New(wsToken, apiToken string) (*PureAPI, error) {
	conn, err := wsconnect.New(wsToken)
	if err != nil {
		return nil, err
	}

	api := webapi.New(apiToken)

	return &PureAPI{
		wsapi:  conn,
		webapi: api,
	}, nil
}

func (a *PureAPI) Close() error {
	if err := a.wsapi.Close(); err != nil {
		return err
	}

	return nil
}

func (a *PureAPI) Like(userId string) error {
	return a.webapi.Like(userId)
}

func (a *PureAPI) GetUsers(lastID int, filters GetUsersFilters) (Users, error) {
	return a.wsapi.GetUsers(lastID, wsconnect.GetUsersFilters(filters))
}
