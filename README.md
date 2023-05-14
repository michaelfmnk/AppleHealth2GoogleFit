# AppleHealth2GoogleFit
[![Go](https://github.com/michaelfmnk/AppleHealth2GoogleFit/actions/workflows/go.yml/badge.svg)](https://github.com/michaelfmnk/AppleHealth2GoogleFit/actions/workflows/go.yml)

This is a command-line interface (CLI) tool that enables you to parse an Apple Health XML file and import its data into Google Fitness using Google's Fitness API. The tool is written in Go and uses the Google APIs Go client library.

## Installation

To use this tool, you need to have Go and Git installed on your computer. Then, you can download the tool by running the following command in your terminal:

```
go get github.com/michaelfmnk/AppleHealth2GoogleFit
```

## Usage

The CLI tool has two commands:

### Parse Apple Health XML file

This command parses an Apple Health XML file and creates a CSV file containing the parsed data. The CSV file can then be used to import the data into Google Fitness. To use this command, run the following command in your terminal:

```
AppleHealth2GoogleFit parse --xml-file [path to Apple Health XML file] --out [path to output CSV file]
```

Replace `[path to Apple Health XML file]` with the path to your Apple Health XML file, and replace `[path to output CSV file]` with the path where you want to save the output CSV file.

### Import Google Fitness data

This command imports data from a CSV file into Google Fitness using Google's Fitness API. To use this command, run the following command in your terminal:

```
AppleHealth2GoogleFit import --client-id [Google Client ID] --client-secret [Google Client Secret] --project-number [Google Project Number] --input [path to input CSV file]
```

Replace `[Google Client ID]` with your Google Client ID, `[Google Client Secret]` with your Google Client Secret, `[Google Project Number]` with your Google Project Number, and `[path to input CSV file]` with the path to your input CSV file.

Note that you need to create a Google Cloud project and enable the Fitness API before you can use this command. You can follow the instructions in Google's [Fitness API Quickstart Guide](https://developers.google.com/fit/rest/v1/get-started) to create a project and enable the Fitness API.


## Contribution

This tool is open-source and contributions are welcome. If you find a bug or have an idea for a new feature, please open an issue or a pull request on the GitHub repository.