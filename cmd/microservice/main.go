package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/ell1jah/linkcompress/internal/microservice/repo"
	"github.com/ell1jah/linkcompress/internal/microservice/service"
	"github.com/ell1jah/linkcompress/internal/microservice/transport"
	"github.com/ell1jah/linkcompress/internal/proto"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

const addr = ":50051"
const postgresStor = "postgres"
const inMemoryStor = "inmem"

func main() {
	var storType = flag.String("stor", postgresStor, "storage type flag")
	flag.Parse()
	if *storType != postgresStor && *storType != inMemoryStor {
		panic(fmt.Errorf("wrong storage type flag value"))
	}

	zapLogger := zap.Must(zap.NewDevelopment())
	logger := zapLogger.Sugar()

	var repository service.LinkRepo
	var err error

	if *storType == postgresStor {
		ctx, finish := context.WithCancel(context.Background())
		defer func() {
			finish()
		}()

		repository, err = repo.NewPostgresRepo(logger, ctx)
		if err != nil {
			logger.Fatalw("NewPostgresRepo err",
				"err", err)
			panic(err)
		}
	} else {
		repository = repo.NewInMemoryRepo(logger)
	}

	service := service.NewLinkService(repository, logger)
	transport := transport.NewLinkTransport(service, logger)

	accLog := func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		reply, err := handler(ctx, req)
		if err != nil {
			logger.Errorw("accLog err",
				"calling", "handler",
				"req", req,
				"error", err)

			return reply, err
		}

		timestamp := time.Now().Unix() - start.Unix()

		p, ok := peer.FromContext(ctx)
		if !ok {
			logger.Errorw("accLog err",
				"calling", "peer.FromContext",
				"req", req,
				"error", err)

			return reply, fmt.Errorf("consumer addr not found")
		}

		logger.Infow("accLog",
			"timestamp", timestamp,
			"method", info.FullMethod,
			"peer addr", p.Addr.String(),
			"req", req,
			"reply", reply,
			"error", err)

		return reply, err
	}

	logger.Infow("starting server",
		"addr", addr)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatalw("cant listen",
			"addr", addr,
			"err", err)

		return
	}

	server := grpc.NewServer(grpc.UnaryInterceptor(accLog))

	proto.RegisterLinkCompresserServer(server, transport)

	logger.Infow("server stoped serving",
		"addr", addr,
		"result", server.Serve(lis))
}
