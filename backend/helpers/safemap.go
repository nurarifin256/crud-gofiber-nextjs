package helpers

import (
	"encoding/json"

	"github.com/iancoleman/orderedmap"
)

func StructToSafeMap(input interface{}, hiddenFields ...string) (*orderedmap.OrderedMap, error) {
	om := orderedmap.New()
	bytes, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bytes, &om); err != nil {
		return nil, err
	}

	for _, field := range hiddenFields {
		om.Delete(field)
	}

	return om, nil
}
