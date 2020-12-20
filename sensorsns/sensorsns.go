package sensorsns

import (
	"strconv"

	"github.com/justinast/golang/sensor"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

type SensorSnsNotifier struct {
	session *session.Session
}

func New(region string, credentials *credentials.Credentials) SensorSnsNotifier {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials,
	})
	if err != nil {
		panic(err)
	}

	s := SensorSnsNotifier{session: sess}

	return s
}

func (n SensorSnsNotifier) PublishSensorStateToSns(state sensor.SensorState) {
	ma := map[string]*sns.MessageAttributeValue{
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
	}

	if state.ValueType == "float" {
		ma["state"] = &sns.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(strconv.FormatFloat(state.ValueF, 'f', 1, 64)),
		}
		ma["type"] = &sns.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String("float"),
		}
	} else if state.ValueType == "bool" {
		s := "true"
		if state.ValueB == false {
			s = "false"
		}
		ma["state"] = &sns.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(s),
		}
		ma["type"] = &sns.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String("bool"),
		}
	} else {
		panic("Unknown value type: " + state.ValueType)
	}

	input := &sns.PublishInput{
		MessageAttributes: ma,
		Message:           aws.String("{\"message\":\"Sensor state\"}"),
		TopicArn:          aws.String("arn:aws:sns:eu-west-1:310819670781:HomeSensorIncomingData"),
	}

	_, err := sns.New(n.session).Publish(input)
	if err != nil {
		panic(err)
	}
}
