package meshctx

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShouldStoreTenantIDInContext(t *testing.T) {
	// Arrange
	tenantID := uuid.New()
	ctx := context.Background()

	// Act
	ctx = WithTenantID(ctx, tenantID)
	result, err := TenantIDFromContext(ctx)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, tenantID, result)
}

func TestShouldReturnErrorWhenTenantIDNotInContext(t *testing.T) {
	// Arrange
	ctx := context.Background()

	// Act
	result, err := TenantIDFromContext(ctx)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, ErrTenantIDMissing, err)
	assert.Equal(t, uuid.Nil, result)
}

func TestShouldStoreStreamPathInContext(t *testing.T) {
	// Arrange
	path := "/tmp/test-stream"
	ctx := context.Background()

	// Act
	ctx = WithStreamPath(ctx, path)
	result, ok := StreamPathFromContext(ctx)

	// Assert
	assert.True(t, ok)
	assert.Equal(t, path, result)
}

func TestShouldReturnFalseWhenStreamPathNotInContext(t *testing.T) {
	// Arrange
	ctx := context.Background()

	// Act
	result, ok := StreamPathFromContext(ctx)

	// Assert
	assert.False(t, ok)
	assert.Empty(t, result)
}

func TestShouldOverwriteExistingTenantID(t *testing.T) {
	// Arrange
	firstID := uuid.New()
	secondID := uuid.New()
	ctx := context.Background()

	// Act
	ctx = WithTenantID(ctx, firstID)
	ctx = WithTenantID(ctx, secondID)
	result, err := TenantIDFromContext(ctx)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, secondID, result)
	assert.NotEqual(t, firstID, result)
}

func TestShouldOverwriteExistingStreamPath(t *testing.T) {
	// Arrange
	firstPath := "/tmp/first"
	secondPath := "/tmp/second"
	ctx := context.Background()

	// Act
	ctx = WithStreamPath(ctx, firstPath)
	ctx = WithStreamPath(ctx, secondPath)
	result, ok := StreamPathFromContext(ctx)

	// Assert
	assert.True(t, ok)
	assert.Equal(t, secondPath, result)
	assert.NotEqual(t, firstPath, result)
}
