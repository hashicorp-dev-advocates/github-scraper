package database

import (
	"database/sql/driver"
	"strings"
)

type StringArray []string

// e.g. ["postgres", "golang"] -> {"postgres","golang"}
func (x StringArray) Value() (driver.Value, error) {
	if len(x) == 0 {
		return "{}", nil
	}

	transform := "{\"" + strings.Join(x, "\", \"") + "\"}"

	return transform, nil
}

// e.g. {"postgres","golang"} -> ["postgres", "golang"]
func (x *StringArray) Scan(src interface{}) error {
	v := string(src.(string))

	replacer := strings.NewReplacer("{", "", "\"", "", "}", "", "\\", "")
	v = replacer.Replace(v)
	transform := strings.Split(v, ",")

	*x = transform
	return nil
}
