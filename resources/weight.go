package resources

import (
	"strconv"
	"strings"
)

var WeightUnitLb int8 = 1
var WeightUnitOz int8 = 2
var WeightUnitG int8 = 3
var WeightUnitKg int8 = 4

var WeightUnits = map[string]int8{
	"lb": WeightUnitLb,
	"oz": WeightUnitOz,
	"g":  WeightUnitG,
	"gr": WeightUnitG,
	"kg": WeightUnitKg,
}

type Weight struct {
	Value     float32
	Unit      int8
	IsDefined bool
}

func (w *Weight) ParseString(raw string) {
	w.Unit = 0
	w.Value = 0
	w.IsDefined = false

	for k, _ := range WeightUnits {
		// Add space before unit if missing
		if strings.Contains(raw, k) && !strings.Contains(raw, " "+k) {

			// Warning for g: Replace only if not kg
			if !strings.Contains(raw, "k"+k) {
				raw = strings.ReplaceAll(raw, k, " "+k)
			}
		}
	}

	parts := strings.Split(strings.ToLower(raw), " ")

	if len(parts[0]) > 0 {
		weight, err := strconv.ParseFloat(parts[0], 16)
		if nil == err {
			w.Value = float32(weight)
			w.IsDefined = true
			w.Unit = WeightUnitKg // Kg by default

			if len(parts) >= 2 {
				found, ok := WeightUnits[parts[1]]
				if ok {
					w.Unit = found
				}
			}
		}
	}
}
