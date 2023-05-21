package service

import (
	"fmt"
	"sync"

	"github.com/ell1jah/linkcompress/internal/microservice/domain"
)

const compressedLen = 10

var compressedLinkChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

type LinkRepo interface {
	GetOriginal(domain.Link) (domain.Link, error)
	GetCompressed(domain.Link) (domain.Link, error)
	AddCompressed(domain.Link, domain.Link) error
	GetLastCompressed() (domain.Link, error)
}

type LinkServiceLogger interface {
	Infow(string, ...interface{})
	Errorw(string, ...interface{})
}

type LinkService struct {
	linkRepo LinkRepo
	logger   LinkServiceLogger
	mu       sync.Mutex
}

func NewLinkService(linkRepo LinkRepo, logger LinkServiceLogger) *LinkService {
	return &LinkService{
		linkRepo: linkRepo,
		logger:   logger,
	}
}

func (ls *LinkService) genFirstCompressed() domain.Link {
	res := make([]rune, compressedLen)

	for i := 0; i < compressedLen; i++ {
		res[i] = rune(compressedLinkChars[0])
	}

	return domain.Link(res)
}

func letterI(char byte) (int, error) {
	for i, c := range compressedLinkChars {
		if rune(char) == c {
			return i, nil
		}
	}

	return 0, fmt.Errorf("unknown symbol")
}

func (ls *LinkService) genNewCompressed(last domain.Link) (domain.Link, error) {
	res := make([]rune, compressedLen)

	overflowFlag := true

	for i := compressedLen - 1; i >= 0; i-- {
		if last[i] == compressedLinkChars[len(compressedLinkChars)-1] {
			res[i] = rune(compressedLinkChars[0])
		} else {
			letI, err := letterI(last[i])
			if err != nil {
				return domain.Link(""), err
			}

			res[i] = rune(compressedLinkChars[letI+1])
			overflowFlag = false
			break
		}
	}

	if overflowFlag {
		return domain.Link(""), fmt.Errorf("link overflow")
	}

	return domain.Link(res), nil
}

func (ls *LinkService) Compress(src domain.Link) (domain.Link, error) {
	existing, err := ls.linkRepo.GetCompressed(src)
	if err != nil {
		ls.logger.Errorw("LinkService err",
			"method", "Compress",
			"calling", "ls.linkRepo.GetCompressed",
			"src link", src,
			"error", err)
		return domain.Link(""), err
	} else if existing != "" {
		return existing, nil
	}

	ls.mu.Lock()
	defer ls.mu.Unlock()

	last, err := ls.linkRepo.GetLastCompressed()
	if err != nil {
		ls.logger.Errorw("LinkService err",
			"method", "Compress",
			"calling", "ls.linkRepo.GetLastCompressed",
			"src link", src,
			"error", err)

		return domain.Link(""), err
	}

	var res domain.Link

	if last == "" {
		res = ls.genFirstCompressed()
	} else {
		res, err = ls.genNewCompressed(src)
		if err != nil {
			ls.logger.Errorw("LinkService err",
				"method", "Compress",
				"calling", "ls.genNewCompressed",
				"src link", src,
				"error", err)

			return domain.Link(""), err
		}
	}

	err = ls.linkRepo.AddCompressed(src, res)
	if err != nil {
		ls.logger.Errorw("LinkService err",
			"method", "Compress",
			"calling", "ls.linkRepo.AddCompressed",
			"src link", src,
			"error", err)

		return domain.Link(""), err
	}

	ls.logger.Infow("LinkService log",
		"method", "Compress",
		"src link", src,
		"compressed link", res)
	return res, nil
}

func (ls *LinkService) Original(src domain.Link) (domain.Link, error) {
	res, err := ls.linkRepo.GetOriginal(src)
	if err != nil {
		ls.logger.Errorw("LinkService err",
			"method", "Original",
			"calling", "ls.linkRepo.GetOriginal",
			"src link", src,
			"error", err)

		return domain.Link(""), err
	}

	if res == "" {
		ls.logger.Infow("LinkService log",
			"method", "Original",
			"src link", src,
			"original link", "not found")
		return domain.Link(""), fmt.Errorf("not found")
	}

	ls.logger.Infow("LinkService log",
		"method", "Original",
		"src link", src,
		"original link", res)
	return res, nil
}
