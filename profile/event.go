package profile

import "github.com/tobyjsullivan/life/signup-svc/store"

const (
    EventType_ProfileCreated = store.EventType("ProfileCreated")
    EventType_NameChanged = store.EventType("NameChanged")
)

func NewProfileCreated(profileId string, expectedVersion uint, firstName, lastName string) *store.Event {
    return &store.Event{
        Version: expectedVersion,
        Type: EventType_ProfileCreated,
        Data: map[string]string {
            "profileId": profileId,
            "firstName": firstName,
            "lastName": lastName,
        },
    }
}

func NewNameChanged(expectedVersion uint, firstName, lastName string) *store.Event {
    return &store.Event{
        Version:expectedVersion,
        Type: EventType_NameChanged,
        Data: map[string]string {
            "firstName": firstName,
            "lastName": lastName,
        },
    }
}
