package app

import (
	"github.com/xeipuuv/gojsonschema"
	"io/ioutil"
)

// Validate inspect struct according json schema.
// ./schemas/trade.json
func (t Trade) Validate() []string {
	errors := []string{}

	documentLoader := gojsonschema.NewGoLoader(t)

	schema, err := ioutil.ReadFile("./schemas/trade.json")
	if err != nil {
		errors = append(errors, err.Error())
		return errors
	}
	schemaLoader := gojsonschema.NewBytesLoader(schema)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		errors = append(errors, err.Error())
		return errors
	}
	if !result.Valid() {
		for _, desc := range result.Errors() {
			errors = append(errors, desc.String())
		}
		return errors
	}

	return nil
}
