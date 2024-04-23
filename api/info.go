package api

import (
	"fmt"
	"time"

	"github.com/gernest/rbf/quantum"
	"github.com/gernest/sql3/dax"
)

const (
	FieldTypeSet       = "set"
	FieldTypeInt       = "int"
	FieldTypeTime      = "time"
	FieldTypeMutex     = "mutex"
	FieldTypeBool      = "bool"
	FieldTypeDecimal   = "decimal"
	FieldTypeTimestamp = "timestamp"
)

// Constants related to timestamp.
const (
	TimeUnitSeconds      = "s"
	TimeUnitMilliseconds = "ms"
	TimeUnitMicroseconds = "µs"
	TimeUnitUSeconds     = "us"
	TimeUnitNanoseconds  = "ns"
)

// IsValidTimeUnit returns true if unit is valid.
func IsValidTimeUnit(unit string) bool {
	switch unit {
	case TimeUnitSeconds, TimeUnitMilliseconds, TimeUnitMicroseconds, TimeUnitUSeconds, TimeUnitNanoseconds:
		return true
	default:
		return false
	}
}

// ValToTimestamp takes a timeunit and an integer value and converts it to time.Time
func ValToTimestamp(unit string, val int64) (time.Time, error) {
	switch unit {
	case TimeUnitSeconds:
		return time.Unix(val, 0).UTC(), nil
	case TimeUnitMilliseconds:
		return time.UnixMilli(val).UTC(), nil
	case TimeUnitMicroseconds, TimeUnitUSeconds:
		return time.UnixMicro(val).UTC(), nil
	case TimeUnitNanoseconds:
		return time.Unix(0, val).UTC(), nil
	default:
		return time.Time{}, fmt.Errorf("unknown time unit: '%v'", unit)
	}
}

type FieldInfo struct {
	Name        string       `json:"name"`
	CreatedAt   int64        `json:"createdAt,omitempty"`
	Owner       string       `json:"owner"`
	Options     FieldOptions `json:"options"`
	Cardinality *uint64      `json:"cardinality,omitempty"`
	Views       []*ViewInfo  `json:"views,omitempty"`
}

type FieldOptions struct {
	Base           int64               `json:"base,omitempty"`
	BitDepth       uint64              `json:"bitDepth,omitempty"`
	Min            int64               `json:"min,omitempty"`
	Max            int64               `json:"max,omitempty"`
	Scale          int64               `json:"scale,omitempty"`
	Keys           bool                `json:"keys"`
	NoStandardView bool                `json:"noStandardView,omitempty"`
	TrackExistence bool                `json:"trackExistence,omitempty"`
	CacheSize      uint32              `json:"cacheSize,omitempty"`
	CacheType      string              `json:"cacheType,omitempty"`
	Type           string              `json:"type,omitempty"`
	TimeUnit       string              `json:"timeUnit,omitempty"`
	TimeQuantum    quantum.TimeQuantum `json:"timeQuantum,omitempty"`
	ForeignIndex   string              `json:"foreignIndex"`
	TTL            time.Duration       `json:"ttl,omitempty"`
}

// ViewInfo represents schema information for a view.
type ViewInfo struct {
	Name string `json:"name"`
}

func FieldToFieldInfo(fld *dax.Field) *FieldInfo {
	var timeUnit string
	var base int64
	min := fld.Options.Min
	max := fld.Options.Max

	return &FieldInfo{
		Name: string(fld.Name),
		Options: FieldOptions{
			Type:           fieldToFieldType(fld),
			Base:           base,
			Min:            min,
			Max:            max,
			Scale:          fld.Options.Scale,
			Keys:           fld.StringKeys(),
			NoStandardView: fld.Options.NoStandardView,
			CacheType:      fld.Options.CacheType,
			CacheSize:      fld.Options.CacheSize,
			TimeUnit:       timeUnit,
			TimeQuantum:    quantum.TimeQuantum(fld.Options.TimeQuantum),
			TTL:            fld.Options.TTL,
			ForeignIndex:   fld.Options.ForeignIndex,
			TrackExistence: fld.Options.TrackExistence,
		},
		Views: nil, // TODO(tlt): do we need views populated?
	}
}

// fieldToFieldType returns the featurebase.FieldType for the given dax.Field.
func fieldToFieldType(f *dax.Field) string {
	switch f.Type {
	case dax.BaseTypeID, dax.BaseTypeString:
		if f.Name == dax.PrimaryKeyFieldName {
			return string(f.Type)
		}
		return "mutex"

	case dax.BaseTypeIDSet, dax.BaseTypeStringSet:
		return "set"

	case dax.BaseTypeIDSetQ, dax.BaseTypeStringSetQ:
		return "time"

	default:
		return string(f.Type)
	}
}
