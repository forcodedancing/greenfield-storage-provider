package bsdb

import (
	"errors"

	"github.com/forbole/juno/v4/common"
	"gorm.io/gorm"
)

// GetPolicyByResourceAndPrincipal get policy by resource type & ID, principal type & value
func (b *BsDBImpl) GetPolicyByResourceAndPrincipal(resourceType, resourceID, principalType, principalValue string) (*Permission, error) {
	var (
		permission *Permission
		err        error
	)

	err = b.db.Table((&Permission{}).TableName()).
		Select("*").
		Where("resource_type = ? and resource_id = ? and principal_type = ? and principal_value = ?", resourceType, resourceID, principalType, principalValue).
		Take(&permission).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return permission, err
}

// GetStatementsByPolicyID get statements info by a policy id
func (b *BsDBImpl) GetStatementsByPolicyID(policyIDList []common.Hash) ([]*Statement, error) {
	var (
		statements []*Statement
		err        error
	)

	err = b.db.Table((&Statement{}).TableName()).
		Select("*").
		Where("policy_id in ?", policyIDList).
		Find(&statements).Error
	return statements, err
}

// GetPermissionsByResourceAndPrincipleType get permission by resource type & ID, principal type
func (b *BsDBImpl) GetPermissionsByResourceAndPrincipleType(resourceType, resourceID, principalType string) ([]*Permission, error) {
	var (
		permissions []*Permission
		err         error
	)

	err = b.db.Table((&Permission{}).TableName()).
		Select("*").
		Where("resource_type = ? and resource_id = ? and principal_type = ?", resourceType, resourceID, principalType).
		Find(&permissions).Error
	return permissions, err
}
