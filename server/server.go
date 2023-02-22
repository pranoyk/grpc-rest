package server

import (
	"context"
	"sync"

	"github.com/gofrs/uuid"
	usersv1 "github.com/pranoyk/grpc-rest/proto/users/v1"
)

// Backend implements the protobuf interface
type Backend struct {
	mu    *sync.RWMutex
	users []*usersv1.User
}

// New initializes a new Backend struct.
func New() *Backend {
	return &Backend{
		mu: &sync.RWMutex{},
	}
}

// AddUser adds a user to the in-memory store.
func (b *Backend) AddUser(ctx context.Context, _ *usersv1.AddUserRequest) (*usersv1.User, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	user := &usersv1.User{
		Id: uuid.Must(uuid.NewV4()).String(),
	}
	b.users = append(b.users, user)

	return user, nil
}

// ListUsers lists all users in the store.
func (b *Backend) ListUsers(_ *usersv1.ListUsersRequest, srv usersv1.UserService_ListUsersServer) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, user := range b.users {
		err := srv.Send(user)
		if err != nil {
			return err
		}
	}

	return nil
}
