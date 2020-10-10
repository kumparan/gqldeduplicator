package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/graph-gophers/graphql-go"

	"github.com/kumparan/gqldeduplicator"
	"github.com/kumparan/gqldeduplicator/example/graph-gophers/resolver"
)

type (
	// Handler represent graphql handler
	Handler struct {
		Schema *graphql.Schema
	}

	request struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}
)

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var request request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := h.Schema.Exec(r.Context(), request.Query, request.OperationName, request.Variables)
	if r.URL.Query().Get("deduplicate") == "1" {
		result, err := gqldeduplicator.Deflate(response.Data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if result.Deflated {
			response.Data = result.Data
			w.Header().Set("GraphQL-Deduplicator", "1")
		}
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(responseJSON)
}

func main() {
	schema := graphql.MustParseSchema(gqlSchema, &resolver.Resolver{})

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(gqlPlaygroundPage)
	}))
	http.Handle("/query", &Handler{Schema: schema})
	log.Println("Running...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var gqlSchema = `
type Child {
    id: ID!
    name: String!
    age: Int!
}

type Parent {
    id: ID!
    name: String!
    child: Child
}

type Query {
    parents: [Parent!]!
}
`

var gqlPlaygroundPage = []byte(`
<!DOCTYPE html>
<html>

<head>
  <meta charset=utf-8/>
  <meta name="viewport" content="user-scalable=no, initial-scale=1.0, minimum-scale=1.0, maximum-scale=1.0, minimal-ui">
  <title>GraphQL Playground</title>
  <link rel="stylesheet" href="//cdn.jsdelivr.net/npm/graphql-playground-react/build/static/css/index.css" />
  <link rel="shortcut icon" href="//cdn.jsdelivr.net/npm/graphql-playground-react/build/favicon.png" />
  <script src="//cdn.jsdelivr.net/npm/graphql-playground-react/build/static/js/middleware.js"></script>
</head>

<body>
  <div id="root">
    <style>
      body {
        background-color: rgb(23, 42, 58);
        font-family: Open Sans, sans-serif;
        height: 90vh;
      }

      #root {
        height: 100%;
        width: 100%;
        display: flex;
        align-items: center;
        justify-content: center;
      }

      .loading {
        font-size: 32px;
        font-weight: 200;
        color: rgba(255, 255, 255, .6);
        margin-left: 20px;
      }

      img {
        width: 78px;
        height: 78px;
      }

      .title {
        font-weight: 400;
      }
    </style>
    <img src='//cdn.jsdelivr.net/npm/graphql-playground-react/build/logo.png' alt=''>
    <div class="loading"> Loading
      <span class="title">GraphQL Playground</span>
    </div>
  </div>
  <script>window.addEventListener('load', function (event) {
      GraphQLPlayground.init(document.getElementById('root'), {
        // options as 'endpoint' belong here
		endpoint: '/query'
      })
    })</script>
</body>

</html>
`)
