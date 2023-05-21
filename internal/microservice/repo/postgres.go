package repo

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ell1jah/linkcompress/internal/microservice/domain"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const dataSource = "host=db port=5432 dbname=links user=docker password=postgres"
const driver = "pgx"

type PostgresRepo struct {
	logger RepoLogger
	ctx    context.Context
	db     *sql.DB
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
		db:     db,
	}, nil
}

func (pr *PostgresRepo) GetOriginal(compressed domain.Link) (domain.Link, error) {
	var res string
	err := pr.db.QueryRow(`SELECT DISTINCT original FROM pairs WHERE compressed = $1;`, compressed).Scan(&res)
	if errors.Is(err, sql.ErrNoRows) {
		pr.logger.Infow("PostgresRepo log",
			"method", "GetOriginal",
			"src link", compressed,
			"original link", "not found")
		return "", nil
	} else if err != nil {
		pr.logger.Errorw("PostgresRepo err",
			"method", "GetOriginal",
			"src link", compressed,
			"err", err)
		return "", err
	}

	pr.logger.Infow("PostgresRepo log",
		"method", "GetOriginal",
		"src link", compressed,
		"original link", res)

	return domain.Link(res), err
}

func (pr *PostgresRepo) GetCompressed(original domain.Link) (domain.Link, error) {
	var res string
	err := pr.db.QueryRow(`SELECT DISTINCT compressed FROM pairs WHERE original = $1;`, original).Scan(&res)
	if errors.Is(err, sql.ErrNoRows) {
		pr.logger.Infow("PostgresRepo log",
			"method", "GetCompressed",
			"src link", original,
			"compressed link", "not found")
		return "", nil
	} else if err != nil {
		pr.logger.Errorw("PostgresRepo err",
			"method", "GetCompressed",
			"src link", original,
			"err", err)
		return "", err
	}

	pr.logger.Infow("PostgresRepo log",
		"method", "GetCompressed",
		"src link", original,
		"compressed link", res)

	return domain.Link(res), err
}

func (pr *PostgresRepo) AddCompressed(original domain.Link, compressed domain.Link) error {
	_, err := pr.db.Exec(`INSERT INTO pairs(original, compressed) values ($1, $2);`, original, compressed)
	if err != nil {
		pr.logger.Errorw("PostgresRepo err",
			"method", "AddCompressed",
			"original link", original,
			"compressed link", compressed,
			"err", err)

		return err
	}

	pr.logger.Infow("PostgresRepo log",
		"method", "AddCompressed",
		"original link", original,
		"compressed link", compressed)

	return nil
}

func (pr *PostgresRepo) GetLastCompressed() (domain.Link, error) {
	var res string
	err := pr.db.QueryRow(`SELECT compressed FROM pairs ORDER BY id DESC LIMIT 1;`).Scan(&res)
	if errors.Is(err, sql.ErrNoRows) {
		pr.logger.Infow("PostgresRepo log",
			"method", "GetLastCompressed",
			"last compressed link", "not found")
		return "", nil
	} else if err != nil {
		pr.logger.Errorw("PostgresRepo err",
			"method", "GetLastCompressed",
			"err", err)
		return "", err
	}

	pr.logger.Infow("PostgresRepo log",
		"method", "GetCompressed",
		"last compressed link", res)

	return domain.Link(res), err
}
