package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/xuri/excelize/v2"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func main() {
	// Open the file
	file, err := os.Open("coso.tmcf")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Create a map to store the values and observationAbout
	observation := []string{}
	place := []string{}

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Skip empty lines
		if line == "" {
			continue
		}

		// Extract the observationAbout value
		if strings.HasPrefix(line, "observationAbout") {
			obsAboutParts := strings.Split(line, "dcid:")
			observation = append(observation, obsAboutParts[1])
		}

		// Extract the value for other observations
		if strings.HasPrefix(line, "value") {
			valueParts := strings.Split(line, "->")
			place = append(place, valueParts[1])
		}
	}

	placeMap := make(map[string]string)

	for i, v := range place {
		placeMap[v] = observation[i]
	}

	f, err := excelize.OpenFile("sample.xlsx")
	if err != nil {
		return
	}
	defer f.Close()

	// change this
	rows, err := f.GetRows("1.6")
	if err != nil {
		fmt.Println("Error opening sheet:", err)
		return
	}
	count := 0
	// change this
	data := [][]string{{"locations", "totalLevel2", "totalUrbanLevel2", "totalLevel3", "totalUrbanLevel3"}}
	// change this
	locationIndex := 1
	for _, row := range rows {
		if len(row) > 1 {
			result := ConvertToUTF8(row[locationIndex])
			if value, ok := placeMap[result]; ok {
				count++
				rw := []string{value}
				for i := locationIndex + 1; i < len(row); i++ {
					rw = append(rw, row[i])
				}
				data = append(data, rw)
			} else {
				fmt.Println("false on", row[locationIndex])
			}
		}

	}
	fmt.Println(count)

	// Create the CSV file
	csvfile, err := os.Create("output.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer csvfile.Close()

	// Create a CSV writer
	writer := csv.NewWriter(csvfile)
	defer writer.Flush()

	// Write data to the CSV file
	for _, row := range data {
		err := writer.Write(row)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println("CSV file created successfully.")
}
func ConvertToUTF8(input string) string {
	// func main() {

	// Create a transformation to remove diacritics
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isDiacritic), norm.NFC)

	// Apply the transformation to the input string
	result, _, _ := transform.String(t, input)

	// Convert to lowercase
	result = strings.ToLower(result)

	// Remove any remaining non-alphanumeric characters and spaces
	result = removeNonAlphaNumeric(result)
	result = strings.ReplaceAll(result, " ", "")
	return result
	// fmt.Println(result) // Output: danang
}

// isDiacritic checks if the rune is a diacritic character
func isDiacritic(r rune) bool {
	return unicode.Is(unicode.Mn, r)
}

// removeNonAlphaNumeric removes any non-alphanumeric characters from the string
func removeNonAlphaNumeric(s string) string {
	var sb strings.Builder

	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			if r == 'đ' {
				r = 'd'
			}
			if r > unicode.MaxASCII {
				fmt.Println("err in", s)
			}

			sb.WriteRune(r)
		}
	}
	return sb.String()
}
