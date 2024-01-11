// Copyright 2019 The go-jjang Authors
// This file is part of the go-jjang library.
//
// The go-jjang library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-jjang library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-jjang library. If not, see <http://www.gnu.org/licenses/>.

package graphql

import (
	"encoding/json"
	"net/http"

	"github.com/jjang-network/go-jjang/eth/filters"
	"github.com/jjang-network/go-jjang/internal/ethapi"
	"github.com/jjang-network/go-jjang/node"
	"github.com/graph-gophers/graphql-go"
)

type handler struct {
	Schema *graphql.Schema
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := h.Schema.Exec(r.Context(), params.Query, params.OperationName, params.Variables)
	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(response.Errors) > 0 {
		w.WriteHeader(http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(responseJSON)
}

// New constructs a new GraphQL service instance.
func New(stack *node.Node, backend ethapi.Backend, filterSystem *filters.FilterSystem, cors, vhosts []string) error {
	_, err := newHandler(stack, backend, filterSystem, cors, vhosts)
	return err
}

// newHandler returns a new `http.Handler` that will answer GraphQL queries.
// It additionally exports an interactive query browser on the / endpoint.
func newHandler(stack *node.Node, backend ethapi.Backend, filterSystem *filters.FilterSystem, cors, vhosts []string) (*handler, error) {
	q := Resolver{backend, filterSystem}

	s, err := graphql.ParseSchema(schema, &q)
	if err != nil {
		return nil, err
	}
	h := handler{Schema: s}
	handler := node.NewHTTPHandlerStack(h, cors, vhosts, nil)

	stack.RegisterHandler("GraphQL UI", "/graphql/ui", GraphiQL{})
	stack.RegisterHandler("GraphQL", "/graphql", handler)
	stack.RegisterHandler("GraphQL", "/graphql/", handler)

	return &h, nil
}