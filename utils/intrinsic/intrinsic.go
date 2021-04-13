package intrinsic

import (
	"fmt"
	"os"
	"sync"
	"io"
	"encoding/json"
	"bytes"
)

var lock sync.Mutex

// CreateDirIfNotExist function to create directory if not exist
func CreateDirIfNotExist(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Errorf("error: %v", err)
		}
	}
}

// Save saves a representation of v to the file at path.
func Save(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	r, err := Marshal(v)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, r)
	return err
}

// Load loads the file at path into v.
func Load(path string, v interface{}) error {
	lock.Lock()
	defer lock.Unlock()
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
		// file does not exist, do something
		return fmt.Errorf("filename: %v; error: %v", path, err)
		} else {
		return fmt.Errorf("error: %v", err)
		}	
	}
	defer f.Close()
	return Unmarshal(f, v)
}

// Marshal is a function that marshals the object into an
// io.Reader.
// By default, it uses the JSON marshaller.
var Marshal = func(v interface{}) (io.Reader, error) {
	b, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

// Unmarshal is a function that unmarshals the data from the
// reader into the specified value.
// By default, it uses the JSON unmarshaller.
var Unmarshal = func(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}
