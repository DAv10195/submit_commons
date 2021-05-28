package http

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const (
	limit		= "limit"
	afterId		= "after_id"

	base		= 10
	bitSize		= 64
)

// params for performing rest pagination
type PagingParams struct {
	Limit 		int64
	AfterId		int64
}

// parse the query string of the given request and return the required paging params if they exist (-1
// value means not given). If multiple values are assigned to a param, the 1st one is taken
func PagingParamsFromRequest(r *http.Request) (*PagingParams, error) {
	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		return nil, err
	}
	var limitVal, afterIdVal int64
	limitVals, ok := params[limit]
	if ok && len(limitVals) > 0 {
		limitVal, err = strconv.ParseInt(limitVals[0], base, bitSize)
		if err != nil {
			return nil, err
		}
		if limitVal <= 0 {
			return nil, fmt.Errorf("%s query param should be positive", limit)
		}
	} else {
		limitVal = 0
	}
	afterIdVals, ok := params[afterId]
	if ok && len(afterIdVals) > 0 {
		afterIdVal, err = strconv.ParseInt(afterIdVals[0], base, bitSize)
		if err != nil {
			return nil, err
		}
		if afterIdVal <= 0 {
			return nil, fmt.Errorf("%s query param should be positive", afterId)
		}
	} else {
		afterIdVal = 0
	}
	return &PagingParams{limitVal, afterIdVal}, nil
}
