package sensor

type SensorState struct {
	Timestamp   int64
	Id          string
	Name        string
	MeasureName string
	ValueType   string
	ValueF      float64
	ValueB      bool
}
