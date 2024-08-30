package pureapi

import (
	"errors"

	"github.com/0x16f/pureapi/src/usecase/webapi"
	"github.com/0x16f/pureapi/src/usecase/wsconnect"
)

var ErrTokenExpired = webapi.ErrTokenExpired

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
	if err := a.webapi.Like(userId); err != nil {
		if errors.Is(err, webapi.ErrTokenExpired) {
			return ErrTokenExpired
		}

		return err
	}

	return nil
}

func (a *PureAPI) GetUsers(lastID int, filters GetUsersFilters) (Users, error) {
	return a.wsapi.GetUsers(lastID, wsconnect.GetUsersFilters(filters))
}

func (a *PureAPI) SetLocation(lat, lng float64) error {
	if err := a.webapi.SetLocation(lat, lng); err != nil {
		if errors.Is(err, webapi.ErrTokenExpired) {
			return ErrTokenExpired
		}

		return err
	}

	return nil
}
