package resources

import "testing"

func TestWeight_ParseString(t *testing.T) {
	weight := Weight{}

	raw := "12kg"
	weight.ParseString(raw)
	validWeight(t, weight, raw, 12, WeightUnitKg)

	raw = "33g"
	weight.ParseString(raw)
	validWeight(t, weight, raw, 33, WeightUnitG)

	raw = "0.1oz"
	weight.ParseString(raw)
	validWeight(t, weight, raw, 0.1, WeightUnitOz)

	raw = "123.456 lb"
	weight.ParseString(raw)
	validWeight(t, weight, raw, 123.456, WeightUnitLb)

	raw = "10.000 gr"
	weight.ParseString(raw)
	validWeight(t, weight, raw, 10.0, WeightUnitG)

	raw = "bla"
	weight.ParseString(raw)
	if weight.IsDefined {
		t.Errorf("Weight should not be defined from empty string")
	}
}

func validWeight(t *testing.T, testValue Weight, raw string, expectedValue float32, expectedUnit int8) {
	if !testValue.IsDefined {
		t.Errorf("Weight should be defined from '%s'", raw)
	}

	if testValue.Value != expectedValue {
		t.Errorf("Weight value invalid from '%s'. Expected: %f, current: %f", raw, expectedValue, testValue.Value)
	}

	if testValue.Unit != expectedUnit {
		t.Errorf("Weight unit invalid from '%s'. Expected: %d, current: %d", raw, expectedUnit, testValue.Unit)
	}
}
