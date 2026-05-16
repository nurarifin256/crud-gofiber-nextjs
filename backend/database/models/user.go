package model

import (
	"time"
)

type User struct {
	ID            int64      `gorm:"column:id;type:bigserial;primaryKey;autoIncrement" json:"id"`
	RoleID        *int32     `gorm:"column:role_id;type:int4" json:"role_id,omitempty"`
	Name          string     `gorm:"column:name;type:varchar(255);not null" json:"name"`
	NIK           string     `gorm:"column:nik;type:varchar(255);not null" json:"nik"`
	Email         *string    `gorm:"column:email;type:varchar(255)" json:"email,omitempty"`
	PhoneNumber   *string    `gorm:"column:phone_number;type:varchar(255)" json:"phone_number,omitempty"`
	Password      string     `gorm:"column:password;type:varchar(255);not null" json:"password"`
	DepartementID *int32     `gorm:"column:departement_id;type:int4" json:"departement_id,omitempty"`
	Picture       *string    `gorm:"column:picture;type:varchar(255)" json:"picture,omitempty"`
	RememberToken *string    `gorm:"column:remember_token;type:varchar(255)" json:"-"`
	CreatedBy     *string    `gorm:"column:created_by;type:varchar(50)" json:"created_by,omitempty"`
	UpdatedBy     *string    `gorm:"column:updated_by;type:varchar(50)" json:"updated_by,omitempty"`
	CreatedAt     *time.Time `gorm:"column:created_at;type:timestamp(0)" json:"created_at,omitempty"`
	UpdatedAt     *time.Time `gorm:"column:updated_at;type:timestamp(0)" json:"updated_at,omitempty"`
}

func (User) TableName() string { return "m_users" }
