package submit_commons

import "github.com/dchest/uniuri"

const UniqueIdLen = uniuri.StdLen

func GenerateUniqueId() string {
	return uniuri.New()
}
