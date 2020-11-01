package query

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func TestQuery(t *testing.T) {

	queryString, _ := url.QueryUnescape("test")
	assert.Equal(t, "test", queryString)

}

func TestQueryMultipleTarget(t *testing.T) {
	linkedQuery := makeLinkedQuery("name,title,email", "mfmf")
	assert.Equal(t, "name like '%mfmf%' OR title like '%mfmf%' OR email like '%mfmf%'", linkedQuery)
}
