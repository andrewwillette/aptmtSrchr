package aptmtSrchr

import (
	"fmt"
	"testing"
)

func TestGetApartments(t *testing.T) {
	apartments := GetApartments()
	for _, apt := range apartments {
		fmt.Printf("%+v\n", apt)
	}
}
