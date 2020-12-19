package sensorsns

import (
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type SensorSnsNotifier struct {
	session *session.Session
}

type SensorState struct {
	Id          string
	Name        string
	MeasureName string
	Value       float64
	Timestamp   int64
}

func New(region string, credentials *credentials.Credentials) SensorSnsNotifier {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("eu-west-1"),
		Credentials: credentials,
	})
	if err != nil {
		panic(err)
	}

	s := SensorSnsNotifier{session: sess}

	return s
}

func (n SensorSnsNotifier) publishSensorStateToSns(state SensorState) {
	input := &sns.PublishInput{
		MessageAttributes: map[string]*sns.MessageAttributeValue{
			"ts": &sns.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(strconv.FormatInt(state.Timestamp, 10)),
			},
			"sensorId": &sns.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(state.Id),
			},
			"sensorName": &sns.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(state.Name),
			},
			"measureName": &sns.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(state.MeasureName),
			},
			"state": &sns.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String(strconv.FormatFloat(state.Value, 'f', 1, 64)),
			},
			"type": &sns.MessageAttributeValue{
				DataType:    aws.String("String"),
				StringValue: aws.String("float"),
			},
		},
		Message:  aws.String("{\"message\":\"Sensor state\"}"),
		TopicArn: aws.String("arn:aws:sns:eu-west-1:310819670781:HomeSensorIncomingData"),
	}

	_, err := sns.New(n.session).Publish(input)
	if err != nil {
		panic(err)
	}
}
