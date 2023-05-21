package service

import (
	"context"

	"github.com/ell1jah/linkcompress/internal/http/domain"
	"github.com/ell1jah/linkcompress/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const addr = ":443"

type MicroserviceClienLogger interface {
	Infow(string, ...interface{})
	Errorw(string, ...interface{})
}

type MicroserviceClient struct {
	logger      MicroserviceClienLogger
	ctx         context.Context
	linkService proto.LinkCompresserClient
}

func NewMicroserviceClient(logger MicroserviceClienLogger, ctx context.Context) *MicroserviceClient {
	grcpConn, err := grpc.Dial(
		addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logger.Errorw("MicroserviceClient err",
			"method", "getGrpcConn",
			"calling", "grpc.Dial",
			"addr", addr,
			"error", err)
	}

	go func() {
		<-ctx.Done()
		grcpConn.Close()

		logger.Infow("grpc connection closed",
			"method", "NewMicroserviceClient",
			"addr", addr)
	}()

	linkService := proto.NewLinkCompresserClient(grcpConn)

	return &MicroserviceClient{
		logger:      logger,
		ctx:         ctx,
		linkService: linkService,
	}
}

func (mc *MicroserviceClient) Compress(src domain.Link) (domain.Link, error) {
	resp, err := mc.linkService.Compress(context.Background(), &proto.Link{Body: string(src)})
	if err != nil {
		mc.logger.Errorw("MicroserviceClient got error",
			"method:", "Compress",
			"src link", src,
			"err", err)
		return domain.Link(""), err
	}

	mc.logger.Infow("MicroserviceClient got responce",
		"method:", "Compress",
		"src link", src,
		"resp", resp.Body)

	return domain.Link(resp.Body), nil
}

func (mc *MicroserviceClient) Original(src domain.Link) (domain.Link, error) {
	resp, err := mc.linkService.Original(context.Background(), &proto.Link{Body: string(src)})
	if err != nil {
		mc.logger.Errorw("MicroserviceClient got error",
			"method:", "Original",
			"src link", src,
			"err", err)
		return domain.Link(""), err
	}

	mc.logger.Infow("MicroserviceClient got responce",
		"method:", "Original",
		"src link", src,
		"resp", resp.Body)

	return domain.Link(resp.Body), nil
}
