package docs

import (
	"github.com/swaggo/swag"
)

var doc = `{
    "swagger": "2.0",
    "info": {
        "description": "This is a Caching Microservice server.",
        "title": "Caching API",
        
        "version": "1.0"
    },
    "host": "127.0.0.1:8000",
	"schemes": ["http"],
    "paths": {
			  "/view/page/{pageno}":{
				"get": {
				"tags": ["CachingApiExplorer"],
				"parameters":[{
								"name": "pageno",
								"in": "path",
								"required": true,
								"type": "integer",
					            "format": "int32",
								"description": "Page No"
                             }],
				"produces": ["application/json"],
				"responses": {
					"200": {
						"description": "OK",
						"schema": {
							"type": "array",
							"items": {
								"type": "string"
							}
						}
					}
				}
			}
		},
			  "/insert/{data}":{
                "put":{
				"tags": ["CachingApiExplorer"],
				"parameters":[{
								"name": "data",
								"in": "path",
								"required": true,
								"type": "string",
								"description": "Data to be added in Cache"
                             }],
                "produces": ["application/json"],
                 "responses": {
					"200": {
						"description": "OK",
						"schema": {
							"type": "string",
							"items": {
								"type": "string"
							}
						}
					}
				}
          }
		},
				"/notify/{msg}":{
                "put":{
				"tags": ["CachingApiExplorer"],
				"parameters":[{
								"name": "msg",
								"in": "path",
								"required": true,
								"type": "string",
								"description": "Pass msg reload to reload cache from DB"
                             }],
                "produces": ["application/json"],
                 "responses": {
					"200": {
						"description": "OK",
						"schema": {
							"type": "string",
							"items": {
								"type": "string"
							}
						}
					}
				}
          }
		}
			

}
}`

type s struct{}

func (s *s) ReadDoc() string {
	return doc
}
func init() {
	swag.Register(swag.Name, &s{})
}
