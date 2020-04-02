package clients

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/Suburbia-io/dashboard/pkg/errors"
)

type TextSearchArgs struct {
	ApiKey    string
	Input     string
	InputType string
	Language  string
}

type TextSearchResBody struct {
	Candidates []struct {
		PlaceId string `json:"place_id"`
	} `json:"candidates"`
}

func GooglePlacesTextSearch(args TextSearchArgs) (placeId *string, err error) {
	uriParams := url.Values{
		"key":          []string{args.ApiKey},
		"input":        []string{args.Input},
		"inputtype":    []string{args.InputType},
		"language":     []string{args.Language},
		"locationbias": []string{"circle:1@0,0"}, // dirty hack to disable location bias
	}

	uri := url.URL{
		Scheme:   "https",
		Host:     "maps.googleapis.com",
		Path:     "maps/api/place/findplacefromtext/json",
		RawQuery: uriParams.Encode(),
	}

	res, err := http.Get(uri.String())
	if err != nil {
		return nil, errors.Unexpected
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var resBodyJson TextSearchResBody

	if err := json.Unmarshal(b, &resBodyJson); err != nil {
		return nil, err
	}

	if len(resBodyJson.Candidates) > 0 {
		return &resBodyJson.Candidates[0].PlaceId, nil
	}

	return nil, nil
}

type PlaceDetailsArgs struct {
	ApiKey   string
	PlaceId  string
	Language string
}

type PlaceDetailResBody struct {
	Result struct {
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
	} `json:"result"`
}

func GooglePlacesDetails(args PlaceDetailsArgs) (city string, country string, err error) {

	uriParams := url.Values{
		"key":      []string{args.ApiKey},
		"placeid":  []string{args.PlaceId},
		"language": []string{args.Language},
	}

	uri := url.URL{
		Scheme:   "https",
		Host:     "maps.googleapis.com",
		Path:     "maps/api/place/details/json",
		RawQuery: uriParams.Encode(),
	}

	res, err := http.Get(uri.String())
	if err != nil {
		return "", "", errors.Unexpected
	}

	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", "", err
	}

	var resBodyJson PlaceDetailResBody
	if err := json.Unmarshal(b, &resBodyJson); err != nil {
		return "", "", err
	}

	for _, addressComponent := range resBodyJson.Result.AddressComponents {
		for _, addressComponentType := range addressComponent.Types {
			if addressComponentType == "locality" {
				city = addressComponent.LongName
				break
			}
		}
		if city != "" {
			break
		}
	}

	for _, addressComponent := range resBodyJson.Result.AddressComponents {
		for _, addressComponentType := range addressComponent.Types {
			if addressComponentType == "country" {
				country = addressComponent.ShortName
				break
			}
		}
		if country != "" {
			break
		}
	}

	return city, country, nil
}
