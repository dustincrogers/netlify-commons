package nconf

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

type GenericConfig struct {
	Type   string
	Config map[string]interface{}
}

func (gc *GenericConfig) Load(into interface{}) error {
	err := mapstructure.Decode(&gc.Config, into)
	if err != nil {
		return errors.Wrapf(err, "Failed to parse config")
	}
	return nil
}
