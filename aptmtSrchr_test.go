package aptmtSrchr

import (
	"fmt"
	"testing"
)

func TestGetApartments(t *testing.T) {
	uliApartmentQueries := []string{
		"https://www.uli.com/residential/apartment-search?field_property_target_id%5B%5D=4&field_property_target_id%5B%5D=1980&field_property_target_id%5B%5D=2133&field_bedrooms_value%5B%5D=studio&field_bedrooms_value%5B%5D=1_bed&field_bedrooms_value%5B%5D=1_bed_den&field_bedrooms_value%5B%5D=2_bed&field_bedrooms_value%5B%5D=2_bed_den&field_available_date_value_1%5Bvalue%5D%5Bdate%5D=July%2C+2022",
		"https://www.uli.com/residential/apartment-search?field_property_target_id%5B%5D=4&field_property_target_id%5B%5D=1980&field_property_target_id%5B%5D=2133&field_bedrooms_value%5B%5D=studio&field_bedrooms_value%5B%5D=1_bed&field_bedrooms_value%5B%5D=1_bed_den&field_bedrooms_value%5B%5D=2_bed&field_bedrooms_value%5B%5D=2_bed_den&field_available_date_value_1%5Bvalue%5D%5Bdate%5D=August%2C+2022",
	}
	apartments := GetUliMadisonAptmts(uliApartmentQueries)
	for _, apt := range apartments {
		fmt.Printf("%+v\n", apt)
	}
}
