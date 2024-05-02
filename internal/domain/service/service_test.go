package domain

import (
	"testing"

	"github.com/google/uuid"
	me "github.com/octoposprime/op-be-user/internal/domain/model/entity"
	mo "github.com/octoposprime/op-be-user/internal/domain/model/object"
)

func TestService_CheckUserNameRules(t *testing.T) {

	type args struct {
		user *me.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid Username",
			args: args{user: &me.User{
				Id: uuid.UUID{},
				User: mo.User{
					UserName: "Qwe123_ee",
				},
			}},
			wantErr: false,
		},
		{
			name: "Username With Invalid Characters",
			args: args{user: &me.User{
				Id: uuid.UUID{},
				User: mo.User{
					UserName: "asd!^+.++dfghdr",
				},
			}},
			wantErr: true,
		},
		{
			name: "Username With Space",
			args: args{user: &me.User{
				Id: uuid.UUID{},
				User: mo.User{
					UserName: "user name",
				},
			}},
			wantErr: true,
		},
		{
			name: "Too Short Username",
			args: args{user: &me.User{
				Id: uuid.UUID{},
				User: mo.User{
					UserName: "short",
				},
			}},
			wantErr: true,
		},
		{
			name: "Too Long Username",
			args: args{user: &me.User{
				Id: uuid.UUID{},
				User: mo.User{
					UserName: "toolongusernameisdefinitelytoolong",
				},
			}},
			wantErr: true,
		},
		{
			name: "Valid Numeric Username",
			args: args{user: &me.User{
				Id: uuid.UUID{},
				User: mo.User{
					UserName: "12345678",
				},
			}},
			wantErr: false,
		},
		{
			name: "Valid Username With Period",
			args: args{user: &me.User{
				Id: uuid.UUID{},
				User: mo.User{
					UserName: "John.Doe",
				},
			}},
			wantErr: false,
		},
		{
			name: "Valid Username With Underscores",
			args: args{user: &me.User{
				Id: uuid.UUID{},
				User: mo.User{
					UserName: "valid_user123",
				},
			}},
			wantErr: false,
		},
		{
			name: "All Underscores",
			args: args{user: &me.User{
				Id: uuid.UUID{},
				User: mo.User{
					UserName: "________",
				},
			}},
			wantErr: true,
		},
		{
			name: "All Periods",
			args: args{user: &me.User{
				Id: uuid.UUID{},
				User: mo.User{
					UserName: "......",
				},
			}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{}
			if err := s.CheckUserNameRules(tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("Service.CheckUserNameRules() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
