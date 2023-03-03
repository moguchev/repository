package domain

import (
	"context"
	"errors"
	"log"

	"github.com/moguchev/repository/internal/app/models"
)

var (
	ErrInvalidItemID     = errors.New("invalid item id")
	ErrInvalidItemFilter = errors.New("invalid item filter")
	ErrNotEnoughItems    = errors.New("not enough items")
	ErrReservationFailed = errors.New("not enough items")
)

type TransactionManager interface {
	RunRepeteableReade(ctx context.Context, f func(ctxTX context.Context) error) error
}

type ItemsRepository interface {
	CreateItem(ctx context.Context, item models.Item) (models.ItemID, error)
	GetItemByID(ctx context.Context, id models.ItemID) (models.Item, error)
	SearchItems(ctx context.Context, filter models.ItemFilter) ([]models.Item, error)

	AddItemStock(ctx context.Context, itemID models.ItemID, count uint32, warehouseID models.WarehouseID) (uint32, error)
	GetItemStocks(ctx context.Context, itemID models.ItemID) (models.ItemStocks, error)
	ReserveStock(ctx context.Context, itemID models.ItemID, count uint32, warehouseID models.WarehouseID) (uint32, error)
}

type Deps struct {
	ItemsRepository
	TransactionManager
}

type ProductService struct {
	Deps
}

func NewProductService(d Deps) *ProductService {
	return &ProductService{d}
}

func (p *ProductService) validateItem(item models.Item) error {
	return nil
}

func (p *ProductService) CreateItem(ctx context.Context, item models.Item) (models.Item, error) {
	if err := p.validateItem(item); err != nil {
		return models.Item{}, nil
	}

	id, err := p.ItemsRepository.CreateItem(ctx, item)
	if err != nil {
		return models.Item{}, nil
	}

	item.ID = id

	return item, nil
}

func (p *ProductService) GetItemByID(ctx context.Context, id models.ItemID) (models.Item, error) {
	if id == 0 {
		return models.Item{}, ErrInvalidItemID
	}

	return p.GetItemByID(ctx, id)
}

func (p *ProductService) validateItemsFilter(item models.ItemFilter) error {
	return nil
}

func (p *ProductService) SearchItems(ctx context.Context, filter models.ItemFilter) ([]models.Item, error) {
	if err := p.validateItemsFilter(filter); err != nil {
		return nil, ErrInvalidItemFilter
	}

	return p.SearchItems(ctx, filter)
}

func (p *ProductService) AddStock(ctx context.Context, itemID models.ItemID, count uint32, warehouseID models.WarehouseID) (uint32, error) {
	// we need transaction !
	if _, err := p.GetItemByID(ctx, itemID); err != nil {
		return 0, err
	}

	newCount, err := p.AddItemStock(ctx, itemID, count, warehouseID)
	if err != nil {
		return 0, err
	}
	return newCount, nil
}

// ReserveStock
func (p *ProductService) ReserveStock(ctx context.Context, itemID models.ItemID, count uint32, warehouseID models.WarehouseID) (models.ItemStocks, error) {
	// some business logic

	var needToReserve = make(models.ItemStocks, 1)

	err := p.TransactionManager.RunRepeteableReade(ctx, func(ctxTX context.Context) error {
		// tx: begin

		stocks, err := p.GetItemStocks(ctxTX, itemID)
		if err != nil {
			return err
		}

		var (
			reservedCount uint32
		)

		for warehouseID, warehouseStock := range stocks {
			left := count - reservedCount
			if left == 0 {
				break
			}
			if warehouseStock >= left {
				needToReserve[warehouseID] = left
				reservedCount += left
			} else {
				needToReserve[warehouseID] = warehouseStock
				reservedCount += warehouseStock
			}
		}

		if reservedCount != count {
			return ErrNotEnoughItems
		}

		for warehouseID, reservation := range needToReserve {
			if _, err := p.ItemsRepository.ReserveStock(ctxTX, itemID, reservation, warehouseID); err != nil {
				return err
			}
		}
		return nil
		// tx: end
	})
	if err != nil {
		log.Println("reservaition failed", err)
		return nil, ErrReservationFailed
	}

	return needToReserve, nil
}
