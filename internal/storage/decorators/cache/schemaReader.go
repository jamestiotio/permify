package cache

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Permify/permify/internal/storage"
	"github.com/Permify/permify/pkg/cache"
	base "github.com/Permify/permify/pkg/pb/base/v1"
)

// SchemaReader - Add cache behaviour to schema reader
type SchemaReader struct {
	delegate storage.SchemaReader
	cache    cache.Cache
}

// NewSchemaReader new instance of SchemaReader
func NewSchemaReader(delegate storage.SchemaReader, cache cache.Cache) *SchemaReader {
	return &SchemaReader{
		delegate: delegate,
		cache:    cache,
	}
}

// ReadSchema  - Read schema from the repository
func (r *SchemaReader) ReadSchema(ctx context.Context, tenantID, version string) (schema *base.SchemaDefinition, err error) {
	return r.delegate.ReadSchema(ctx, tenantID, version)
}

// ReadEntityDefinition - Read entity definition from the repository
func (r *SchemaReader) ReadEntityDefinition(ctx context.Context, tenantID, entityName, version string) (definition *base.EntityDefinition, v string, err error) {
	var s interface{}
	found := false
	if version != "" {
		s, found = r.cache.Get(fmt.Sprintf("%s|%s|%s", tenantID, entityName, version))
	}
	if !found {
		definition, version, err = r.delegate.ReadEntityDefinition(ctx, tenantID, entityName, version)
		if err != nil {
			return nil, "", err
		}
		size := reflect.TypeOf(definition).Size()
		r.cache.Set(fmt.Sprintf("%s|%s|%s", tenantID, entityName, version), definition, int64(size))
		return definition, version, nil
	}
	def, ok := s.(*base.EntityDefinition)
	if !ok {
		return nil, "", errors.New(base.ErrorCode_ERROR_CODE_SCAN.String())
	}
	return def, "", err
}

// ReadRuleDefinition - Read rule definition from the repository
func (r *SchemaReader) ReadRuleDefinition(ctx context.Context, tenantID, ruleName, version string) (definition *base.RuleDefinition, v string, err error) {
	var s interface{}
	found := false
	if version != "" {
		s, found = r.cache.Get(fmt.Sprintf("%s|%s|%s", tenantID, ruleName, version))
	}
	if !found {
		definition, version, err = r.delegate.ReadRuleDefinition(ctx, tenantID, ruleName, version)
		if err != nil {
			return nil, "", err
		}
		size := reflect.TypeOf(definition).Size()
		r.cache.Set(fmt.Sprintf("%s|%s|%s", tenantID, ruleName, version), definition, int64(size))
		return definition, version, nil
	}
	def, ok := s.(*base.RuleDefinition)
	if !ok {
		return nil, "", errors.New(base.ErrorCode_ERROR_CODE_SCAN.String())
	}
	return def, "", err
}

// HeadVersion - Finds the latest version of the schema.
func (r *SchemaReader) HeadVersion(ctx context.Context, tenantID string) (version string, err error) {
	return r.delegate.HeadVersion(ctx, tenantID)
}
