package aptmtSrchr

import (
	"fmt"
	"testing"
)

func TestGetApartments(t *testing.T) {
	apartments := GetApartments()
	fmt.Println(apartments)	
}
