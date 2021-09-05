package nagiosplugin

import (
	"fmt"
	"math"
	"strconv"
)

type PerfDatumValue interface {
	fmt.Stringer
}

type FloatPerfDatumValue float64

func NewFloatPerfDatumValue(f float64) (FloatPerfDatumValue, error) {
	if math.IsInf(f, 0) || math.IsNaN(f) {
		return FloatPerfDatumValue(f), fmt.Errorf("Perfdata value may not be infinity or NaN: %v.", f)
	}
	return FloatPerfDatumValue(f), nil
}

func (f FloatPerfDatumValue) String() string {
	return fmtPerfFloat(float64(f))
}

type UndeterminedPerfDatumValue struct{}

const (
	// value may be a literal "U" instead
	// this would indicate that the actual value couldn't be determined
	UndeterminedPerfDatumStr = "U"
)

func NewUndeterminedPerfDatumValue() UndeterminedPerfDatumValue {
	return UndeterminedPerfDatumValue{}
}

func (u UndeterminedPerfDatumValue) String() string {
	return UndeterminedPerfDatumStr
}

// PerfDatum represents one metric to be reported as part of a check
// result.
type PerfDatum struct {
	label string
	value PerfDatumValue
	unit  string
	min   *float64
	max   *float64
	warn  *Range
	crit  *Range
}

// fmtPerfFloat returns a string representation of n formatted in the
// typical /\d+(\.\d+)/ pattern. The difference from %f is that it
// removes any trailing zeroes (like %g except it never returns
// values in scientific notation).
func fmtPerfFloat(n float64) string {
	return strconv.FormatFloat(n, 'f', -1, 64)
}

// NewPerfDatum returns a PerfDatum object suitable to use in a check
// result.
//
// Zero to four thresholds may be supplied: min, max, warn and crit.
// Thresholds may be positive infinity, negative infinity, or NaN, in
// which case they will be omitted in check output.
func NewPerfDatum(label string, unit string, value PerfDatumValue, warn, crit *Range, min, max *float64) (*PerfDatum, error) {
	datum := new(PerfDatum)
	datum.label = label
	datum.value = value
	datum.unit = unit
	datum.warn = warn
	datum.crit = crit
	datum.min = min
	datum.max = max
	return datum, nil
}

// String returns the string representation of a PerfDatum, suitable for
// check output.
func (p PerfDatum) String() string {
	value := fmt.Sprintf("'%s'=%s%s", p.label, p.value.String(), p.unit)

	var warn, crit string
	if p.warn != nil {
		warn = p.warn.String()
	}
	if p.crit != nil {
		crit = p.crit.String()
	}

	var min, max string
	if p.min != nil && !math.IsInf(*p.min, -1) {
		min = fmtPerfFloat(*p.min)
	}
	if p.max != nil && !math.IsInf(*p.max, 1) {
		max = fmtPerfFloat(*p.max)
	}

	value += fmt.Sprintf(";%s;%s", warn, crit)
	value += fmt.Sprintf(";%s;%s", min, max)
	return value
}

// RenderPerfdata accepts a slice of PerfDatum objects and returns their
// concatenated string representations in a form suitable to append to
// the first line of check output.
func RenderPerfdata(perfdata []PerfDatum) string {
	value := ""
	if len(perfdata) == 0 {
		return value
	}
	// Demarcate start of perfdata in check output.
	value += " |"
	for _, datum := range perfdata {
		value += fmt.Sprintf(" %v", datum)
	}
	return value
}
