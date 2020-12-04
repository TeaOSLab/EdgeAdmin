package userutils

import (
	"errors"
	"github.com/TeaOSLab/EdgeAdmin/internal/web/actions/actionutils"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
	"github.com/iwind/TeaGo/maps"
)

// 查找用户基本信息
func InitUser(p *actionutils.ParentAction, userId int64) error {
	resp, err := p.RPC().UserRPC().FindEnabledUser(p.AdminContext(), &pb.FindEnabledUserRequest{UserId: userId})
	if err != nil {
		return err
	}
	if resp.User == nil {
		return errors.New("not found user")
	}
	p.Data["user"] = maps.Map{
		"id":       userId,
		"fullname": resp.User.Fullname,
		"username": resp.User.Username,
	}
	return nil
}
