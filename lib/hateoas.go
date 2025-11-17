package lib

import (
	"fmt"
	"net/url"
	"strconv"
)

func Hateoas(baseURL string, path string, page int, limit int, totalPage int, extraQuery url.Values) map[string]string {
	makeQuery := func(pageValue int) string {
		q := url.Values{}
		for key, val := range extraQuery {
			q[key] = val
		}
		q.Set("page", strconv.Itoa(pageValue))
		q.Set("limit", strconv.Itoa(limit))
		return q.Encode()
	}

	var nextURL, prevURL string

	if page > 1 {
		prevURL = fmt.Sprintf("%s%s?%s", baseURL, path, makeQuery(page-1))
	}

	if page < totalPage {
		nextURL = fmt.Sprintf("%s%s?%s", baseURL, path, makeQuery(page+1))
	}

	return map[string]string{
		"page":       strconv.Itoa(page),
		"limit":      strconv.Itoa(limit),
		"total_page": strconv.Itoa(totalPage),
		"prev":       prevURL,
		"next":       nextURL,
	}
}
