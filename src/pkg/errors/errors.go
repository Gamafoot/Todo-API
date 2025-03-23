package errors

import errLib "errors"

func MatchIn(err error, errs ...error) bool {
	for _, e := range errs {
		if errLib.Is(err, e) {
			return true
		}
	}
	return false
}
