// Copyright 2015 Zalando SE
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sed

import (
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/zalando/skipper/filters"
)

type sedType int

const (
	request sedType = iota
	response
)

type sed struct {
	regex   *regexp.Regexp
	replace string
}

/*
Substitutes the contents in the body matching a pattern with a given
replacement string. Think of it as using the Unix 'sed' utility with a typical
substitution command of the form 's/regexp/replacement/g'.

The substitution can be applied to a request or response body.

Name: "sed"
*/
func NewSed() filters.Spec { return &sed{} }

// Returns the name of this filter.
func (spec *sed) Name() string {
	return SedName
}

// Creates a new sed filter with the parameters specified in config.
func (spec *sed) CreateFilter(config []interface{}) (filters.Filter, error) {
	if len(config) != 2 {
		return nil, filters.ErrInvalidFilterParameters
	}

	expr, ok := config[0].(string)
	if !ok {
		return nil, filters.ErrInvalidFilterParameters
	}

	replace, ok := config[1].(string)
	if !ok {
		return nil, filters.ErrInvalidFilterParameters
	}

	regex, err := regexp.Compile(expr)
	if err != nil {
		return nil, err
	}

	return &sed{regex: regex, replace: replace}, nil
}

// Intentionally left with no implementation.
func (_ *sed) Request(_ filters.FilterContext) {}

// Applies this filter's regex to the response body and replaces what was
// matched with the provided replacement string.
func (f *sed) Response(ctx filters.FilterContext) {

	body, err := ioutil.ReadAll(ctx.Response().Body)
	if err != nil {
		log.Println(err)
		return
	}

	transformed := f.regex.ReplaceAllString(string(body), f.replace)

	ctx.Response().Body = ioutil.NopCloser(strings.NewReader(transformed))
	ctx.Response().ContentLength = int64(len(transformed))
}
