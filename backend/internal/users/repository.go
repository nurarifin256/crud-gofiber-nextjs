package user

import (
	"context"

	model "simple-crud/database/models"

	"gorm.io/gorm"
)

type Repository interface {
	FindUserByNik(ctx context.Context, nik string) (*model.User, error)
	SubmitUser(ctx context.Context, req SubmitUserRequest, picture string, password string) (*model.User, error)
}

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) Repository {
	return &UserRepositoryImpl{
		DB: DB,
	}
}

func (r *UserRepositoryImpl) FindUserByNik(ctx context.Context, nik string) (*model.User, error) {
	var user model.User
	query := `
		SELECT
			u.id,
			u.role_id,
			u.name,
			u.nik,
			u.email,
			u.phone_number,
			u.password,
			u.departement_id,
			u.picture,
			u.remember_token,
			u.created_by,
			u.updated_by,
			TO_CHAR(u.created_at, 'YYYY-MM-DD HH24:MI:SS') as created_at,
			TO_CHAR(u.updated_at, 'YYYY-MM-DD HH24:MI:SS') as updated_at
		FROM
			m_users u
		WHERE
			u.nik = ?
	`
	if err := r.DB.WithContext(ctx).Raw(query, nik).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) SubmitUser(
	ctx context.Context,
	req SubmitUserRequest,
	picture string,
	password string,
) (*model.User, error) {

	tx := r.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	user := &model.User{
		RoleID:       req.RoleID,
		Name:         req.Name,
		NIK:          req.NIK,
		Email:        req.Email,
		PhoneNumber:  req.PhoneNumber,
		Password:     password,
		DepartmentID: req.DepartmentID,
		Picture:      &picture,
	}

	if err := tx.Create(user).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return user, nil
}
