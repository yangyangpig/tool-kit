package goods

// 怎么解决贫血模型
type GoodsDomain struct {
	Name string
	Id   uint64
	Price float64
}

func NewGoodsDomain() *GoodsDomain { return nil}


