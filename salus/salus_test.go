package salus

import (
	"os"
	"strconv"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func TestGetTemperature(t *testing.T) {
	salus := New(getCredentials(), getDynamoDB(), 10)

	expTemp, err := strconv.ParseFloat(os.Getenv("EXP_TEMP"), 64)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	salusTemp := salus.GetTemperature()
	if salusTemp != expTemp {
		t.Errorf("Temperature incorrect, got: %f, want: %f.", salusTemp, expTemp)
	}
}

func TestGetSetPoint(t *testing.T) {
	salus := New(getCredentials(), getDynamoDB(), 10)

	expSP, err := strconv.ParseFloat(os.Getenv("EXP_SET_POINT"), 64)
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	setPoint := salus.GetSetPoint()
	if setPoint != expSP {
		t.Errorf("Set point incorrect, got: %f, want: %f.", setPoint, expSP)
	}
}

func TestIsHeating(t *testing.T) {
	salus := New(getCredentials(), getDynamoDB(), 10)

	expIsHeating := false
	if os.Getenv("EXP_IS_HEATING") == "1" {
		expIsHeating = true
	} else if os.Getenv("EXP_IS_HEATING") != "0" {
		panic("Unknown expectation for heating")
	}

	heating := salus.GetIsHeating()
	if heating != expIsHeating {
		t.Errorf("Heater status incorrect")
	}
}

func getCredentials() Credentials {
	return Credentials{Email: os.Getenv("SALUS_EMAIL"), Password: os.Getenv("SALUS_PASSWORD")}
}

func getDynamoDB() *dynamodb.DynamoDB {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("eu-west-1"),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_KEY"), os.Getenv("AWS_SECRET"), ""),
	})
	if err != nil {
		panic(err)
	}

	return dynamodb.New(sess)
}
