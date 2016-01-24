package models

import (
	"gopkg.in/bluesuncorp/validator.v8"
)

var V = validator.New(&validator.Config{TagName: "v"})
