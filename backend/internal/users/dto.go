package user

import "errors"

var (
	ErrInvalidCredentials = errors.New("incorrect NIK or password")
	ErrUserNotFound       = errors.New("user not found")
)

type User struct {
	ID            int64   `json:"id" gorm:"column:id"`
	RoleID        *int32  `json:"role_id" gorm:"column:role_id"`
	RoleName      *string `json:"role_name" gorm:"column:role_name"`
	Name          string  `json:"name" gorm:"column:name"`
	NIK           string  `json:"nik" gorm:"column:nik"`
	Email         *string `json:"email" gorm:"column:email"`
	PhoneNumber   *string `json:"phone_number" gorm:"column:phone_number"`
	Password      string  `json:"-" gorm:"column:password"`
	Departement   *string `json:"departement" gorm:"column:departement"`
	Picture       *string `json:"picture" gorm:"column:picture"`
	RememberToken *string `json:"-" gorm:"column:remember_token"`
	CreatedNik    *string `json:"created_nik" gorm:"column:created_nik"`
	UpdatedNik    *string `json:"updated_nik" gorm:"column:updated_nik"`
	CreatedBy     *string `json:"created_by" gorm:"column:created_by"`
	UpdatedBy     *string `json:"updated_by" gorm:"column:updated_by"`
	CreatedAt     *string `json:"created_at" gorm:"column:created_at"` // format yyyy-mm-dd hh:mm:ss
	UpdatedAt     *string `json:"updated_at" gorm:"column:updated_at"` // format yyyy-mm-dd hh:mm:ss
}
