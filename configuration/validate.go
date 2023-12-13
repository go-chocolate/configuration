package configuration

type Validator interface {
	Validate() error
}
