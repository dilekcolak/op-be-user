package presentation

import (
	"fmt"

	pb "github.com/octoposprime/op-be-shared/pkg/proto/pb/user"
	tuuid "github.com/octoposprime/op-be-shared/tool/uuid"
	me "github.com/octoposprime/op-be-user/internal/domain/model/entity"
	mo "github.com/octoposprime/op-be-user/internal/domain/model/object"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserWithPassword is a struct that represents the dto of a user with password values.
type UserWithPassword struct {
	proto *pb.UserWithPassword
}

// NewUserWithPassword creates a new *UserWithPassword.
func NewUserWithPassword(pb *pb.UserWithPassword) *UserWithPassword {
	return &UserWithPassword{
		proto: pb,
	}
}

// String returns a string representation of the UserWithPassword.
func (s *UserWithPassword) String() string {
	return fmt.Sprintf(
		"User: %v, "+
			"Password: %v",
		s.proto.User,
		"***",
	)
}

func NewUserWithPasswordFromEntity(entity me.User) *UserWithPassword {
	return &UserWithPassword{
		&pb.UserWithPassword{
			User: &pb.User{
				Id:         entity.Id.String(),
				Username:   entity.UserName,
				Email:      entity.Email,
				Role:       entity.Role,
				UserType:   pb.UserType(entity.UserType),
				UserStatus: pb.UserStatus(entity.UserStatus),
				FirstName:  entity.FirstName,
				LastName:   entity.LastName,
				CreatedAt:  timestamppb.New(entity.CreatedAt), // Not sure if this needs to be converted
				UpdatedAt:  timestamppb.New(entity.UpdatedAt), // Not sure if this needs to be converted
			},
		},
	}
}

// ToPb returns a protobuf representation of the UserWithPassword.
func (s *UserWithPassword) ToPb() *pb.UserWithPassword {
	return s.proto
}

// ToEntity returns a entity representation of the UserWithPassword.
func (s *UserWithPassword) ToEntity() (*me.User, *me.UserPassword) {
	return &me.User{
			Id: tuuid.FromString(s.proto.User.Id),
			User: mo.User{
				UserName:   s.proto.User.Username,
				Email:      s.proto.User.Email,
				Role:       s.proto.User.Role,
				UserType:   mo.UserType(s.proto.User.UserType),
				UserStatus: mo.UserStatus(s.proto.User.UserStatus),
				Tags:       s.proto.User.Tags,
				FirstName:  s.proto.User.FirstName,
				LastName:   s.proto.User.LastName,
			},
			CreatedAt: s.proto.User.CreatedAt.AsTime(), // Not sure if this needs to be converted
			UpdatedAt: s.proto.User.UpdatedAt.AsTime(), // Not sure if this needs to be converted
		}, &me.UserPassword{
			Id:     tuuid.FromString(s.proto.UserPassword.Id),
			UserId: tuuid.FromString(s.proto.UserPassword.UserId),
			UserPassword: mo.UserPassword{
				Password:       s.proto.UserPassword.Password,
				PasswordStatus: mo.PasswordStatus(s.proto.UserPassword.PasswordStatus),
			},
		}
}

type UserWithPasswords []*UserWithPassword

// NewUserWithPasswordsFromEntities creates a new []*UserWithPassword from entities.
func NewUserWithPasswordsFromEntities(entities []me.User) UserWithPasswords {
	userWithPasswords := make([]*UserWithPassword, len(entities))
	for i, entity := range entities {
		userWithPasswords[i] = NewUserWithPasswordFromEntity(entity)
	}
	return userWithPasswords
}

// NewEmptyUserWithPassword creates a new *UserWithPassword with empty values.
func NewEmptyUserWithPassword() *UserWithPassword {
	return &UserWithPassword{
		&pb.UserWithPassword{
			User: &pb.User{
				Id:         "",
				Username:   "",
				Email:      "",
				Role:       "",
				UserType:   pb.UserType_UserTypeNONE,
				UserStatus: pb.UserStatus_UserStatusNONE,
				Tags:       []string{},
				FirstName:  "",
				LastName:   "",
				CreatedAt:  timestamppb.Now(), // Not sure if this needs to be converted
				UpdatedAt:  timestamppb.Now(), // Not sure if this needs to be converted
			},
			UserPassword: &pb.UserPassword{
				Id:             "",
				UserId:         "",
				Password:       "",
				PasswordStatus: pb.PasswordStatus_PasswordStatusNONE,
			},
		},
	}
}

// ToPbs returns a protobuf representation of the UserWithPasswords
func (s UserWithPasswords) ToPbs() []*pb.UserWithPassword {
	userWithPasswords := make([]*pb.UserWithPassword, len(s))
	for i, userWithPassword := range s {
		userWithPasswords[i] = userWithPassword.proto
	}
	return userWithPasswords
}
