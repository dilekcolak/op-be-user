package application

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	smodel "github.com/octoposprime/op-be-shared/pkg/model"
	pb_logging "github.com/octoposprime/op-be-shared/pkg/proto/pb/logging"
	me "github.com/octoposprime/op-be-user/internal/domain/model/entity"
	mo "github.com/octoposprime/op-be-user/internal/domain/model/object"
	"golang.org/x/crypto/bcrypt"
)

// GetUsersByFilter returns the users that match the given filter.
func (a *Service) GetUsersByFilter(ctx context.Context, userFilter me.UserFilter) (me.Users, error) {
	return a.DbPort.GetUsersByFilter(ctx, userFilter)
}

// CreateUser sends the given user to the repository of the infrastructure layer for creating a new user.
func (a *Service) CreateUser(ctx context.Context, user me.User) (me.User, error) {
	user.Id = uuid.UUID{}
	if err := a.ValidateUser(&user); err != nil {
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "CreateUser", userId, err.Error()))
		return me.User{}, err
	}
	if err := a.CheckUserNameRules(&user); err != nil {
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "CreateUser", userId, err.Error()))
		return me.User{}, err
	}
	if err := a.CheckEmailRules(&user); err != nil {
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "CreateUser", userId, err.Error()))
		return me.User{}, err
	}
	var userEmailCheckFilter me.UserFilter
	userEmailCheckFilter.Email = user.Email
	emailExistsUsers, err := a.GetUsersByFilter(ctx, userEmailCheckFilter)
	if err != nil {
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "CreateUser", userId, err.Error()))
		return me.User{}, err
	}
	if emailExistsUsers.TotalRows > 0 {
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		err := mo.ErrorUserEmailIsExists
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "CreateUser", userId, err.Error()))
		return me.User{}, err
	}
	var userNameCheckFilter me.UserFilter
	userNameCheckFilter.UserName = user.UserName
	nameExistsUsers, err := a.GetUsersByFilter(ctx, userNameCheckFilter)
	if err != nil {
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "CreateUser", userId, err.Error()))
		return me.User{}, err
	}
	if nameExistsUsers.TotalRows > 0 {
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		err := mo.ErrorUserUsernameIsExists
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "CreateUser", userId, err.Error()))
		return me.User{}, err
	}
	if user.UserStatus == mo.UserStatusNONE {
		user.UserStatus = mo.UserStatusACTIVE
	}
	return a.DbPort.SaveUser(ctx, user)
}

// UpdateUserBase sends the given base values of the user to the repository of the infrastructure layer for updating base values of user data.
func (a *Service) UpdateUserBase(ctx context.Context, user me.User) (me.User, error) {
	if user.Id.String() == "" || user.Id == (uuid.UUID{}) {
		err := mo.ErrorUserIdIsEmpty
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "UpdateUserStatus", userId, err.Error()))
		return me.User{}, err
	}
	var userFilter me.UserFilter
	userFilter.Id = user.Id
	users, err := a.GetUsersByFilter(ctx, userFilter)
	if err != nil {
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "UpdateUserStatus", userId, err.Error()))
		return me.User{}, err
	}
	if users.TotalRows > 0 {
		dbUser := users.Users[0]
		dbUser.Tags = user.Tags
		dbUser.FirstName = user.FirstName
		dbUser.LastName = user.LastName
		if err := a.ValidateUser(&dbUser); err != nil {
			userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
			go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "UpdateUserStatus", userId, err.Error()))
			return me.User{}, err
		}
		return a.DbPort.SaveUser(ctx, dbUser)
	} else {
		return user, mo.ErrorUserNotFound
	}
}

// UpdateUserStatus sends the given status value of the user to the repository of the infrastructure layer for updating status of user data.
func (a *Service) UpdateUserStatus(ctx context.Context, user me.User) (me.User, error) {
	if user.Id.String() == "" || user.Id == (uuid.UUID{}) {
		err := mo.ErrorUserIdIsEmpty
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "UpdateUserStatus", userId, err.Error()))
		return me.User{}, err
	}
	var userFilter me.UserFilter
	userFilter.Id = user.Id
	users, err := a.GetUsersByFilter(ctx, userFilter)
	if err != nil {
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "UpdateUserStatus", userId, err.Error()))
		return me.User{}, err
	}
	if users.TotalRows > 0 {
		dbUser := users.Users[0]
		dbUser.UserStatus = user.UserStatus
		if err := a.ValidateUser(&dbUser); err != nil {
			userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
			go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "UpdateUserStatus", userId, err.Error()))
			return me.User{}, err
		}
		return a.DbPort.SaveUser(ctx, dbUser)
	} else {
		return user, mo.ErrorUserNotFound
	}
}

// UpdateUserRole sends the given type value of the user to the repository of the infrastructure layer for updating role of user data.
func (a *Service) UpdateUserRole(ctx context.Context, user me.User) (me.User, error) {
	if user.Id.String() == "" || user.Id == (uuid.UUID{}) {
		err := mo.ErrorUserIdIsEmpty
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "UpdateUserRole", userId, err.Error()))
		return me.User{}, err
	}
	var userFilter me.UserFilter
	userFilter.Id = user.Id
	users, err := a.GetUsersByFilter(ctx, userFilter)
	if err != nil {
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "UpdateUserRole", userId, err.Error()))
		return me.User{}, err
	}
	if users.TotalRows > 0 {
		dbUser := users.Users[0]
		dbUser.Role = user.Role
		if err := a.ValidateUser(&dbUser); err != nil {
			userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
			go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "UpdateUserRole", userId, err.Error()))
			return me.User{}, err
		}
		return a.DbPort.SaveUser(ctx, dbUser)
	} else {
		return user, mo.ErrorUserNotFound
	}
}

// DeleteUser sends the given user to the repository of the infrastructure layer for deleting data.
func (a *Service) DeleteUser(ctx context.Context, user me.User) (me.User, error) {
	var err error
	user, err = a.DbPort.DeleteUser(ctx, user)
	if err != nil {
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "DeleteUser", userId, err.Error()))
		return me.User{}, err
	}

	err = a.RedisPort.DeleteUserPasswordByUserId(ctx, user.Id)
	if err != nil {
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "DeleteUser", userId, err.Error()))
		return me.User{}, err
	}
	return user, err
}

// ChangePassword sends the given user password to the repository of the infrastructure layer for changing user password.
func (a *Service) ChangePassword(ctx context.Context, userPassword me.UserPassword) error {
	if userPassword.UserId.String() == "" || userPassword.UserId == (uuid.UUID{}) {
		err := mo.ErrorUserIdIsEmpty
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "ChangePassword", userId, err.Error()))
		return err
	}
	if err := a.ValidatePassword(&userPassword); err != nil {
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "ChangePassword", userId, err.Error()))
		return err
	}
	if err := a.CheckPasswordRules(&userPassword); err != nil {
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "ChangePassword", userId, err.Error()))
		return err
	}
	passByte, err := bcrypt.GenerateFromPassword([]byte(userPassword.Password), 4)
	if err != nil {
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "ChangePassword", userId, err.Error()))
		return err
	}
	userPassword.Password = string(passByte)
	userPassword.PasswordStatus = mo.PasswordStatusACTIVE
	userPassword, err = a.DbPort.ChangePassword(ctx, userPassword)
	if err != nil {
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "ChangePassword", userId, err.Error()))
		return err
	}
	err = a.RedisPort.ChangePassword(ctx, userPassword)
	if err != nil {
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "ChangePassword", userId, err.Error()))
		return err
	}
	return err
}

// GetUserPasswordByUserId returns active password of the given user.
func (a *Service) GetUserPasswordByUserId(ctx context.Context, userId uuid.UUID) (me.UserPassword, error) {
	if userId.String() == "" || userId == (uuid.UUID{}) {
		err := mo.ErrorUserIdIsEmpty
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "GetUserPasswordByUserId", userId, err.Error()))
		return me.UserPassword{}, err
	}
	// if the userPassword is cached in the redis repository return it
	userPassword, err := a.RedisPort.GetUserPasswordByUserId(ctx, userId)
	if err == nil && userPassword.UserId == userId {
		return userPassword, err
	}
	// else the userPassword is not cached in the redis repository get and return the userPassword from db
	userPassword, err = a.DbPort.GetUserPasswordByUserId(ctx, userId)
	if err != nil {
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "GetUserPasswordByUserId", userId, err.Error()))
		return userPassword, err
	}
	// and also write it to redis.
	err = a.RedisPort.ChangePassword(ctx, userPassword)
	if err != nil {
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "GetUserPasswordByUserId", userId, err.Error()))
		return userPassword, err
	}
	return userPassword, err
}

// CheckUserPassword sends the given user password to the repository of the infrastructure layer for checking user password.
func (a *Service) CheckUserPassword(ctx context.Context, user me.User, userPassword me.UserPassword) (me.User, error) {
	fmt.Println(user.String(), "user.String")
	var users me.Users
	var err error
	returnErr := errors.New("Request failed")
	if user.UserName == "" {
		if user.Email == "" {
			// err := mo.ErrorUserEmailAndUsernameIsEmpty // TODO : Implement this error
			err = fmt.Errorf("user email and username is empty")
			userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
			go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "CheckUserPassword", userId, err.Error()))
			return me.User{}, returnErr
		}
		users, err = a.GetUsersByFilter(ctx, me.UserFilter{Email: user.Email})
		if err != nil {
			userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
			go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "CheckUserPassword", userId, err.Error()))
			return me.User{}, returnErr
		}
	}
	users, err = a.GetUsersByFilter(ctx, me.UserFilter{UserName: user.UserName})
	if err != nil {
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "CheckUserPassword", userId, err.Error()))
		return me.User{}, returnErr
	}
	if users.TotalRows == 0 {
		err := mo.ErrorUserNotFound
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "CheckUserPassword", userId, err.Error()))
		return me.User{}, returnErr
	}
	if users.TotalRows > 1 {
		// err := mo.ErrorMultipleUserFound //TODO : Implement this error
		err := fmt.Errorf("multiple user found")
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "CheckUserPassword", userId, err.Error()))
		return me.User{}, returnErr
	}
	if users.TotalRows == 1 {
		user = users.Users[0]
	}

	if user.Id.String() == "" || user.Id == (uuid.UUID{}) {
		err := mo.ErrorUserIdIsEmpty
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "CheckUserPassword", userId, err.Error()))
		return me.User{}, returnErr
	}
	// if userPassword.UserId.String() == "" || userPassword.UserId == (uuid.UUID{}) {
	// 	err := mo.ErrorUserIdIsEmpty
	// 	userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
	// 	go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "CheckUserPassword", userId, err.Error()))
	// 	return me.User{}, err
	// }
	// Not sure if this Validation is necessary or right approach.
	// if err := a.ValidatePassword(&userPassword); err != nil {
	// 	userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
	// 	go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "CheckUserPassword", userId, err.Error()))
	// 	return me.User{}, err
	// }
	// This may be cause problem in the future if the password rules are changed.
	// if err := a.CheckPasswordRules(&userPassword); err != nil {
	// 	userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
	// 	go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "CheckUserPassword", userId, err.Error()))
	// 	return me.User{}, err
	// }
	fmt.Println(user.String(), "user.String")

	// dbUser, err := a.DbPort.CheckUserPassword(ctx, user, userPassword)
	dbPassword, err := a.GetUserPasswordByUserId(ctx, user.Id)

	if dbPassword.PasswordStatus != mo.PasswordStatusACTIVE {
		// err := mo.ErrorUserPasswordIsNotActive //TODO : Implement this error
		err := fmt.Errorf("user password is not active")
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "CheckUserPassword", userId, err.Error()))
		return me.User{}, returnErr
	}

	isPasswordMatch := dbPassword.ComparePass(userPassword.Password)

	if !isPasswordMatch {
		// err := mo.ErrorUserPasswordIsNotMatch //TODO : Implement this error
		err := fmt.Errorf("user password is not match")
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "CheckUserPassword", userId, err.Error()))
		return me.User{}, returnErr
	}

	if err != nil {
		// err := mo.ErrorUserPasswordIsNotMatch //TODO : Implement this error
		err := fmt.Errorf("user password is not match")
		userId, _ := ctx.Value(smodel.QueryKeyUid).(string)
		go a.Log(context.Background(), me.NewLogData().GenerateLogData(pb_logging.LogType_LogTypeERROR, "CheckUserPassword", userId, err.Error()))
		return me.User{}, returnErr
	}
	return user, nil
}
