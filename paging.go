package webkit

import "strings"

// Paging ...
type Paging struct {
	Current int    `json:"current"`
	Total   int    `json:"total"`
	Link    string `json:"link"`
}

// NewPaging ...
func NewPaging(link string, current, count, total int) *Paging {
	if strings.Index(link, "?") < 0 {
		link += "?1=1"
	}
	if count <= 0 || total <= 0 {
		return &Paging{
			Link:    link,
			Total:   1,
			Current: 1,
		}
	}

	ret := &Paging{
		Link: link,
	}

	ret.Total = total / count
	if total%count != 0 {
		ret.Total++
	}

	if current <= 0 {
		current = 1
	}

	if current >= ret.Total {
		current = ret.Total
	}

	ret.Current = current

	return ret
}
