package authentication

import (
	"github.com/gofrs/uuid"
	"strings"
)

//
// UUID
// @Description:
// @return string
//
func UUID() string {
	v4, err := uuid.NewV4()
	if err != nil {
		return ""
	}
	return strings.ReplaceAll(v4.String(), "-", "")
}
