package app

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/shahulsonhal/store-service/internal/data"
)

var storeLocation = map[string][]data.StoreDetails{
	"DE": {
		{
			StoreID:     "DE1",
			Name:        "De store1",
			Country:     "Germany",
			CountryCode: "DE",
			Location: data.LocationData{
				Lat: 1.4,
				Lng: 11.2,
			},
			SlowService: true,
		},
		{
			StoreID:     "DE2",
			Name:        "De store2",
			Country:     "Germany",
			CountryCode: "DE",
			Location: data.LocationData{
				Lat: 2.4,
				Lng: 14.2,
			},
			SlowService: false,
		},
	},
	"FR": {
		{
			StoreID:     "FR1",
			Name:        "Fr store1",
			Country:     "France",
			CountryCode: "FR",
			Location: data.LocationData{
				Lat: 15.4,
				Lng: 25.2,
			},
			SlowService: false,
		},
		{
			StoreID:     "FR2",
			Name:        "Fr store2",
			Country:     "France",
			CountryCode: "FR",
			Location: data.LocationData{
				Lat: 10.4,
				Lng: 20.2,
			},
			SlowService: false,
		},
	},
}

func TestServer_handleGetStoreV1(t *testing.T) {
	type expect struct {
		apiStatus  string
		statusCode int
		response   []data.StoreDetails
	}

	type fields struct {
		max,
		country string
		repo data.Repo
	}

	tests := []struct {
		name   string
		fields fields
		want   expect
	}{
		{
			name: "success case: initial case for Germany",
			fields: fields{
				country: "DE",
				repo: &data.MockRepo{
					MockGetStore: func(max int, country string) ([]data.StoreDetails, error) {
						return storeLocation["DE"], nil
					},
				},
			},
			want: expect{
				apiStatus:  statusOK,
				statusCode: http.StatusOK,
				response:   storeLocation["DE"],
			},
		},
		{
			name: "success case: initial case for France",
			fields: fields{
				country: "FR",
				repo: &data.MockRepo{
					MockGetStore: func(max int, country string) ([]data.StoreDetails, error) {
						return storeLocation["FR"], nil
					},
				},
			},
			want: expect{
				apiStatus:  statusOK,
				statusCode: http.StatusOK,
				response:   storeLocation["FR"],
			},
		},
		{
			name: "success case: with max value less than items",
			fields: fields{
				country: "FR",
				max:     "1",
				repo: &data.MockRepo{
					MockGetStore: func(max int, country string) ([]data.StoreDetails, error) {
						return storeLocation["FR"], nil
					},
				},
			},
			want: expect{
				apiStatus:  statusOK,
				statusCode: http.StatusOK,
				response:   storeLocation["FR"],
			},
		},
		{
			name: "failure case: invalid max",
			fields: fields{
				country: "FR",
				max:     "abc",
				repo: &data.MockRepo{
					MockGetStore: func(max int, country string) ([]data.StoreDetails, error) {
						return storeLocation["FR"], nil
					},
				},
			},
			want: expect{
				apiStatus:  statusFail,
				statusCode: http.StatusBadRequest,
			},
		},
		{
			name: "failure case: country record not found",
			fields: fields{
				country: "FR",
				max:     "1",
				repo: &data.MockRepo{
					MockGetStore: func(max int, country string) ([]data.StoreDetails, error) {
						return nil, data.ErrResourceNotFound
					},
				},
			},
			want: expect{
				apiStatus:  statusFail,
				statusCode: http.StatusNotFound,
			},
		},
		{
			name: "failure case: fail to list store (500)",
			fields: fields{
				country: "FR",
				max:     "1",
				repo: &data.MockRepo{
					MockGetStore: func(max int, country string) ([]data.StoreDetails, error) {
						return nil, fmt.Errorf("repo error")
					},
				},
			},
			want: expect{
				apiStatus:  statusFail,
				statusCode: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				repo: tt.fields.repo,
			}
			url := "/"
			r, err := http.NewRequest("GET", url, nil)
			if err != nil {
				t.Fatal(err)
			}

			if tt.fields.max != "" {
				q := r.URL.Query()
				q.Add("max", tt.fields.max)
				r.URL.RawQuery = q.Encode()
			}

			if tt.fields.country != "" {
				q := r.URL.Query()
				q.Add("country", tt.fields.country)
				r.URL.RawQuery = q.Encode()
			}

			w := httptest.NewRecorder()
			s.initStoreRouter().ServeHTTP(w, r)

			if tt.want.statusCode != w.Result().StatusCode {
				t.Fatalf(
					"expected http status code %d, got %d",
					tt.want.statusCode,
					w.Result().StatusCode,
				)
			}

			if tt.want.apiStatus == statusOK {
				response, err := readSuccess(w)
				if err != nil {
					t.Fatal(err)
				}

				var result []data.StoreDetails
				if err := json.Unmarshal(response.Result, &result); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(result, tt.want.response) {
					t.Errorf("expect %v got %v", tt.want.response, result)
				}
			} else {
				response, err := readFailure(w)
				if err != nil {
					t.Fatal(err)
				}
				if tt.want.statusCode != response.Error.Code {
					t.Errorf(
						"expected error code %d but got %d",
						tt.want.statusCode,
						response.Error.Code,
					)
				}
			}
		})
	}
}

func readSuccess(w *httptest.ResponseRecorder) (*Response, error) {
	r, err := readResponse(w)
	if err != nil {
		return nil, err
	}
	if r.Error != nil {
		return r, err
	}
	if r.Status != statusOK {
		return r, err
	}
	return r, nil
}

func readFailure(w *httptest.ResponseRecorder) (*Response, error) {
	r, err := readResponse(w)
	if err != nil {
		return nil, err
	}
	if r.Status != statusFail {
		return r, err
	}
	if r.Error == nil {
		return r, err
	}
	return r, nil
}

func readResponse(w *httptest.ResponseRecorder) (*Response, error) {
	r := &Response{}
	if err := json.Unmarshal(w.Body.Bytes(), r); err != nil {
		return nil, err
	}
	return r, nil
}
