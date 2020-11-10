package models

import (
	"context"
	"github.com/TeaOSLab/EdgeAdmin/internal/rpc"
	"github.com/TeaOSLab/EdgeCommon/pkg/rpc/pb"
)

var SharedLogDAO = NewLogDAO()

type LogDAO struct {
}

func NewLogDAO() *LogDAO {
	return &LogDAO{}
}

func (this *LogDAO) CreateAdminLog(ctx context.Context, level string, action string, description string, ip string) error {
	client, err := rpc.SharedRPC()
	if err != nil {
		return err
	}
	_, err = client.LogRPC().CreateLog(ctx, &pb.CreateLogRequest{
		Level:       level,
		Description: description,
		Action:      action,
		Ip:          ip,
	})
	return err
}
