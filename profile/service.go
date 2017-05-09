package profile

import (
    "github.com/tobyjsullivan/life/signup-svc/store"
    "github.com/satori/go.uuid"
    "errors"
    "fmt"
)

type Service struct {
}

func NewService() *Service {
    s := &Service{    }

    return s
}

func (s *Service) CreateProfile(firstName, lastName string) error {
    id := "profile-"+uuid.NewV4().String()

    st, err := store.NewStore(id)
    if err != nil {
        return err
    }

    e := NewProfileCreated(id, 0, firstName, lastName)

    if err := st.Commit([]*store.Event{e}); err != nil {
        return err
    }

    logger.Println(fmt.Sprintf("Profile created. (%s)", id))
    return nil
}

func (s *Service) ChangeName(profileId, firstName, lastName string) error {
    st, err := store.NewStore(profileId)
    if err != nil {
        return err
    }

    aggregate := NewProfileFromHistory(st.ReadAll())

    // Checks
    if aggregate.version == 0 {
        return errors.New("Profile doesn't exist")
    }

    e := NewNameChanged(aggregate.version, firstName, lastName)

    if err := st.Commit([]*store.Event{e}); err != nil {
        return err
    }

    logger.Println(fmt.Sprintf("Name changed. (%s)", profileId))
    return nil
}
