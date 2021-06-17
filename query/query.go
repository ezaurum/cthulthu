package query

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"net/url"
	"reflect"
	"strings"
)

func Query(q Param, db *gorm.DB, out interface{}) (*Response, error) {
	// 기본 제한수를 200개로 둔다
	if q.Limit == 0 {
		q.Limit = 200
	}
	return query(q, db, out)
}

const (
	exactPattern  = "%s = '%v'"
	queryPattern  = "%s like '%%%s%%'"
	queryPattern2 = " OR " + queryPattern
)

func makeLinkedQuery(queryTarget, queryString string) string {
	split := strings.Split(queryTarget, ",")
	qs := fmt.Sprintf(queryPattern, split[0], queryString)
	for _, s := range split[1:] {
		qs += fmt.Sprintf(queryPattern2, s, queryString)
	}
	return qs
}

func query(q Param, w *gorm.DB, out interface{}) (*Response, error) {
	orderedResultQueryString, orderedWhere, unorderedWhere := MakeQuery(q, w)

	return MakeResponse(q, orderedWhere, unorderedWhere, orderedResultQueryString, out)
}

//MakeResponse 만든 쿼리를 가지고 Response 객체를 만든다
func MakeResponse(q Param, orderedWhere *gorm.DB, unorderedWhere *gorm.DB, orderedResultQueryString string, out interface{}) (*Response, error) {
	if f := orderedWhere.Find(out); nil != f.Error && f.Error != gorm.ErrRecordNotFound {
		return nil, f.Error
	}

	var count int
	if c := unorderedWhere.Model(out).Count(&count); nil != c.Error && c.Error != gorm.ErrRecordNotFound {
		return nil, c.Error
	}

	r, err := response(q, orderedResultQueryString, out)
	if nil != r {
		r.Total = count
	}
	return r, err
}

// MakeQuery 전체적으로 쿼리를 만들어준다
func MakeQuery(q Param, rw *gorm.DB) (string, *gorm.DB, *gorm.DB) {
	resultLinkQueryString, unorderedWhere := MakeRawQuery(q, rw)
	orderedWhere, orderedResultString := MakeOrderQuery(q, unorderedWhere, resultLinkQueryString)
	return orderedResultString, orderedWhere, unorderedWhere
}

// MakeOrderQuery 이전 쿼리에다 정렬 관련 쿼리를 붙여준다
func MakeOrderQuery(q Param, w *gorm.DB, resultLinkQueryString string) (*gorm.DB, string) {
	if len(q.OrderBy) > 0 {
		split := strings.Split(q.OrderBy, ",")
		for _, s := range split {
			s = strings.ToLower(s)
			replace := strings.ReplaceAll(s, ".", " ")
			w = w.Order(replace)
		}
		resultLinkQueryString += fmt.Sprintf("&o=%s", q.OrderBy)
	}
	if q.Limit > 0 {
		w = w.Limit(q.Limit + 1)
	}
	if q.Start > 0 {
		w = w.Offset(q.Start)
	}
	if q.After != 0 {
		w = w.Where("id > ?", q.After)
	}
	if q.Before != 0 {
		w = w.Where("id < ?", q.Before)
	}
	return w, resultLinkQueryString
}

// MakeRawQuery 정렬 전 쿼리를 만들어준다
func MakeRawQuery(q Param, w *gorm.DB) (string, *gorm.DB) {
	resultLinkQueryString := ""
	queryString, _ := url.QueryUnescape(q.QueryString)
	if len(queryString) > 0 && len(q.QueryTarget) > 0 {
		lqs := makeLinkedQuery(q.QueryTarget, queryString)
		w = w.Where(lqs)
		resultLinkQueryString += fmt.Sprintf("&q=%s&qt=%s", q.QueryString, q.QueryTarget)
	}
	if len(q.QueryExact) > 0 {
		split := strings.Split(q.QueryExact, ",")
		for _, s := range split {
			pair := strings.Split(s, ".")
			w = w.Where(fmt.Sprintf(exactPattern, strings.ToLower(pair[0]), pair[1]))
		}
		resultLinkQueryString += fmt.Sprintf("&qx=%s", q.QueryExact)
	}
	if q.QueryValues != nil {
		for k, v := range q.QueryValues {
			w = w.Where(fmt.Sprintf(exactPattern, k, v))
			resultLinkQueryString += fmt.Sprintf("&%s=%v", k, v)
		}
	}
	return resultLinkQueryString, w
}

func response(q Param, resultLinkQueryString string, out interface{}) (*Response, error) {
	// 파라마터는 &[]domainObject 형태로 넘어오므로 실제 slice를 가져온다
	of := reflect.ValueOf(out)
	v := of.Elem()
	count := v.Len()
	var links NavigationLinks
	if count > q.Limit && q.Limit > 0 {
		count = count - 1
		v = v.Slice(0, count)
		// 자른 걸 실제 슬라이드에 반영
		of.Elem().Set(v)
		links.Next = fmt.Sprintf("l=%d&s=%d", q.Limit, q.Start+count) + resultLinkQueryString
	}
	if q.Start >= q.Limit && q.Limit > 0 {
		links.Prev = fmt.Sprintf("l=%d&s=%d", q.Limit, q.Start-q.Limit) + resultLinkQueryString
	}

	links.Self = fmt.Sprintf("l=%d&s=%d", q.Limit, q.Start) + resultLinkQueryString

	return &Response{
		Param:  q,
		Size:   count,
		Links:  links,
		Result: v.Interface(),
	}, nil
}
