package domain

// PaginationQuery 分页查询值对象
type PaginationQuery struct {
	Skip  int64 // 跳过条数
	Limit int64 // 每页条数
}

// SearchQuery 搜索查询值对象
type SearchQuery struct {
	Search string // 搜索关键词
}

// StudentFiltersQuery 学生筛选查询值对象
type StudentFiltersQuery struct {
	RegisterDateFrom  string // 注册日期起始
	RegisterDateTo    string // 注册日期截止
	LastVisitDateFrom string // 最后登录日期起始
	LastVisitDateTo   string // 最后登录日期截止
	Verified          *bool  // 是否验证
}

// GetStudentsQuery 获取学生列表查询值对象
type GetStudentsQuery struct {
	PaginationQuery
	SearchQuery
	StudentFiltersQuery
}

// OrdersFiltersQuery 订单筛选查询值对象
type OrdersFiltersQuery struct {
	DateFrom string // 日期起始
	DateTo   string // 日期截止
	Status   string // 订单状态
}

// GetOrdersQuery 获取订单列表查询值对象
type GetOrdersQuery struct {
	PaginationQuery
	SearchQuery
	OrdersFiltersQuery
}
