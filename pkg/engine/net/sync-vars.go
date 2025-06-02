package net

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
	return getTypedVal[bool](v.data, key)
}

func (v *SyncVars) GetInt32(key string) (int32, bool) {
	return getTypedVal[int32](v.data, key)
}

func (v *SyncVars) GetFloat32(key string) (float32, bool) {
	return getTypedVal[float32](v.data, key)
}

func (v *SyncVars) GetString(key string) (string, bool) {
	return getTypedVal[string](v.data, key)
}

func getTypedVal[T any](data map[string]any, key string) (T, bool) {
	val, ok := data[key]
	if !ok {
		return *new(T), false
	}
	casted, ok := val.(T)
	return casted, ok
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
