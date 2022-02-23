package keys

import (
	_ "embed"
	"encoding/json"
)

//go:embed auth.dev.json
var devAuth []byte

func GetDevAuth() (map[string][]byte, error) {
	data := map[string]string{}
	if err := json.Unmarshal(devAuth, &data); err != nil {
		return nil, err
	}

	creds := map[string][]byte{}
	for k, v := range data {
		creds[k] = []byte(v)
	}
	return creds, nil
}
