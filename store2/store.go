package store2

import (
    "sync"
    "strings"
    "errors"
    "fmt"
    "log"
    "os"
    "path"
    "io/ioutil"
    "strconv"
    "bufio"
    "encoding/json"
)

var (
    logger *log.Logger

    storesDir string
    stores map[string]*Store
    storesMx sync.RWMutex

    ext = "v2"
)

type Store struct {
    StreamID string
    Events []*Event
    mx *sync.RWMutex
}

type EventType string

type Event struct {
    Type EventType          `json:"type"`
    Version uint            `json:"version"`
    Data map[string]string  `json:"data"`
}

func init() {
    logger = log.New(os.Stdout, "[store2] ", 0)

    stores = make(map[string]*Store)

    pwd, err := os.Getwd()
    if err != nil {
        logger.Panicln("Failed to get working directory: "+err.Error())
    }

    storesDir = path.Join(pwd, "stores")
}

func OpenStore(streamId string) (*Store, error) {
    if strings.Contains(streamId, ".") {
        return nil, errors.New(fmt.Sprintf("Stream ID includes a period. (%s)", streamId))
    }

    storesMx.Lock()
    defer storesMx.Unlock()

    s := stores[streamId]

    if s == nil {
        s = &Store{
            StreamID: streamId,
            Events: make([]*Event, 0),
            mx: &sync.RWMutex{},
        }

        f, err := findLatestFile(streamId)
        if err != nil {
            return nil, err
        }

        if f != nil {
            filePath := path.Join(storesDir, f.Name())
            fHandle, err := os.Open(filePath)
            defer fHandle.Close()
            if err != nil {
                return nil, err
            }

            scanner := bufio.NewScanner(fHandle)
            for scanner.Scan() {
                var e Event
                if err := json.Unmarshal(scanner.Bytes(), &e); err != nil {
                    logger.Println("Error while unmarshalling event: "+err.Error())
                    return nil, err
                }

                s.Events = append(s.Events, &e)
            }
        }

        stores[streamId] = s
    }

    return s, nil
}

func (s *Store) Commit(e *Event) error {
    s.mx.Lock()
    defer s.mx.Unlock()

    nextVersion := uint(len(s.Events))

    if e.Version != nextVersion {
        return errors.New("Event out of sequence")
    }

    prevFile, err := findLatestFile(s.StreamID)
    if err != nil {
        return err
    }

    filename := fmt.Sprintf("%s.%d.%s", s.StreamID, nextVersion, ext)

    fWriter, err := os.OpenFile(path.Join(storesDir, filename), os.O_CREATE|os.O_WRONLY, 0600)
    if err != nil {
        return err
    }
    defer fWriter.Close()

    encoder := json.NewEncoder(fWriter)
    for _, curEvent := range s.Events {
        if wrErr := encoder.Encode(curEvent); wrErr != nil {
            err = wrErr
            break
        }
    }
    if err != nil {
        defer os.Remove(filename)
        return err
    }
    if err := encoder.Encode(e); err != nil {
        defer os.Remove(filename)
        return err
    }

    if prevFile != nil {
        prevFilename := path.Join(storesDir, prevFile.Name())
        defer os.Remove(prevFilename)
    }

    s.Events = append(s.Events, e)
    return nil
}

func (s *Store) ReadAll() []*Event {
    s.mx.RLock()
    defer s.mx.RUnlock()

    out := make([]*Event, len(s.Events))

    copy(out, s.Events)

    return out
}

func findLatestFile(streamId string) (os.FileInfo, error) {
    files, err := ioutil.ReadDir(storesDir)
    if err != nil {
        return nil, err
    }

    var highest uint = 0
    var highestFile os.FileInfo = nil

    for _, f := range files {
        name := f.Name()
        if !strings.HasPrefix(name, streamId) {
            continue
        }

        if !validFilename(name) {
            continue
        }

        version, err := parseVersionNumber(streamId, name)
        if err != nil {
            return nil, err
        }

        if highestFile == nil || version > highest {
            highest = version
            highestFile = f
        }
    }

    return highestFile, nil
}

func validFilename(filename string) bool {
    return strings.HasSuffix(filename, ext)
}

func parseVersionNumber(streamId, filename string) (uint, error) {
    sNum := strings.TrimPrefix(strings.TrimSuffix(filename, "."+ext), streamId+".")
    ui, err := strconv.ParseUint(sNum, 10, 32)
    if err != nil {
        return 0, err
    }
    return uint(ui), nil
}