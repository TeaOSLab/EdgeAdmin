package userutils

import (
	"context"
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/TeaOSLab/EdgeCommon/pkg/userconfigs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

var ErrUserNotFound = errors.New("not found user")

// InitUser 查找用户基本信息
func InitUser(p *actionutils.ParentAction, userId int64) error {
	resp, err := p.RPC().UserRPC().FindEnabledUser(p.AdminContext(), &pb.FindEnabledUserRequest{UserId: userId})
	if err != nil {
		return err
	}
	if resp.User == nil {
		return ErrUserNotFound
	}

	// AccessKey数量
	countAccessKeysResp, err := p.RPC().UserAccessKeyRPC().CountAllEnabledUserAccessKeys(p.AdminContext(), &pb.CountAllEnabledUserAccessKeysRequest{
		AdminId: 0,
		UserId:  userId,
	})
	if err != nil {
		return err
	}

	// 是否有实名认证
	hasNewIndividualIdentity, hasNewEnterpriseIdentity, identityTag, err := CheckUserIdentity(p.RPC(), p.AdminContext(), userId)
	if err != nil {
		return err
	}

	p.Data["user"] = maps.Map{
		"id":                       userId,
		"fullname":                 resp.User.Fullname,
		"username":                 resp.User.Username,
		"countAccessKeys":          countAccessKeysResp.Count,
		"hasNewIndividualIdentity": hasNewIndividualIdentity,
		"hasNewEnterpriseIdentity": hasNewEnterpriseIdentity,
		"identityTag":              identityTag,
	}
	return nil
}

// CheckUserIdentity 实名认证信息
func CheckUserIdentity(rpcClient *rpc.RPCClient, ctx context.Context, userId int64) (hasNewIndividualIdentity bool, hasNewEnterpriseIdentity bool, identityTag string, err error) {
	var tags = []string{}

	// 个人
	individualIdentityResp, err := rpcClient.UserIdentityRPC().FindEnabledUserIdentityWithOrgType(ctx, &pb.FindEnabledUserIdentityWithOrgTypeRequest{
		UserId:  userId,
		OrgType: userconfigs.UserIdentityOrgTypeIndividual,
	})
	if err != nil {
		return false, false, "", err
	}
	var individualIdentity = individualIdentityResp.UserIdentity
	hasNewIndividualIdentity = individualIdentity != nil && individualIdentity.Status == userconfigs.UserIdentityStatusSubmitted

	if individualIdentity != nil && individualIdentity.Status == userconfigs.UserIdentityStatusVerified {
		tags = append(tags, "个人")
	}

	// 企业
	enterpriseIdentityResp, err := rpcClient.UserIdentityRPC().FindEnabledUserIdentityWithOrgType(ctx, &pb.FindEnabledUserIdentityWithOrgTypeRequest{
		UserId:  userId,
		OrgType: userconfigs.UserIdentityOrgTypeEnterprise,
	})
	if err != nil {
		return false, false, "", err
	}
	var enterpriseIdentity = enterpriseIdentityResp.UserIdentity
	hasNewEnterpriseIdentity = enterpriseIdentity != nil && enterpriseIdentity.Status == userconfigs.UserIdentityStatusSubmitted

	if enterpriseIdentity != nil && enterpriseIdentity.Status == userconfigs.UserIdentityStatusVerified {
		tags = append(tags, "企业")
	}

	identityTag = strings.Join(tags, "+")

	return
}
