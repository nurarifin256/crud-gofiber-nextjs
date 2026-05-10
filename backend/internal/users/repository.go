package user

import (
	"context"

	"gorm.io/gorm"
)

type Repository interface {
	FindUserByNik(ctx context.Context, nik string) (*User, error)
}

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) Repository {
	return &UserRepositoryImpl{
		DB: DB,
	}
}

func (r *UserRepositoryImpl) FindUserByNik(ctx context.Context, nik string) (*User, error) {
	var user User
	query := `
		SELECT
			u.id,
			u.role_id,
			r.role as role_name,
			u.name,
			u.nik,
			u.email,
			u.phone_number,
			u.password,
			u.departement,
			u.picture,
			u.picture_sign,
			u.remember_token,
			u.created_nik,
			u.updated_nik,
			u.created_by,
			u.updated_by,
			TO_CHAR(u.created_at, 'YYYY-MM-DD HH24:MI:SS') as created_at,
			TO_CHAR(u.updated_at, 'YYYY-MM-DD HH24:MI:SS') as updated_at
		FROM
			m_user u
		LEFT JOIN m_role r ON u.role_id = r.id
		WHERE
			u.nik = ? AND u.deleted_at IS NULL
	`
	if err := r.DB.WithContext(ctx).Raw(query, nik).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}
