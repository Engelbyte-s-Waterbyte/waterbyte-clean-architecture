package logic

import (
	"errors"
	"strings"

	"github.com/Engelbyte-s-Waterbyte/waterbyte-clean-architecture/models"
	jwt "github.com/dgrijalva/jwt-go"
)

type FetchUserDataFunc func(token string) (models.User, error)

// SelectNextUsernameFunc: A function that takes a username and then queries the database to esnure it is unique.
// If it isn't unique, a number is added at the end to make it unique.
type SelectNextUsernameFunc func(username string) (nextUsername string, err error)

var (
	jwtSecret = []byte("jwt_secret")
)

func SignIn(googleToken *string, appleToken *string, selectUserByThirdPartyID func(*models.User) (found bool, err error), fetchUserDataFromGoogle, fetchUserDataFromApple FetchUserDataFunc, selectNextUsername SelectNextUsernameFunc, insertUser func(*models.User) error) (token *string, err error) {
	generateTokenForUser := func(user *models.User) string {
		claims := jwt.MapClaims{}
		claims["id"] = user.ID
		at := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
		token, _ := at.SignedString(jwtSecret)
		return token
	}

	var user *models.User
	if appleToken != nil {
		rUser, err := fetchUserDataFromApple(*appleToken)
		if err != nil {
			return nil, err
		}
		user = &rUser
	}
	if googleToken != nil {
		rUser, err := fetchUserDataFromGoogle(*googleToken)
		if err != nil {
			return nil, err
		}
		user = &rUser
	}
	if user == nil {
		return nil, errors.New("No third party token provided")
	}
	if user.Username == "" {
		user.Username = strings.Split(user.Email, "@")[0]
	}

	found, err := selectUserByThirdPartyID(user)
	if err != nil {
		return nil, err
	}
	if found {
		token := generateTokenForUser(user)
		return &token, err
	}

	nextUsername, err := selectNextUsername(user.Username)
	if err != nil {
		return nil, err
	}
	user.Username = nextUsername
	if err := insertUser(user); err != nil {
		return nil, err
	}
	generatedToken := generateTokenForUser(user)
	return &generatedToken, nil
}

func AuthenticatedUser(token string, selectUserByID func(*models.User) (found bool, err error)) (*models.User, error) {
	decodeTokenToID := func(tokenString string) (*int, error) {
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Unexpected Token signing method")
			}
			return jwtSecret, nil
		})
		if err != nil {
			return nil, err
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok && !token.Valid {
			return nil, errors.New("Token not OK or not Valid")
		}
		uId := int(claims["id"].(float64))
		return &uId, nil
	}

	user := &models.User{}
	userID, err := decodeTokenToID(token)
	if err != nil {
		return nil, err
	}
	user.ID = uint(*userID)

	found, err := selectUserByID(user)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, errors.New("User does not exist")
	}

	return user, nil
}
