package helpers

import (
	"errors"
	"net/http"

	"goauthentik.io/api/v3"
)

// Generic interface that mimics a generated request by the API client
// Requires mainly `Treq` which will be the actual request type, and
// `Tres` which is the response type
type PaginatorRequest[Treq any, Tres any] interface {
	Page(page int32) Treq
	PageSize(size int32) Treq
	Execute() (Tres, *http.Response, error)
}

// Generic interface that mimics a generated response by the API client
type PaginatorResponse[Tobj any] interface {
	GetResults() []Tobj
	GetPagination() api.Pagination
}

// Paginator options for page size
type PaginatorOptions struct {
	PageSize int
}

// Automatically fetch all objects from an API endpoint using the pagination
// data received from the server.
func Paginator[Tobj any, Treq any, Tres PaginatorResponse[Tobj]](
	req PaginatorRequest[Treq, Tres],
	opts PaginatorOptions,
) ([]Tobj, *http.Response, error) {
	var bfreq, cfreq any
	fetchOffset := func(page int32) (Tres, *http.Response, error) {
		bfreq = req.Page(page)
		cfreq = bfreq.(PaginatorRequest[Treq, Tres]).PageSize(int32(opts.PageSize))
		res, hres, err := cfreq.(PaginatorRequest[Treq, Tres]).Execute()
		if err != nil {
			if hres != nil && hres.StatusCode >= 400 && hres.StatusCode < 500 {
				return res, hres, err
			}
		}
		return res, hres, err
	}
	var page int32 = 1
	errs := make([]error, 0)
	objects := make([]Tobj, 0)
	for {
		apiObjects, hr, err := fetchOffset(page)
		if err != nil {
			if page == 1 {
				return objects, hr, err
			}
			errs = append(errs, err)
			continue
		}
		objects = append(objects, apiObjects.GetResults()...)
		if apiObjects.GetPagination().Next > 0 {
			page += 1
		} else {
			break
		}
	}
	return objects, nil, errors.Join(errs...)
}
