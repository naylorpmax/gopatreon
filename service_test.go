package gopatreon

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPatreon_AuthenticateUser(t *testing.T) {
	isPaused := true
	isPausedAlt := false

	type (
		in struct {
			fetchUser    func() (*User, error)
			fetchPledges func(string) ([]*Pledge, error)
		}
		exp struct {
			fullName string
			err      error
		}
	)
	cases := []struct {
		name     string
		input    in
		expected exp
	}{
		{
			name: "happy-creator",
			input: in{
				fetchUser: func() (*User, error) {
					return &User{
						ID: CreatorUserID,
						Attributes: Attributes{
							FirstName: "max",
							LastName:  "naylor",
						},
					}, nil
				},
			},
			expected: exp{
				fullName: "max naylor",
			},
		},
		{
			name: "happy-patron",
			input: in{
				fetchUser: func() (*User, error) {
					return &User{
						ID: "not-creator-id",
						Attributes: Attributes{
							FirstName:   "max",
							LastName:    "naylor",
							IsDeleted:   false,
							IsNuked:     false,
							IsSuspended: false,
						},
					}, nil
				},
				fetchPledges: func(campaignID string) ([]*Pledge, error) {
					isPaused := false
					return []*Pledge{
						{
							AmountCents:    500,
							IsPaused:       &isPaused,
							PatronPaysFees: true,
						},
					}, nil
				},
			},
			expected: exp{
				fullName: "max naylor",
				err:      nil,
			},
		},
		{
			name: "sad-patron-unable-to-fetch-user",
			input: in{
				fetchUser: func() (*User, error) {
					return nil, errors.New("oh no!")
				},
			},
			expected: exp{
				err: errors.New("unable to fetch user: oh no!"),
			},
		},
		{
			name: "sad-patron-unable-to-fetch-pledges",
			input: in{
				fetchUser: func() (*User, error) {
					return &User{
						ID: "not-creator-id",
					}, nil
				},
				fetchPledges: func(campaignID string) ([]*Pledge, error) {
					return nil, errors.New("oh no!")
				},
			},
			expected: exp{
				err: errors.New("unable to fetch user's pledges: oh no!"),
			},
		},
		{
			name: "sad-patron-not-enough-dough",
			input: in{
				fetchUser: func() (*User, error) {
					return &User{
						ID: "not-creator-id",
						Attributes: Attributes{
							IsSuspended: false,
							IsDeleted:   false,
							IsNuked:     false,
						},
					}, nil
				},
				fetchPledges: func(campaignID string) ([]*Pledge, error) {
					return []*Pledge{
						{
							AmountCents:    MinUserAmountCents/2 - 1,
							PatronPaysFees: true,
						},
						{
							AmountCents:    MinUserAmountCents/2 - 1,
							PatronPaysFees: true,
						},
					}, nil
				},
			},
			expected: exp{
				err: errors.New("patron level not high enough to access content"),
			},
		},
		{
			name: "sad-patron-suspended",
			input: in{
				fetchUser: func() (*User, error) {
					return &User{
						ID: "not-creator-id",
						Attributes: Attributes{
							IsSuspended: true,
							IsDeleted:   false,
							IsNuked:     false,
						},
					}, nil
				},
				fetchPledges: func(campaignID string) ([]*Pledge, error) {
					return []*Pledge{
						{
							AmountCents:    MinUserAmountCents,
							PatronPaysFees: true,
						},
						{
							AmountCents:    MinUserAmountCents,
							PatronPaysFees: true,
						},
					}, nil
				},
			},
			expected: exp{
				err: errors.New("user is not in good standing with this campaign: user is suspended"),
			},
		},
		{
			name: "sad-patron-deleted",
			input: in{
				fetchUser: func() (*User, error) {
					return &User{
						ID: "not-creator-id",
						Attributes: Attributes{
							IsSuspended: false,
							IsDeleted:   true,
							IsNuked:     false,
						},
					}, nil
				},
				fetchPledges: func(campaignID string) ([]*Pledge, error) {
					return []*Pledge{
						{
							AmountCents:    MinUserAmountCents,
							PatronPaysFees: true,
						},
						{
							AmountCents:    MinUserAmountCents,
							PatronPaysFees: true,
						},
					}, nil
				},
			},
			expected: exp{
				err: errors.New("user is not in good standing with this campaign: user is deleted"),
			},
		},
		{
			name: "sad-patron-nuked",
			input: in{
				fetchUser: func() (*User, error) {
					return &User{
						ID: "not-creator-id",
						Attributes: Attributes{
							IsSuspended: false,
							IsDeleted:   false,
							IsNuked:     true,
						},
					}, nil
				},
				fetchPledges: func(campaignID string) ([]*Pledge, error) {
					return []*Pledge{
						{
							AmountCents:    MinUserAmountCents,
							PatronPaysFees: true,
						},
						{
							AmountCents:    MinUserAmountCents,
							PatronPaysFees: true,
						},
					}, nil
				},
			},
			expected: exp{
				err: errors.New("user is not in good standing with this campaign: user is nuked"),
			},
		},
		{
			name: "sad-patron-paused",
			input: in{
				fetchUser: func() (*User, error) {
					return &User{
						ID: "not-creator-id",
						Attributes: Attributes{
							IsSuspended: false,
							IsDeleted:   false,
							IsNuked:     false,
						},
					}, nil
				},
				fetchPledges: func(campaignID string) ([]*Pledge, error) {
					return []*Pledge{
						{
							AmountCents:    MinUserAmountCents,
							PatronPaysFees: true,
							IsPaused:       &isPausedAlt,
						},
						{
							AmountCents:    MinUserAmountCents,
							PatronPaysFees: true,
							IsPaused:       &isPaused,
						},
					}, nil
				},
			},
			expected: exp{
				err: errors.New("user is not in good standing with this campaign: user is paused"),
			},
		},
	}

	for _, test := range cases {
		client := mockClient{
			fetchUserFn:    test.input.fetchUser,
			fetchPledgesFn: test.input.fetchPledges,
		}
		patreonService, _ := NewService(&client)

		actualFullName, actualErr := patreonService.AuthenticateUser()

		assert.Equal(t, test.expected.fullName, actualFullName)
		assert.Equal(t, test.expected.err, actualErr)
	}
}

type mockClient struct {
	fetchUserFn    func() (*User, error)
	fetchPledgesFn func(string) ([]*Pledge, error)
}

func (m *mockClient) FetchUser() (*User, error) {
	return m.fetchUserFn()
}

func (m *mockClient) FetchPledges(campaignID string) ([]*Pledge, error) {
	return m.fetchPledgesFn(campaignID)
}
