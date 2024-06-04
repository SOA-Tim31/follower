package handlers

import (
	"Rest/data"
	"Rest/domain"
	follower "Rest/proto"
	"context"
	"fmt"
)

type UserHandlergRPC struct {
	repo *data.UserRepository
	follower.UnimplementedFollowerServiceServer
}

func NewgRPCUserHandler(r *data.UserRepository) *UserHandlergRPC {
	return &UserHandlergRPC{repo: r}
}

func (uHandlergRPC *UserHandlergRPC) CreateUser(ctx context.Context, req *follower.UserRequest) (*follower.UserResponse, error) {
	person := domain.User{
		Id: int32(req.Id),
		Username: req.Username,
	}


	err := uHandlergRPC.repo.WriteUser(&person)
	if err != nil {
		fmt.Printf("Database exception: %s\n", err)
		return nil, err
	}

	return &follower.UserResponse{
		MessageConfirmation: "SUCCESS",
	}, nil
}

func (uHandlergRPC *UserHandlergRPC) FollowUser(ctx context.Context, req *follower.FollowRequest) (*follower.FollowResponse, error) {
	var users []domain.User
	for _, username := range req.Username {
		users = append(users, domain.User{
			Username: username,
		})
	}

	if len(users) < 2{
		fmt.Println("Expecting at least 2 users")
		return nil, fmt.Errorf("send 2 users minimum")
	}

	err := uHandlergRPC.repo.FollowUser(&users[0], &users[1])
	if err != nil {
		fmt.Printf("Database exception: %s\n", err)
		return nil, err
	}

	return &follower.FollowResponse{MessageConfirmation: "SUCCESS"}, nil
}