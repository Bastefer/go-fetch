package service

import (
	"context"
	"errors"
	"log/slog"
	"service-parser/internal/app/domain"
	"service-parser/internal/app/dto"
	"service-parser/internal/app/external"
	"service-parser/internal/app/repository"
	"service-parser/internal/db/wrapper"
	"service-parser/internal/logger/sl"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type FetchService struct {
	log *slog.Logger
	db  *pgxpool.Pool

	sourceClient external.SourceClient

	productRepo       repository.ProductRepository
	clientRepo        repository.ClientRepository
	brandRepo         repository.BrandRepository
	categoryRepo      repository.CategoryRepository
	clientProductRepo repository.ClientProductRepository
	taskRepo          repository.TaskRepository

	productURLs []string
	clientURL   string
	mu          sync.Mutex

	activeTaskID int
	cancel       context.CancelFunc
	done         chan struct{}
}

func NewFetchService(
	log *slog.Logger,
	db *pgxpool.Pool,
	sourceClient external.SourceClient,

	productRepo repository.ProductRepository,
	clientRepo repository.ClientRepository,

	brandRepo repository.BrandRepository,
	categoryRepo repository.CategoryRepository,

	clientProductRepo repository.ClientProductRepository,

	taskRepo repository.TaskRepository,

	productURLs []string,
	clientURL string,
) *FetchService {
	return &FetchService{
		log:          log,
		db:           db,
		sourceClient: sourceClient,

		productRepo: productRepo,
		clientRepo:  clientRepo,

		brandRepo:    brandRepo,
		categoryRepo: categoryRepo,

		clientProductRepo: clientProductRepo,

		taskRepo: taskRepo,

		productURLs: productURLs,
		clientURL:   clientURL,
	}
}

func (s *FetchService) Fetch(ctx context.Context, taskID int) error {

	tx, err := s.db.Begin(ctx)
	if err != nil {
		s.log.Error("fail start transaction", slog.Any(sl.Error, err))
		return err
	}

	defer tx.Rollback(ctx)

	products, clients, err := s.fetchData(ctx)
	if err != nil {
		s.log.Error("fail fetch urls", slog.Any(sl.Error, err))
		return err
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	if err := s.saveProducts(ctx, tx, products); err != nil {
		s.log.Error("fail save products", slog.Any(sl.Error, err))
		return err
	}

	if err := s.saveClients(ctx, tx, clients); err != nil {
		s.log.Error("fail save clients", slog.Any(sl.Error, err))
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		s.log.Error("fail commit transaction", slog.Any(sl.Error, err))
		return err
	}

	return nil
}

func (s *FetchService) fetchData(ctx context.Context) ([]dto.SourceProduct, []dto.SourceClient, error) {

	var (
		wg sync.WaitGroup

		productMu sync.Mutex

		allProducts []dto.SourceProduct

		clients []dto.SourceClient

		errChan = make(chan error, len(s.productURLs)+1)
	)

	for _, url := range s.productURLs {

		wg.Add(1)

		go func() {
			defer wg.Done()

			products, err := s.sourceClient.GetProducts(ctx, url)
			if err != nil {
				errChan <- err
				return
			}

			productMu.Lock()
			allProducts = append(allProducts, products...)
			productMu.Unlock()
		}()
	}

	wg.Add(1)

	go func() {
		defer wg.Done()

		result, err := s.sourceClient.GetClients(ctx, s.clientURL)
		if err != nil {
			errChan <- err
			return
		}

		clients = result
	}()

	wg.Wait()
	close(errChan)

	var errs []error
	for err := range errChan {
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return nil, nil, errors.Join(errs...)
	}

	return allProducts, clients, nil
}

func (s *FetchService) saveProducts(ctx context.Context, tx wrapper.DB, products []dto.SourceProduct) error {
	for _, p := range products {

		if err := ctx.Err(); err != nil {
			return err
		}

		brandID, err := s.brandRepo.GetOrCreate(
			ctx,
			tx,
			p.Brand,
		)
		if err != nil {
			return err
		}

		categoryID, err := s.categoryRepo.GetOrCreate(
			ctx,
			tx,
			p.Category,
		)
		if err != nil {
			return err
		}

		price := parsePrice(p.Price)

		err = s.productRepo.Upsert(
			ctx,
			tx,
			domain.Product{
				ID:         p.ID,
				Name:       p.Name,
				BrandID:    brandID,
				CategoryID: categoryID,
				Price:      price,
				Stock:      p.Stock,
			},
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *FetchService) saveClients(ctx context.Context, tx wrapper.DB, clients []dto.SourceClient) error {
	for _, c := range clients {
		if err := ctx.Err(); err != nil {
			return err
		}
		err := s.clientRepo.Upsert(
			ctx,
			tx,
			domain.Client{
				ID:        c.ID,
				FirstName: c.FirstName,
				LastName:  c.LastName,
			},
		)

		if err != nil {
			return err
		}

		err = s.clientProductRepo.ReplaceProducts(
			ctx,
			tx,
			c.ID,
			c.Products,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *FetchService) StartDownload() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cancel != nil {
		s.cancel()

		if s.done != nil {
			<-s.done
		}
		now := time.Now()

		err := s.taskRepo.UpdateStatus(
			context.Background(),
			s.db,
			s.activeTaskID,
			domain.TaskCompleted,
			&now,
		)
		if err != nil {
			s.log.Error("fail update cancelled task status", slog.Any(sl.Error, err))
		}
	}

	taskID, err := s.taskRepo.Create(
		context.Background(),
		s.db,
	)
	if err != nil {
		s.log.Error("fail create task", slog.Any(sl.Error, err))
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})

	s.cancel = cancel
	s.done = done
	s.activeTaskID = taskID

	go func() {
		defer close(done)

		err = s.Fetch(ctx, taskID)
		if err != nil {

			if errors.Is(err, context.Canceled) {
				s.log.Info(
					"task cancelled",
					slog.Any("task_id", taskID),
				)
			} else {
				s.log.Error(
					"task failed",
					slog.Any(sl.Error, err),
				)
			}
		}
		now := time.Now()

		err := s.taskRepo.UpdateStatus(
			context.Background(),
			s.db,
			taskID,
			domain.TaskCompleted,
			&now,
		)
		if err != nil {
			s.log.Error("fail update cancelled task status", slog.Any(sl.Error, err))
		}
	}()

	return nil
}
