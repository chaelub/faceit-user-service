package models

type (
	EventType uint8

	Event struct {
		Type    EventType
		Message []byte
	}
)

func (et EventType) String() string {
	if et <= 0 {

	}
	s, got := eventType2String[et]
	if et <= 0 || !got {
		return "Unknown"
	}
	return s
}

const (
	UnknownEvent EventType = iota
	NewUser
	UserUpdated
	UserDeleted
)

var (
	eventType2String = map[EventType]string{
		UnknownEvent: "Unknown",
		NewUser:      "NewUser",
		UserUpdated:  "UserUpdated",
		UserDeleted:  "UserDeleted",
	}
)
