package common

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/savsgio/atreugo/v11"
)

func BindQueryToStruct[T any](rc *atreugo.RequestCtx) (*T, error) {
	queryParams := make(map[string]string)
	rc.QueryArgs().VisitAll(func(key, value []byte) {
		queryParams[string(key)] = string(value)
	})

	jsonBytes, err := json.Marshal(queryParams)
	if err != nil {
		return nil, err
	}

	var result T
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func BindBodyToStruct[T any](rc *atreugo.RequestCtx) (*T, error) {
	body := make(map[string]string)
	bodyData, ok := rc.UserValue("body").([]byte)
	if !ok {
		return nil, errors.New("failed to get body from request context")
	}
	bodyStream := bytes.NewReader(bodyData)
	if err := json.NewDecoder(bodyStream).Decode(&body); err != nil {
		return nil, err
	}
	rc.Request.CloseBodyStream()

	jsonBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	var result T
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func GetParams(rc *atreugo.RequestCtx, keys []string) map[string]string {
	params := make(map[string]string)
	for _, key := range keys {
		val, ok := rc.UserValue(key).(string)
		if !ok {
			continue
		}
		params[key] = val
	}

	return params
}
