package common

import "github.com/samber/lo"

type LabelItem struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type FilterType string

const (
	FilterTypeString FilterType = "string"
	FilterTypeNumber FilterType = "number"
	FilterTypeList   FilterType = "list"
	FilterTypeDate   FilterType = "date"
	FilterTypeBool   FilterType = "bool"
)

type Filter struct {
	Type  FilterType  `json:"type"`
	Label string      `json:"label"`
	Name  string      `json:"name"`
	Value string      `json:"value,omitempty"`
	Items []LabelItem `json:"items,omitempty"`
}

type Columns struct {
	Columns []TableColumnData `json:"columns,omitempty"`
	Hidden  []string          `json:"hidden,omitempty"`
}

func NewColumns(columns []TableColumnData) Columns {
	return Columns{
		Columns: lo.Map(columns, func(item TableColumnData, _ int) TableColumnData {
			return item.Format()
		}),
		Hidden: lo.FilterMap(columns, func(item TableColumnData, _ int) (string, bool) {
			if item.Hidden {
				return item.DataIndex, true
			}
			return "", false
		}),
	}
}
type TableColumnData struct {
	DataIndex         string             `json:"dataIndex,omitempty"`
	Title             string             `json:"title,omitempty"` // RenderFunction 类型被忽略
	Width             *int               `json:"width,omitempty"`
	Align             *string            `json:"align,omitempty"`
	Fixed             *string            `json:"fixed,omitempty"`
	Ellipsis          *bool              `json:"ellipsis,omitempty"`
	Tooltip           interface{}        `json:"tooltip,omitempty"` // 使用 interface{} 来表示 boolean | Record<string, any>
	Sortable          *TableSortable     `json:"sortable,omitempty"`
	Filterable        *TableFilterable   `json:"filterable,omitempty"`
	Children          *[]TableColumnData `json:"children,omitempty"`
	CellClass         *string            `json:"cellClass,omitempty"`
	HeaderCellClass   *string            `json:"headerCellClass,omitempty"`
	BodyCellClass     *string            `json:"bodyCellClass,omitempty"`
	SummaryCellClass  *string            `json:"summaryCellClass,omitempty"`
	CellStyle         *string            `json:"cellStyle,omitempty"`
	HeaderCellStyle   *string            `json:"headerCellStyle,omitempty"`
	BodyCellStyle     *string            `json:"bodyCellStyle,omitempty"`
	SummaryCellStyle  *string            `json:"summaryCellStyle,omitempty"`
	SlotName          *string            `json:"slotName,omitempty"`
	TitleSlotName     *string            `json:"titleSlotName,omitempty"`
	IsLastLeftFixed   *bool              `json:"isLastLeftFixed,omitempty"`
	IsFirstRightFixed *bool              `json:"isFirstRightFixed,omitempty"`
	ColSpan           *int               `json:"colSpan,omitempty"`
	RowSpan           *int               `json:"rowSpan,omitempty"`
	Index             *int               `json:"index,omitempty"`
	Parent            *TableColumnData   `json:"parent,omitempty"`
	ResizeWidth       *int               `json:"resizeWidth,omitempty"`
	Hidden            bool               `json:"hidden,omitempty"`
}

func (t TableColumnData) Format() TableColumnData {
	if t.SlotName == nil {
		t.SlotName = lo.ToPtr(t.DataIndex)
	}
	return t
}

type TableSortable struct {
	SortDirections   []string    `json:"sortDirections"`
	Sorter           interface{} `json:"sorter"` // 因为它可以是函数或布尔值，所以使用 interface{}
	SortOrder        *string     `json:"sortOrder,omitempty"`
	DefaultSortOrder *string     `json:"defaultSortOrder,omitempty"`
}

type TableFilterData struct {
	Text  interface{} `json:"text"` // 使用 interface{} 是因为它可以是 string 或 RenderFunction
	Value string      `json:"value"`
}

type TableFilterable struct {
	Filters              *[]TableFilterData `json:"filters,omitempty"`
	Filter               interface{}        `json:"filter"` // 使用 interface{} 是因为它是一个函数
	Multiple             *bool              `json:"multiple,omitempty"`
	FilteredValue        *[]string          `json:"filteredValue,omitempty"`
	DefaultFilteredValue *[]string          `json:"defaultFilteredValue,omitempty"`
	RenderContent        interface{}        `json:"renderContent,omitempty"` // 使用 interface{} 是因为它是一个函数
	Icon                 interface{}        `json:"icon,omitempty"`          // 使用 interface{} 是因为它是 RenderFunction
	TriggerProps         interface{}        `json:"triggerProps,omitempty"`  // TriggerProps 的实际 Go 类型需要进一步定义
	AlignLeft            *bool              `json:"alignLeft,omitempty"`
	SlotName             *string            `json:"slotName,omitempty"`
}
