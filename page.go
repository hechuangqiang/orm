package orm

type Page struct {
	PageNo      int           `json:"pageNo"`
	PageSize    int           `json:"pageSize"`
	ResultCount int           `json:"resultCount"`
	List        []interface{} `json:"list"`
}

func NewPage(pageNo, pageSize int) *Page {
	p := Page{}
	if pageNo == 0 {
		p.PageNo = 1
	} else {
		p.PageNo = pageNo
	}

	if pageSize == 0 {
		p.PageSize = 15
	} else {
		p.PageSize = pageSize
	}

	return &p
}

func (p *Page) StartRow() int {
	return p.PageSize * (p.PageNo - 1)
}

func (p *Page) PageCount() int {
	pageCount := 0
	if p.ResultCount%p.PageSize == 0 {
		pageCount = p.ResultCount / p.PageSize
	} else {
		pageCount = p.ResultCount/p.PageSize + 1
	}
	return pageCount
}
