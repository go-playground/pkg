package strconvext

import (
	"strconv"
)

// ParseBool returns the boolean value represented by the string. It extends the std library parse bool with a few more
// valid options.
//
// It accepts 1, t, T, true, TRUE, True, on, yes, ok as true values and 0, f, F, false, FALSE, False, off, no as false.
func ParseBool(str string) (bool, error) {
	switch str {
	case "1", "t", "T", "true", "TRUE", "True", "on", "yes", "ok":
		return true, nil
	case "", "0", "f", "F", "false", "FALSE", "False", "off", "no":
		return false, nil
	}
	// strconv.NumError mimicking exactly the strconv.ParseBool(..) error and type
	// to ensure compatibility with std library.
	return false, &strconv.NumError{Func: "ParseBool", Num: string([]byte(str)), Err: strconv.ErrSyntax}
}
