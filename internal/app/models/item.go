package models

type ItemID uint32

type Price uint32 // RUB коп.

type Item struct {
	ID    ItemID
	Name  string
	Price Price
	Tags
	Size
}

type Size struct {
	Length, Width, Height uint32
}

type Tag uint32

type Tags []Tag

type ItemFilter struct {
	Name   string
	ByTags struct {
		Tags
		Use bool
	}
	ByPrice struct {
		From, To Price
		Use      bool
	}
}

type ItemStocks map[WarehouseID]uint32
