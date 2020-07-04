package query

import (
	"github.com/ezaurum/cthulthu/jsonconv"
)

type Param struct {
	// id.asc or id.desc,created_at.asc
	OrderBy     string `json:"order_by,omitempty" query:"o"`
	QueryString string `json:"query_string,omitempty" query:"q"`
	// qt=name,position
	QueryTarget string                 `json:"query_target,omitempty" query:"qt"`
	Start       int                    `json:"start,omitempty" query:"s"`
	Limit       int                    `json:"limit,omitempty" query:"l"`
	QueryExact  string                 `json:"query_exact,omitempty" query:"qx"`
	QueryValues map[string]interface{} `json:"-"`
	After       int64                  `json:"after_id,string,omitempty" query:"afi"`
	Before      int64                  `json:"before_id,string,omitempty" query:"bfi"`
}

type NavigationLinks struct {
	//"base": "http://localhost:8080/confluence",
	Base    string `json:"base,omitempty"`
	Context string `json:"context,omitempty"`
	//"next": "/rest/api/space/ds/content/page?limit=5&start=5",
	Next string `json:"next,omitempty"`
	//  "prev": "/rest/api/space/ds/content/page?limit=5&start=0",
	Prev string `json:"prev,omitempty"`
	//"self": "http://localhost:8080/confluence/rest/api/space/ds/content/page"
	Self string `json:"self,omitempty"`
}

type Response struct {
	Param
	// 다른 link와 헷갈리지 않도록 _를 붙여둔다.
	Links  NavigationLinks `json:"_links"`
	Result interface{}     `json:"result"`
	// 결과 크기
	Size int `json:"size"`
}

type BulkActionRequest struct {
	Action  string                    `json:"action"`
	Targets jsonconv.Int64StringSlice `json:"target_list,omitempty"`
}

func (req *BulkActionRequest) Valid() bool {
	if len(req.Targets) > 0 {
		return true
	}
	return false
}
