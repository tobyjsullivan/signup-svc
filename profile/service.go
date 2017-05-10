package profile

import (
    "github.com/satori/go.uuid"
    "errors"
    "fmt"
    "github.com/tobyjsullivan/life/signup-svc/store2"
)

type Service struct {
}

func NewService() *Service {
    s := &Service{}

    return s
}

func (s *Service) CreateProfile(firstName, lastName string) error {
    id := "profile-"+uuid.NewV4().String()

    st, err := store2.OpenStore(id)
    if err != nil {
        return err
    }

    e := NewProfileCreated(id, 0, firstName, lastName)

    if err := st.Commit(e); err != nil {
        logger.Println("Error on commit: "+err.Error())
        return err
    }

    logger.Println(fmt.Sprintf("Profile created. (%s)", id))
    return nil
}

func (s *Service) ChangeName(profileId, firstName, lastName string) error {
    st, err := store2.OpenStore(profileId)
    if err != nil {
        return err
    }

    aggregate := NewProfileFromHistory(st.ReadAll())

    // Checks
    if aggregate.version == 0 {
        return errors.New("Profile doesn't exist")
    }

    e := NewNameChanged(aggregate.version, firstName, lastName)

    if err := st.Commit(e); err != nil {
        logger.Println("Error on commit: "+err.Error())
        return err
    }

    logger.Println(fmt.Sprintf("Name changed. (%s)", profileId))
    return nil
}
