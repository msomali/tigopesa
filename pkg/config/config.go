package config

type (

	//Validator validates the configurations and return nil if all
	//is good or return an error
	Validator interface {
		Validate() error
	}
)
