package events

type Fetcher interface {
	Fetch(limit int) ([]Event, error)
}

type Processor interface {
	Process(e Event) error
}

type EventsType int

const (
	Unknown EventsType = iota
	Message
)

type Event struct {
	EventsType EventsType
	Text       string
	Meta       interface{}
}
