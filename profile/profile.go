package profile

import (
    "fmt"
    "log"
    "os"
    "github.com/tobyjsullivan/life/signup-svc/store"
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

func (p *Profile) apply(event *store.Event) {
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
func (p *Profile) String() string {
    return fmt.Sprintf("(Profile, %s, %s, %s)", p.profileID, p.firstName, p.lastName)
}

func NewProfileFromHistory(events []*store.Event) *Profile {
    p := &Profile{}
    for _, e := range events {
        p.apply(e)
    }
    return p
}
