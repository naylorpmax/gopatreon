package gopatreon

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
	mxpv "gopkg.in/mxpv/patreon-go.v1"
)

type (
	Client interface {
		FetchUser() (*User, error)
		FetchPledges(string) ([]*Pledge, error)
	}

	standardClient struct {
		client mxpvClientWrapper
	}

	mxpvClientWrapper interface {
		FetchUser() (*mxpv.UserResponse, error)
		FetchPledges(string) (*mxpv.PledgeResponse, error)
	}

	mxpvClient struct {
		client *mxpv.Client
	}

	User struct {
		ID string
		Attributes
	}

	Attributes struct {
		FirstName   string
		LastName    string
		IsSuspended bool
		IsDeleted   bool
		IsNuked     bool
	}

	Pledge struct {
		AmountCents    int
		PatronPaysFees bool
		IsPaused       *bool
	}
)

func New(ctx context.Context, code string, oauth2Config *oauth2.Config) (Client, error) {
	tok, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("unable to create Patreon client: %v", err.Error())
	}

	client := oauth2Config.Client(ctx, tok)

	return &standardClient{
		client: &mxpvClient{mxpv.NewClient(client)},
	}, nil
}

func (c *standardClient) FetchUser() (*User, error) {
	user, err := c.client.FetchUser()
	if err != nil {
		return nil, err
	}
	return &User{
		ID: user.Data.ID,
		Attributes: Attributes{
			FirstName:   user.Data.Attributes.FirstName,
			LastName:    user.Data.Attributes.LastName,
			IsSuspended: user.Data.Attributes.IsSuspended,
			IsDeleted:   user.Data.Attributes.IsDeleted,
			IsNuked:     user.Data.Attributes.IsNuked,
		},
	}, nil
}

func (c *standardClient) FetchPledges(campaignID string) ([]*Pledge, error) {
	pledgesResp, err := c.client.FetchPledges(campaignID)
	if err != nil {
		return nil, err
	}

	pledges := make([]*Pledge, 0)
	for _, pledgeResp := range pledgesResp.Data {
		pledge := &Pledge{
			AmountCents:    pledgeResp.Attributes.AmountCents,
			PatronPaysFees: pledgeResp.Attributes.PatronPaysFees,
		}
		if pledgeResp.Attributes.IsPaused != nil {
			pledge.IsPaused = pledgeResp.Attributes.IsPaused
		}
		pledges = append(pledges, pledge)
	}
	return pledges, nil
}

func (m *mxpvClient) FetchPledges(campaignID string) (*mxpv.PledgeResponse, error) {
	return m.client.FetchPledges(campaignID)
}

func (m *mxpvClient) FetchUser() (*mxpv.UserResponse, error) {
	return m.client.FetchUser()
}
