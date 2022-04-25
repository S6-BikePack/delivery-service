// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
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
        "/api/deliveries": {
            "get": {
                "description": "gets all deliveries in the system",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "get all deliveries",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/domain.Delivery"
                            }
                        }
                    }
                }
            },
            "post": {
                "description": "creates a new delivery",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "create delivery",
                "parameters": [
                    {
                        "description": "Add delivery",
                        "name": "rider",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.BodyCreateDelivery"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.ResponseCreateDelivery"
                        }
                    }
                }
            }
        },
        "/api/deliveries/radius/{latlon}": {
            "get": {
                "description": "gets a delivery from the system based on the distance to the given point",
                "produces": [
                    "application/json"
                ],
                "summary": "get delivery by distance",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Latitude,Longitude",
                        "name": "latlon",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "radius of search in meters (default = 1000)",
                        "name": "radius",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.Delivery"
                        }
                    }
                }
            }
        },
        "/api/deliveries/{id}": {
            "get": {
                "description": "gets a delivery from the system by its ID",
                "produces": [
                    "application/json"
                ],
                "summary": "get delivery",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Delivery id",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.Delivery"
                        }
                    }
                }
            }
        },
        "/api/deliveries/{id}/complete": {
            "get": {
                "description": "completes a delivery",
                "produces": [
                    "application/json"
                ],
                "summary": "complete delivery",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.Delivery"
                        }
                    }
                }
            }
        },
        "/api/deliveries/{id}/rider": {
            "post": {
                "description": "assigns a rider to a delivery",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "assign rider",
                "parameters": [
                    {
                        "description": "Assign rider",
                        "name": "delivery",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/dto.BodyAssignRider"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/dto.ResponseAssignRider"
                        }
                    }
                }
            }
        },
        "/api/deliveries/{id}/start": {
            "get": {
                "description": "starts a delivery",
                "produces": [
                    "application/json"
                ],
                "summary": "start delivery",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/domain.Delivery"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.Customer": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                }
            }
        },
        "domain.Delivery": {
            "type": "object",
            "properties": {
                "customer": {
                    "$ref": "#/definitions/domain.Customer"
                },
                "destination": {
                    "$ref": "#/definitions/domain.TimeAndPlace"
                },
                "id": {
                    "type": "string"
                },
                "parcel": {
                    "$ref": "#/definitions/domain.Parcel"
                },
                "pickup": {
                    "$ref": "#/definitions/domain.TimeAndPlace"
                },
                "rider": {
                    "$ref": "#/definitions/domain.Rider"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "domain.Dimensions": {
            "type": "object",
            "properties": {
                "depth": {
                    "type": "integer"
                },
                "height": {
                    "type": "integer"
                },
                "width": {
                    "type": "integer"
                }
            }
        },
        "domain.Location": {
            "type": "object",
            "properties": {
                "latitude": {
                    "type": "number"
                },
                "longitude": {
                    "type": "number"
                }
            }
        },
        "domain.Parcel": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "serviceArea": {
                    "type": "integer"
                },
                "size": {
                    "$ref": "#/definitions/domain.Dimensions"
                },
                "weight": {
                    "type": "integer"
                }
            }
        },
        "domain.Rider": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string"
                },
                "serviceArea": {
                    "type": "integer"
                }
            }
        },
        "domain.TimeAndPlace": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "coordinates": {
                    "$ref": "#/definitions/domain.Location"
                },
                "time": {
                    "type": "string"
                }
            }
        },
        "dto.BodyAssignRider": {
            "type": "object",
            "properties": {
                "riderId": {
                    "type": "string"
                }
            }
        },
        "dto.BodyCreateDelivery": {
            "type": "object",
            "properties": {
                "destination": {
                    "$ref": "#/definitions/dto.BodyCreateDeliveryDestination"
                },
                "ownerId": {
                    "type": "string"
                },
                "parcelId": {
                    "type": "string"
                },
                "pickup": {
                    "$ref": "#/definitions/dto.BodyCreateDeliveryPickup"
                }
            }
        },
        "dto.BodyCreateDeliveryDestination": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "coordinates": {
                    "$ref": "#/definitions/domain.Location"
                }
            }
        },
        "dto.BodyCreateDeliveryPickup": {
            "type": "object",
            "properties": {
                "address": {
                    "type": "string"
                },
                "coordinates": {
                    "$ref": "#/definitions/domain.Location"
                },
                "time": {
                    "type": "integer"
                }
            }
        },
        "dto.ResponseAssignRider": {
            "type": "object",
            "properties": {
                "customer": {
                    "$ref": "#/definitions/domain.Customer"
                },
                "destination": {
                    "$ref": "#/definitions/domain.TimeAndPlace"
                },
                "id": {
                    "type": "string"
                },
                "parcel": {
                    "$ref": "#/definitions/domain.Parcel"
                },
                "pickup": {
                    "$ref": "#/definitions/domain.TimeAndPlace"
                },
                "rider": {
                    "$ref": "#/definitions/domain.Rider"
                },
                "status": {
                    "type": "integer"
                }
            }
        },
        "dto.ResponseCreateDelivery": {
            "type": "object",
            "properties": {
                "destination": {
                    "$ref": "#/definitions/domain.TimeAndPlace"
                },
                "owner": {
                    "$ref": "#/definitions/domain.Customer"
                },
                "parcel": {
                    "$ref": "#/definitions/domain.Parcel"
                },
                "pickup": {
                    "$ref": "#/definitions/domain.TimeAndPlace"
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
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
