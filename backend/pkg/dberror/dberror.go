package dberror

import (
	"errors"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FriendlyErrors struct{}

func NewFriendlyErrors() *FriendlyErrors { return &FriendlyErrors{} }
func (p *FriendlyErrors) Name() string   { return "friendly_errors" }

func (p *FriendlyErrors) Initialize(db *gorm.DB) error {
	// Untuk semua jalur, kalau tx.Error terisi, maka ditangani.
	// Create / Update / Delete / Query
	db.Callback().Create().After("gorm:create").Register("friendly_errors:create", wrapTxErr)
	db.Callback().Update().After("gorm:update").Register("friendly_errors:update", wrapTxErr)
	db.Callback().Delete().After("gorm:delete").Register("friendly_errors:delete", wrapTxErr)
	db.Callback().Query().After("gorm:query").Register("friendly_errors:query", wrapTxErr)
	return nil
}

func wrapTxErr(tx *gorm.DB) {
	if tx.Error == nil {
		return
	}

	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return
	}

	orig := tx.Error
	friendly := translate(orig)

	zap.S().Error(orig)

	if friendly.Error() != orig.Error() {
		zap.S().Error(friendly)
	}

	tx.Error = friendly
}

// ====== Mapping logic (bisa kamu ubah wording-nya) ======

var (
	msgNotFound       = "Record not found."
	msgDuplicate      = "Duplicate value. This field must be unique."
	msgForeignKey     = "Operation blocked: this record is referenced by other data (foreign key)."
	msgNotNull        = "Required field cannot be empty."
	msgCheckViolation = "Validation failed: data does not meet the required constraints."
	msgTooLong        = "Value is too long for the target column."
	msgInvalidInput   = "Invalid input."
	msgInternal       = "Internal Server Error."
)

func translate(err error) error {
	// Pertahankan not found supaya middleware kamu bisa map ke 404
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return gorm.ErrRecordNotFound
	}

	// pgx
	var pgxE *pgconn.PgError
	if errors.As(err, &pgxE) {
		return errors.New(mapPgErr(pgxE.Code, pgxE.ConstraintName, pgxE.ColumnName))
	}

	// Fallback lain → jangan bocorkan detail
	return errors.New(msgInternal)
}

func mapPgErr(code, constraint, column string) string {
	switch code {
	case pgerrcode.UniqueViolation: // 23505
		return label(msgDuplicate, constraint, "")
	case pgerrcode.ForeignKeyViolation: // 23503
		return msgForeignKey
	case pgerrcode.NotNullViolation: // 23502
		return label(msgNotNull, "", column)
	case pgerrcode.CheckViolation: // 23514
		return msgCheckViolation
	case pgerrcode.StringDataRightTruncationWarning: // 22001
		return label(msgTooLong, "", column)
	case pgerrcode.InvalidTextRepresentation: // 22P02
		return msgInvalidInput
	case pgerrcode.NoDataFound: // 22P02
		return msgNotFound
	default:
		return msgInternal
	}
}

func label(base, constraint, column string) string {
	if column != "" {
		return fmt.Sprintf("%s: %s", base, column)
	}
	return base
}
