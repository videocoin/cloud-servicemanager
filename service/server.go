package service

import (
	"context"

	"github.com/sirupsen/logrus"
	svcmgr "github.com/videocoin/videocoinapis-admin/videocoin/admin/api/servicemanagement/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	"github.com/videocoin/common/api/resource/project"
	"github.com/videocoin/common/api/resource/service"
	"github.com/videocoin/go-service-manager/datastore"
)

// Server implements the ServiceManagerServer interface.
type Server struct {
	ds     datastore.DataStore
	logger *logrus.Entry
}

// NewServer creates a service manager server.
func NewServer(logger *logrus.Entry, ds datastore.DataStore) *Server {
	return &Server{
		ds:     ds,
		logger: logger,
	}
}

// CreateService creates a new managed service.
func (srv *Server) CreateService(ctx context.Context, req *svcmgr.CreateServiceRequest) (*svcmgr.ManagedService, error) {
	if req.Service == nil {
		return nil, status.Error(codes.InvalidArgument, "CreateServiceRequest.Service required")
	}
	if ok := service.IsValidName(req.Service.ServiceName); !ok {
		return nil, status.Error(codes.InvalidArgument, service.ErrInvalidName.Error())
	}

	svc, err := srv.ds.CreateService(&datastore.Service{
		ID:   uuid.New().String(),
		Name: req.Service.ServiceName,
	})
	if err != nil {
		return nil, err
	}

	return &svcmgr.ManagedService{ServiceName: svc.Name}, nil
}

// ListServices lists managed services.
func (srv *Server) ListServices(ctx context.Context, req *svcmgr.ListServicesRequest) (*svcmgr.ListServicesResponse, error) {
	var (
		svcs []*datastore.Service
		err  error
	)

	if req.ConsumerId != "" {
		if ok := project.IsValidID(req.ConsumerId); !ok {
			return nil, status.Error(codes.InvalidArgument, project.ErrInvalidID.Error())
		}

		svcs, err = srv.ds.ListConsumerServices(req.ConsumerId)
		if err != nil {
			return nil, err
		}
	} else {
		svcs, err = srv.ds.ListServices()
		if err != nil {
			return nil, err
		}
	}

	svcsPB := make([]*svcmgr.ManagedService, 0, len(svcs))
	for _, svc := range svcs {
		svcsPB = append(svcsPB, &svcmgr.ManagedService{
			ServiceName: svc.Name,
		})
	}

	return &svcmgr.ListServicesResponse{Services: svcsPB}, nil
}

// GetService gets a managed service.
func (srv *Server) GetService(ctx context.Context, req *svcmgr.GetServiceRequest) (*svcmgr.ManagedService, error) {
	if ok := service.IsValidName(req.ServiceName); !ok {
		return nil, status.Error(codes.InvalidArgument, service.ErrInvalidName.Error())
	}
	svc, err := srv.ds.GetService(req.ServiceName)
	if err != nil {
		return nil, err
	}
	return &svcmgr.ManagedService{ServiceName: svc.Name}, nil
}

// DeleteService deletes a managed service.
func (srv *Server) DeleteService(ctx context.Context, req *svcmgr.DeleteServiceRequest) (*empty.Empty, error) {
	if ok := service.IsValidName(req.ServiceName); !ok {
		return nil, status.Error(codes.InvalidArgument, service.ErrInvalidName.Error())
	}
	return new(empty.Empty), srv.ds.DeleteService(req.ServiceName)
}

// EnableService enables a service for a project, so it can be used for the
// project.
func (srv *Server) EnableService(ctx context.Context, req *svcmgr.EnableServiceRequest) (*empty.Empty, error) {
	if ok := service.IsValidName(req.ServiceName); !ok {
		return nil, status.Error(codes.InvalidArgument, service.ErrInvalidName.Error())
	}
	if ok := project.IsValidID(req.ConsumerId); !ok {
		return nil, status.Error(codes.InvalidArgument, project.ErrInvalidID.Error())
	}

	return new(empty.Empty), srv.ds.CreateServiceConsumer(req.ServiceName, req.ConsumerId)
}

// DisableService disables a service for a project, so it can no longer be be used for the
// project. It prevents security leaks.
func (srv *Server) DisableService(ctx context.Context, req *svcmgr.DisableServiceRequest) (*empty.Empty, error) {
	if ok := service.IsValidName(req.ServiceName); !ok {
		return nil, status.Error(codes.InvalidArgument, service.ErrInvalidName.Error())
	}
	if ok := project.IsValidID(req.ConsumerId); !ok {
		return nil, status.Error(codes.InvalidArgument, project.ErrInvalidID.Error())
	}

	return new(empty.Empty), srv.ds.DeleteServiceConsumer(req.ServiceName, req.ConsumerId)
}
