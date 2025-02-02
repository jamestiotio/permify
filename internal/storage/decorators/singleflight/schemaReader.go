package singleflight

import (
	"context"

	"resenje.org/singleflight"

	"github.com/Permify/permify/internal/storage"
	base "github.com/Permify/permify/pkg/pb/base/v1"
)

// SchemaReader - Add singleflight behaviour to schema reader
type SchemaReader struct {
	delegate storage.SchemaReader
	group    singleflight.Group[string, string]
}

// NewSchemaReader - Add singleflight behaviour to new schema reader
func NewSchemaReader(delegate storage.SchemaReader) *SchemaReader {
	return &SchemaReader{delegate: delegate}
}

// ReadSchema - Read schema from repository
func (r *SchemaReader) ReadSchema(ctx context.Context, tenantID, version string) (*base.SchemaDefinition, error) {
	return r.delegate.ReadSchema(ctx, tenantID, version)
}

// ReadEntityDefinition - Read entity definition from repository
func (r *SchemaReader) ReadEntityDefinition(ctx context.Context, tenantID, entityName, version string) (*base.EntityDefinition, string, error) {
	return r.delegate.ReadEntityDefinition(ctx, tenantID, entityName, version)
}

// ReadRuleDefinition - Read rule definition from repository
func (r *SchemaReader) ReadRuleDefinition(ctx context.Context, tenantID, ruleName, version string) (*base.RuleDefinition, string, error) {
	return r.delegate.ReadRuleDefinition(ctx, tenantID, ruleName, version)
}

// HeadVersion - Finds the latest version of the schema.
func (r *SchemaReader) HeadVersion(ctx context.Context, tenantID string) (version string, err error) {
	rev, _, err := r.group.Do(ctx, "", func(ctx context.Context) (string, error) {
		return r.delegate.HeadVersion(ctx, tenantID)
	})
	return rev, err
}
