package main

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"os"
)

const (
	BodyWeightKey = "HKQuantityTypeIdentifierBodyMass"
)

type ParseAppleHealthXML struct {
	Out     string `help:"Output file" name:"out" type:"path" short:"o"`
	XMLFile string `help:"Apple Health XML file to parse" name:"xml-file" type:"existingfile" short:"i"`
}

type HealthData struct {
	RecordType string  `xml:"type,attr"`
	StartDate  string  `xml:"startDate,attr"`
	EndDate    string  `xml:"endDate,attr"`
	Value      float64 `xml:"value,attr"`
}

func (a *ParseAppleHealthXML) Run() error {
	// Open XML file
	file, err := os.Open(a.XMLFile)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	// Create output file
	outFile, err := os.Create(a.Out)
	if err != nil {
		panic(err)
	}
	defer func(outFile *os.File) {
		_ = outFile.Close()
	}(outFile)

	// Create CSV writer
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	if err := a.writeContents(err, writer, file); err != nil {
		_ = os.Remove(a.Out)
		return err
	}

	return nil
}

func (a *ParseAppleHealthXML) writeContents(err error, writer *csv.Writer, file *os.File) error {
	err = writer.Write([]string{"Date", "Weight"})
	if err != nil {
		return err
	}

	decoder := xml.NewDecoder(file)
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}

		if se, ok := token.(xml.StartElement); ok && se.Name.Local == "Record" {
			var hd HealthData
			err := decoder.DecodeElement(&hd, &se)
			if err != nil {
				return err
			}

			if hd.RecordType == BodyWeightKey {
				// Write weight data to CSV
				_ = writer.Write([]string{hd.StartDate, fmt.Sprintf("%f", hd.Value)})
			}
		}
	}
	return nil
}
