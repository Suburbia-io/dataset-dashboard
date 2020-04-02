package application

import (
	"encoding/csv"
	"encoding/json"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Suburbia-io/dashboard/pkg/database"
	"github.com/Suburbia-io/dashboard/pkg/database/tables"
	"github.com/Suburbia-io/dashboard/pkg/errors"
	"github.com/Suburbia-io/dashboard/pkg/helpers/clients"
	"github.com/Suburbia-io/dashboard/pkg/helpers/postalcodes"
	"github.com/Suburbia-io/dashboard/pkg/helpers/validate"
	"github.com/Suburbia-io/dashboard/shared"
	"github.com/rainycape/countries"
)

func (app *App) AdminApiLocationUpsertCSV(w http.ResponseWriter, r *http.Request, s shared.SharedSession) error {
	datasetID := r.URL.Query().Get("datasetID")
	if err := validate.UUID(datasetID); err != nil {
		return err
	}
	reader := csv.NewReader(r.Body)
	if err := app.DBAL.LocationUpsertCSV(datasetID, reader); err != nil {
		return errors.HttpBadRequestArgs
	}
	defer app.DBAL.AuditTrailByUserInsertAsync(s.Session, tables.Locations.Table(), datasetID, "locationUpsertCSV", "")
	go app.SearchGeonamesForLocations(datasetID)
	app.respondApi(w, r, true, nil)
	return nil
}

var emptyJSON = []byte("{}")

func (app *App) AdminApiLocationApprove(w http.ResponseWriter, r *http.Request, s shared.SharedSession) error {
	args := struct {
		DatasetID    string `json:"datasetID"`
		LocationHash string `json:"locationHash"`
		Approved     *bool  `json:"approved"`
	}{}
	if err := app.decodeRequest(r, &args); err != nil {
		return err
	}

	location, err := app.DBAL.LocationGet(args.DatasetID, args.LocationHash)
	if err != nil {
		return err
	}

	location.Approved = args.Approved

	if location.Approved != nil && *location.Approved && location.GeonameID == nil {
		return errors.Unexpected.Wrap("cannot approve a location without a geoname id")
	}

	if err = app.DBAL.LocationUpdate(&location); err != nil {
		return err
	}

	// Refresh the location object before returning.
	location, err = app.DBAL.LocationGet(args.DatasetID, args.LocationHash)
	if err != nil {
		return err
	}

	defer app.DBAL.AuditTrailByUserInsertAsync(
		s.Session, tables.Locations.Table(), location.LocationHash, "approve", location.Approved)
	app.respondApi(w, r, location, nil)
	return nil
}

func (app *App) AdminApiLocationSetGeonameID(w http.ResponseWriter, r *http.Request, s shared.SharedSession) (err error) {
	args := struct {
		DatasetID    string `json:"datasetID"`
		LocationHash string `json:"locationHash"`
		GeonameID    *int   `json:"geonameID"`
	}{}

	if err := app.decodeRequest(r, &args); err != nil {
		return err
	}

	location, err := app.DBAL.LocationGet(args.DatasetID, args.LocationHash)
	if err != nil {
		return err
	}

	if args.GeonameID != nil {
		urlParams := url.Values{
			"geonameId": []string{strconv.Itoa(*args.GeonameID)},
		}
		// We need to use "get" to check the geoname id because "hierarchy" defaults to the geoname id for "Earth" for
		// nonexistent geoname ids.
		_, err =
			clients.GeonamesApiCall(app.Config.GeonamesApiUsername, app.Config.GeonamesApiToken, "get", urlParams)
		if err != nil {
			return err
		}

		location.GeonamesHierarchy, err =
			clients.GeonamesApiCall(app.Config.GeonamesApiUsername, app.Config.GeonamesApiToken, "hierarchy", urlParams)
		if err != nil {
			return err
		}
	} else {
		location.GeonamesHierarchy = emptyJSON
	}
	location.GeonameID = args.GeonameID

	if err = app.DBAL.LocationUpdate(&location); err != nil {
		return err
	}

	// Refresh the location object before returning.
	location, err = app.DBAL.LocationGet(args.DatasetID, args.LocationHash)
	if err != nil {
		return err
	}

	defer app.DBAL.AuditTrailByUserInsertAsync(
		s.Session, tables.Locations.Table(), location.LocationHash, "setGeonameID", location.GeonameID)
	app.respondApi(w, r, location, nil)
	return nil
}

// Addition mapping found in the raw strings. If this comes up a lot, we can add multi-lingual support to the
// country string library we're using to avoid having to maintain this manually.
var countryMapping = map[string]string{
	"NEDERLAND":   "NL",
	"BELGIE":      "BE",
	"SPANJE":      "ES",
	"DEUTSCHLAND": "DE",
}

func extractCountryCode(locationString string) (countryCode string, citiesPostalCodes []string) {
	for _, part := range strings.Split(locationString, ",") {
		part = strings.TrimSpace(part)
		if part == "" || strings.ToLower(part) == "undefined" {
			continue
		}

		if countryCode == "" {
			country, _ := countries.Parse(part)
			if country != nil {
				countryCode = country.ISO2
				continue
			}
			if result, ok := countryMapping[strings.ToUpper(part)]; ok {
				countryCode = result
				continue
			}
		}
		citiesPostalCodes = append(citiesPostalCodes, part)
	}
	return
}

func (app *App) SearchGeonamesForLocations(datasetID string) error {
	status := "pending"
	locations, err := app.DBAL.LocationList(datasetID, database.LocationListArgs{Status: &status, Limit: math.MaxInt32})
	if err != nil {
		return err
	}

	for _, location := range locations {
		// The first step is to extract country information from location string.
		var citiesPostalCodes []string
		location.ParsedCountryCode, citiesPostalCodes = extractCountryCode(location.LocationString)

		// The second step is to extract the postal codes and then cites from the remaining parts of the location string.
		var cities []string
		for _, part := range citiesPostalCodes {
			// Search for the postal code in Geonames.
			if len(location.GeonamesPostalCodes) == len(emptyJSON) || string(location.GeonamesPostalCodes) == `{"postalCodes":[]}` {
				if location.ParsedCountryCode != "" {
					if re, ok := postalcodes.PostalCodeRegex[location.ParsedCountryCode]; ok {
						postalCode := re.FindString(part)
						if postalCode != "" {
							location.ParsedPostalCode = strings.ReplaceAll(strings.ToUpper(postalCode), " ", "")

							// Remove postal code from the city string.
							part = strings.TrimSpace(strings.ReplaceAll(part, postalCode, ""))

							rawResponse, err :=
								clients.GeonamesPostalCodeSearch(app.Config.GeonamesApiUsername, app.Config.GeonamesApiToken, postalCode, location.ParsedCountryCode)
							if err != nil {
								return err
							}
							var response clients.GeonamesPostalCodeSearchResponse
							if err := json.Unmarshal(rawResponse, &response); err != nil {
								return err
							}
							if len(response.PostalCodes) > 0 {
								location.GeonamesPostalCodes = rawResponse
							}
						}
					}
				} else if part != "" { // CountryCode == ""
					// Fallback to a dumb search with geonames postal code search without using the regex to verify the postal
					// code first.
					rawResponse, err :=
						clients.GeonamesPostalCodeSearch(app.Config.GeonamesApiUsername, app.Config.GeonamesApiToken, part, "")
					if err != nil {
						return err
					}
					// Reset the raw response if the search has no results.
					var response clients.GeonamesPostalCodeSearchResponse
					if err := json.Unmarshal(rawResponse, &response); err != nil {
						return err
					}
					if len(response.PostalCodes) > 0 {
						// Found a result in geonames postal code search so we don't add part to the cites slice.
						location.ParsedPostalCode = strings.ReplaceAll(strings.ToUpper(part), " ", "")
						location.GeonamesPostalCodes = rawResponse
						// TODO Could check all results for the same country instead of only setting the country if there's 1 result.
						if len(response.PostalCodes) == 1 {
							location.ParsedCountryCode = response.PostalCodes[0].CountryCode
						}
						continue
					}
				}
			}

			if part != "" {
				cities = append(cities, part)
			}
		}

		err = app.DBAL.LocationUpdate(&location)
		if err != nil {
			return err
		}

		geonameID, err := app.SearchGeonamesForCityStrings(cities, location.ParsedCountryCode)
		if err != nil {
			return err
		}
		if geonameID == nil {
			continue
		}

		// Update the hierarchy if needed.
		if location.GeonameID == nil || *location.GeonameID != *geonameID {
			location.GeonameID = geonameID

			urlParams := url.Values{
				"geonameId": []string{strconv.Itoa(*location.GeonameID)},
			}
			location.GeonamesHierarchy, err =
				clients.GeonamesApiCall(app.Config.GeonamesApiUsername, app.Config.GeonamesApiToken, "hierarchy", urlParams)
			if err != nil {
				return err
			}

			err = app.DBAL.LocationUpdate(&location)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (app *App) SearchGeonamesForCityStrings(cities []string, countryCode string) (geonameID *int, err error) {
	countryCodes := []string{countryCode}

	// Add an empty country code at the end of the country code list so that the city is searched without a country if
	// there are no matches with the country (or there are no countries found).
	if countryCode != "" {
		countryCodes = append(countryCodes, "")
	}

	// Search geonames for various combinations of the extracted information.
	for _, city := range cities {
		for _, countryCode := range countryCodes {
			geonameID, err =
				clients.GeonamesSearchForCity(app.Config.GeonamesApiUsername, app.Config.GeonamesApiToken, city, countryCode)
			if err != nil || geonameID != nil {
				return geonameID, err
			}
		}

		// Substrings of the city.
		if strings.Contains(city, " ") {
			for _, citySubstring := range strings.Split(city, " ") {
				// Substrings of the city with all countries.
				for _, countryCode := range countryCodes {
					geonameID, err =
						clients.GeonamesSearchForCity(app.Config.GeonamesApiUsername, app.Config.GeonamesApiToken, citySubstring, countryCode)
					if err != nil || geonameID != nil {
						return geonameID, err
					}
				}
			}
		}
	}

	// Fallback to country search.
	for _, countryCode := range countryCodes {
		geonameID, err =
			clients.GeonamesSearchForCountry(app.Config.GeonamesApiUsername, app.Config.GeonamesApiToken, countryCode)
		if err != nil || geonameID != nil {
			return geonameID, err
		}
	}

	return geonameID, nil
}
