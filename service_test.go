package gopatreon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPatreon_AuthenticateUser(t *testing.T) {
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
			name: "happy-patron",
			input: in{
				fetchUser: func() (*User, error) {
					return &User{
						ID: "123456",
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
