package service

import (
	"context"

	"github.com/bnb-chain/greenfield/types/resource"
	permissiontypes "github.com/bnb-chain/greenfield/x/permission/types"
	"github.com/forbole/juno/v4/common"

	"github.com/bnb-chain/greenfield-storage-provider/pkg/log"
	metatypes "github.com/bnb-chain/greenfield-storage-provider/service/metadata/types"
)

// VerifyPermission Verify the input items permission.
func (metadata *Metadata) VerifyPermission(ctx context.Context, req *metatypes.QueryVerifyPermissionRequest) (resp *metatypes.QueryVerifyPermissionResponse, err error) {

	ctx = log.Context(ctx, req)

	// TODO Verify Bucket
	object, err := metadata.bsDB.GetObjectByObjectAndBucketName(req.ObjectName, req.BucketName)
	if err != nil {
		log.CtxErrorw(ctx, "failed to get object by object name and bucket name", "error", err)
		return
	}
	if object == nil {
		log.CtxError(ctx, "failed to get object by object name and bucket name since this object is not exist")
		return
	}

	// TODO Verify Group
	// 1. get permission list and retrieve principle value(group id) from result
	// input resourceID/resourceType & priciple type = 2(group)
	// verify bucket,
	// VerifyPolicy(ctx sdk.Context, resourceID math.Uint, resourceType resource.ResourceType, operator sdk.AccAddress,
	//		action permtypes.ActionType, opts *permtypes.VerifyOptions) permtypes.Effect
	// effect := k.permKeeper.VerifyPolicy(ctx, bucketInfo.Id, gnfdresource.RESOURCE_TYPE_BUCKET, operator, action, options)
	// objectEffect := k.permKeeper.VerifyPolicy(ctx, objectInfo.Id, gnfdresource.RESOURCE_TYPE_OBJECT, operator, action, nil)
	permissions, err := metadata.bsDB.GetPermissionsByResourceAndPrincipleType(object.ObjectID.String(), resource.RESOURCE_TYPE_OBJECT.String(), permissiontypes.PRINCIPAL_TYPE_GNFD_GROUP.String())
	if err != nil {
		log.CtxErrorw(ctx, "failed to get permission by resource and principle type", "error", err)
		return
	}
	if permissions == nil || len(permissions) == 0 {
		log.CtxError(ctx, "failed to get permission by resource and principle type since the permissions list is empty")
		return
	}

	groupIDList := make([]common.Hash, len(permissions))
	for i, permission := range permissions {
		groupIDList[i] = common.HexToHash(permission.PrincipalValue)
	}

	// filter group id by account
	groups, err := metadata.bsDB.GetGroupsByGroupIDAndAccount(groupIDList, common.HexToHash(req.Operator))
	if err != nil {
		log.CtxErrorw(ctx, "failed to get groups by group id and account", "error", err)
		return
	}
	if groups == nil || len(groups) == 0 {
		log.CtxError(ctx, "failed to get groups by group id and account since the group list is empty")
		return
	}

	groupIDMap := make(map[common.Hash]bool)
	policyIDList := make([]common.Hash, 0)
	for _, group := range groups {
		groupIDMap[group.GroupID] = true
	}
	for _, permission := range permissions {
		_, ok := groupIDMap[common.HexToHash(permission.PrincipalValue)]
		if ok {
			policyIDList = append(policyIDList, permission.PolicyID)
		}
	}

	// statements
	statements, err := metadata.bsDB.GetStatementsByPolicyID(policyIDList)
	if err != nil {
		log.CtxErrorw(ctx, "failed to get statements by policy id", "error", err)
		return
	}
	if statements == nil || len(statements) == 0 {
		log.CtxError(ctx, "failed to get statements by policy id  since the statements list is empty")
		return
	}

	// TODO check the statement

	//resp = &metatypes.QueryVerifyPermissionResponse{Effect: }
	log.CtxInfow(ctx, "succeed to get payment by bucket id")
	return resp, nil
}
