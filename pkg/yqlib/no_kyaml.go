//go:build yq_nokyaml

package yqlib

func NewKYamlEncoder(_ KYamlPreferences) Encoder {
	return nil
}
