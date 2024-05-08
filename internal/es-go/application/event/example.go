package event

type TestEvent struct {
	Id            string
	ImportantData string
}

func (e TestEvent) EventName() string {
	return "TestEvent"
}
