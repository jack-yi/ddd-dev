package rpc

import (
	"context"

	"github.com/yangboyi/ddd-dev/user-center/internal/application"
	"github.com/yangboyi/ddd-dev/user-center/proto/pb"
)

type UserCenterServer struct {
	pb.UnimplementedUserCenterServer
	userApp *application.UserApp
}

func NewUserCenterServer(app *application.UserApp) *UserCenterServer {
	return &UserCenterServer{userApp: app}
}

func (s *UserCenterServer) VerifyToken(ctx context.Context, req *pb.VerifyTokenReq) (*pb.VerifyTokenResp, error) {
	claims, err := s.userApp.VerifyToken(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	return &pb.VerifyTokenResp{
		UserId: claims.UserID,
		Email:  claims.Email,
		Name:   claims.Name,
		Roles:  claims.Roles,
	}, nil
}

func (s *UserCenterServer) CheckRole(ctx context.Context, req *pb.CheckRoleReq) (*pb.CheckRoleResp, error) {
	has, err := s.userApp.CheckRole(ctx, req.UserId, req.Role)
	if err != nil {
		return nil, err
	}
	return &pb.CheckRoleResp{HasRole: has}, nil
}

func (s *UserCenterServer) GetUserInfo(ctx context.Context, req *pb.GetUserInfoReq) (*pb.UserInfo, error) {
	user, err := s.userApp.GetUserInfo(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &pb.UserInfo{
		Id:     user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Avatar: user.Avatar,
		Status: user.Status,
		Roles:  user.RoleNames(),
	}, nil
}
