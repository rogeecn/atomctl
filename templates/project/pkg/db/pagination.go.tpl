package db

type Pagination struct {
	Offset   string `json:"offset,omitempty"`
	OffsetID int64  `json:"-"`
	Action   int    `json:"action"` // action: 0 :加载更多 1:刷新
}
