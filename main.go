package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// BondPrefixes what bond numbers begin with, used to search and filter file
// results
var BondPrefixes = []string{"F1", "J1", "X1"}

// ExtractFileContent opens a file, copying all content into an array splitting
// the lines into an array []string{"line1","line2","line3"...}
func ExtractFileContent(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}
	defer file.Close()

	return text, nil
}

// TextToArray cleans up text files removing all spaces, returning every word
// as indexes of a single array
func TextToArray(text []string) ([]string, []string, error) {
	lines := []string{}
	words := []string{}

	// cleans up each line, getting rid of spaces, only appeding non-blank lines
	for _, line := range text {
		line = strings.TrimSpace(line)

		// cleans up each line, getting rid of spaces, appending only lines
		// contianing bond numbers and addresses
		w := strings.Split(line, " ")
		if len(line) > 0 {
			if FirstWordContains(w[0], BondPrefixes) {
				lines = append(lines, line)
			}
			if FirstWordIsNumber(w[0]) {
				lines = append(lines, line)
			}
		}

		// cleans up each word, getting rid of spaces, only appending non-blank
		// words
		for _, trimedWord := range w {
			strings.TrimSpace(trimedWord)
			if len(trimedWord) > 0 {
				words = append(words, trimedWord)
			}
		}
	}
	// further cleaning of empty space
	lines, err := TrimSpaceFromLines(lines)
	if err != nil {
		log.Fatal(err)
	}
	// puts a delimiter between indexes
	lines, err = toCsv(lines)
	if err != nil {
		log.Fatal(err)
	}
	// apends indexes together
	lines, err = CollapseSlice(lines)
	if err != nil {
		log.Fatal(err)
	}
	// changes a section of the strings delimiter
	lines, err = ChangeDelimiter(lines, 5, "/", ";", " ")
	if err != nil {
		log.Fatal(err)
	}
	// changes a section of the strings delimiter
	lines, err = ChangeAddressDelimiter(lines, ";", " ")
	if err != nil {
		log.Fatal(err)
	}
	return lines, words, nil
}

// TrimSpaceFromLines trims all space between words in slice of strings
func TrimSpaceFromLines(lines []string) ([]string, error) {
	trimedLines := []string{}

	for _, l := range lines {
		tl := strings.Split(l, " ")
		ntl := []string{}

		for _, word := range tl {

			if len(word) > 0 {
				ntl = append(ntl, strings.TrimSpace(word))
			}
		}

		if len(l) > 0 {
			trimedLines = append(trimedLines, strings.Join(ntl, " "))
		}
	}

	return trimedLines, nil
}

// FirstWordContains checks the first word of a line and compairs it to substrings,
// returns false if no substrings match
func FirstWordContains(word string, subString []string) bool {
	return IsMatch(word, subString)
}

// FirstWordIsNumber splits a string and checks if the fist word in that string
// is a number
func FirstWordIsNumber(line string) bool {
	words := strings.Split(line, " ")
	_, err := strconv.Atoi(words[0])
	if err != nil {
		return false
	}
	return true
}

// IsMatch is a wrapper for strings.Contains that checks if any of multiple
// substrings is contained in s
func IsMatch(s string, subStrings []string) bool {

	is := false

	for _, sub := range subStrings {
		is = strings.Contains(s, sub)
		if is {
			return true
		}
	}
	return false
}

// FindIndexWithSubstring returns the index of string that contains substring
func FindIndexWithSubstring(line string, sep string, subString string) (int, error) {
	words := strings.Split(line, sep)

	for i, w := range words {
		if IsMatch(w, []string{subString}) {
			return i, nil
		}

	}
	return 0, errors.New("could not find index")
}

// FindAddressNumberIndex .
func FindAddressNumberIndex(line string, sep string) (int, error) {
	words := strings.Split(line, sep)

	for i := len(words) - 1; i > 0; i-- {
		_, err := strconv.Atoi(words[i])
		if err == nil {
			if i != len(words)-1 {
				return i, nil
			}
		}
	}

	return 0, errors.New("error finding address number index")
}

// ChangeDelimiter replaces delemiters for a slice of the string
func ChangeDelimiter(lines []string, startIndex int, subString string,
	fromSep string, toSep string) ([]string, error) {
	newLines := make([]string, len(lines))

	for i, l := range lines {
		matchedIndex, err := FindIndexWithSubstring(l, fromSep, subString)
		if err != nil {
			return []string{},
				errors.New("error while attempting to change delimiter")
		}

		words := strings.Split(l, fromSep)
		newLines[i] = strings.Join(words[:startIndex+1], fromSep) + " " +
			strings.Join(words[startIndex+1:matchedIndex], toSep) + fromSep +
			strings.Join(words[matchedIndex:], fromSep)
	}

	return newLines, nil
}

// ChangeAddressDelimiter replaces delemiters for a slice of the string
func ChangeAddressDelimiter(lines []string, fromSep string,
	toSep string) ([]string, error) {
	newLines := make([]string, len(lines))

	for i, l := range lines {
		matchedIndex, err := FindAddressNumberIndex(l, fromSep)
		if err != nil {
			return []string{},
				errors.New("error while attempting to change delimiter")
		}

		words := strings.Split(l, fromSep)
		fmt.Println(matchedIndex)
		fmt.Println(words[matchedIndex])
		newLines[i] = strings.Join(words[:matchedIndex], fromSep) + fromSep +
			strings.Join(words[matchedIndex:], toSep)
	}

	return newLines, nil
}

// WriteLinesToFile .
func WriteLinesToFile(lines []string) error {

	f, err := os.Create("./new.csv")
	defer f.Close()

	if err != nil {
		log.Fatal(err)
	}

	for _, l := range lines {
		f.WriteString(l + "\n")
	}
	return nil
}

func toCsv(lines []string) ([]string, error) {
	csv := []string{}

	for _, l := range lines {
		csv = append(csv, strings.Join(strings.Split(l, " "), ";"))
	}

	return csv, nil
}

// CollapseSlice collapses the given array by appending the odd indexes to the
// even. ex. converts []string{"0","1","2","3"} to []string{"01","23"}
func CollapseSlice(lines []string) ([]string, error) {
	collapsed := []string{}
	even := []string{}
	odd := []string{}

	for i, l := range lines {
		switch i % 2 {
		case 0:
			even = append(even, l)
		case 1:
			odd = append(odd, l)
		}
	}

	for i := 0; i < len(even); i++ {
		next := i + 1
		if next < len(even)-1 {
			if len(odd[i+1]) < 60 {
				collapsed = append(collapsed, even[i]+";"+odd[i])
			}
		}
	}

	return collapsed, nil
}

func main() {
	text, err := ExtractFileContent("Attachment_A.txt")
	if err != nil {
		log.Fatal(err)
	}

	lines, _, err := TextToArray(text)
	if err != nil {
		log.Fatal(err)
	}

	WriteLinesToFile(lines)

}
