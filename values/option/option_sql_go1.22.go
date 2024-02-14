//go:build go1.22

package optionext

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"time"
)

var (
	scanType      = reflect.TypeFor[sql.Scanner]()
	valuerType    = reflect.TypeFor[driver.Valuer]()
	byteSliceType = reflect.TypeFor[[]byte]()
	timeType      = reflect.TypeFor[time.Time]()
	stringType    = reflect.TypeFor[string]()
	int64Type     = reflect.TypeFor[int64]()
	float64Type   = reflect.TypeFor[float64]()
	boolType      = reflect.TypeFor[bool]()
)

// Value implements the driver.Valuer interface.
//
// This honours the `driver.Valuer` interface if the value implements it.
// It also supports custom types of the std types and treats all else as []byte
func (o Option[T]) Value() (driver.Value, error) {
	if o.IsNone() {
		return nil, nil
	}
	val := reflect.ValueOf(o.value)

	if val.Type().Implements(valuerType) {
		return val.Interface().(driver.Valuer).Value()
	}
	switch val.Kind() {
	case reflect.String:
		return val.Convert(stringType).Interface(), nil
	case reflect.Bool:
		return val.Convert(boolType).Interface(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return val.Convert(int64Type).Interface(), nil
	case reflect.Float64:
		return val.Convert(float64Type).Interface(), nil
	case reflect.Slice, reflect.Array:
		if val.Type().ConvertibleTo(byteSliceType) {
			return val.Convert(byteSliceType).Interface(), nil
		}
		return json.Marshal(val.Interface())
	case reflect.Struct:
		if val.CanConvert(timeType) {
			return val.Convert(timeType).Interface(), nil
		}
		return json.Marshal(val.Interface())
	case reflect.Map:
		return json.Marshal(val.Interface())
	default:
		return o.value, nil
	}
}

// Scan implements the sql.Scanner interface.
func (o *Option[T]) Scan(value any) error {

	if value == nil {
		*o = None[T]()
		return nil
	}

	val := reflect.ValueOf(&o.value)

	if val.Type().Implements(scanType) {
		err := val.Interface().(sql.Scanner).Scan(value)
		if err != nil {
			return err
		}
		o.isSome = true
		return nil
	}

	val = val.Elem()

	switch val.Kind() {
	case reflect.String, reflect.Bool, reflect.Uint8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float64:
		var v sql.Null[T]
		if err := v.Scan(value); err != nil {
			return err
		}
		*o = Some(reflect.ValueOf(v.V).Convert(val.Type()).Interface().(T))
	case reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		v := reflect.ValueOf(value)
		if v.Type().ConvertibleTo(val.Type()) {
			*o = Some(reflect.ValueOf(v.Convert(val.Type()).Interface()).Interface().(T))
		} else {
			return fmt.Errorf("value %T not convertable to %T", value, o.value)
		}
	case reflect.Float32:
		var v sql.Null[float64]
		if err := v.Scan(value); err != nil {
			return err
		}
		*o = Some(reflect.ValueOf(v.V).Convert(val.Type()).Interface().(T))
	case reflect.Int:
		var v sql.Null[int64]
		if err := v.Scan(value); err != nil {
			return err
		}
		if v.V > math.MaxInt || v.V < math.MinInt {
			return fmt.Errorf("value %d out of range for int", v.V)
		}
		*o = Some(reflect.ValueOf(v.V).Convert(val.Type()).Interface().(T))
	case reflect.Int8:
		var v sql.Null[int64]
		if err := v.Scan(value); err != nil {
			return err
		}
		if v.V > math.MaxInt8 || v.V < math.MinInt8 {
			return fmt.Errorf("value %d out of range for int8", v.V)
		}
		*o = Some(reflect.ValueOf(v.V).Convert(val.Type()).Interface().(T))
	case reflect.Interface:
		*o = Some(reflect.ValueOf(value).Convert(val.Type()).Interface().(T))
	case reflect.Struct:
		if val.CanConvert(timeType) {
			switch t := value.(type) {
			case string:
				tm, err := time.Parse(time.RFC3339Nano, t)
				if err != nil {
					return err
				}
				*o = Some(reflect.ValueOf(tm).Convert(val.Type()).Interface().(T))

			case []byte:
				tm, err := time.Parse(time.RFC3339Nano, string(t))
				if err != nil {
					return err
				}
				*o = Some(reflect.ValueOf(tm).Convert(val.Type()).Interface().(T))

			default:
				var v sql.Null[time.Time]
				if err := v.Scan(value); err != nil {
					return err
				}
				*o = Some(reflect.ValueOf(v.V).Convert(val.Type()).Interface().(T))
			}
			return nil
		}
		fallthrough

	default:
		switch val.Kind() {
		case reflect.Struct, reflect.Slice, reflect.Map:
			v := reflect.ValueOf(value)

			if v.Type().ConvertibleTo(byteSliceType) {
				if val.Kind() == reflect.Slice && val.Type().Elem().Kind() == reflect.Uint8 {
					*o = Some(reflect.ValueOf(v.Convert(val.Type()).Interface()).Interface().(T))
				} else {
					if err := json.Unmarshal(v.Convert(byteSliceType).Interface().([]byte), &o.value); err != nil {
						return err
					}
				}
				o.isSome = true
				return nil
			}
		}
		return fmt.Errorf("unsupported Scan, storing driver.Value type %T into type %T", value, o.value)
	}
	return nil
}
