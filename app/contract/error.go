// This package hold enum values or custom type of wallet service response code and status
package contract

import "strings"

type (
	// Error hold error of wallet service.
	Error struct {
		Code StatusCode
		Raw  error

		Custom       string
		AppendFormat []string
		CustomAppend []string
	}
)

func (e *Error) String() (s string) {
	switch {
	case e == nil:
		break

	case len(strings.TrimSpace(e.Custom)) > 0:
		s = e.Code.String(e.Custom)

	case len(e.AppendFormat) > 0:
		s = e.Code.FormatedString(e.AppendFormat...)

	case e.Raw != nil:
		s = e.Raw.Error()

	case len(e.CustomAppend) > 0:
		s = e.Code.String() + " " + strings.Join(e.CustomAppend, " ")

	default:
		s = e.Code.String()
	}

	return
}
