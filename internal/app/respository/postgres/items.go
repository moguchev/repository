package respository

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/moguchev/repository/internal/app/models"
	"github.com/moguchev/repository/internal/app/respository/postgres/transactor"
	"github.com/moguchev/repository/internal/app/respository/schema"
)

type ItemsRepo struct {
	transactor.QueryEngineProvider
}

func NewItemsRepo(provider transactor.QueryEngineProvider) *ItemsRepo {
	return &ItemsRepo{
		QueryEngineProvider: provider,
	}
}

var (
	itemsColumns = []string{"id", "name", "price", "tags", "length", "width", "height"}
)

const (
	itemsTable = "items"
)

func (r *ItemsRepo) SearchItems(ctx context.Context, filter models.ItemFilter) ([]models.Item, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	query := sq.Select(itemsColumns...).
		From(itemsTable).
		Where("name LIKE ?", fmt.Sprint("%", filter.Name, "%"))

	if filter.ByTags.Use {
		query = query.Where("tags && ?", filter.ByTags.Tags) // WE NEED CREATE GIN INDEX FOR TABLE items(tags)
	}
	if filter.ByPrice.Use {
		query = query.Where(
			sq.GtOrEq{"price": filter.ByPrice.From})
		if max := filter.ByPrice.To; max != 0 {
			query = query.Where(
				sq.LtOrEq{"price": max})
		}
	}

	rawQuery, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	// делаем SELECT в пару строк без циклов и сканирования)
	var items []schema.Item
	if err := pgxscan.Select(ctx, db, &items, rawQuery, args...); err != nil {
		return nil, err
	}

	return bindSchemaItemsToModelsItems(items), nil
}

func bindSchemaItemsToModelsItems(items []schema.Item) []models.Item {
	result := make([]models.Item, len(items))
	for i := range items {
		result[i].ID = models.ItemID(items[i].ID)
		result[i].Name = items[i].Name
		result[i].Price = models.Price(items[i].Price)
		result[i].Tags = bindInts64ToTags(items[i].Tags)
		result[i].Length = uint32(items[i].Length.Int64)
		result[i].Width = uint32(items[i].Width.Int64)
		result[i].Height = uint32(items[i].Height.Int64)
	}
	return result
}

func bindInts64ToTags(tags []int64) models.Tags {
	result := make(models.Tags, len(tags))
	for i := range tags {
		result[i] = models.Tag(tags[i])
	}
	return result
}

func (r *ItemsRepo) GetItemStocks(ctx context.Context, itemID models.ItemID) (models.ItemStocks, error) {
	db := r.QueryEngineProvider.GetQueryEngine(ctx)

	// build query
	const query = ""

	rows, err := db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(models.ItemStocks)

	// построчная обработка
	_ = pgxscan.NewRowScanner(rows)
	for rows.Next() {
		// var st Student
		// if err := rs.Scan(&st); err != nil {
		// 	log.Fatal(err)
		// }
		// do something here
		// students[st.ID] = st
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *ItemsRepo) ReserveStock(ctx context.Context, itemID models.ItemID, count uint32, warehouseID models.WarehouseID) (uint32, error) {
	return 0, nil
}
