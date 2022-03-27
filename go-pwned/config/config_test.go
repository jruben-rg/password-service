package config

import (
	"bytes"
	"testing"

	"gopkg.in/yaml.v2"
)

type (
	TestTypeOne struct {
		TestTypeTwo TestTypeTwo `yaml:"config"`
	}

	TestTypeTwo struct {
		Enabled bool   `yaml:"enabled"`
		Value   string `yaml:"value"`
	}
)

func TestShouldUnmarshallYAMLIntoDataType(t *testing.T) {

	tests := []struct {
		scenario  string
		dataTypes TestTypeOne
	}{
		{
			scenario: "Should unmarshall nested files and read configuration for Enabled: false",
			dataTypes: TestTypeOne{
				TestTypeTwo: TestTypeTwo{
					Enabled: false,
					Value:   "NotEnabled",
				},
			},
		},
		{
			scenario: "Should unmarshall nested files and read configuration for Enabled: true",
			dataTypes: TestTypeOne{
				TestTypeTwo: TestTypeTwo{
					Enabled: true,
					Value:   "Enabled",
				},
			},
		},
	}

	for _, test := range tests {

		wantSerializedYml, err := yaml.Marshal(&test.dataTypes)
		if err != nil {
			t.Error("failed to serialize application config")
		}

		result := TestTypeOne{}
		var buffer bytes.Buffer
		buffer.Write(wantSerializedYml)

		err = unmarshallYml(&buffer, &result)
		if err != nil {
			t.Errorf("Scenario '%s' got error: %s.\n", test.scenario, err)
		}

		if result.TestTypeTwo.Enabled != test.dataTypes.TestTypeTwo.Enabled {
			t.Errorf("Scenario: '%s', was expecting enabled to be '%t', got '%t'\n", test.scenario,
				test.dataTypes.TestTypeTwo.Enabled, result.TestTypeTwo.Enabled)
		}

		if result.TestTypeTwo.Value != test.dataTypes.TestTypeTwo.Value {
			t.Errorf("Scenario: '%s', was expecting value to be '%s', got '%s'\n", test.scenario,
				test.dataTypes.TestTypeTwo.Value, result.TestTypeTwo.Value)
		}

	}

}
