package service

import (
	"delivery/errs"
	"delivery/internal/models"
	repository "delivery/internal/repositories"
	"delivery/utils"
	"errors"
	"fmt"
)

func GetAllUsers() (users []models.User, err error) {
	users, err = repository.GetAllUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetUserByID(id uint) (user models.User, err error) {
	user, err = repository.GetUserByID(id)
	if err != nil {
		if errors.Is(err, errs.ErrRecordNotFound) {
			return user, errs.ErrUserNotFound
		}
		return user, err
	}

	return user, nil
}

func CreateUser(user models.User) error {
	// 1. Check username uniqueness
	userFromDB, err := repository.GetUserByUsername(user.Username)
	if err != nil && !errors.Is(err, errs.ErrRecordNotFound) {
		return err
	}

	if userFromDB.ID > 0 {
		return errs.ErrUsernameUniquenessFailed
	}

	user.Role = "viewer"

	// 2. Generate password hash
	user.Password = utils.GenerateHash(user.Password)

	// 3. Repository call
	err = repository.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}

func EditUserByID(id uint, userInput models.User) (*models.User, error) {
	_, err := repository.GetUserByID(id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %v", err)
	}

	userInput.ID = id

	updatedUser, err := repository.EditUserByID(&userInput)
	if err != nil {
		return nil, fmt.Errorf("could not update user: %v", err)
	}

	return updatedUser, nil
}

func DeleteUserByID(id uint) error {
	user, err := repository.GetUserByID(id)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	if err := repository.DeleteUserByID(&user); err != nil {
		return fmt.Errorf("could not delete user: %v", err)
	}

	return nil
}
