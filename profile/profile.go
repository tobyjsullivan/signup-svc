package profile

import (
    "log"
    "os"
    "github.com/tobyjsullivan/life/signup-svc/store2"
)

var logger *log.Logger

func init() {
    logger = log.New(os.Stdout, "[profile] ", 0)
}

// Aggregates
type Profile struct {
    profileID       string
    version uint
    firstName       string
    lastName        string
}

func (p *Profile) apply(event *store2.Event) {
    switch event.Type {
    case EventType_ProfileCreated:
        p.profileID = event.Data["profileId"]
        p.firstName = event.Data["firstName"]
        p.lastName = event.Data["lastName"]
    case EventType_NameChanged:
        p.firstName = event.Data["firstName"]
        p.lastName = event.Data["lastName"]
    }

    p.version++
}

func NewProfileFromHistory(events []*store2.Event) *Profile {
    p := &Profile{}
    for _, e := range events {
        p.apply(e)
    }
    return p
}
