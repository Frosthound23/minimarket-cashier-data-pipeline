package helper

import (
	"database/sql"
	"time"
)

func NullableStringSQL(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}

	return &value.String
}

func NullStringToPointer(value interface{}) interface{} {
	switch v := value.(type) {
	case interface {
		Value() (driverValue interface{}, err error)
	}:
		driverValue, err := v.Value()
		if err != nil {
			return nil
		}

		return driverValue
	default:
		return value
	}
}

func NullableTimeSQL(value sql.NullTime) *time.Time {
	if !value.Valid {
		return nil
	}

	return &value.Time
}

func NullableInt32SQL(value sql.NullInt64) *int32 {
	if !value.Valid {
		return nil
	}

	v := int32(value.Int64)
	return &v
}
