package secret

import (
	"errors"

	cipher "worachanon226/gophercises/secret-api-cli/encrypt"
)

func Memory(encodingKey string) Vault {
	return Vault{
		encodingKey: encodingKey,
		keyVaules:   make(map[string]string),
	}
}

type Vault struct {
	encodingKey string
	keyVaules   map[string]string
}

func (v *Vault) Get(key string) (string, error) {
	hex, ok := v.keyVaules[key]
	if !ok {
		return "", errors.New("secret: no value for that key")
	}
	ret, err := cipher.Decrypt(v.encodingKey, hex)
	if err != nil {
		return "", err
	}
	return ret, nil
}

func (v *Vault) Set(key, value string) error {
	encryptedValue, err := cipher.Encrypt(v.encodingKey, value)
	if err != nil {
		return err
	}
	v.keyVaules[key] = encryptedValue
	return nil
}
