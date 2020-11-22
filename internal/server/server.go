package server

import (
	"context"
	"github.com/maglink/products-fetcher/pkg/messages"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"strconv"
	"time"
)

const requestTimeout = 5 * time.Second

type Server struct {
	messages.UnimplementedProductsFetcherServer
	mongoCollection *mongo.Collection
	cfg             Config
	ctx             context.Context
}

func New(cfg Config, ctx context.Context) *Server {
	s := &Server{
		cfg: cfg,
		ctx: ctx,
	}

	return s
}

func (s *Server) Run() {
	var err error
	s.mongoCollection, err = mongoConnect(s.cfg.Mongo, s.ctx)
	if err != nil {
		log.Fatalf("failed to connect to mongo: %v", err)
	}

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(s.cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("start listening on tcp port", s.cfg.Port)

	grpcSrv := grpc.NewServer()
	messages.RegisterProductsFetcherServer(grpcSrv, s)
	go func() {
		if err := grpcSrv.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	defer grpcSrv.GracefulStop()

	<-s.ctx.Done()
}

func (s *Server) Fetch(ctx context.Context, in *messages.FetchRequest) (*messages.FetchResponse, error) {
	if in.Url == "" {
		return nil, status.Error(codes.InvalidArgument, "url is empty")
	}

	ctx, cancel := context.WithTimeout(s.ctx, requestTimeout)
	defer cancel()

	productEntries, err := FetchCsv(ctx, in.Url)
	if err != nil {
		return nil, errors.Wrap(err, "fetch csv")
	}

	err = mongoUpsertProducts(s.mongoCollection, productEntries, ctx)
	if err != nil {
		return nil, errors.Wrap(err, "upsert products")
	}

	return &messages.FetchResponse{
		Status: messages.Status_OK,
	}, nil
}

func (s *Server) List(ctx context.Context, in *messages.ListRequest) (*messages.ListResponse, error) {
	ctx, cancel := context.WithTimeout(s.ctx, requestTimeout)
	defer cancel()

	list, err := mongoGetProductsList(in, s.mongoCollection, ctx)
	if err != nil {
		return nil, err
	}

	return &messages.ListResponse{
		Status: messages.Status_OK,
		List:   list,
	}, nil
}
