package service

import (
	"context"
	"time"

	"github.com/micro/go-micro/auth"
	pb "github.com/micro/go-micro/auth/service/proto"
	"github.com/micro/go-micro/client"
)

// svc is the implementation of the Auth interface
type svc struct {
	options auth.Options
	auth    pb.AuthService
}

// Generate a new auth ServiceAccount
func (s *svc) Generate(sa *auth.ServiceAccount) (*auth.ServiceAccount, error) {
	// format the roles and resources
	roles := make([]*pb.Role, len(sa.Roles))
	for i, r := range sa.Roles {
		roles[i] = &pb.Role{
			Name: r.Name,
		}

		if r.Resource != nil {
			roles[i].Resource = &pb.Resource{
				Id:   r.Resource.Id,
				Type: r.Resource.Type,
			}
		}
	}

	// construct the request
	req := &pb.GenerateRequest{
		ServiceAccount: &pb.ServiceAccount{
			Roles:    roles,
			Metadata: sa.Metadata,
			Parent: &pb.Resource{
				Id:   sa.Parent.Id,
				Type: sa.Parent.Type,
			},
		},
	}

	// execute the request
	resp, err := s.auth.Generate(context.Background(), req)
	if err != nil {
		return nil, err
	}

	// format the response
	sa = &auth.ServiceAccount{
		Token:    resp.ServiceAccount.Token,
		Created:  time.Unix(resp.ServiceAccount.Created, 0),
		Expiry:   time.Unix(resp.ServiceAccount.Expiry, 0),
		Metadata: resp.ServiceAccount.Metadata,
	}
	if resp.ServiceAccount.Parent != nil {
		sa.Parent = &auth.Resource{
			Id:   resp.ServiceAccount.Parent.Id,
			Type: resp.ServiceAccount.Parent.Type,
		}
	}

	sa.Roles = make([]*auth.Role, len(resp.ServiceAccount.Roles))
	for i, r := range resp.ServiceAccount.Roles {
		sa.Roles[i] = &auth.Role{
			Name: r.Name,
		}

		if r.Resource != nil {
			sa.Roles[i].Resource = &auth.Resource{
				Id:   r.Resource.Id,
				Type: r.Resource.Type,
			}
		}
	}

	return sa, nil
}

// Revoke an authorization ServiceAccount
func (s *svc) Revoke(sa *auth.ServiceAccount) error {
	// contruct the request
	req := &pb.RevokeRequest{
		ServiceAccount: &pb.ServiceAccount{Token: sa.Token},
	}

	// execute the request
	_, err := s.auth.Revoke(context.Background(), req)
	return err
}

// Validate a service account token
func (s *svc) Validate(token string) (*auth.ServiceAccount, error) {
	return nil, nil
}

// NewAuth returns a new instance of the Auth service
func NewAuth(opts ...auth.Option) auth.Auth {
	options := auth.Options{}

	for _, o := range opts {
		o(&options)
	}

	client := client.DefaultClient
	srv := pb.NewAuthService("go.micro.srv.auth", client)

	return &svc{options, srv}
}