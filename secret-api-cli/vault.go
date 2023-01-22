package secret

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"worachanon226/gophercises/secret-api-cli/cipher"
)

func File(encodingKey, filepath string) *Vault {
	return &Vault{
		encodingKey: encodingKey,
		filepath:    filepath,
		keyVaules:   make(map[string]string),
	}
}

type Vault struct {
	encodingKey string
	filepath    string
	mutex       sync.Mutex
	keyVaules   map[string]string
}

func (v *Vault) loadKeyValues() error {
	f, err := os.Open(v.filepath)
	if err != nil {
		v.keyVaules = make(map[string]string)
		return nil
	}
	defer f.Close()
	r, err := cipher.DecryptReader(v.encodingKey, f)
	if err != nil {
		return err
	}
	err = v.readKeyValues(r)
	if err != nil {
		return err
	}

	return nil
}

func (v *Vault) readKeyValues(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(&v.keyVaules)
}

func (v *Vault) writeKeyValues(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(v.keyVaules)
}

func (v *Vault) saveKeyValues() error {
	var sb strings.Builder
	enc := json.NewEncoder(&sb)
	err := enc.Encode(v.keyVaules)
	if err != nil {
		return err
	}

	encryptedJSON, err := cipher.Encrypt(v.encodingKey, sb.String())
	if err != nil {
		return err
	}

	f, err := os.OpenFile(v.filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer f.Close()
	_, err = fmt.Fprint(f, encryptedJSON)
	if err != nil {
		return err
	}

	return nil
}

func (v *Vault) Get(key string) (string, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	err := v.loadKeyValues()
	if err != nil {
		return "", err
	}
	value, ok := v.keyVaules[key]
	if !ok {
		return "", errors.New("secret: no value for that key")
	}
	return value, nil
}

func (v *Vault) Set(key, value string) error {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	err := v.loadKeyValues()
	if err != nil {
		return err
	}
	v.keyVaules[key] = value
	err = v.saveKeyValues()
	return err
}
