package repo

import (
	"context"
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const dataSource = "host=localhost port=8080 dbname=payments user=postgres password=postgres"
const driver = "pgx"

type PostgresRepo struct {
	logger RepoLogger
	ctx    context.Context
}

func NewPostgresRepo(logger RepoLogger, ctx context.Context) (*PostgresRepo, error) {
	db, err := sql.Open(driver, dataSource)
	if err != nil {
		logger.Errorw("PostgresRepo err",
			"method", "NewPostgresRepo",
			"calling", "sql.Open",
			"driver", driver,
			"dataSource", dataSource,
			"error", err)
		return nil, err
	}

	logger.Infow("DB connection opened",
		"method", "NewPostgresRepo",
		"driver", driver,
		"dataSource", dataSource)

	go func() {
		<-ctx.Done()
		db.Close()

		logger.Infow("DB connection closed",
			"method", "NewPostgresRepo",
			"driver", driver,
			"dataSource", dataSource)
	}()

	return &PostgresRepo{
		logger: logger,
		ctx:    ctx,
	}, nil
}

// func (pr *PostgresRepo) GetOriginal(compressed domain.Link) (domain.Link, error) {

// }

// func (pr *PostgresRepo) GetCompressed(original domain.Link) (domain.Link, error) {

// }

// func (pr *PostgresRepo) AddCompressed(original domain.Link, compressed domain.Link) error {

// }

// func (pr *PostgresRepo) GetLastCompressed() (domain.Link, error) {

// }
