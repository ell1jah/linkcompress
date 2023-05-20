package transport

import (
	"context"

	"github.com/ell1jah/linkcompress/internal/microservice/domain"
	"github.com/ell1jah/linkcompress/internal/proto"
)

type LinkService interface {
	Compress(domain.Link) (domain.Link, error)
	Original(domain.Link) (domain.Link, error)
}

type LinkTransportLogger interface {
	Debugw(string, ...interface{})
	Infow(string, ...interface{})
	Errorw(string, ...interface{})
}

type LinkTransport struct {
	proto.UnimplementedLinkCompresserServer
	linkService LinkService
	logger      LinkTransportLogger
}

func NewLinkService(linkService LinkService, logger LinkTransportLogger) *LinkTransport {
	return &LinkTransport{
		linkService: linkService,
		logger:      logger,
	}
}

func (lt *LinkTransport) Compress(ctx context.Context, link *proto.Link) (*proto.Link, error) {
	dLink := domain.Link(link.Body)

	resp, err := lt.linkService.Compress(dLink)

	if err != nil {
		lt.logger.Infow("LinkTransport err",
			"method", "Compress",
			"src link", link.Body,
			"error", err)

		return &proto.Link{}, err
	}

	lt.logger.Infow("LinkTransport log",
		"method", "Compress",
		"src link", link.Body,
		"compressed link", resp)

	return &proto.Link{Body: string(resp)}, err
}

func (lt *LinkTransport) Original(ctx context.Context, link *proto.Link) (*proto.Link, error) {
	dLink := domain.Link(link.Body)

	resp, err := lt.linkService.Original(dLink)

	if err != nil {
		lt.logger.Infow("LinkTransport err",
			"method", "Original",
			"src link", link.Body,
			"error", err)

		return &proto.Link{}, err
	}

	lt.logger.Infow("LinkTransport log",
		"method", "Original",
		"src link", link.Body,
		"original link", resp)

	return &proto.Link{Body: string(resp)}, err
}
