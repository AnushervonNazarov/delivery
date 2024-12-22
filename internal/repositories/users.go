package repository

import (
	"blogging_platform/internal/models"
	"blogging_platform/logger"
	"blogging_platform/pkg/db"
)

func GetAllUsers() (users []models.User, err error) {
	err = db.GetDBConn().Find(&users).Error
	if err != nil {
		logger.Error.Println("[repository.GetAllUsers] cannot get all users. Error is:", err.Error())
		return nil, translateError(err)
	}

	return users, nil
}

func GetUserByID(id uint) (user models.User, err error) {
	err = db.GetDBConn().Where("id = ?", id).First(&user).Error
	if err != nil {
		logger.Error.Println("[repository.GetUserByID] cannot get user by id. Error is:", err.Error())
		return user, translateError(err)
	}

	return user, nil
}

func GetUserByUsername(username string) (user models.User, err error) {
	err = db.GetDBConn().Where("username = ?", username).First(&user).Error
	if err != nil {
		logger.Error.Println("[repository.GetUserByUsername] cannot get user by username. Error is:", err.Error())
		return user, translateError(err)
	}

	return user, nil
}

func GetUserByUsernameAndPassword(username string, password string) (user models.User, err error) {
	err = db.GetDBConn().Where("username = ? AND password = ?", username, password).First(&user).Error
	if err != nil {
		logger.Error.Println("[repository.GetUserByUsernameAndPassword] cannot get user by username and password. Error is:", err.Error())
		return user, translateError(err)
	}

	return user, nil
}

func CreateUser(user models.User) (err error) {
	if err = db.GetDBConn().Create(&user).Error; err != nil {
		logger.Error.Println("[repository.CreateUser] cannot create user. Error is:", err.Error())
		return translateError(err)
	}

	return nil
}

func EditUserByID(user *models.User) (*models.User, error) {
	if err := db.GetDBConn().Omit("password").Save(&user).Error; err != nil {
		logger.Error.Println("[repository.EditUserByID] cannot edit user. Error is:", err.Error())
		return nil, translateError(err)
	}
	return user, nil
}

func DeleteUserByID(user *models.User) error {
	if err := db.GetDBConn().Delete(user).Error; err != nil {
		logger.Error.Println("[repository.DeleteUserByID] cannot delete user. Error is:", err.Error())
		return translateError(err)
	}
	return nil
}
