package usecase

import (
	"encoding/json"
	"github.com/google/uuid"
	"test2/constanta"
	"test2/model"
	"test2/model/entity"
	"test2/repository/cacheRepository"
	"test2/repository/sqlRepository"
	"test2/util"
	"time"
)

type Auths interface {
	Register(u model.UserRequest) *model.Error
	Login(u model.UserRequest) (model.UserLoginResponse, *model.Error)
}

type Auth struct {
	repository sqlRepository.Users
	sqlRepo    sqlRepository.SqlDB
	cache      cacheRepository.Cache
	jwtUtil    util.JwtToken
}

func NewAuth(
	repository sqlRepository.Users,
	sqlRepo sqlRepository.SqlDB,
	cache cacheRepository.Cache,
	jwt util.JwtToken,
) Auths {
	return Auth{
		repository: repository,
		sqlRepo:    sqlRepo,
		cache:      cache,
		jwtUtil:    jwt,
	}
}

func (a Auth) Register(u model.UserRequest) *model.Error {
	// validate
	if err := a.validateRegister(u); err != nil {
		return err
	}

	_, err := a.validateUserRegistered(u.PhoneNumber, false)
	if err != nil {
		return err
	}

	tx := a.sqlRepo.DBBeginTX()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// insert into db
	err = a.repository.Insert(tx, &entity.User{
		FirstName:   u.FirstName,
		LastName:    u.LastName,
		PhoneNumber: u.PhoneNumber,
		Pin:         util.GenerateHmacSHA256(u.Pin),
		Address:     u.Address,
	})
	return err
}

func (a Auth) Login(u model.UserRequest) (model.UserLoginResponse, *model.Error) {
	// validate
	if err := a.validateLogin(u); err != nil {
		return model.UserLoginResponse{}, err
	}

	usr, err := a.validateUserRegistered(u.PhoneNumber, true)
	if err != nil {
		return model.UserLoginResponse{}, err
	}

	if usr.Pin != util.GenerateHmacSHA256(u.Pin) {
		return model.UserLoginResponse{}, model.NewError(401, "phone number or pin miss match", nil)
	}

	// generate token
	accessToken, err := a.jwtUtil.GenerateToken(usr.ID.String(), 1*time.Hour)
	if err != nil {
		return model.UserLoginResponse{}, err
	}

	refreshToken, err := a.jwtUtil.GenerateToken(usr.ID.String(), 7*24*time.Hour)
	if err != nil {
		return model.UserLoginResponse{}, err
	}

	authInfo, _ := json.Marshal(model.UserAccess{
		UserID:    usr.ID.String(),
		FirstName: usr.FirstName,
		LastName:  usr.LastName,
	})

	// set data into redis
	err = a.cache.Set(accessToken, string(authInfo), 1*time.Hour)
	if err != nil {
		return model.UserLoginResponse{}, err
	}

	// set response
	return model.UserLoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a Auth) validateRegister(r model.UserRequest) *model.Error {
	if err := a.validateLogin(r); err != nil {
		return err
	}
	return nil
}

func (a Auth) validateLogin(u model.UserRequest) *model.Error {
	if util.IsEmptyStringWithTrimSpace(&u.PhoneNumber) {
		return model.NewError(400, "phone number can't empty", nil)
	}
	if util.IsEmptyString(u.Pin) {
		return model.NewError(400, "pin can't empty", nil)
	}
	if valid, _ := util.ValidateRegex(u.PhoneNumber, constanta.RegexPhoneNumber); !valid {
		return model.NewError(400, "invalid format phone number", nil)
	}
	if len(u.Pin) != 6 {
		return model.NewError(400, "pin must 6 number", nil)
	}
	if valid, _ := util.ValidateRegex(u.Pin, constanta.RegexPinNumber); !valid {
		return model.NewError(400, "invalid format pin", nil)
	}
	return nil
}

func (a Auth) validateUserRegistered(phoneNumber string, isLogin bool) (entity.User, *model.Error) {
	usr, err := a.repository.GetByPhoneNumber(phoneNumber)
	if err != nil {
		return entity.User{}, err
	}
	if !isLogin && usr.ID != uuid.Nil {
		return entity.User{}, model.NewError(400, "user has been registered", nil)
	}

	if isLogin && usr.ID == uuid.Nil {
		return entity.User{}, model.NewError(400, "user not registered", nil)
	}
	return usr, nil
}
