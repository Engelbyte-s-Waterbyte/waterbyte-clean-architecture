package db

import (
	"errors"
	"fmt"

	"github.com/Engelbyte-s-Waterbyte/waterbyte-clean-architecture/models"
	"gorm.io/gorm"
)

func SelectUserByThirdPartyID(user *models.User) (found bool, err error) {
	var result models.User
	var ctx *gorm.DB
	if user.GoogleID != nil {
		ctx = db.First(&result, "google_id = ?", user.GoogleID)
	}
	if user.AppleID != nil {
		ctx = db.First(&result, "apple_id = ?", user.AppleID)
	}
	if ctx == nil {
		return false, errors.New("No third party ID provided")
	}
	if ctx.Error != nil {
		if errors.Is(ctx.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, ctx.Error
	}
	*user = result
	return true, nil
}

func InsertUser(user *models.User) error {
	result := db.Create(user)
	return result.Error
}

func SelectNextUsername(username string) (nextUsername string, err error) {
	var selectNextUsername func(username string, run int) (nextUsername string, err error)
	selectNextUsername = func(username string, run int) (nextUsername string, err error) {
		var queryUsername string
		if run == -1 {
			queryUsername = username
		} else {
			queryUsername = fmt.Sprint(username, run)
		}
		result := db.First(&models.User{}, "username = ?", queryUsername)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				return username, nil
			} else {
				return username, result.Error
			}
		}
		return selectNextUsername(username, run+1)
	}
	return selectNextUsername(username, -1)
}

func SelectUserByID(user *models.User) (found bool, err error) {
	var result models.User
	ctx := db.First(&result, user.ID)
	if ctx.Error != nil {
		if errors.Is(ctx.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	*user = result
	return true, nil
}
