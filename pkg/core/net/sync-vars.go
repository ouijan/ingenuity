package net

import (
	"fmt"
	"strconv"
)

type SyncVars struct {
	data map[string]any
}

func (v *SyncVars) Data() map[string]any {
	return v.data
}

func (v *SyncVars) IsEmpty() bool {
	return len(v.data) == 0
}

func (v *SyncVars) Set(key string, value any) {
	v.data[key] = value
}

func (v *SyncVars) Get(key string) (any, bool) {
	val, ok := v.data[key]
	return val, ok
}

func (v *SyncVars) GetBool(key string) (bool, bool) {
	val, ok := v.Get(key)
	if !ok {
		return false, false
	}
	switch val := val.(type) {
	case bool:
		return val, true
	case int:
		return val > 0, true
	case uint:
		return val > 0, true
	case float32:
		return val > 0, true
	case float64:
		return val > 0, true
	case string:
		strVal, err := strconv.ParseBool(val)
		if err == nil {
			return strVal, true
		}
		return false, false
	case []byte:
		strVal, err := strconv.ParseBool(string(val))
		if err == nil {
			return strVal, true
		}
		return false, false
	}
	return false, false
}

func (v *SyncVars) GetInt32(key string) (int32, bool) {
	val, ok := v.Get(key)
	if !ok {
		return 0, false
	}
	switch val := val.(type) {
	case int:
		return int32(val), true
	case uint:
		return int32(val), true
	case float32:
		return int32(val), true
	case float64:
		return int32(val), true
	case bool:
		if val {
			return 1, true
		}
		return 0, true

	case string:
	case []byte:
		intVal, err := strconv.ParseInt(string(val), 10, 32)
		if err == nil {
			return int32(intVal), true
		}
		return 0, false
	case nil:
		return 0, true
	}
	return 0, false
}

func (v *SyncVars) GetFloat32(key string) (float32, bool) {
	val, ok := v.Get(key)
	if !ok {
		return 0, false
	}
	switch val := val.(type) {
	case int:
		return float32(val), true
	case uint:
		return float32(val), true
	case float32:
		return val, true
	case float64:
		return float32(val), true
	case bool:
		if val {
			return 1.0, true
		}
		return 0.0, true

	case string:
		fVal, err := strconv.ParseFloat(val, 32)
		if err == nil {
			return float32(fVal), true
		}
		return 0.0, false
	case []byte:
		fVal, err := strconv.ParseFloat(string(val), 32)
		if err == nil {
			return float32(fVal), true
		}
		return 0.0, false
	case nil:
		return 0.0, true
	}
	return 0, false
}

func (v *SyncVars) GetString(key string) (string, bool) {
	val, ok := v.Get(key)
	if !ok {
		return "", false
	}
	switch val := val.(type) {
	case string:
		return val, true
	case []byte:
		return string(val), true
	case int:
		return fmt.Sprintf("%d", val), true
	case uint:
		return fmt.Sprintf("%d", val), true
	case float32:
		return fmt.Sprintf("%f", val), true
	case float64:
		return fmt.Sprintf("%f", val), true
	case bool:
		if val {
			return "true", true
		}
		return "false", true
	}
	return "", false
}

func (v *SyncVars) IsEq(o SyncVars) bool {
	if len(v.data) != len(o.data) {
		return false
	}

	for k, v := range v.data {
		if w, ok := o.data[k]; !ok || v != w {
			return false
		}
	}

	return true
}

func NewSyncVars() SyncVars {
	return SyncVars{
		data: map[string]any{},
	}
}

func SyncVarsFromMap(data map[string]any) SyncVars {
	return SyncVars{
		data: data,
	}
}
