package clients

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

type GeonamesSearchResponse struct {
	Geonames []struct {
		GeonameID   int    `json:"geonameId"`
		Name        string `json:"toponymName"`
		Population  int64  `json:"population"`
		CountryID   string `json:"countryId"`
		CountryCode string `json:"countryCode"`
	} `json:"geonames"`
}

type GeonamesPostalCodeSearchResponse struct {
	PostalCodes []struct {
		PlaceName   string `json:"placeName"`
		AdminName2  string `json:"adminName2"`
		AdminName1  string `json:"adminName1"`
		CountryCode string `json:"countryCode"`
		PostalCode  string `json:"postalCode"`
		ISO31662    string `json:"ISO3166-2"`
	} `json:"postalCodes"`
}

func GeonamesSearchForCity(username, token, city, countryCode string) (geonameID *int, err error) {
	urlParams := url.Values{
		"name":         []string{city},
		"featureClass": []string{"P"}, // city, village, etc.
		"maxRows":      []string{"1"},
	}
	if countryCode != "" {
		urlParams["countryBias"] = []string{countryCode}
	}
	geonamesResponse, err := GeonamesSearchApiCall(username, token, urlParams)
	if err != nil {
		return nil, err
	}
	if len(geonamesResponse.Geonames) != 1 {
		// Try search again without specifying the featureClass to other features like municipalities.
		delete(urlParams, "featureClass")
		geonamesResponse, err = GeonamesSearchApiCall(username, token, urlParams)
		if err != nil {
			return nil, err
		}
		if len(geonamesResponse.Geonames) != 1 {
			return nil, nil
		}
	}
	return &geonamesResponse.Geonames[0].GeonameID, nil
}

func GeonamesSearchForCountry(username, token, countryCode string) (geonameID *int, err error) {
	urlParams := url.Values{
		"country":      []string{countryCode},
		"featureClass": []string{"A"},    // country, state, region, etc.
		"featureCode":  []string{"PCLI"}, // independent political entity
		"maxRows":      []string{"1"},
	}
	geonamesResponse, err := GeonamesSearchApiCall(username, token, urlParams)
	if err != nil {
		return nil, err

	}
	if len(geonamesResponse.Geonames) != 1 {
		return nil, nil
	}
	return &geonamesResponse.Geonames[0].GeonameID, nil
}

func GeonamesPostalCodeSearch(username, token, postalCode, countryCode string) (response json.RawMessage, err error) {
	urlParams := url.Values{
		"postalcode": []string{postalCode},
		"maxRows":    []string{"10"},
		"isReduced":  []string{"false"}, // Don't use backwards compatible API.
	}
	if countryCode != "" {
		urlParams["country"] = []string{countryCode}

		// Some adjustments for the Postal Code restrictions in GeoNames:
		// - For Canada we have only the first letters of the full postal codes (for copyright reasons)
		// - For Ireland we have only the first letters of the full postal codes (for copyright reasons)
		// - For Malta we have only the first letters of the full postal codes (for copyright reasons)
		// - The Argentina data file contains 4-digit postal codes which were replaced with a new system in 1999.
		// - For Brazil only major postal codes are available (only the codes ending with -000 and the major code per
		//   municipality).
		if countryCode == "CA" || countryCode == "IE" || countryCode == "MT" {
			postalCode = postalCode[0:3]
		}
	}
	return GeonamesApiCall(username, token, "postalCodeSearch", urlParams)
}

func GeonamesSearchApiCall(username, token string, urlParams url.Values) (*GeonamesSearchResponse, error) {
	rawResponse, err := GeonamesApiCall(username, token, "search", urlParams)
	if err != nil {
		return nil, err
	}

	var response GeonamesSearchResponse
	if err := json.Unmarshal(rawResponse, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func GeonamesApiCall(username string, token string, webService string, uriParams url.Values) (response json.RawMessage, err error) {
	// List of available web services: https://www.geonames.org/export/ws-overview.html
	uriParams["username"] = []string{username}
	uriParams["token"] = []string{token}

	uri := url.URL{
		Scheme:   "https",
		Host:     "secure.geonames.net",
		Path:     webService + "JSON",
		RawQuery: uriParams.Encode(),
	}

	resp, err := http.Get(uri.String())
	if err != nil {
		return nil, errors.Unexpected.Wrap("geonames.org api error", err)
	}

	rawResponse, err := ioutil.ReadAll(resp.Body)
	if (http.StatusBadRequest < resp.StatusCode && resp.StatusCode < http.StatusInternalServerError) ||
		strings.HasPrefix(string(rawResponse), `{"status":{"message":`) {
		// Application level error in the Geonames API.
		var errorResponse struct {
			Status struct {
				Message string `json:"message"`
			} `json:"status"`
		}
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(rawResponse, &errorResponse); err != nil {
			return nil, err
		}
		return nil, errors.Unexpected.Wrap("%s", errorResponse.Status.Message)

	} else if resp.StatusCode != http.StatusOK {
		// Unexpected error in the Geonames API.
		return nil, errors.Unexpected.Wrap(
			fmt.Sprintf("unexpected response status code from geonames.org api: %d", resp.StatusCode))
	}

	return rawResponse, nil
}
