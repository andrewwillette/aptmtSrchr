package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/gocolly/colly"
	"github.com/spf13/cobra"
)

// tsafeApartments apartments map with mutex
// for concurrency use
type tsafeApartments struct {
	apartments map[apartment]bool
	mut        sync.Mutex
}

type apartment struct {
	availDate string
	unitTitle string
	bedrooms  int
	sqFootage int
	rent      int
}

func (aptmts *tsafeApartments) insertApartment(apt apartment) {
	aptmts.mut.Lock()
	defer aptmts.mut.Unlock()
	aptmts.apartments[apt] = true
}

func newApartmentsSet() *tsafeApartments {
	aptmtSet := map[apartment]bool{}
	return &tsafeApartments{apartments: aptmtSet}
}

func main() {
	handleCliConfigs()
	apartments := newApartmentsSet()
	c := colly.NewCollector(
		colly.AllowedDomains("www.uli.com"),
		colly.Async(true),
	)

	c.OnHTML(".unit-result-item", func(e *colly.HTMLElement) {
		temp := apartment{}
		temp.availDate = getAvailableDate(e.ChildText(".avail-date"))
		temp.unitTitle = e.ChildText(".unit-title")
		temp.sqFootage = getSqFootage(e.ChildText(".sq-footage"))
		temp.rent = getRent(e.ChildText(".rent"))
		temp.bedrooms = getBedrooms(e.ChildText(".bedrooms"))
		apartments.insertApartment(temp)
	})

	c.Limit(&colly.LimitRule{
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	// Crawl all reddits the user passes in
	uliApartmentQueries := []string{
		"https://www.uli.com/residential/apartment-search?field_property_target_id%5B%5D=4&field_property_target_id%5B%5D=1980&field_property_target_id%5B%5D=2133&field_bedrooms_value%5B%5D=studio&field_bedrooms_value%5B%5D=1_bed&field_bedrooms_value%5B%5D=1_bed_den&field_bedrooms_value%5B%5D=2_bed&field_bedrooms_value%5B%5D=2_bed_den&field_available_date_value_1%5Bvalue%5D%5Bdate%5D=July%2C+2022",
		"https://www.uli.com/residential/apartment-search?field_property_target_id%5B%5D=4&field_property_target_id%5B%5D=1980&field_property_target_id%5B%5D=2133&field_bedrooms_value%5B%5D=studio&field_bedrooms_value%5B%5D=1_bed&field_bedrooms_value%5B%5D=1_bed_den&field_bedrooms_value%5B%5D=2_bed&field_bedrooms_value%5B%5D=2_bed_den&field_available_date_value_1%5Bvalue%5D%5Bdate%5D=August%2C+2022",
	}

	for _, aptQuery := range uliApartmentQueries {
		c.Visit(aptQuery)
	}

	c.Wait()
	displayAptmts(apartments)
}

func displayAptmts(apartments *tsafeApartments) {
	aptmts := []apartment{}
	for aptmt := range apartments.apartments {
		aptmts = append(aptmts, aptmt)
	}
	aptmts = sortFromCliConfig(aptmts)
	for _, aptmt := range aptmts {
		fmt.Printf("%+v\n", aptmt)
	}
}

func sortFromCliConfig(apts []apartment) []apartment {
	sort.SliceStable(apts, func(i, j int) bool {
		switch sortedInput {
		case rent:
			return apts[i].rent < apts[j].rent
		case sqFeet:
			return apts[i].sqFootage < apts[j].sqFootage
		case availDate:
			return apts[i].availDate < apts[j].availDate
		default:
			return true
		}
	})
	return apts
}

func getBedrooms(html string) int {
	r, _ := regexp.Compile(`\d{1,2}`)
	result := r.FindString(html)
	intVar, _ := strconv.Atoi(result)
	return intVar
}

func getRent(html string) int {
	r, _ := regexp.Compile(`\d{1,4}`)
	result := r.FindString(html)
	intVar, _ := strconv.Atoi(result)
	return intVar
}

func getSqFootage(html string) int {
	r, _ := regexp.Compile(`\d{0,4}`)
	result := r.FindString(html)
	intVar, _ := strconv.Atoi(result)
	return intVar
}

func getAvailableDate(html string) string {
	r, _ := regexp.Compile(`\d{1,2}/\d{1,2}/\d{1,4}`)
	return r.FindString(html)
}

type aptmtSortable string

const (
	rent      aptmtSortable = "r"
	availDate aptmtSortable = "d"
	sqFeet    aptmtSortable = "s"
)

// String is used both by fmt.Print and by Cobra in help text
func (e *aptmtSortable) String() string {
	return string(*e)
}

// Set must have pointer receiver so it doesn't change the value of a copy
func (e *aptmtSortable) Set(v string) error {
	switch v {
	case "r", "d", "s":
		*e = aptmtSortable(v)
		return nil
	default:
		return errors.New(`must be one of "r", "d", or "s"`)
	}
}

func (e *aptmtSortable) Type() string {
	return "aptmtSortable"
}

var verbose bool
var sortedInput aptmtSortable

var rootCmd = &cobra.Command{
	Use:   "get apartment details",
	Short: "apartment buying quickly",
	Long:  `Get apartments quickly but longer text`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func handleCliConfigs() {
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.Flags().VarP(&sortedInput, "sort", "s", `sort by partcular column: "r": rent, "d": availDate, "s": square feet`)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
