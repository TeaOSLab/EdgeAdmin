package ttlcache

type OptionInterface interface {
}

type PiecesOption struct {
	Count int
}

func NewPiecesOption(count int) *PiecesOption {
	return &PiecesOption{Count: count}
}

type MaxItemsOption struct {
	Count int
}

func NewMaxItemsOption(count int) *MaxItemsOption {
	return &MaxItemsOption{Count: count}
}
