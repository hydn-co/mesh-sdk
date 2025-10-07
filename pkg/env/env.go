package env

import (
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// Environment variable names used across the project.
const (
	// Azure
	AzureKeyVaultURL = "AZURE_KEY_VAULT_URL"

	// Mesh-specific
	MeshJWTKeySecretName = "MESH_JWT_KEY_SECRET_NAME"
	MeshPortalURL        = "MESH_PORTAL_URL"
	MeshStreamURL        = "MESH_STREAM_URL"
	MeshClientID         = "MESH_CLIENT_ID"
	MeshClientSecret     = "MESH_CLIENT_SECRET"
	MeshSeed             = "MESH_SEED"
)

// ─── String ─────────────────────────────────────────────────────

func GetEnvOrDefaultStr(key string, defaultValue string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return defaultValue
}

func TryGetEnvStr(key string) (string, bool) {
	v, ok := os.LookupEnv(key)
	return v, ok
}

func MustGetEnvStr(key string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	panic("required env var not set: " + key)
}

func TryGetEnvStrSlice(key string) ([]string, bool) {
	if v, ok := os.LookupEnv(key); ok {
		return strings.Split(v, ","), true
	}
	return nil, false
}

// ─── Int ────────────────────────────────────────────────────────

func GetEnvOrDefaultInt(key string, defaultValue int) int {
	if v, ok := TryGetEnvInt(key); ok {
		return v
	}
	return defaultValue
}

func TryGetEnvInt(key string) (int, bool) {
	if v, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(v); err == nil {
			return i, true
		}
	}
	return 0, false
}

func MustGetEnvInt(key string) int {
	if i, ok := TryGetEnvInt(key); ok {
		return i
	}
	panic("invalid or missing int env var: " + key)
}

func TryGetEnvIntSlice(key string) ([]int, bool) {
	raw, ok := os.LookupEnv(key)
	if !ok {
		return nil, false
	}
	parts := strings.Split(raw, ",")
	result := make([]int, 0, len(parts))
	for _, p := range parts {
		if i, err := strconv.Atoi(strings.TrimSpace(p)); err == nil {
			result = append(result, i)
		} else {
			return nil, false
		}
	}
	return result, true
}

// ─── Int32 ──────────────────────────────────────────────────────

func GetEnvOrDefaultInt32(key string, defaultValue int32) int32 {
	if v, ok := TryGetEnvInt32(key); ok {
		return v
	}
	return defaultValue
}

func TryGetEnvInt32(key string) (int32, bool) {
	if v, ok := os.LookupEnv(key); ok {
		if i, err := strconv.ParseInt(v, 10, 32); err == nil {
			return int32(i), true
		}
	}
	return 0, false
}

func MustGetEnvInt32(key string) int32 {
	if i, ok := TryGetEnvInt32(key); ok {
		return i
	}
	panic("invalid or missing int32 env var: " + key)
}

func TryGetEnvInt32Slice(key string) ([]int32, bool) {
	raw, ok := os.LookupEnv(key)
	if !ok {
		return nil, false
	}
	parts := strings.Split(raw, ",")
	result := make([]int32, 0, len(parts))
	for _, p := range parts {
		if i, err := strconv.ParseInt(strings.TrimSpace(p), 10, 32); err == nil {
			result = append(result, int32(i))
		} else {
			return nil, false
		}
	}
	return result, true
}

// ─── Int64 ──────────────────────────────────────────────────────

func GetEnvOrDefaultInt64(key string, defaultValue int64) int64 {
	if v, ok := TryGetEnvInt64(key); ok {
		return v
	}
	return defaultValue
}

func TryGetEnvInt64(key string) (int64, bool) {
	if v, ok := os.LookupEnv(key); ok {
		if i, err := strconv.ParseInt(v, 10, 64); err == nil {
			return i, true
		}
	}
	return 0, false
}

func MustGetEnvInt64(key string) int64 {
	if i, ok := TryGetEnvInt64(key); ok {
		return i
	}
	panic("invalid or missing int64 env var: " + key)
}

func TryGetEnvInt64Slice(key string) ([]int64, bool) {
	raw, ok := os.LookupEnv(key)
	if !ok {
		return nil, false
	}
	parts := strings.Split(raw, ",")
	result := make([]int64, 0, len(parts))
	for _, p := range parts {
		if i, err := strconv.ParseInt(strings.TrimSpace(p), 10, 64); err == nil {
			result = append(result, i)
		} else {
			return nil, false
		}
	}
	return result, true
}

// ─── Float32 ────────────────────────────────────────────────────

func GetEnvOrDefaultFloat32(key string, defaultValue float32) float32 {
	if v, ok := TryGetEnvFloat32(key); ok {
		return v
	}
	return defaultValue
}

func TryGetEnvFloat32(key string) (float32, bool) {
	if v, ok := os.LookupEnv(key); ok {
		if f, err := strconv.ParseFloat(v, 32); err == nil {
			return float32(f), true
		}
	}
	return 0, false
}

func MustGetEnvFloat32(key string) float32 {
	if f, ok := TryGetEnvFloat32(key); ok {
		return f
	}
	panic("invalid or missing float32 env var: " + key)
}

func TryGetEnvFloat32Slice(key string) ([]float32, bool) {
	raw, ok := os.LookupEnv(key)
	if !ok {
		return nil, false
	}
	parts := strings.Split(raw, ",")
	result := make([]float32, 0, len(parts))
	for _, p := range parts {
		if f, err := strconv.ParseFloat(strings.TrimSpace(p), 32); err == nil {
			result = append(result, float32(f))
		} else {
			return nil, false
		}
	}
	return result, true
}

// ─── Float64 ────────────────────────────────────────────────────

func GetEnvOrDefaultFloat64(key string, defaultValue float64) float64 {
	if v, ok := TryGetEnvFloat64(key); ok {
		return v
	}
	return defaultValue
}

func TryGetEnvFloat64(key string) (float64, bool) {
	if v, ok := os.LookupEnv(key); ok {
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f, true
		}
	}
	return 0, false
}

func MustGetEnvFloat64(key string) float64 {
	if f, ok := TryGetEnvFloat64(key); ok {
		return f
	}
	panic("invalid or missing float64 env var: " + key)
}

func TryGetEnvFloat64Slice(key string) ([]float64, bool) {
	raw, ok := os.LookupEnv(key)
	if !ok {
		return nil, false
	}
	parts := strings.Split(raw, ",")
	result := make([]float64, 0, len(parts))
	for _, p := range parts {
		if f, err := strconv.ParseFloat(strings.TrimSpace(p), 64); err == nil {
			result = append(result, f)
		} else {
			return nil, false
		}
	}
	return result, true
}

// ─── Bool ───────────────────────────────────────────────────────

func GetEnvOrDefaultBool(key string, defaultValue bool) bool {
	if v, ok := TryGetEnvBool(key); ok {
		return v
	}
	return defaultValue
}

func TryGetEnvBool(key string) (bool, bool) {
	if v, ok := os.LookupEnv(key); ok {
		if b, err := strconv.ParseBool(v); err == nil {
			return b, true
		}
	}
	return false, false
}

func MustGetEnvBool(key string) bool {
	if b, ok := TryGetEnvBool(key); ok {
		return b
	}
	panic("invalid or missing bool env var: " + key)
}

func TryGetEnvBoolSlice(key string) ([]bool, bool) {
	raw, ok := os.LookupEnv(key)
	if !ok {
		return nil, false
	}
	parts := strings.Split(raw, ",")
	result := make([]bool, 0, len(parts))
	for _, p := range parts {
		if b, err := strconv.ParseBool(strings.TrimSpace(p)); err == nil {
			result = append(result, b)
		} else {
			return nil, false
		}
	}
	return result, true
}

// ─── UUID ───────────────────────────────────────────────────────

func GetEnvOrDefaultUUID(key string, defaultValue uuid.UUID) uuid.UUID {
	if v, ok := TryGetEnvUUID(key); ok {
		return v
	}
	return defaultValue
}

func TryGetEnvUUID(key string) (uuid.UUID, bool) {
	if v, ok := os.LookupEnv(key); ok {
		if id, err := uuid.Parse(v); err == nil {
			return id, true
		}
	}
	return uuid.Nil, false
}

func MustGetEnvUUID(key string) uuid.UUID {
	if id, ok := TryGetEnvUUID(key); ok {
		return id
	}
	panic("invalid or missing UUID env var: " + key)
}

func TryGetEnvUUIDSlice(key string) ([]uuid.UUID, bool) {
	raw, ok := os.LookupEnv(key)
	if !ok {
		return nil, false
	}
	parts := strings.Split(raw, ",")
	result := make([]uuid.UUID, 0, len(parts))
	for _, p := range parts {
		if id, err := uuid.Parse(strings.TrimSpace(p)); err == nil {
			result = append(result, id)
		} else {
			return nil, false
		}
	}
	return result, true
}
