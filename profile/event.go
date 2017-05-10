package profile

import (
    "github.com/tobyjsullivan/life/signup-svc/store2"
)

const (
    EventType_ProfileCreated = store2.EventType("ProfileCreated")
    EventType_NameChanged = store2.EventType("NameChanged")
)

func NewProfileCreated(profileId string, expectedVersion uint, firstName, lastName string) *store2.Event {
    return &store2.Event{
        Version: expectedVersion,
        Type: EventType_ProfileCreated,
        Data: map[string]string {
            "profileId": profileId,
            "firstName": firstName,
            "lastName": lastName,
        },
    }
}

func NewNameChanged(expectedVersion uint, firstName, lastName string) *store2.Event {
    return &store2.Event{
        Version:expectedVersion,
        Type: EventType_NameChanged,
        Data: map[string]string {
            "firstName": firstName,
            "lastName": lastName,
        },
    }
}
