package pager

const (
	// DefaultMaxPerPage максимальное кол-во записей на странице по умолчанию
	DefaultMaxPerPage = 1000
)

// Pager — структура, встраиваемая в структуру фильтра
// для удобной работы с LIMIT N OFFSET N
type Pager struct {
	page, perPage, perPageMax uint32
	sort                      []string
}

// New вернет новый пагинатор
func New() *Pager {
	return &Pager{}
}

// NewPagePer пагинатор с предустановленной страницей и кол-вом записей на странице
func NewPagePer(page, perPage uint32) *Pager {
	return New().SetPage(page).SetPerPage(perPage)
}

// SetPage задать текущую страницу
func (p *Pager) SetPage(val uint32) *Pager {
	p.page = val
	return p
}

// SetPerPage задать кол-во записей на странице
func (p *Pager) SetPerPage(val uint32) *Pager {
	p.perPage = val
	return p
}

// SetPerPageMax задать максимальное кол-во записей на странице
func (p *Pager) SetPerPageMax(val uint32) *Pager {
	p.perPageMax = val
	return p
}

// SetSortColumns задать поля для OrderBy
func (p *Pager) SetSortColumns(columns ...string) {
	p.sort = columns
}

// SetReverseOrderSort установить обратный порядок для заданных полей
func (p *Pager) SetReverseOrderSort(columns ...string) {
	for _, str := range columns {
		for i := range p.sort {
			if str == p.sort[i] {
				p.sort[i] += " DESC"
				break
			}
		}
	}
}

// GetSort вернуть строку для OrderBy
func (p *Pager) GetSort() []string {
	return p.sort
}

// GetLimit вернет SQL LIMIT
func (p *Pager) GetLimit() uint32 {
	if p.perPage == 0 || p.perPage > p.getPerPageMax() {
		return p.getPerPageMax()
	}

	return p.perPage
}

// GetOffset вернет для SQL OFFSET
func (p *Pager) GetOffset() uint32 {
	if p.page == 0 {
		return 0
	}

	return (p.page - 1) * p.GetLimit()
}

// GetLimit64 вернет SQL LIMIT для squirrel
func (p *Pager) GetLimit64() uint64 {
	return uint64(p.GetLimit())
}

// GetOffset64 вернет SQL OFFSET для squirrel
func (p *Pager) GetOffset64() uint64 {
	return uint64(p.GetOffset())
}

// GetPage вернет page (для тестирования мапперов)
func (p *Pager) GetPage() uint32 {
	return p.page
}

// GetPerPage вернет perPage (для тестирования мапперов)
func (p *Pager) GetPerPage() uint32 {
	return p.perPage
}

func (p *Pager) getPerPageMax() uint32 {
	if p.perPageMax == 0 {
		return DefaultMaxPerPage
	}

	return p.perPageMax
}
