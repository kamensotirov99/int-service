package app

import (
	"context"
	pb "int-service/_proto"
	transport_grpc "int-service/grpc"
	"int-service/repository"
	"int-service/service"
	"net"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	logger *logrus.Logger
}

func Initialize(grpcPort string, logger *logrus.Logger) {
	a := App{}
	a.logger = logger

	client, err := a.connectMongo()
	if err != nil {
		panic("Error while connecting to Mongo database")
	}
	a.createGprcServer(grpcPort, client)
}

func (a *App) createGprcServer(serverPort string, client *mongo.Client) {
	service := service.NewSvc(a.logger, repository.NewMongoDB(client))
	grpcServer := transport_grpc.NewSvc(service, a.logger)

	listen, err := net.Listen("tcp", ":"+serverPort)
	if err != nil {
		a.logger.WithError(err).Fatal("Error while starting grpc server")
	}

	s := grpc.NewServer()

	pb.RegisterArticleSvcServer(s, grpcServer)
	pb.RegisterCelebritySvcServer(s, grpcServer)
	pb.RegisterEpisodeSvcServer(s, grpcServer)
	pb.RegisterShowSvcServer(s, grpcServer)
	pb.RegisterGenreSvcServer(s, grpcServer)
	pb.RegisterSeasonSvcServer(s, grpcServer)
	pb.RegisterJournalistSvcServer(s, grpcServer)
	reflection.Register(s)
	a.logger.Info("GRPC server listening on port: " + serverPort)
	err = s.Serve(listen)
	if err != nil {
		a.logger.WithError(err).Fatal("Error while serving grpc server")
	}
}

func (a *App) connectMongo() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI("mongodb://ArgoXInterns:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, errors.Wrap(err, "Error while connecting to Mongo database")
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Error while verifying the connection between client and database")
	}
	return client, nil
}
