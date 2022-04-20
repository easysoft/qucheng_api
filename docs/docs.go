// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "QuCheng Pangu Team"
        },
        "license": {
            "name": "Z PUBLIC LICENSE 1.2"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/cne/app/install": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "安装接口",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "应用管理"
                ],
                "summary": "安装接口",
                "parameters": [
                    {
                        "type": "string",
                        "description": "jwtToken",
                        "name": "Authorization",
                        "in": "header"
                    },
                    {
                        "type": "string",
                        "description": "staticToken",
                        "name": "X-Auth-Token",
                        "in": "header"
                    },
                    {
                        "description": "meta",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.AppCreateModel"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/router.response2xx"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/router.response5xx"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.AppCreateModel": {
            "type": "object",
            "required": [
                "chart",
                "cluster",
                "name",
                "namespace"
            ],
            "properties": {
                "chart": {
                    "type": "string"
                },
                "cluster": {
                    "type": "string"
                },
                "name": {
                    "type": "string"
                },
                "namespace": {
                    "type": "string"
                }
            }
        },
        "router.response2xx": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "object"
                },
                "message": {
                    "type": "string"
                },
                "pagination": {
                    "type": "object"
                },
                "success": {
                    "type": "boolean"
                }
            }
        },
        "router.response5xx": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "object"
                },
                "message": {
                    "type": "string"
                },
                "pagination": {
                    "type": "object"
                },
                "success": {
                    "type": "boolean",
                    "default": false
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0.0",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "CNE API",
	Description: "CNE API.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
