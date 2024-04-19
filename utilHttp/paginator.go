// Copyright 2013 wetalk authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package utilHttp

import (
	"github.com/hilaoyu/go-utils/utilConvert"
	"math"
	"net/http"
	"net/url"
	"strconv"
)

type Paginator struct {
	Request     *http.Request `json:"-"`
	PerPage     int           `json:"per_page"`
	CurrentPage int           `json:"current_page"`
	MaxPages    int           `json:"-"`

	Total     int64 `json:"total"`
	PageRange []int `json:"page_range"`
	PageNums  int   `json:"page_nums"`
}

func NewPaginator(req *http.Request, pageSize int, total interface{}) *Paginator {
	p := Paginator{}
	p.Request = req
	if pageSize <= 0 {
		pageSize = 10
	}
	p.PerPage = pageSize
	p.SetTotal(total)
	p.GetPages()
	return &p
}
func (p *Paginator) GetTotal() int64 {
	return p.Total
}
func (p *Paginator) GetPageNums() int {
	if p.PageNums != 0 {
		return p.PageNums
	}
	pageNums := math.Ceil(float64(p.Total) / float64(p.PerPage))
	if p.MaxPages > 0 {
		pageNums = math.Min(pageNums, float64(p.MaxPages))
	}
	p.PageNums = int(pageNums)
	return p.PageNums
}

func (p *Paginator) SetTotal(total interface{}) {
	p.Total, _ = utilConvert.ToInt64(total)
	p.GetPages()
}

func (p *Paginator) GetCurrentPage() int {
	if p.CurrentPage != 0 {
		return p.CurrentPage
	}
	if p.Request.Form == nil {
		p.Request.ParseForm()
	}
	p.CurrentPage, _ = strconv.Atoi(p.Request.Form.Get("pager_page"))
	/*if p.CurrentPage > p.GetPageNums() {
		p.CurrentPage = p.GetPageNums()
	}*/
	if p.CurrentPage <= 0 {
		p.CurrentPage = 1
	}

	return p.CurrentPage
}

func (p *Paginator) GetPages() []int {
	if p.PageRange == nil && p.Total > 0 {
		var pages []int
		pageNums := p.GetPageNums()
		page := p.GetCurrentPage()
		switch {
		case page >= pageNums-4 && pageNums > 9:
			start := pageNums - 9 + 1
			pages = make([]int, 9)
			for i, _ := range pages {
				pages[i] = start + i
			}
		case page >= 5 && pageNums > 9:
			start := page - 5 + 1
			pages = make([]int, int(math.Min(9, float64(page+4+1))))
			for i, _ := range pages {
				pages[i] = start + i
			}
		default:
			pages = make([]int, int(math.Min(9, float64(pageNums))))
			for i, _ := range pages {
				pages[i] = i + 1
			}
		}
		p.PageRange = pages
	}
	return p.PageRange
}

func (p *Paginator) PageLink(page int) string {
	link, _ := url.ParseRequestURI(p.Request.RequestURI)
	values := link.Query()
	if page == 1 {
		values.Del("pager_page")
	} else {
		values.Set("pager_page", strconv.Itoa(page))
	}
	link.RawQuery = values.Encode()
	return link.String()
}

func (p *Paginator) PageLinkPrev() (link string) {
	if p.HasPrev() {
		link = p.PageLink(p.GetCurrentPage() - 1)
	}
	return
}

func (p *Paginator) PageLinkNext() (link string) {
	if p.HasNext() {
		link = p.PageLink(p.GetCurrentPage() + 1)
	}
	return
}

func (p *Paginator) PageLinkFirst() (link string) {
	return p.PageLink(1)
}

func (p *Paginator) PageLinkLast() (link string) {
	return p.PageLink(p.GetPageNums())
}

func (p *Paginator) HasPrev() bool {
	return p.GetCurrentPage() > 1
}

func (p *Paginator) HasNext() bool {
	return p.GetCurrentPage() < p.GetPageNums()
}

func (p *Paginator) IsActive(page int) bool {
	return p.GetCurrentPage() == page
}

func (p *Paginator) Offset() int {
	return (p.GetCurrentPage() - 1) * p.PerPage
}

func (p *Paginator) HasPages() bool {
	return p.GetPageNums() > 1
}
