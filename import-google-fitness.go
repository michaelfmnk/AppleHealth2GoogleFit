package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"golang.org/x/oauth2"
	fitness "google.golang.org/api/fitness/v1"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	UserID       = "me"
	DateLayout   = "2006-01-02 15:04:05 -0700"
	DataSourceID = "raw:com.google.weight:%s:Mykhailo Fomenko:AppleHealth2GoogleFit:AP2GF"
)

var AppDataSource = &fitness.DataSource{
	Type: "raw",
	Application: &fitness.Application{
		Name: "AppleHealth2GoogleFit",
	},
	DataType: &fitness.DataType{
		Name: "com.google.weight",
		Field: []*fitness.DataTypeField{
			{
				Format: "floatPoint",
				Name:   "weight",
			},
		},
	},
	Device: &fitness.Device{
		Type:         "unknown",
		Manufacturer: "Mykhailo Fomenko",
		Model:        "AppleHealth2GoogleFit",
		Uid:          "AP2GF",
		Version:      "1.0",
	},
}

type ImportGoogleFitness struct {
	ClientID      string `help:"Google Client ID" name:"client-id" type:"string" short:"c"`
	ClientSecret  string `help:"Google Client Secret" name:"client-secret" type:"string" short:"s"`
	ProjectNumber string `help:"Google Project Number" name:"project-number" type:"string" short:"p"`
	Input         string `help:"Input CSV file" name:"input" type:"string" short:"i"`

	client *fitness.Service
}

func (a *ImportGoogleFitness) Run() error {
	client, err := a.createFitnessService()
	if err != nil {
		return err
	}
	a.client = client

	dataSource, err := a.prepareDataSource()
	if err != nil {
		return err
	}

	err = a.InsertData(dataSource)
	if err != nil {
		return err
	}

	return nil
}

func (a *ImportGoogleFitness) prepareDataSource() (*fitness.DataSource, error) {
	result, err := a.createDataSource()
	if shouldRaise(err) {
		return nil, err
	}

	id := fmt.Sprintf(DataSourceID, a.ProjectNumber)
	result, err = a.getDataSource(id)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *ImportGoogleFitness) createFitnessService() (*fitness.Service, error) {
	ctx := context.Background()
	config := &oauth2.Config{
		ClientID:     a.ClientID,
		ClientSecret: a.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://accounts.google.com/o/oauth2/auth",
			TokenURL:  "https://oauth2.googleapis.com/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
		RedirectURL: "https://developers.google.com/oauthplayground",
		Scopes:      []string{fitness.FitnessBodyWriteScope},
	}
	url := config.AuthCodeURL("")
	fmt.Printf("Visit the URL for the auth dialog: %v\n", url)
	fmt.Printf("Enter the code: ")

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return nil, err
	}
	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	client, err := fitness.NewService(ctx, option.WithTokenSource(config.TokenSource(ctx, token)))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func shouldRaise(err error) bool {
	if e, ok := err.(*googleapi.Error); ok {
		if strings.Contains(e.Error(), "alreadyExists") {
			return false
		}
	}
	return true
}

func (a *ImportGoogleFitness) createDataSource() (*fitness.DataSource, error) {
	result, err := a.client.Users.DataSources.Create(UserID, AppDataSource).Do()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *ImportGoogleFitness) getDataSource(id string) (*fitness.DataSource, error) {
	dataSource, err := a.client.Users.DataSources.Get(UserID, id).Do()
	if err != nil {
		return nil, err
	}
	return dataSource, nil
}

func (a *ImportGoogleFitness) InsertData(ds *fitness.DataSource) error {
	csvFile, err := os.Open(a.Input)
	if err != nil {
		return err
	}
	defer func(csvFile *os.File) {
		_ = csvFile.Close()
	}(csvFile)

	records, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		return err
	}

	points, err := a.prepareDataPoints(records, ds)
	if err != nil {
		return err
	}

	datasetID := fmt.Sprintf("%d-%d", points[0].StartTimeNanos, points[len(points)-1].EndTimeNanos)
	_, err = a.client.Users.DataSources.Datasets.Patch(UserID, ds.DataStreamId, datasetID, &fitness.Dataset{
		DataSourceId:   ds.DataStreamId,
		MinStartTimeNs: points[0].StartTimeNanos,
		MaxEndTimeNs:   points[len(points)-1].EndTimeNanos,
		Point:          points,
	}).Do()
	if err != nil {
		return err
	}

	return nil
}

func (a *ImportGoogleFitness) prepareDataPoints(records [][]string, ds *fitness.DataSource) ([]*fitness.DataPoint, error) {
	var points []*fitness.DataPoint
	for i, record := range records {
		if i == 0 {
			continue
		}

		recordTime, err := time.Parse(DateLayout, record[0])
		if err != nil {
			return nil, err
		}
		val, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			continue
		}
		points = append(points, &fitness.DataPoint{
			DataTypeName:   "com.google.weight",
			StartTimeNanos: recordTime.UnixNano(),
			EndTimeNanos:   recordTime.UnixNano(),
			Value: []*fitness.Value{
				{
					FpVal: val,
				},
			},
			OriginDataSourceId: ds.DataStreamId,
		})
	}

	sort.Slice(points, func(i, j int) bool {
		return points[i].StartTimeNanos < points[j].StartTimeNanos
	})
	return points, nil
}
