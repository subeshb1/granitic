// Copyright 2016 Granitic. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be found in the LICENSE file at the root of this project.

package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/graniticio/granitic/instance"
	"github.com/graniticio/granitic/logging"
	"io/ioutil"
	"net/http"
)

const jsonMergerComponentName string = instance.FrameworkPrefix + "JsonMerger"

// NewJSONMerger creates a JSONMerger with a Logger
func NewJSONMerger(flm *logging.ComponentLoggerManager) *JSONMerger {
	jm := new(JSONMerger)

	jm.Logger = flm.CreateLogger(jsonMergerComponentName)

	return jm
}

// A JSONMerger can merge a sequence of JSON configuration files (from a filesystem or HTTP URL) into a single
// view of configuration that will be used to configure Grantic's facilities and the user's IoC components. See the top
// of this page for a brief explanation of how merging works.
type JSONMerger struct {
	// Logger used by Granitic framework components. Automatically injected.
	Logger logging.Logger

	// True if arrays should be joined when merging; false if the entire conetnts of the array should be overwritten.
	MergeArrays bool
}

// LoadAndMergeConfig takes a list of file paths or URIs to JSON files and merges them into a single in-memory object representation.
// See the top of this page for a brief explanation of how merging works. Returns an error if a remote URI returned a 4xx or 5xx response code,
// a file or folder could not be accessed or if two files could not be merged dued to JSON parsing errors.
func (jm *JSONMerger) LoadAndMergeConfig(files []string) (map[string]interface{}, error) {

	var mergedConfig map[string]interface{}
	var jsonData []byte
	var err error

	for index, fileName := range files {

		if isURL(fileName) {
			jm.Logger.LogTracef("Acessing URL %s", fileName)

			jsonData, err = jm.loadFromURL(fileName)

		} else {
			jm.Logger.LogTracef("Reading file %s", fileName)

			jsonData, err = ioutil.ReadFile(fileName)
		}

		if err != nil {
			m := fmt.Sprintf("Problem reading data from file/URL %s: %s", fileName, err)
			return nil, errors.New(m)
		}

		var loadedConfig interface{}
		err = json.Unmarshal(jsonData, &loadedConfig)

		if err != nil {
			m := fmt.Sprintf("Problem parsing data from a file or URL (%s) as JSON : %s", fileName, err)
			return nil, errors.New(m)
		}

		additionalConfig := loadedConfig.(map[string]interface{})

		if index == 0 {
			mergedConfig = additionalConfig
		} else {
			mergedConfig = jm.merge(mergedConfig, additionalConfig)
		}

	}

	return mergedConfig, nil
}

func (jm *JSONMerger) loadFromURL(url string) ([]byte, error) {

	r, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	if r.StatusCode >= 400 {
		m := fmt.Sprintf("HTTP %d", r.StatusCode)
		return nil, errors.New(m)
	}

	var b bytes.Buffer

	b.ReadFrom(r.Body)
	r.Body.Close()

	return b.Bytes(), nil
}

func (jm *JSONMerger) merge(base, additional map[string]interface{}) map[string]interface{} {

	for key, value := range additional {

		if existingEntry, ok := base[key]; ok {

			existingEntryType := JsonType(existingEntry)
			newEntryType := JsonType(value)

			if existingEntryType == JsonMap && newEntryType == JsonMap {
				jm.merge(existingEntry.(map[string]interface{}), value.(map[string]interface{}))
			} else if jm.MergeArrays && existingEntryType == JsonArray && newEntryType == JsonArray {
				base[key] = jm.mergeArrays(existingEntry.([]interface{}), value.([]interface{}))
			} else {
				base[key] = value
			}
		} else {
			jm.Logger.LogTracef("Adding %s", key)

			base[key] = value
		}

	}

	return base
}

func (jm *JSONMerger) mergeArrays(a []interface{}, b []interface{}) []interface{} {
	return append(a, b...)
}
