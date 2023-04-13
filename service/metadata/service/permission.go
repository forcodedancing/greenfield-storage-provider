package service

import (
	"context"
	"time"

	"cosmossdk.io/math"
	chaintypes "github.com/bnb-chain/greenfield/types"
	"github.com/bnb-chain/greenfield/types/resource"
	gnfdresource "github.com/bnb-chain/greenfield/types/resource"
	permtypes "github.com/bnb-chain/greenfield/x/permission/types"
	"github.com/bnb-chain/greenfield/x/storage/keeper"
	storagetypes "github.com/bnb-chain/greenfield/x/storage/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/forbole/juno/v4/common"

	"github.com/bnb-chain/greenfield-storage-provider/model/errors"
	"github.com/bnb-chain/greenfield-storage-provider/pkg/log"
	"github.com/bnb-chain/greenfield-storage-provider/store/bsdb"
)

// VerifyPermission Verify the input items permission.
func (metadata *Metadata) VerifyPermission(ctx context.Context, req *storagetypes.QueryVerifyPermissionRequest) (resp *storagetypes.QueryVerifyPermissionResponse, err error) {
	log.Debugf("VerifyPermission request: operator:%s, bucket_name: %s, object_name: %s, action_type: %d", req.Operator, req.BucketName, req.ObjectName, req.ActionType)
	var (
		operator   sdk.AccAddress
		bucketInfo *bsdb.Bucket
		objectInfo *bsdb.Object
		effect     permtypes.Effect
	)

	ctx = log.Context(ctx, req)

	if req == nil {
		log.CtxErrorw(ctx, "invalid request", "error", err)
		return nil, errors.ErrInvalidParams
	}

	operator, err = sdk.AccAddressFromHexUnsafe(req.Operator)
	if err != nil {
		log.CtxErrorw(ctx, "failed to creates an AccAddress from a HEX-encoded string", "error", err)
		return nil, err
	}
	log.Debugf("AccAddressFromHexUnsafe result: operator: %s", operator.String())

	if req.BucketName == "" {
		log.Errorw("failed to check bucket name", "bucket_name", req.BucketName, "error", err)
		return nil, errors.ErrInvalidBucketName
	}

	bucketInfo, err = metadata.bsDB.GetBucketByName(req.BucketName, true)
	if err != nil || bucketInfo == nil {
		log.CtxErrorw(ctx, "failed to get bucket info", "error", err)
		return nil, errors.ErrNoSuchBucket
	}
	log.Debugf("GetBucketByName result: bucketInfo: %v", bucketInfo)

	if req.ObjectName == "" {
		effect, err = metadata.VerifyBucketPermission(ctx, bucketInfo, operator, req.ActionType, nil)
		if err != nil {
			log.CtxErrorw(ctx, "failed to verify bucket permission", "error", err)
			return nil, err
		}
		log.Debugf("VerifyBucketPermission result: effect: %s", effect.String())
	} else {
		objectInfo, err = metadata.bsDB.GetObjectInfo(req.BucketName, req.ObjectName)
		if err != nil || objectInfo == nil {
			log.CtxErrorw(ctx, "failed to get object info", "error", err)
			return nil, errors.ErrNoSuchObject
		}
		effect = metadata.VerifyObjectPermission(ctx, bucketInfo, objectInfo, operator, req.ActionType)
		log.Debugf("VerifyObjectPermission result: effect: %s", effect.String())
	}

	resp = &storagetypes.QueryVerifyPermissionResponse{Effect: effect}
	log.CtxInfow(ctx, "succeed to get payment by bucket id")
	return resp, nil
}

// VerifyBucketPermission Bucket permissions checks are divided into three steps:
// First, if the bucket is a public bucket and the action is a read-only action, it returns "allow".
// Second, if the operator is the owner of the bucket, it returns "allow", as the owner has the highest permission.
// Third, verify the policy corresponding to the bucket and the operator.
//  1. If the policy is evaluated as "allow", return "allow" to the user.
//  2. If it is evaluated as "deny" or "unspecified", return "deny".
func (metadata *Metadata) VerifyBucketPermission(ctx context.Context, bucketInfo *bsdb.Bucket, operator sdk.AccAddress,
	action permtypes.ActionType, options *permtypes.VerifyOptions) (permtypes.Effect, error) {
	log.Debugf("VerifyBucketPermission request: bucketInfo:%v, operator: %s  action:%s, options: %v", bucketInfo, operator.String(), action.String(), options)
	var (
		owner  sdk.AccAddress
		err    error
		effect permtypes.Effect
	)

	// if bucket is public, anyone can read but can not write it.
	if bucketInfo.Visibility == storagetypes.VISIBILITY_TYPE_PUBLIC_READ.String() && keeper.PublicReadBucketAllowedActions[action] {
		log.Debugf("check bucket visibility and public read actions is true")
		return permtypes.EFFECT_ALLOW, nil
	}

	owner, err = sdk.AccAddressFromHexUnsafe(bucketInfo.Owner.String())
	if err != nil {
		log.CtxErrorw(ctx, "failed to creates an AccAddress from a HEX-encoded string", "error", err)
		return permtypes.EFFECT_DENY, err
	}

	// The owner has full permissions
	if operator.Equals(owner) {
		log.Debugf("The owner has full permissions, operator:%s, owner:%s", operator.String(), owner.String())
		return permtypes.EFFECT_ALLOW, nil
	}

	// verify policy
	effect, err = metadata.VerifyPolicy(ctx, math.NewUintFromBigInt(bucketInfo.BucketID.Big()), gnfdresource.RESOURCE_TYPE_BUCKET, operator, action, options)
	if err != nil {
		log.CtxErrorw(ctx, "failed to verify bucket policy", "error", err)
		return permtypes.EFFECT_DENY, err
	}
	if effect == permtypes.EFFECT_ALLOW {
		return permtypes.EFFECT_ALLOW, nil
	}
	return permtypes.EFFECT_DENY, nil
}

// VerifyObjectPermission verify object permission
func (metadata *Metadata) VerifyObjectPermission(ctx context.Context, bucketInfo *bsdb.Bucket, objectInfo *bsdb.Object,
	operator sdk.AccAddress, action permtypes.ActionType) permtypes.Effect {
	var (
		visibility   bool
		err          error
		ownerAcc     sdk.AccAddress
		opts         *permtypes.VerifyOptions
		bucketEffect permtypes.Effect
		objectEffect permtypes.Effect
	)

	if objectInfo.Visibility == storagetypes.VISIBILITY_TYPE_PUBLIC_READ.String() ||
		(objectInfo.Visibility == storagetypes.VISIBILITY_TYPE_INHERIT.String() && bucketInfo.Visibility == storagetypes.VISIBILITY_TYPE_PUBLIC_READ.String()) {
		visibility = true
	}
	if visibility && keeper.PublicReadObjectAllowedActions[action] {
		return permtypes.EFFECT_ALLOW
	}
	// The owner has full permissions
	ownerAcc, err = sdk.AccAddressFromHexUnsafe(objectInfo.Owner.String())
	if err != nil {
		log.CtxErrorw(ctx, "failed to creates an AccAddress from a HEX-encoded string", "error", err)
		return permtypes.EFFECT_DENY
	}
	log.Debugf("AccAddressFromHexUnsafe result: ownerAcc: %s", ownerAcc.String())
	if ownerAcc.Equals(operator) {
		log.Debugf("ownerAcc.Equals(operator)")
		return permtypes.EFFECT_ALLOW
	}

	// verify policy
	opts = &permtypes.VerifyOptions{
		Resource: chaintypes.NewObjectGRN(objectInfo.BucketName, objectInfo.ObjectName).String(),
	}
	bucketEffect, err = metadata.VerifyPolicy(ctx, math.NewUintFromBigInt(bucketInfo.BucketID.Big()), gnfdresource.RESOURCE_TYPE_BUCKET, operator, action, opts)
	if err != nil || bucketEffect == permtypes.EFFECT_DENY {
		log.CtxErrorw(ctx, "failed to verify object policy", "error", err)
		return permtypes.EFFECT_DENY
	}
	log.Debugf("VerifyPolicy result: bucketEffect: %s", bucketEffect.String())

	objectEffect, err = metadata.VerifyPolicy(ctx, math.NewUintFromBigInt(objectInfo.ObjectID.Big()), gnfdresource.RESOURCE_TYPE_OBJECT, operator, action,
		nil)
	if err != nil || objectEffect == permtypes.EFFECT_DENY {
		log.CtxErrorw(ctx, "failed to verify object policy", "error", err)
		return permtypes.EFFECT_DENY
	}
	log.Debugf("VerifyPolicy result: objectEffect: %s", objectEffect.String())

	if bucketEffect == permtypes.EFFECT_ALLOW || objectEffect == permtypes.EFFECT_ALLOW {
		log.Debugf("bucketEffect == permtypes.EFFECT_ALLOW || objectEffect == permtypes.EFFECT_ALLOW,bucketEffect = %s,objectEffect = %s", bucketEffect, objectEffect)
		return permtypes.EFFECT_ALLOW
	}
	return permtypes.EFFECT_DENY
}

// VerifyPolicy verify policy of permission
func (metadata *Metadata) VerifyPolicy(ctx context.Context, resourceID math.Uint, resourceType resource.ResourceType,
	operator sdk.AccAddress, action permtypes.ActionType, opts *permtypes.VerifyOptions) (permtypes.Effect, error) {
	log.Debugf("VerifyPolicy request: resourceID:%s, resourceType: %s  operator:%s, action: %v, opts:%v", resourceID.String(), resourceType.String(), operator.String(), action.String(), opts)
	var (
		err                    error
		effect                 permtypes.Effect
		allowed                bool
		permission             *bsdb.Permission
		permissions            []*bsdb.Permission
		groups                 []*bsdb.Group
		statements             []*bsdb.Statement
		groupIDMap             = make(map[common.Hash]bool)
		accountPolicyID        = make([]common.Hash, 0)
		policyIDList           = make([]common.Hash, 0)
		filteredPermissionList = make([]*bsdb.Permission, 0)
	)

	// verify policy which grant permission to account
	permission, err = metadata.bsDB.GetPermissionByResourceAndPrincipal(resourceType.String(), resourceID.String(), permtypes.PRINCIPAL_TYPE_GNFD_ACCOUNT.String(), operator.String())
	if err != nil || permission == nil {
		log.CtxErrorw(ctx, "failed to get permission by resource and principal", "error", err)
		return permtypes.EFFECT_DENY, err
	}
	log.Debugf("GetPermissionByResourceAndPrincipal result: permission: %v", permission)

	accountPolicyID = append(accountPolicyID, permission.PolicyID)
	statements, err = metadata.bsDB.GetStatementsByPolicyID(accountPolicyID)
	if err != nil || statements == nil {
		log.CtxErrorw(ctx, "failed to get statements by policy id", "error", err)
		return permtypes.EFFECT_DENY, err
	}
	log.Debugf("GetStatementsByPolicyID result: statements: %v", statements)
	effect = permission.Eval(action, time.Now(), opts, statements)
	log.Debugf("Permission Eval Result: effect:%s", effect.String())
	if effect != permtypes.EFFECT_UNSPECIFIED {
		return effect, nil
	}
	log.Debugf("account permission.Eval result: effect: %s", effect.String())
	// verify policy which grant permission to group
	log.Debugf("GetPermissionsByResourceAndPrincipleType request: %s,%s,%s", resourceType.String(), resourceID.String(), permtypes.PRINCIPAL_TYPE_GNFD_GROUP.String())
	permissions, err = metadata.bsDB.GetPermissionsByResourceAndPrincipleType(resourceType.String(), resourceID.String(), permtypes.PRINCIPAL_TYPE_GNFD_GROUP.String())
	if err != nil || permissions == nil {
		log.CtxErrorw(ctx, "failed to get permission by resource and principle type", "error", err)
		return permtypes.EFFECT_DENY, err
	}
	log.Debugf("GetPermissionsByResourceAndPrincipleType response: %v", permissions)

	groupIDList := make([]common.Hash, len(permissions))
	for i, perm := range permissions {
		groupIDList[i] = common.HexToHash(perm.PrincipalValue)
	}

	// filter group id by account
	groups, err = metadata.bsDB.GetGroupsByGroupIDAndAccount(groupIDList, common.HexToHash(operator.String()))
	if err != nil || groups == nil {
		log.CtxErrorw(ctx, "failed to get groups by group id and account", "error", err)
		return permtypes.EFFECT_DENY, err
	}
	log.Debugf("GetGroupsByGroupIDAndAccount result: group: %v", groups)

	for _, group := range groups {
		groupIDMap[group.GroupID] = true
	}

	for _, perm := range permissions {
		_, ok := groupIDMap[common.HexToHash(perm.PrincipalValue)]
		if ok {
			policyIDList = append(policyIDList, perm.PolicyID)
			filteredPermissionList = append(filteredPermissionList, perm)
		}
	}
	log.Debugf("GetGroupsByGroupIDAndAccount result: len policyIDList: %d, filteredPermissionList: %d", len(policyIDList), len(filteredPermissionList))

	// statements
	statements, err = metadata.bsDB.GetStatementsByPolicyID(policyIDList)
	if err != nil || statements == nil {
		log.CtxErrorw(ctx, "failed to get statements by policy id", "error", err)
		return permtypes.EFFECT_DENY, err
	}
	log.Debugf("GetStatementsByPolicyID result: statements: %v", statements)
	for _, perm := range filteredPermissionList {
		effect = perm.Eval(action, time.Now(), opts, statements)
		log.Debugf("group permission eval result: effect: %s", effect.String())
		if effect != permtypes.EFFECT_UNSPECIFIED {
			if effect == permtypes.EFFECT_ALLOW {
				allowed = true
			} else if effect == permtypes.EFFECT_DENY {
				return permtypes.EFFECT_DENY, nil
			}
		}
	}
	log.Debugf("if allowed %t", allowed)
	if allowed {
		return permtypes.EFFECT_ALLOW, nil
	}
	return permtypes.EFFECT_UNSPECIFIED, nil
}
