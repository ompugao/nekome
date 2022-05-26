package api

import (
	"context"

	"github.com/arrow2nd/nekome/oauth"
	"github.com/g8rswimmer/go-twitter/v2"
)

// AuthUserLookup 現在のユーザの情報を取得
func (a *API) AuthUserLookup() (*twitter.UserObj, error) {
	return a.AuthUserLookupFromToken(a.token)
}

// AuthUserLookupFromToken トークンに紐づいたユーザの情報を取得
func (a *API) AuthUserLookupFromToken(token *oauth.Token) (*twitter.UserObj, error) {
	client, err := a.newClient(token)
	if err != nil {
		return nil, err
	}

	opts := twitter.UserLookupOpts{}
	userResponse, err := client.AuthUserLookup(context.Background(), opts)
	if err != nil {
		return nil, err
	}

	return userResponse.Raw.Users[0], nil
}
