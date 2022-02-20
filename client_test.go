package gopatreon

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	mxpv "gopkg.in/mxpv/patreon-go.v1"
)

func TestClient_FetchUser(t *testing.T) {
	testCases := []struct {
		name         string
		client       mxpvClientWrapper
		expectedUser *User
		expectedErr  error
	}{
		{
			name:   "happy",
			client: &mockMXPVClient{testCase: "happy"},
			expectedUser: &User{
				ID: "123456789",
				Attributes: Attributes{
					FirstName:   "max",
					LastName:    "naylor",
					IsSuspended: false,
					IsDeleted:   false,
					IsNuked:     true,
				},
			},
		},
		{
			name:        "sad",
			client:      &mockMXPVClient{testCase: "sad", expectedErr: errors.New("connection lost")},
			expectedErr: errors.New("connection lost"),
		},
	}

	for _, testCase := range testCases {
		client := standardClient{client: testCase.client}
		user, err := client.FetchUser()

		assert.Equal(t, testCase.expectedUser, user)
		assert.Equal(t, testCase.expectedErr, err)
	}
}

func TestClient_FetchPledges(t *testing.T) {
	isPaused := true

	testCases := []struct {
		name            string
		client          mxpvClientWrapper
		campaignID      string
		expectedPledges []*Pledge
		expectedErr     error
	}{
		{
			name:       "happy",
			client:     &mockMXPVClient{testCase: "happy"},
			campaignID: "888",
			expectedPledges: []*Pledge{
				{
					AmountCents:    1000,
					IsPaused:       &isPaused,
					PatronPaysFees: false,
				},
				{
					AmountCents:    2000,
					PatronPaysFees: true,
				},
			},
		},
		{
			name:        "sad",
			client:      &mockMXPVClient{testCase: "sad", expectedErr: errors.New("connection lost")},
			campaignID:  "999",
			expectedErr: errors.New("connection lost"),
		},
	}

	for _, testCase := range testCases {
		client := standardClient{client: testCase.client}
		user, err := client.FetchPledges(testCase.campaignID)

		assert.Equal(t, testCase.expectedPledges, user)
		assert.Equal(t, testCase.expectedErr, err)
	}
}

type mockMXPVClient struct {
	testCase    string
	expectedErr error
}

func (m *mockMXPVClient) FetchPledges(campaignID string) (*mxpv.PledgeResponse, error) {
	switch m.testCase {
	case "happy":
		resp := &mxpv.PledgeResponse{
			Data: []mxpv.Pledge{
				{
					Type: "paused",
					ID:   "123456789",
				},
				{
					Type: "not paused",
					ID:   "987654321",
				},
			},
		}

		isPaused := true
		resp.Data[0].Attributes.AmountCents = 1000
		resp.Data[0].Attributes.IsPaused = &isPaused
		resp.Data[0].Attributes.PatronPaysFees = false

		resp.Data[1].Attributes.AmountCents = 2000
		resp.Data[1].Attributes.PatronPaysFees = true

		return resp, nil
	case "sad":
		return nil, m.expectedErr
	}
	return nil, errors.New("unrecognized test case")
}

func (m *mockMXPVClient) FetchUser() (*mxpv.UserResponse, error) {
	switch m.testCase {
	case "happy":
		resp := &mxpv.UserResponse{
			Data: mxpv.User{
				ID:   "123456789",
				Type: "patron",
			},
		}
		resp.Data.Attributes.FirstName = "max"
		resp.Data.Attributes.LastName = "naylor"
		resp.Data.Attributes.Gender = -1
		resp.Data.Attributes.About = "hi i'm max and i create software"
		resp.Data.Attributes.IsDeleted = false
		resp.Data.Attributes.IsNuked = true
		resp.Data.Attributes.IsSuspended = false
		return resp, nil
	case "sad":
		return nil, m.expectedErr
	}
	return nil, errors.New("unrecognized test case")
}
