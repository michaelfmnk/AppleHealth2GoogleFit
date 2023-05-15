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
		return err
	}
	defer func(outFile *os.File) {
		_ = outFile.Close()
	}(outFile)

	// Create CSV writer
	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	if err := writeContents(file, writer); err != nil {
		_ = os.Remove(a.Out)
		return err
	}

	return nil
}

func writeContents(inputFile *os.File, outWriter *csv.Writer) error {
	err := outWriter.Write([]string{"Date", "Weight"})
	if err != nil {
		return err
	}

	decoder := xml.NewDecoder(inputFile)
	for {
		token, err := decoder.Token()
		if err != nil {
			break
		}

		elem, ok := token.(xml.StartElement)
		if ok && elem.Name.Local == "Record" {
			var hd HealthData
			err := decoder.DecodeElement(&hd, &elem)
			if err != nil {
				return err
			}

			if hd.RecordType == BodyWeightKey {
				// Write weight data to CSV
				_ = outWriter.Write([]string{hd.StartDate, fmt.Sprintf("%f", hd.Value)})
			}
		}
	}
	return nil
}
