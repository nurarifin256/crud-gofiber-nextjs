package user

import "errors"

var (
	ErrInvalidCredentials = errors.New("incorrect NIK or password")
	ErrUserNotFound       = errors.New("user not found")
)

type User struct {
	ID            int64   `json:"id" gorm:"column:id"`
	RoleID        *int32  `json:"role_id" gorm:"column:role_id"`
	Name          string  `json:"name" gorm:"column:name"`
	NIK           string  `json:"nik" gorm:"column:nik"`
	Email         *string `json:"email" gorm:"column:email"`
	PhoneNumber   *string `json:"phone_number" gorm:"column:phone_number"`
	Password      string  `json:"-" gorm:"column:password"`
	DepartmentID  *int32  `json:"department_id" gorm:"column:department_id"`
	Picture       *string `json:"picture" gorm:"column:picture"`
	RememberToken *string `json:"-" gorm:"column:remember_token"`
	CreatedBy     *string `json:"created_by" gorm:"column:created_by"`
	UpdatedBy     *string `json:"updated_by" gorm:"column:updated_by"`
	CreatedAt     *string `json:"created_at" gorm:"column:created_at"` // format yyyy-mm-dd hh:mm:ss
	UpdatedAt     *string `json:"updated_at" gorm:"column:updated_at"` // format yyyy-mm-dd hh:mm:ss
}

type SubmitUserRequest struct {
	RoleID       *int32  `form:"role_id" validate:"required"`
	Name         string  `form:"name" validate:"required"`
	NIK          string  `form:"nik" validate:"required"`
	Email        *string `form:"email" validate:"omitempty,email"`
	PhoneNumber  *string `form:"phone_number" validate:"omitempty"`
	Password     string  `form:"password" validate:"required,min=6"`
	DepartmentID *int32  `form:"department_id" validate:"required"`
}
