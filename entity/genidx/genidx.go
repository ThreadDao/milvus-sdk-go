//go:build ignore
// +build ignore

// Copyright (C) 2019-2021 Zilliz. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License
// is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
// or implied. See the License for the specific language governing permissions and limitations under the License.

// This program generates entity/indexes_gen.go. Invoked by go generate
package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/milvus-io/milvus-sdk-go/v2/entity"
)

var indexTemplate = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT
// This file is generated by go generate 

package entity

import (
	"errors"
	"encoding/json"
	"fmt"
)

{{ range .Indexes }}{{with .}}

var _ Index = &Index{{.IdxName}}{}

// Index{{.IdxName}} idx type for {{.IdxType}}
type Index{{.IdxName}} struct { //auto generated fields{{range .ConstructParams}}{{with .}}
	{{.Name}} {{.ParamType}}{{end}}{{end}}
	metricType MetricType
}

// Name returns index type name, implementing Index interface
func(i *Index{{.IdxName}}) Name() string {
	return "{{.IdxName}}"
}

// IndexType returns IndexType, implementing Index interface
func(i *Index{{.IdxName}}) IndexType() IndexType {
	return IndexType("{{.IdxType}}")
}

// SupportBinary returns whether index type support binary vector
func(i *Index{{.IdxName}}) SupportBinary() bool {
	return {{.VectorSupport}} & 2 > 0
}

// Params returns index construction params, implementing Index interface
func(i *Index{{.IdxName}}) Params() map[string]string {
	params := map[string]string {//auto generated mapping {{range .ConstructParams}}{{with .}}
		"{{.Name}}": fmt.Sprintf("%v",i.{{.Name}}),{{end}}{{end}}
	}
	bs, _ := json.Marshal(params)
	return map[string]string {
		"params": string(bs),
		"index_type": string(i.IndexType()),
		"metric_type": string(i.metricType),
	}
}

// NewIndex{{.IdxName}} create index with construction parameters
func NewIndex{{.IdxName}}(metricType MetricType, {{range .ConstructParams}}{{with .}}
	{{.Name}} {{.ParamType}},
{{end}}{{end}}) (*Index{{.IdxName}}, error) {
	// auto generate parameters validation code, if any{{range .ConstructParams}}{{with .}}
	{{.ValidationCode}}
	{{end}}{{end}}
	return &Index{{.IdxName}}{ {{range .ConstructParams}}{{with .}}
	//auto generated setting
	{{.Name}}: {{.Name}},{{end}}{{end}}
	metricType: metricType,
	}, nil
}
{{end}}{{end}}
`))

var indexTestTemplate = template.Must(template.New("").Funcs(template.FuncMap{
	"SContains": func(s, sub string) bool { return strings.Contains(s, sub) },
}).Parse(`// Code generated by go generate; DO NOT EDIT
// This file is generated by go generate

package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

{{range $idx := .Indexes}}{{with.}}
func TestIndex{{.IdxName}}(t *testing.T){
	{{range .ConstructParams}}{{with.}}
	var {{.Name}} {{.ParamType}}{{end}}{{end}}

	{{if SContains $idx.IdxName "Bin" }}mt := HAMMING{{else}}mt := L2{{end}}
	

	t.Run("valid usage case", func(t *testing.T){
		{{range $i, $example := .ValidExamples}}
		{{$example}}
		idx{{$i}}, err := NewIndex{{$idx.IdxName}}(mt, {{range $idx.ConstructParams}}{{with.}}
			{{.Name}},{{end}}{{end}}
		)
		assert.Nil(t, err)
		assert.NotNil(t, idx{{$i}})
		assert.Equal(t, "{{$idx.IdxName}}", idx{{$i}}.Name())
		assert.EqualValues(t, "{{$idx.IdxType}}", idx{{$i}}.IndexType())
		assert.NotNil(t, idx{{$i}}.Params())
		{{if SContains $idx.IdxName "Bin" }}assert.True(t, idx{{$i}}.SupportBinary()){{else}}assert.False(t, idx{{$i}}.SupportBinary()){{end}}
		{{end}}
	})

	t.Run("invalid usage case", func(t *testing.T){
		{{range $i, $example := .InvalidExamples}}
		{{$example}}
		idx{{$i}}, err := NewIndex{{$idx.IdxName}}(mt, {{range $idx.ConstructParams}}{{with.}}
			{{.Name}},{{end}}{{end}}
		)
		assert.NotNil(t, err)
		assert.Nil(t, idx{{$i}})
		{{end}}
	})
}
{{end}}{{end}}
`))

var indexSearchParamTemplate = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT
// This file is generated by go generate 

package entity

import (
	"errors"
)

{{ range .Indexes }}{{with.}}
var _ SearchParam = &Index{{.IdxName}}SearchParam{}

// Index{{.IdxName}}SearchParam search param struct for index type {{.IdxType}}
type Index{{.IdxName}}SearchParam struct { //auto generated fields
	baseSearchParams
	{{range .SearchParams}}{{with.}}
	{{.Name}} {{.ParamType}}{{end}}{{end}}
}

// NewIndex{{.IdxName}}SearchParam create index search param
func NewIndex{{.IdxName}}SearchParam({{range .SearchParams}}{{with .}}
	{{.Name}} {{.ParamType}},
{{end}}{{end}}) (*Index{{.IdxName}}SearchParam, error) {
	// auto generate parameters validation code, if any{{range .SearchParams}}{{with .}}
	{{.ValidationCode}}
	{{end}}{{end}}
	sp := &Index{{.IdxName}}SearchParam{
		baseSearchParams: newBaseSearchParams(),
	}
	
	//auto generated setting{{range .SearchParams}}{{with .}}
	sp.params["{{.Name}}"] = {{.Name}}{{end}}{{end}}

	return sp, nil
}
{{end}}{{end}}
`))

var indexSearchParamTestTemplate = template.Must(template.New("").Parse(`// Code generated by go generate; DO NOT EDIT
// This file is generated by go generate

package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

{{range $idx := .Indexes}}{{with.}}
func TestIndex{{.IdxName}}SearchParam(t *testing.T) {
	{{range .SearchParams}}{{with.}}
	var {{.Name}} {{.ParamType}}{{end}}{{end}}

	t.Run("valid usage case", func(t *testing.T){
		{{range $i, $example := .ValidSearchParams}}
		{{$example}}
		idx{{$i}}, err := NewIndex{{$idx.IdxName}}SearchParam({{range $idx.SearchParams}}{{with.}}
			{{.Name}},{{end}}{{end}}
		)
		assert.Nil(t, err)
		assert.NotNil(t, idx{{$i}})
		assert.NotNil(t, idx{{$i}}.Params())
		{{end}}
	})
	{{ if .InvalidSearchParams }}
	t.Run("invalid usage case", func(t *testing.T){
		{{range $i, $example := .InvalidSearchParams}}
		{{$example}}
		idx{{$i}}, err := NewIndex{{$idx.IdxName}}SearchParam({{range $idx.SearchParams}}{{with.}}
			{{.Name}},{{end}}{{end}}
		)
		assert.NotNil(t, err)
		assert.Nil(t, idx{{$i}})
		{{end}}
	})
	{{end}}
}
{{end}}{{end}}
`))

type idxDef struct {
	IdxName             string
	IdxType             entity.IndexType
	VectorSupport       int8
	ConstructParams     []idxParam
	SearchParams        []idxParam
	ValidExamples       []string // valid value examples, used in tests
	InvalidExamples     []string // invalid value examples, used in tests
	ValidSearchParams   []string
	InvalidSearchParams []string
}

type idxParam struct {
	Name           string
	Type           string // default int
	ValidationRule string // support format `[min, max]`, `in (...)` `{other param}===0 (mod self)`
}

// ParamType param definition type
func (ip idxParam) ParamType() string {
	if ip.Type == "" {
		return "int"
	}
	return ip.Type
}

// ValidationCode auto generate validation code by validation rule
func (ip idxParam) ValidationCode() string {
	if ip.Type == "" {
		ip.Type = "int"
	}
	switch {
	case strings.HasPrefix(ip.ValidationRule, "in ("): // TODO change to regex
		raw := ip.ValidationRule[4 : len(ip.ValidationRule)-1]
		return fmt.Sprintf(`validRange := []%s{%s}
	%sOk := false
	for _, v := range validRange {
		if v == %s {
			%sOk = true
			break
		}
	}
	if !%sOk {
		return nil, errors.New("%s not valid")
	}`, ip.Type, raw, ip.Name, ip.Name, ip.Name, ip.Name, ip.Name)
	case strings.Contains(ip.ValidationRule, "==="): // TODO change to regex
		//TODO NOT IMPLEMENT YET
	case strings.Contains(ip.ValidationRule, "[") &&
		strings.Contains(ip.ValidationRule, "]"): // TODO change to regex
		// very restrict format contstraint to `[min, max]`
		// not extra space allowed
		raw := ip.ValidationRule[1 : len(ip.ValidationRule)-1]
		parts := strings.Split(raw, ",")
		min := strings.TrimSpace(parts[0])
		max := strings.TrimSpace(parts[1])
		return fmt.Sprintf(`if %s < %s {
		return nil, errors.New("%s has to be in range [%s, %s]")
	}
	if %s > %s {
		return nil, errors.New("%s has to be in range [%s, %s]")
	}`, ip.Name, min, ip.Name, min, max, ip.Name, max, ip.Name, min, max)
	default:
		return ""
	}

	return ""
}

type vectorTypeSupport int8

const (
	floatVectorSupport  vectorTypeSupport = 1
	binaryVectorSupport vectorTypeSupport = 1 << 1
)

func main() {
	f, err := os.OpenFile("indexes_gen.go", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer f.Close()
	ft, err := os.OpenFile("indexes_gen_test.go", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer ft.Close()
	fp, err := os.OpenFile("indexes_search_param_gen.go", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer fp.Close()
	fpt, err := os.OpenFile("indexes_search_param_gen_test.go", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer fpt.Close()

	settings := struct {
		Indexes []idxDef
	}{
		Indexes: []idxDef{
			// FLAT
			{
				IdxName:         "Flat",
				IdxType:         entity.Flat,
				ConstructParams: []idxParam{},
				SearchParams:    []idxParam{},
				ValidExamples: []string{
					"",
				},
				InvalidExamples: []string{},
				ValidSearchParams: []string{
					"",
				},
				InvalidSearchParams: []string{},
			},
			// BIN_FLAT
			{
				IdxName:       "BinFlat",
				IdxType:       entity.BinFlat,
				VectorSupport: int8(binaryVectorSupport),
				ConstructParams: []idxParam{
					{
						Name:           "nlist",
						ValidationRule: "[1, 65536]",
					},
				},
				SearchParams: []idxParam{
					{
						Name:           "nprobe",
						ValidationRule: "[1, 65536]", //[1, nlist]
					},
				},
				ValidExamples: []string{
					"nlist = 10",
				},
				InvalidExamples: []string{
					"nlist = 0",
					"nlist = 65537",
				},
				ValidSearchParams: []string{
					"nprobe = 10",
				},
				InvalidSearchParams: []string{
					"nprobe = 0",
					"nprobe = 65537",
				},
			},
			// IVF_FLAT
			{
				IdxName: "IvfFlat",
				IdxType: entity.IvfFlat,
				ConstructParams: []idxParam{
					{
						Name:           "nlist",
						ValidationRule: "[1, 65536]",
					},
				},
				SearchParams: []idxParam{
					{
						Name:           "nprobe",
						ValidationRule: "[1, 65536]", // [1, nlist], refer to index construct param, not supported yet
					},
				},
				ValidExamples: []string{
					"nlist = 10",
				},
				InvalidExamples: []string{
					"nlist = 0",
					"nlist = 65537",
				},
				ValidSearchParams: []string{
					"nprobe = 10",
				},
				InvalidSearchParams: []string{
					"nprobe = 0",
					"nprobe = 65537",
				},
			},
			// BIN_IVF_FLAT
			{
				IdxName:       "BinIvfFlat",
				IdxType:       entity.BinIvfFlat,
				VectorSupport: int8(binaryVectorSupport),
				ConstructParams: []idxParam{
					{
						Name:           "nlist",
						ValidationRule: "[1, 65536]",
					},
				},
				SearchParams: []idxParam{
					{
						Name:           "nprobe",
						ValidationRule: "[1, 65536]", // [1, nlist], refer to index construct param, not supported yet
					},
				},
				ValidExamples: []string{
					"nlist = 10",
				},
				InvalidExamples: []string{
					"nlist = 0",
					"nlist = 65537",
				},
				ValidSearchParams: []string{
					"nprobe = 10",
				},
				InvalidSearchParams: []string{
					"nprobe = 0",
					"nprobe = 65537",
				},
			},
			// IVF_SQ8
			{
				IdxName: "IvfSQ8",
				IdxType: entity.IvfSQ8,
				ConstructParams: []idxParam{
					{
						Name:           "nlist",
						ValidationRule: "[1, 65536]",
					},
				},
				SearchParams: []idxParam{
					{
						Name:           "nprobe",
						ValidationRule: "[1, 65536]", // [1, nlist], refer to index construct param, not supported yet
					},
				},
				ValidExamples: []string{
					"nlist = 10",
				},
				InvalidExamples: []string{
					"nlist = 0",
					"nlist = 65537",
				},
				ValidSearchParams: []string{
					"nprobe = 10",
				},
				InvalidSearchParams: []string{
					"nprobe = 0",
					"nprobe = 65537",
				},
			},
			// IVF_PQ
			{
				IdxName: "IvfPQ",
				IdxType: entity.IvfPQ,
				ConstructParams: []idxParam{
					{
						Name:           "nlist",
						ValidationRule: "[1, 65536]",
					},
					{
						Name:           "m",
						ValidationRule: "dim===0 (mod self)",
					},
					{
						Name:           "nbits",
						ValidationRule: "[1, 16]",
					},
				},
				SearchParams: []idxParam{
					{
						Name:           "nprobe",
						ValidationRule: "[1, 65536]", // [1, nlist], refer to index construct param, not supported yet
					},
				},
				ValidExamples: []string{
					"nlist, m, nbits = 10, 8, 8",
				},
				InvalidExamples: []string{
					"nlist, m, nbits = 0, 8, 8",
					"nlist, m, nbits = 65537, 8, 8",
					"nlist, m, nbits = 10, 8, 0",
					"nlist, m, nbits = 10, 8, 17",
				},
				ValidSearchParams: []string{
					"nprobe = 10",
				},
				InvalidSearchParams: []string{
					"nprobe = 0",
					"nprobe = 65537",
				},
			},
			// HNSW
			{
				IdxName: "HNSW",
				IdxType: entity.HNSW,
				ConstructParams: []idxParam{
					{
						Name:           "M",
						ValidationRule: "[4, 64]",
					},
					{
						Name:           "efConstruction",
						ValidationRule: "[8, 512]",
					},
				},
				SearchParams: []idxParam{
					{
						Name:           "ef",
						ValidationRule: "[1, 32768]", // [topK, 32768], refer to index construct param, not supported yet
					},
				},
				ValidExamples: []string{
					"M, efConstruction = 16, 40",
				},
				InvalidExamples: []string{
					"M, efConstruction = 3, 40",
					"M, efConstruction = 65, 40",
					"M, efConstruction = 16, 7",
					"M, efConstruction = 16, 513",
				},
				ValidSearchParams: []string{
					"ef = 16",
				},
				InvalidSearchParams: []string{
					"ef = 0",
					"ef = 32769",
				},
			},
			// IVF_HNSW
			{
				IdxName: "IvfHNSW",
				IdxType: entity.IvfHNSW,
				ConstructParams: []idxParam{
					{
						Name:           "nlist",
						ValidationRule: "[1, 65536]",
					},
					{
						Name:           "M",
						ValidationRule: "[4, 64]",
					},
					{
						Name:           "efConstruction",
						ValidationRule: "[8, 512]",
					},
				},
				SearchParams: []idxParam{
					{
						Name:           "nprobe",
						ValidationRule: "[1, 65536]", //[1, nlist]
					},
					{
						Name:           "ef",
						ValidationRule: "[1, 32768]", // [topK, 32768], refer to index construct param, not supported yet
					},
				},
				ValidExamples: []string{
					"nlist, M, efConstruction = 10, 16, 40",
				},
				InvalidExamples: []string{
					"nlist, M, efConstruction = 0, 16, 40",
					"nlist, M, efConstruction = 65537, 16, 40",
					"nlist, M, efConstruction = 10, 3, 40",
					"nlist, M, efConstruction = 10, 65, 40",
					"nlist, M, efConstruction = 10, 16, 7",
					"nlist, M, efConstruction = 10, 16, 513",
				},
				ValidSearchParams: []string{
					"nprobe, ef = 10, 16",
				},
				InvalidSearchParams: []string{
					"nprobe, ef = 0, 16",
					"nprobe, ef = 65537, 16",
					"nprobe, ef = 10, 0",
					"nprobe, ef = 10, 32769",
				},
			},
			{
				IdxName:         "DISKANN",
				IdxType:         entity.DISKANN,
				ConstructParams: []idxParam{},
				SearchParams: []idxParam{
					{
						Name:           "search_list",
						ValidationRule: "[1, 65535]",
					},
				},
				ValidExamples: []string{
					"",
				},
				InvalidExamples: []string{},
				ValidSearchParams: []string{
					"search_list = 30",
				},
				InvalidSearchParams: []string{
					"search_list = 0",
					"search_list = 65537",
				},
			},
			{
				IdxName:         "AUTOINDEX",
				IdxType:         entity.AUTOINDEX,
				ConstructParams: []idxParam{},
				SearchParams: []idxParam{
					{
						Name: "level",
						Type: "interface{}",
					},
				},
				ValidExamples: []string{
					"",
				},
				InvalidExamples: []string{},
				ValidSearchParams: []string{
					"level = 1",
				},
			},
			{
				IdxName: "GPUIvfFlat",
				IdxType: entity.GPUIvfFlat,
				ConstructParams: []idxParam{
					{
						Name:           "nlist",
						ValidationRule: "[1, 65536]",
					},
				},
				SearchParams: []idxParam{
					{
						Name:           "nprobe",
						ValidationRule: "[1, 65536]", // [1, nlist], refer to index construct param, not supported yet
					},
				},
				ValidExamples: []string{
					"nlist = 10",
				},
				InvalidExamples: []string{
					"nlist = 0",
					"nlist = 65537",
				},
				ValidSearchParams: []string{
					"nprobe = 10",
				},
				InvalidSearchParams: []string{
					"nprobe = 0",
					"nprobe = 65537",
				},
			},
			{
				IdxName: "GPUIvfPQ",
				IdxType: entity.GPUIvfPQ,
				ConstructParams: []idxParam{
					{
						Name:           "nlist",
						ValidationRule: "[1, 65536]",
					},
					{
						Name:           "m",
						ValidationRule: "dim===0 (mod self)",
					},
					{
						Name:           "nbits",
						ValidationRule: "[1, 64]",
					},
				},
				SearchParams: []idxParam{
					{
						Name:           "nprobe",
						ValidationRule: "[1, 65536]", // [1, nlist], refer to index construct param, not supported yet
					},
				},
				ValidExamples: []string{
					"nlist, m, nbits = 10, 8, 8",
				},
				InvalidExamples: []string{
					"nlist, m, nbits = 0, 8, 8",
					"nlist, m, nbits = 65537, 8, 8",
					"nlist, m, nbits = 10, 8, 0",
					"nlist, m, nbits = 10, 8, 65",
				},
				ValidSearchParams: []string{
					"nprobe = 10",
				},
				InvalidSearchParams: []string{
					"nprobe = 0",
					"nprobe = 65537",
				},
			},
			{
				IdxName: "SCANN",
				IdxType: entity.SCANN,
				ConstructParams: []idxParam{
					{
						Name:           "nlist",
						ValidationRule: "[1, 65536]",
					},
					{
						Name:           "with_raw_data",
						Type:           "bool",
						ValidationRule: "in (false, true)",
					},
				},
				SearchParams: []idxParam{
					{
						Name:           "nprobe",
						ValidationRule: "[1, 65536]", // [1, nlist], refer to index construct param, not supported yet
					},
					{
						Name:           "reorder_k",
						ValidationRule: "[1, 9223372036854775807]", // [topk, MAX_INT], refer to index construct param, not supported yet
					},
				},
				ValidExamples: []string{
					"nlist = 100",
					"with_raw_data = true",
				},
				InvalidExamples: []string{
					"nlist = 0",
					"nlist = 65537",
				},
				ValidSearchParams: []string{
					"nprobe, reorder_k = 10, 200",
				},
				InvalidSearchParams: []string{
					"nprobe, reorder_k = 0, 200",
					"nprobe, reorder_k = 65537, 200",
					"nprobe, reorder_k = 10, -1",
				},
			},
		},
	}

	indexTemplate.Execute(f, settings)
	indexTestTemplate.Execute(ft, settings)
	indexSearchParamTemplate.Execute(fp, settings)
	indexSearchParamTestTemplate.Execute(fpt, settings)
}
