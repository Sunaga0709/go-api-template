package model

type BookPrice struct {
	uint
}

func NewBookPrice(price uint) BookPrice {
	return BookPrice{price}
}

func (b *BookPrice) Uint() uint {
	return b.uint
}

func (b *BookPrice) Int() int {
	return int(b.uint)
}
