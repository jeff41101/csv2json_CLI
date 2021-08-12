package main

import (
	"encoding/csv"
	"flag"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func Test_getFile(t *testing.T) {
	tests := []struct {
		name    string    // Name of the test
		want    inputFile // Wanted result
		wantErr bool      // Want error
		osArgs  []string  // Command arguments
	}{
		{"Default parameters", inputFile{"test.csv", false}, false, []string{"cmd", "test.csv"}},
		{"No parameters", inputFile{}, true, []string{"cmd"}},
		{"Enable abspath", inputFile{"test.csv", true}, false, []string{"cmd", "-abspath", "test.csv"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Saving original os.Args
			actualOsArgs := os.Args
			defer func() {
				os.Args = actualOsArgs // Restore
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			}()

			os.Args = tt.osArgs   // Setting the command line
			got, err := getFile() // Run getFile()
			if (err != nil) != tt.wantErr {
				t.Errorf("getFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_csv2json(t *testing.T) {
	fakeData := [][]string{
		{"1", "Tiger", "Chone", "TigerChone@gmail.com", "tall", "manager", "0912897987"},
		{"2", "Tiger", "Chone", "TigerChone@gmail.com", "tall", "manager", "0912897987"},
		{"3", "Tiger", "Chone", "TigerChone@gmail.com", "tall", "manager", "0912897987"},
	}

	file, err := os.Create("test.csv")
	fatal(err)
	defer file.Close()
	writer := csv.NewWriter(file)

	for _, value := range fakeData {
		err := writer.Write(value)
		fatal(err)
	}
	// Write the data into disc
	writer.Flush()

	testFileData := inputFile{
		filepath: "test.csv",
		abspath:  false,
	}
	csv2json(testFileData) // Transforming csv to json
	// Getting the text from JSON file created by the previous func
	testOutput, err := ioutil.ReadFile("test.json")
	fatal(err)
	// Getting the text from the JSON file with expected data
	wantOutput, err := ioutil.ReadFile("testSampleFiles/sample.json")
	fatal(err)
	// Check the want ouput and the test output
	if (string(testOutput)) != (string(wantOutput)) {
		t.Errorf("writeJSONFile() = %v, want %v", string(testOutput), string(wantOutput))
	}

}
