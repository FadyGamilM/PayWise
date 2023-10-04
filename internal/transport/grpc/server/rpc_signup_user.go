package server

import (
	"context"
	"paywise/internal/core/dtos"
	"paywise/internal/transport/grpc/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (gs *grpcServer) SignupUser(ctx context.Context, req *pb.SignupRequest) (*pb.SignupResponse, error) {
	// 1. in grpc we don't have to map the request body because its already handled by grpc genereated go code

	// 2. we need to call the signup service which will call the repo
	// TODO => i have to customize the errors to know what to return the the user
	user, err := gs.services.AuthService.Signup(ctx, &dtos.CreateUserDto{
		Username: req.GetUsername(),
		FullName: req.GetFullName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error trying to signup a new user : %v", err)
	}

	// 3. return the response to the client
	return convert_signup_service_layer_response_into_grpc_response(user), nil
}

func convert_signup_service_layer_response_into_grpc_response(slr *dtos.LoginRes) *pb.SignupResponse {
	return &pb.SignupResponse{
		Response: &pb.AuthResponse{
			SessionId:             slr.SessionID.String(),
			AccessToken:           slr.AccessToken,
			AccessTokenExpiresAt:  timestamppb.New(slr.AccessTokenExpiration),
			RefreshToken:          slr.RefreshToken,
			RefreshTokenExpiresAt: timestamppb.New(slr.RefreshTokenExpiration),
			User: &pb.UserResponse{
				Id:       slr.User.ID,
				Username: slr.User.Username,
				FullName: slr.User.FullName,
				Email:    slr.User.Email,
			},
		},
	}

}
