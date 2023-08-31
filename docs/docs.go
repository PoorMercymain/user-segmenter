// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/api/reports/{filename}": {
            "get": {
                "description": "Запрос для получения отчета по истории сегментов пользователя в формате csv",
                "produces": [
                    "text/csv"
                ],
                "tags": [
                    "Reports"
                ],
                "summary": "Запрос чтения отчета по истории сегментов пользователя",
                "parameters": [
                    {
                        "type": "string",
                        "example": "report12345.csv",
                        "description": "report filename",
                        "name": "filename",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/api/segment": {
            "post": {
                "description": "Запрос для создания сегмента по уникальному названию",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Segments"
                ],
                "summary": "Запрос для создания нового сегмента",
                "parameters": [
                    {
                        "description": "segment info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_PoorMercymain_user-segmenter_internal_domain.Slug"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "202": {
                        "description": "Accepted"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "409": {
                        "description": "Conflict"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            },
            "delete": {
                "description": "Запрос для удаления сегмента из списка существующих сегментов по уникальному названию",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Segments"
                ],
                "summary": "Запрос для удаления сегмента",
                "parameters": [
                    {
                        "description": "segment info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_PoorMercymain_user-segmenter_internal_domain.SlugNoPercent"
                        }
                    }
                ],
                "responses": {
                    "202": {
                        "description": "Accepted"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/api/user": {
            "post": {
                "description": "Запрос для обновления списка сегментов пользователя",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Запрос обновления сегментов пользователя",
                "parameters": [
                    {
                        "description": "user segment info",
                        "name": "input",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/github_com_PoorMercymain_user-segmenter_internal_domain.UserUpdate"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/api/user-history/{id}": {
            "get": {
                "description": "Запрос для создания отчета по истории сегментов пользователя в формате csv",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Reports"
                ],
                "summary": "Запрос формирования отчета по истории сегментов пользователя",
                "parameters": [
                    {
                        "type": "string",
                        "example": "1",
                        "description": "user id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "example": "2023-9",
                        "description": "start date",
                        "name": "start",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "example": "2023-9",
                        "description": "end date",
                        "name": "end",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "example": "2023-9",
                        "description": "exact date",
                        "name": "exact",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "400": {
                        "description": "Bad Request"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        },
        "/api/user/{id}": {
            "get": {
                "description": "Запрос для получения списка сегментов пользователя",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Запрос чтения сегментов пользователя",
                "parameters": [
                    {
                        "type": "string",
                        "example": "1",
                        "description": "user id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "204": {
                        "description": "No Content"
                    },
                    "404": {
                        "description": "Not Found"
                    },
                    "500": {
                        "description": "Internal Server Error"
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_PoorMercymain_user-segmenter_internal_domain.Slug": {
            "type": "object",
            "properties": {
                "percent": {
                    "type": "integer",
                    "example": 10
                },
                "slug": {
                    "type": "string",
                    "example": "SEGMENT_NAME"
                }
            }
        },
        "github_com_PoorMercymain_user-segmenter_internal_domain.SlugNoPercent": {
            "type": "object",
            "properties": {
                "slug": {
                    "type": "string",
                    "example": "SEGMENT_NAME"
                }
            }
        },
        "github_com_PoorMercymain_user-segmenter_internal_domain.UserUpdate": {
            "type": "object",
            "properties": {
                "slugs_to_add": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "SEGMENT_NAME"
                    ]
                },
                "slugs_to_delete": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "ANOTHER_SEGMENT_NAME"
                    ]
                },
                "ttl": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "2023-09-30T20:19:05+03:00"
                    ]
                },
                "user_id": {
                    "type": "string",
                    "example": "1"
                }
            }
        }
    },
    "tags": [
        {
            "description": "Группа запросов для управления списком существующих сегментов",
            "name": "Segments"
        },
        {
            "description": "Группа запросов для управления сегментами пользователя",
            "name": "Users"
        },
        {
            "description": "Группа запросов для работы с отчетами по истории сегментов пользователя",
            "name": "Reports"
        }
    ]
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "UserSegmenter API",
	Description:      "Сервис динамического сегментирования пользователей",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
