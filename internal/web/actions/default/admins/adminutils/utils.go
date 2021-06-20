package adminutils

import (
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

// InitAdmin 查找管理员基本信息
func InitAdmin(p *actionutils.ParentAction, adminId int64) error {
	// 管理员信息
	resp, err := p.RPC().AdminRPC().FindEnabledAdmin(p.AdminContext(), &pb.FindEnabledAdminRequest{AdminId: adminId})
	if err != nil {
		return err
	}
	if resp.Admin == nil {
		return errors.New("not found admin")
	}
	var admin = resp.Admin

	// AccessKey数量
	countAccessKeysResp, err := p.RPC().UserAccessKeyRPC().CountAllEnabledUserAccessKeys(p.AdminContext(), &pb.CountAllEnabledUserAccessKeysRequest{
		AdminId: adminId,
		UserId:  0,
	})
	if err != nil {
		return err
	}

	p.Data["admin"] = maps.Map{
		"id":              admin.Id,
		"username":        admin.Username,
		"fullname":        admin.Fullname,
		"countAccessKeys": countAccessKeysResp.Count,
	}
	return nil
}
