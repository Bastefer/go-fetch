package service

import (
	"context"
	"log/slog"
	"service-parser/internal/app/dto"
	"service-parser/internal/app/repository"
	"service-parser/internal/logger/sl"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StatisticService struct {
	log *slog.Logger
	db  *pgxpool.Pool

	productRepo  repository.ProductRepository
	clientRepo   repository.ClientRepository
	brandRepo    repository.BrandRepository
	categoryRepo repository.CategoryRepository
}

func NewStatisticService(
	log *slog.Logger,
	db *pgxpool.Pool,
	productRepo repository.ProductRepository,
	clientRepo repository.ClientRepository,
	brandRepo repository.BrandRepository,
	categoryRepo repository.CategoryRepository,
) *StatisticService {
	return &StatisticService{
		log:          log,
		db:           db,
		productRepo:  productRepo,
		clientRepo:   clientRepo,
		brandRepo:    brandRepo,
		categoryRepo: categoryRepo,
	}
}
func (s *StatisticService) Stats(ctx context.Context) (*dto.Stats, error) {
	productsCount, err := s.productRepo.GetCount(ctx, s.db)
	if err != nil {
		s.log.Error("fail get productsCount", slog.Any(sl.Error, err))
		return nil, err
	}
	clientsCount, err := s.clientRepo.GetCount(ctx, s.db)
	if err != nil {
		s.log.Error("fail get clientsCount", slog.Any(sl.Error, err))
		return nil, err
	}
	brandsCount, err := s.brandRepo.GetCount(ctx, s.db)
	if err != nil {
		s.log.Error("fail get brandsCount", slog.Any(sl.Error, err))
		return nil, err
	}
	categoriesCount, err := s.categoryRepo.GetCount(ctx, s.db)
	if err != nil {
		s.log.Error("fail get categoriesCount", slog.Any(sl.Error, err))
		return nil, err
	}
	return &dto.Stats{
		Products:   productsCount,
		Clients:    clientsCount,
		Brands:     brandsCount,
		Categories: categoriesCount,
	}, nil
}
