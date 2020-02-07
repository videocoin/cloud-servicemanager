package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	svcmgr "github.com/videocoin/cloud-api/servicemanager/v1"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	"github.com/videocoin/cloud-pkg/grpcutil"
	"github.com/videocoin/cloud-pkg/logger"
	"github.com/videocoin/cloud-pkg/tracer"
	"github.com/videocoin/go-service-manager/datastore"
	"github.com/videocoin/go-service-manager/service"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
)

var (
	// ServiceName is the service name.
	ServiceName = "servicemanager"

	// Version is the application version.
	Version = "dev"
)

// Config is the global config.
type Config struct {
	RPCAddr string `default:"0.0.0.0:5000"`
	DBURI   string `default:"root:@tcp(127.0.0.1:3306)/videocoin?charset=utf8&parseTime=True&loc=Local"`
}

func main() {
	logger.Init(ServiceName, Version)
	log := logrus.NewEntry(logrus.New())
	log.Logger.SetReportCaller(true)
	log = logrus.WithFields(logrus.Fields{
		"service": ServiceName,
		"version": Version,
	})

	closer, err := tracer.NewTracer(ServiceName)
	if err != nil {
		log.Info(err.Error())
	} else {
		defer closer.Close()
	}

	cfg := new(Config)
	if err := envconfig.Process(ServiceName, cfg); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stop)

	errgrp, ctx := errgroup.WithContext(ctx)

	healthSrv := health.NewServer()
	var grpcSrv *grpc.Server
	errgrp.Go(func() error {
		DB, err := datastore.Open(cfg.DBURI)
		if err != nil {
			return err
		}
		defer DB.Close()

		grpcSrv = grpc.NewServer(grpcutil.DefaultServerOpts(log)...)
		svcmgr.RegisterServiceManagerServer(grpcSrv, service.NewServer(log, DB))
		healthpb.RegisterHealthServer(grpcSrv, healthSrv)

		healthSrv.SetServingStatus(fmt.Sprintf("grpc.health.v1.%s", ServiceName), healthpb.HealthCheckResponse_SERVING)

		lis, err := net.Listen("tcp", cfg.RPCAddr)
		if err != nil {
			return err
		}
		return grpcSrv.Serve(lis)
	})

	select {
	case <-stop:
		break
	case <-ctx.Done():
		break
	}

	cancel()

	healthSrv.SetServingStatus(fmt.Sprintf("grpc.health.v1.%s", ServiceName), healthpb.HealthCheckResponse_NOT_SERVING)

	if grpcSrv != nil {
		grpcSrv.GracefulStop()
	}

	if err = errgrp.Wait(); err != nil {
		log.Fatal(err)
	}
}
