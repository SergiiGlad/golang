package controllers

import (
	"database/sql"
	"errors"
	"go-team-room/conf"
	"go-team-room/models/dao"
	"go-team-room/models/dao/interfaces"
	"go-team-room/models/dto"
	"gopkg.in/hlandau/passlib.v1/hash/bcrypt"
	"regexp"
	"strings"
)

// Get instance of logger (Formatter, Hook，Level, Output ).
// If you want to use only your log message  It will need use own call logs example
var log = conf.GetLog()

type UserService struct {
	Dao interfaces.UserDao
}

var _ UserServiceInterface = &UserService{}

func (us *UserService) CreateUser(userDto *dto.RequestUserDto) (dto.ResponseUserDto, error) {

	userEntity := dto.RequestUserDtoToEntity(userDto)
	var respUserDto dto.ResponseUserDto

	err := checkUniqueEmail(userEntity.Email, us.Dao)

	if err != nil && err != sql.ErrNoRows {
		return respUserDto, err
	}

	err = checkUniquePhone(userEntity.Phone, us.Dao)

	if err != nil && err != sql.ErrNoRows {
		return respUserDto, err
	}

	if validPasswordLength(userEntity.CurrentPass) == false {
		return respUserDto, errors.New("Password too short.")
	}

	hashPass, err := bcrypt.Crypter.Hash(userEntity.CurrentPass)

	if err != nil {
		return respUserDto, err
	}

	userEntity.CurrentPass = hashPass
	nameLetterToUppep(&userEntity)

	_, err = us.Dao.AddUser(&userEntity)

	if err != nil {
		log.Println(err)
		return respUserDto, err
	}

	respUserDto = dto.UserEntityToResponseDto(&userEntity)

	return respUserDto, nil
}

func (us *UserService) UpdateUser(id int64, userDto *dto.RequestUserDto) (dto.ResponseUserDto, error) {

	newUserData := dto.RequestUserDtoToEntity(userDto)
	oldUserData, err := us.Dao.FindUserById(id)
	var responseUserDto dto.ResponseUserDto

	if err != nil {
		return responseUserDto, err
	}

	if len(newUserData.FirstName) == 0 {
		newUserData.FirstName = oldUserData.FirstName
	}

	if len(newUserData.SecondName) == 0 {
		newUserData.SecondName = oldUserData.SecondName
	}

	if len(newUserData.Email) != 0 {
		err = checkUniqueEmail(newUserData.Email, us.Dao)

		if err != nil && err != sql.ErrNoRows {
			return responseUserDto, err
		}
	} else {
		newUserData.Email = oldUserData.Email
	}

	if len(newUserData.Phone) != 0 {
		err = checkUniquePhone(newUserData.Phone, us.Dao)

		if err != nil && err != sql.ErrNoRows {
			return responseUserDto, err
		}
	} else {
		newUserData.Phone = oldUserData.Phone
	}

	if len(newUserData.CurrentPass) != 0 {
		if validPasswordLength(newUserData.CurrentPass) == false {
			return responseUserDto, errors.New("Password too short.")
		}

		hashPass, err := bcrypt.Crypter.Hash(newUserData.CurrentPass)

		if err != nil {
			log.Println(err)
			return responseUserDto, err
		}

		newUserData.CurrentPass = hashPass
	} else {
		newUserData.CurrentPass = oldUserData.CurrentPass
	}

	nameLetterToUppep(&newUserData)

	err = us.Dao.UpdateUser(id, &newUserData)
	if err != nil {
		log.Println(err)
		return responseUserDto, err
	}

	newUserData.ID = id
	responseUserDto = dto.UserEntityToResponseDto(&newUserData)
	responseUserDto.Friends, _ = us.GetUserFriends(id)

	return responseUserDto, nil
}

func (us *UserService) DeleteUser(id int64) (dto.ResponseUserDto, error) {

	var responseUserDto dto.ResponseUserDto

	userEntity, err := us.Dao.FindUserById(id)

	if err != nil {
		return responseUserDto, err
	}

	responseUserDto = dto.UserEntityToResponseDto(userEntity)
	responseUserDto.Friends, _ = us.GetUserFriends(id)

	return responseUserDto, us.Dao.DeleteUser(id)
}

func (us *UserService) GetUserFriends(id int64) ([]int64, error) {
	_, err := us.Dao.FindUserById(id)

	if err != nil {
		return nil, err
	}

	return us.Dao.FriendsByUserID(id)
}

func checkUniqueEmail(email string, dao interfaces.UserDao) error {

	if validEmail(email) == false {
		return errors.New("Invalid email format.")
	} else {

		_, err := dao.FindUserByEmail(email)

		switch err {
		case sql.ErrNoRows:
			return err

		case nil:
			return errors.New("There is user with such email.")

		default:
			log.Println(err)
			return err
		}
	}

	return nil
}

func checkUniquePhone(phone string, dao interfaces.UserDao) error {
	if len(phone) > 0 {

		if validPhone(phone) == false {
			return errors.New("Invalid phone number format.")
		} else {

			_, err := dao.FindUserByPhone(phone)

			switch err {
			case sql.ErrNoRows:
				return err

			case nil:
				return errors.New("There is user with such phone.")

			default:
				log.Println(err)
				return err
			}
		}
	}

	return nil
}

func validRegexItem(item string, pattern string) bool {
	itemRegex := regexp.MustCompile(pattern)

	if isItemOk := itemRegex.MatchString(item); isItemOk == false {
		return false
	}

	return true
}

func validEmail(email string) bool {
	return validRegexItem(email, "^[a-z0-9]+@[a-z]+[.][a-z]+$")
}

func validPhone(phone string) bool {
	return validRegexItem(phone, "^[+][0-9]{12}")
}

func validCyrillicName(name string) bool {
	return validRegexItem(name, "^[А-Я][а-я]{1,49}")
}

func validLatinName(name string) bool {
	return validRegexItem(name, "^[A-Z][a-z]{1,49}")
}

func validPasswordLength(password string) bool {
	if len(password) < 6 {
		return false
	}

	return true
}

func nameLetterToUppep(user *dao.User) {
	user.FirstName = strings.ToUpper(string([]rune(user.FirstName)[0])) + string([]rune(user.FirstName)[1:])
	user.SecondName = strings.ToUpper(string([]rune(user.SecondName)[0])) + string([]rune(user.SecondName)[1:])
}
