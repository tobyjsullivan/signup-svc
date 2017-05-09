package store

import (
    "log"
    "os"
    "bufio"
    "sync"
    "errors"
    "fmt"
    "bytes"
    "strings"
    "encoding/json"
    "path"
)

var logger *log.Logger
var storesDir string

func init() {
    logger = log.New(os.Stdout, "[store] ", 0)

    wd, err := os.Getwd()
    if err != nil {
        panic(err)
    }

    storesDir = path.Join(wd, "./stores")

    // Ensure the ./stores/ directory exists
    os.MkdirAll(storesDir, 0700)
}

type Store struct {
    file      *os.File
    events    []*Event
    mx        *sync.RWMutex
}

func NewStore(entityId string) (*Store, error) {
    var store Store
    filename := path.Join(storesDir, entityId)

    wrHandle, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
    if err != nil {
        logger.Println("Error while opening store file for append: "+err.Error())
        return nil, err
    }
    store.file = wrHandle

    store.events = make([]*Event, 0)

    // Check if the file exists and, if so, count the lines in it.
    if _, err := os.Stat(filename); err == nil {
        rdHandle, err := os.OpenFile(filename, os.O_RDONLY, 0400)
        if err != nil {
            logger.Println("Error while opening store file for read: "+err.Error())
            return nil, err
        }

        scanner := bufio.NewScanner(rdHandle)
        for scanner.Scan() {
            parsed, err := deserialize(scanner.Text())
            if err != nil {
                logger.Println("Error while deserializing event: "+err.Error())
                return nil, err
            }

            store.events = append(store.events, parsed)
        }

        if err := scanner.Err(); err != nil {
            logger.Println("Error while scanning store file: "+err.Error())
            return nil, err
        }
    }

    store.mx = &sync.RWMutex{}

    return &store, nil
}


func (s *Store) Commit(events []*Event) error {
    s.mx.Lock()
    defer s.mx.Unlock()

    curVer := len(s.events)
    var eventStrings bytes.Buffer
    for i, e := range events {
        if uint(curVer + i) > e.Version {
            return errors.New(fmt.Sprintf("Unexepected version. Expected %d. Actual %d.", e.Version, curVer))
        }

        sJson, err := e.serialize()
        if err != nil {
            logger.Println("Error serializing event: "+err.Error())
            return err
        }

        if _, err := eventStrings.WriteString(strings.TrimSpace(sJson)+"\n"); err != nil {
            logger.Println("Error appending events to buffer: "+err.Error())
            return err
        }
    }

    if _, err := s.file.Write(eventStrings.Bytes()); err != nil {
        logger.Println("Error appending events to store file: "+err.Error())
        return err
    }

    s.events = append(s.events, events...)

    return nil
}

func (s *Store) ReadAll() []*Event {
    s.mx.RLock()
    defer s.mx.RUnlock()

    events := make([]*Event, len(s.events))
    copy(events, s.events)

    return events
}

type EventType string

type Event struct{
    Version uint                 `json:"version"`
    Type    EventType           `json:"type"`
    Data    map[string]string   `json:"data"`
}

func (e *Event) serialize() (string, error) {
    var buf bytes.Buffer
    encoder := json.NewEncoder(&buf)
    if err := encoder.Encode(e); err != nil {
        return "", err
    }

    return buf.String(), nil
}

func deserialize(sJson string) (*Event, error) {
    e := &Event{}
    if err := json.Unmarshal([]byte(sJson), e); err != nil {
        return nil, err
    }
    return e, nil
}