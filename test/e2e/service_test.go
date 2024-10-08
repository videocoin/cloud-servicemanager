package e2e_test

import (
	"context"
	"testing"

	svcmgr "github.com/videocoin/cloud-api/servicemanager/v1"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestIAMService(t *testing.T) {
	conn, err := grpc.Dial(":5000", grpc.WithInsecure())
	require.NoError(t, err)
	require.NotNil(t, conn)
	defer conn.Close()

	cli := svcmgr.NewServiceManagerClient(conn)
	require.NotNil(t, cli)

	ctx := context.Background()
	require.NotNil(t, ctx)

	userID := uuid.New().String()
	projID := userID

	// create service
	svcName := "symphony.videocoin.network"
	req := &svcmgr.CreateServiceRequest{
		Service: &svcmgr.ManagedService{
			ServiceName: "symphony.videocoin.network",
		},
	}

	svc, err := cli.CreateService(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, svc)
	require.Equal(t, svcName, svc.ServiceName)

	// duplicate
	svc2, err := cli.CreateService(ctx, req)
	require.Error(t, err)
	require.Nil(t, svc2)

	// get service
	svc3, err := cli.GetService(ctx, &svcmgr.GetServiceRequest{ServiceName: svcName})
	require.NoError(t, err)
	require.NotNil(t, svc3)
	require.Equal(t, svc, svc3)

	// list services
	svcs, err := cli.ListServices(ctx, &svcmgr.ListServicesRequest{})
	require.NoError(t, err)
	require.NotNil(t, svcs)
	require.NotNil(t, svcs.Services)
	require.Len(t, svcs.Services, 1)
	require.Equal(t, svcs.Services[0], svc)

	// delete service
	empty, err := cli.DeleteService(ctx, &svcmgr.DeleteServiceRequest{ServiceName: svcName})
	require.NoError(t, err)
	require.NotNil(t, empty)

	// list services
	svcs, err = cli.ListServices(ctx, &svcmgr.ListServicesRequest{})
	require.NoError(t, err)
	require.NotNil(t, svcs)
	require.Nil(t, svcs.Services)

	// create service once again
	svc, err = cli.CreateService(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, svc)
	require.Equal(t, svcName, svc.ServiceName)

	// enable service for a consumer
	_, err = cli.EnableService(ctx, &svcmgr.EnableServiceRequest{
		ServiceName: svcName,
		ConsumerId:  projID,
	})
	require.NoError(t, err)

	// list consumer services
	svcs, err = cli.ListServices(ctx, &svcmgr.ListServicesRequest{
		ConsumerId: projID,
	})
	require.NoError(t, err)
	require.NotNil(t, svcs)
	require.NotNil(t, svcs.Services)
	require.Len(t, svcs.Services, 1)
	require.Equal(t, svcs.Services[0], svc)

	// disable service for a consumer
	_, err = cli.DisableService(ctx, &svcmgr.DisableServiceRequest{
		ServiceName: svcName,
		ConsumerId:  projID,
	})
	require.NoError(t, err)

	svcs, err = cli.ListServices(ctx, &svcmgr.ListServicesRequest{
		ConsumerId: projID,
	})
	require.NoError(t, err)
	require.NotNil(t, svcs)
	require.Nil(t, svcs.Services)
}
