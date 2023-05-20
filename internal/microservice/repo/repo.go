package repo

import (
	"sync"

	"github.com/ell1jah/linkcompress/internal/microservice/domain"
)

type RepoLogger interface {
	Infow(string, ...interface{})
	Errorw(string, ...interface{})
}

type InMemoryRepo struct {
	links          map[domain.Link]domain.Link
	lastCompressed domain.Link
	mu             sync.RWMutex
	logger         RepoLogger
}

func NewInMemoryRepo(logger RepoLogger) *InMemoryRepo {
	return &InMemoryRepo{
		logger: logger,
		links:  map[domain.Link]domain.Link{},
	}
}

func (imr *InMemoryRepo) GetOriginal(compressed domain.Link) (domain.Link, error) {
	imr.mu.RLock()
	defer imr.mu.RUnlock()

	for original, foundCompressed := range imr.links {
		if foundCompressed == compressed {
			return original, nil
		}
	}

	return domain.Link(""), nil
}

func (imr *InMemoryRepo) GetCompressed(original domain.Link) (domain.Link, error) {
	imr.mu.RLock()
	defer imr.mu.RUnlock()

	comp, ok := imr.links[original]
	if ok {
		return comp, nil
	}

	return domain.Link(""), nil
}

func (imr *InMemoryRepo) AddCompressed(original domain.Link, compressed domain.Link) error {
	imr.mu.Lock()
	defer imr.mu.Unlock()

	imr.links[original] = compressed
	imr.lastCompressed = compressed

	return nil
}

func (imr *InMemoryRepo) GetLastCompressed() (domain.Link, error) {
	imr.mu.RLock()
	defer imr.mu.RUnlock()

	return imr.lastCompressed, nil
}
