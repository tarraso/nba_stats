// Package docs Code generated by swaggo/swag. DO NOT EDIT
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
        "/add-players": {
            "post": {
                "description": "Add a new player to the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "players"
                ],
                "summary": "Add a new player",
                "parameters": [
                    {
                        "description": "Player",
                        "name": "player",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.Player"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.Player"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/add-stat": {
            "post": {
                "description": "Add a new game stat to the database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "stats"
                ],
                "summary": "Add a new game stat",
                "parameters": [
                    {
                        "description": "Game Stat",
                        "name": "stat",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/models.GameStat"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/models.GameStat"
                        }
                    },
                    "400": {
                        "description": "Bad request",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/players": {
            "get": {
                "description": "Get a list of all players",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "players"
                ],
                "summary": "List all players",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.Player"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/stat/players/{playerId}": {
            "get": {
                "description": "Get a list of all players",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "players"
                ],
                "summary": "player stats",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "PlayerId",
                        "name": "playerId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.AvgStat"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/stat/teams/{teamId}": {
            "get": {
                "description": "Get a list of all players",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "players"
                ],
                "summary": "team stats",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "teamId",
                        "name": "teamId",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/models.AvgStat"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "models.AvgStat": {
            "type": "object",
            "properties": {
                "avg_assists": {
                    "type": "number"
                },
                "avg_blocks": {
                    "type": "number"
                },
                "avg_fouls": {
                    "type": "number"
                },
                "avg_game_date": {
                    "type": "string"
                },
                "avg_minutes_played": {
                    "type": "number"
                },
                "avg_points": {
                    "type": "number"
                },
                "avg_rebounds": {
                    "type": "number"
                },
                "avg_steals": {
                    "type": "number"
                },
                "avg_turnovers": {
                    "type": "number"
                },
                "player_id": {
                    "type": "integer"
                }
            }
        },
        "models.GameStat": {
            "type": "object",
            "properties": {
                "assists": {
                    "type": "integer"
                },
                "blocks": {
                    "type": "integer"
                },
                "fouls": {
                    "type": "integer"
                },
                "game_date": {
                    "type": "string"
                },
                "minutes_played": {
                    "type": "number"
                },
                "player_id": {
                    "type": "integer"
                },
                "points": {
                    "type": "integer"
                },
                "rebounds": {
                    "type": "integer"
                },
                "steals": {
                    "type": "integer"
                },
                "turnovers": {
                    "type": "integer"
                }
            }
        },
        "models.Player": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                },
                "team_id": {
                    "description": "New field for foreign key",
                    "type": "integer"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "",
	Host:             "",
	BasePath:         "",
	Schemes:          []string{},
	Title:            "",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
