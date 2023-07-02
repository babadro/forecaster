// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "swagger": "2.0",
  "info": {
    "description": "This is a sample server for a poll management API.",
    "title": "Poll API",
    "version": "1.0.0"
  },
  "paths": {
    "/polls": {
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "Create poll",
        "operationId": "createPoll",
        "parameters": [
          {
            "description": "Poll object that needs to be added.",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/Poll"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Poll created"
          },
          "400": {
            "description": "Invalid input"
          }
        }
      }
    },
    "/polls/{pollId}": {
      "get": {
        "produces": [
          "application/json"
        ],
        "summary": "Get poll by ID",
        "operationId": "getPollById",
        "parameters": [
          {
            "type": "integer",
            "format": "int32",
            "description": "ID of poll to return",
            "name": "pollId",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "successful operation",
            "schema": {
              "$ref": "#/definitions/Poll"
            }
          },
          "400": {
            "description": "Invalid ID supplied"
          },
          "404": {
            "description": "Poll not found"
          }
        }
      },
      "put": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "Update an existing poll",
        "operationId": "updatePoll",
        "parameters": [
          {
            "type": "integer",
            "format": "int32",
            "description": "ID of poll that needs to be updated",
            "name": "pollId",
            "in": "path",
            "required": true
          },
          {
            "description": "Updated poll object",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/Poll"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Poll updated"
          },
          "400": {
            "description": "Invalid input"
          },
          "404": {
            "description": "Poll not found"
          }
        }
      },
      "delete": {
        "produces": [
          "application/json"
        ],
        "summary": "Deletes a poll",
        "operationId": "deletePoll",
        "parameters": [
          {
            "type": "integer",
            "format": "int32",
            "description": "Poll id to delete",
            "name": "pollId",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "400": {
            "description": "Invalid ID supplied"
          },
          "404": {
            "description": "Poll not found"
          }
        }
      }
    }
  },
  "definitions": {
    "Poll": {
      "type": "object",
      "properties": {
        "description": {
          "type": "string"
        },
        "finish": {
          "type": "string",
          "format": "date-time"
        },
        "id": {
          "type": "integer",
          "format": "int32"
        },
        "start": {
          "type": "string",
          "format": "date-time"
        },
        "title": {
          "type": "string"
        }
      },
      "xml": {
        "name": "Poll"
      }
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "swagger": "2.0",
  "info": {
    "description": "This is a sample server for a poll management API.",
    "title": "Poll API",
    "version": "1.0.0"
  },
  "paths": {
    "/polls": {
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "Create poll",
        "operationId": "createPoll",
        "parameters": [
          {
            "description": "Poll object that needs to be added.",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/Poll"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Poll created"
          },
          "400": {
            "description": "Invalid input"
          }
        }
      }
    },
    "/polls/{pollId}": {
      "get": {
        "produces": [
          "application/json"
        ],
        "summary": "Get poll by ID",
        "operationId": "getPollById",
        "parameters": [
          {
            "type": "integer",
            "format": "int32",
            "description": "ID of poll to return",
            "name": "pollId",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "description": "successful operation",
            "schema": {
              "$ref": "#/definitions/Poll"
            }
          },
          "400": {
            "description": "Invalid ID supplied"
          },
          "404": {
            "description": "Poll not found"
          }
        }
      },
      "put": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "summary": "Update an existing poll",
        "operationId": "updatePoll",
        "parameters": [
          {
            "type": "integer",
            "format": "int32",
            "description": "ID of poll that needs to be updated",
            "name": "pollId",
            "in": "path",
            "required": true
          },
          {
            "description": "Updated poll object",
            "name": "body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/Poll"
            }
          }
        ],
        "responses": {
          "200": {
            "description": "Poll updated"
          },
          "400": {
            "description": "Invalid input"
          },
          "404": {
            "description": "Poll not found"
          }
        }
      },
      "delete": {
        "produces": [
          "application/json"
        ],
        "summary": "Deletes a poll",
        "operationId": "deletePoll",
        "parameters": [
          {
            "type": "integer",
            "format": "int32",
            "description": "Poll id to delete",
            "name": "pollId",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "400": {
            "description": "Invalid ID supplied"
          },
          "404": {
            "description": "Poll not found"
          }
        }
      }
    }
  },
  "definitions": {
    "Poll": {
      "type": "object",
      "properties": {
        "description": {
          "type": "string"
        },
        "finish": {
          "type": "string",
          "format": "date-time"
        },
        "id": {
          "type": "integer",
          "format": "int32"
        },
        "start": {
          "type": "string",
          "format": "date-time"
        },
        "title": {
          "type": "string"
        }
      },
      "xml": {
        "name": "Poll"
      }
    }
  }
}`))
}
