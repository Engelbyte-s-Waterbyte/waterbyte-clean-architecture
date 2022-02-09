package logic

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Engelbyte-s-Waterbyte/waterbyte-clean-architecture/models"
)

func TestAuthenticatedUser(t *testing.T) {
	userWithID1 := (func() *models.User {
		user := &models.User{}
		user.ID = 1
		user.Name = "Lambo Tim"
		return user
	})()

	selectUserByID := func(user *models.User) (found bool, err error) {
		if user.ID == 1 {
			*user = *userWithID1
			return true, nil
		}
		return false, errors.New("User does not exist")
	}
	type args struct {
		token          string
		selectUserByID func(*models.User) (found bool, err error)
	}
	tests := []struct {
		name    string
		args    args
		want    *models.User
		wantErr bool
	}{
		{name: "Empty Token", args: args{token: "", selectUserByID: selectUserByID}, wantErr: true},
		{name: "Invalid Token", args: args{token: "invalid_token", selectUserByID: selectUserByID}, wantErr: true},
		{name: "Valid Token (ID=1)", args: args{token: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpZCI6MX0.yITqM43KV6ssuysiSg1byCq3jU088nL7gs9AyX0NPW8ZHButL7ynGaXq4wI0Lic7bLuFfDr0wn3Rf1Lsr9WHxw", selectUserByID: selectUserByID}, want: userWithID1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AuthenticatedUser(tt.args.token, tt.args.selectUserByID)
			if (err != nil) != tt.wantErr {
				t.Errorf("AuthenticatedUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthenticatedUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func pointStr(a string) *string {
	return &a
}

func TestSignIn(t *testing.T) {
	var users []models.User
	selectUserByThirdPartyID := func(inputUser *models.User) (bool, error) {
		for _, user := range users {
			if inputUser.GoogleID != nil && inputUser.GoogleID == user.GoogleID {
				*inputUser = user
				return true, nil
			}
			if inputUser.AppleID != nil && inputUser.AppleID == user.AppleID {
				*inputUser = user
				return true, nil
			}
		}
		return false, nil
	}

	var googleUsers []models.User
	googleUsers = append(googleUsers, models.User{Name: "zeppes hawara", Email: "zeppelin@gmail.com", GoogleID: pointStr("hawara_zepp")})
	fetchUserDataFromGoogle := func(googleId string) (models.User, error) {
		for _, googleUser := range googleUsers {
			if *googleUser.GoogleID == googleId {
				return googleUser, nil
			}
		}
		return models.User{}, errors.New("")
	}

	var appleUsers []models.User
	appleUsers = append(appleUsers, models.User{Name: "iphone zeppes", Email: "hawara@icloud.com", AppleID: pointStr("zeppi_hawara")})
	fetchUserDataFromApple := func(appleId string) (models.User, error) {
		for _, appleUser := range appleUsers {
			if *appleUser.AppleID == appleId {
				return appleUser, nil
			}
		}
		return models.User{}, errors.New("")
	}

	selectNextUsername := func(username string) (string, error) {
		return username, nil
	}
	insertUser := func(user *models.User) error {
		user.ID = uint(len(users))
		users = append(users, *user)
		return nil
	}
	type args struct {
		googleToken              *string
		appleToken               *string
		selectUserByThirdPartyID func(*models.User) (found bool, err error)
		fetchUserDataFromGoogle  FetchUserDataFunc
		fetchUserDataFromApple   FetchUserDataFunc
		selectNextUsername       SelectNextUsernameFunc
		insertUser               func(*models.User) error
	}
	tests := []struct {
		name      string
		args      args
		wantToken string
		wantErr   bool
	}{
		{name: "Zeppelin Google SignIn 1", args: args{googleToken: pointStr("hawara_zepp"), selectUserByThirdPartyID: selectUserByThirdPartyID, fetchUserDataFromGoogle: fetchUserDataFromGoogle, fetchUserDataFromApple: fetchUserDataFromApple, selectNextUsername: selectNextUsername, insertUser: insertUser}, wantToken: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpZCI6MH0.K2DH0owl4uSdueZWVFFfTgEHXh91UqxfL1kWE5v6qN_LJ_aG7p94r3tvfHq7janh0aFlg6pC0g8qTaQc1wxH3A"},
		{name: "Zeppelin Google SignIn 2", args: args{googleToken: pointStr("hawara_zepp"), selectUserByThirdPartyID: selectUserByThirdPartyID, fetchUserDataFromGoogle: fetchUserDataFromGoogle, fetchUserDataFromApple: fetchUserDataFromApple, selectNextUsername: selectNextUsername, insertUser: insertUser}, wantToken: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpZCI6MH0.K2DH0owl4uSdueZWVFFfTgEHXh91UqxfL1kWE5v6qN_LJ_aG7p94r3tvfHq7janh0aFlg6pC0g8qTaQc1wxH3A"},
		{name: "Zeppelin Apple SignIn", args: args{appleToken: pointStr("zeppi_hawara"), selectUserByThirdPartyID: selectUserByThirdPartyID, fetchUserDataFromGoogle: fetchUserDataFromGoogle, fetchUserDataFromApple: fetchUserDataFromApple, selectNextUsername: selectNextUsername, insertUser: insertUser}, wantToken: "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpZCI6MX0.yITqM43KV6ssuysiSg1byCq3jU088nL7gs9AyX0NPW8ZHButL7ynGaXq4wI0Lic7bLuFfDr0wn3Rf1Lsr9WHxw"},
		{name: "Error", args: args{selectUserByThirdPartyID: selectUserByThirdPartyID, fetchUserDataFromGoogle: fetchUserDataFromGoogle, fetchUserDataFromApple: fetchUserDataFromApple, selectNextUsername: selectNextUsername, insertUser: insertUser}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := SignIn(tt.args.googleToken, tt.args.appleToken, tt.args.selectUserByThirdPartyID, tt.args.fetchUserDataFromGoogle, tt.args.fetchUserDataFromApple, tt.args.selectNextUsername, tt.args.insertUser)
			if (err != nil) != tt.wantErr {
				t.Errorf("SignIn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if *gotToken != tt.wantToken {
				t.Errorf("SignIn() = %v, want %v", *gotToken, tt.wantToken)
			}
		})
	}
}
