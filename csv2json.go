package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Employee struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Description string `json:"description"`
	Role        string `json:"role"`
	Phone       string `json:"phone"`
}

type inputFile struct {
	filepath string // filepath of the inputfile
	abspath  bool   // filepath fo the destination
}

func main() {
	// Showing useful information when the user enters the --help option
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options] <csvFile>\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}

	fileData, err := getFile()
	fatal(err)

	csv2json(fileData)
	if fileData.abspath == true {
		showAbspath(fileData)
	}
}

func fatal(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func getFile() (inputFile, error) {
	if len(os.Args) < 2 {
		return inputFile{}, errors.New("A filepath argument is required")
	}

	// Defining option flags
	abspath := flag.Bool("abspath", false, "Print out absolute path for the json file")

	// Parse the flag
	flag.Parse()
	// flag.Arg returns non-flag command-line arguments
	fileLocation := flag.Arg(0)

	return inputFile{fileLocation, *abspath}, nil

}

func csv2json(fileData inputFile) {
	fmt.Println("Writing JSON file...")
	// Open file for read
	file, err := os.Open(fileData.filepath)
	fatal(err)
	// Close file
	defer file.Close()

	// CSV reader
	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	// Read all
	records, err := reader.ReadAll()
	fatal(err)

	var emp Employee
	var employees []Employee

	// Get all the records from employee
	for _, rec := range records {
		emp.ID, _ = strconv.Atoi(rec[0])
		emp.FirstName = rec[1]
		emp.LastName = rec[2]
		emp.Email = rec[3]
		emp.Description = rec[4]
		emp.Role = rec[5]
		emp.Phone = rec[6]
		employees = append(employees, emp)
	}

	jsonData, err := json.Marshal(employees)
	fatal(err)

	jsonName := fmt.Sprintf("%s.json", strings.TrimSuffix(filepath.Base(fileData.filepath), ".csv")) // Declare the JSON file name
	jsonFile, err := os.Create(jsonName)                                                             // Create JSON file
	fatal(err)

	defer jsonFile.Close()

	jsonFile.Write(jsonData)
	jsonFile.Close()
	fmt.Println("Completed")
}

func showAbspath(fileData inputFile) {
	// Get file path
	jsonDir := filepath.Dir(fileData.filepath)
	jsonName := fmt.Sprintf("%s.json", strings.TrimSuffix(filepath.Base(fileData.filepath), ".csv")) // Declare the JSON file name
	finalLocation := filepath.Join(jsonDir, jsonName)
	// Get abs location
	abspath, err := filepath.Abs(finalLocation)
	fatal(err)
	fmt.Printf("Your JSON file will be at :\n %v", abspath)
}
