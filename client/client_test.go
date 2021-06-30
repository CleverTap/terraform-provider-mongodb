package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetUser(t *testing.T) {
	tesCases := []struct {
		testName     string
		userName     string
		expectErr    bool
		expectedResp *User
	}{
		{
			testName:  "user exist",
			userName:  "existUserID",
			expectErr: false,
			expectedResp: &User{
				Country:      "Country",
				EmailAddress: "existUserID",
				FirstName:    "FirstName",
				ID:           "USerID",
				LastName:     "LastName",
				Roles: []Role{
					{
						OrgID:    "ORGID",
						RoleName: "ORGROLE",
					},
				},
				TeamIds:  []interface{}{},
				Username: "existUserID",
			},
		}, {
			testName:     "user does not exist",
			userName:     "example@gmail.com",
			expectErr:    true,
			expectedResp: nil,
		},
	}

	for _, tc := range tesCases {
		t.Run(tc.testName, func(t *testing.T) {
			client := NewClient("public", "private", "orgID")

			item, err := client.GetUser(tc.userName)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResp, item)
		})
	}
}

func TestClient_AddNewUser(t *testing.T) {
	testCases := []struct {
		testName     string
		newUser      *NewUser
		expectErr    bool
		expectedResp *NewReturnUser
	}{
		{
			testName: "new user created",
			newUser: &NewUser{
				Roles:    []string{"ORG_MEMBER"},
				Username: "newUserID",
			},
			expectErr: false,
			expectedResp: &NewReturnUser{
				InviterUsername: "InviterName",
				OrgID:           "ORGID",
				OrgName:         "ORGNAME",
				TeamIds:         []interface{}{},
				Username:        "newUserID",
			},
		},
		{
			testName: "user already exist",
			newUser: &NewUser{
				Roles:    []string{"ORG_MEMBER"},
				Username: "ExistUserID",
			},
			expectErr:    true,
			expectedResp: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			client := NewClient("public", "private", "orgID")

			item, err := client.AddNewUser(tc.newUser)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResp, item)
		})
	}
}

func TestClient_UpdateUser(t *testing.T) {
	testCases := []struct {
		testName     string
		userId       string
		updatedUser  []string
		expectErr    bool
		expectedResp *User
	}{
		{
			testName:    "user exist",
			userId:      "UserID",
			updatedUser: []string{"ORG_MEMBER"},
			expectErr:   false,
			expectedResp: &User{
				Country:      "Country",
				EmailAddress: "ExistUserEmail",
				FirstName:    "FirstName",
				ID:           "UserID",
				LastName:     "LastName",
				Roles: []Role{
					{
						OrgID:    "ORGID",
						RoleName: "ORG_MEMBER",
					},
				},
				TeamIds:  []interface{}{},
				Username: "ExistUserEmail",
			},
		},
		{
			testName:     "user not exist",
			userId:       "00000",
			updatedUser:  []string{"ORG_MEMBER"},
			expectErr:    true,
			expectedResp: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			client := NewClient("public", "private", "OrgID")

			item, err := client.UpdateUser(tc.updatedUser, tc.userId)
			if tc.expectErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedResp, item)
		})
	}
}
