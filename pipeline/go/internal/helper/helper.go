package helper

import "database/sql"

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
