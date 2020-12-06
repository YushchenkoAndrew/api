exports.swaggerDocument = {
  info: {
    title: "NodeJS API",
    description: "Node.js API with Express, Sequelize & MySQL",
    version: "0.0.1",
  },

  host: `${process.env.API_HOST || "127.0.0.1"}:${process.env.API_PORT || 31337}`,
  basePath: "/",
  swagger: "2.0",

  tags: [
    {
      name: "Visitor",
      description: "Store Countries",
    },
    {
      name: "Views",
      description: "Store amount of views spent on the website",
    },
    {
      name: "Github",
      description: "Store repos info from Github API",
    },
    {
      name: "UsedLanguages",
      description: "Used languages in repos defined by GithubId",
    },
  ],

  schemes: ["http", "https"],
  consumes: ["application/json"],
  produces: ["application/json"],

  paths: {
    "/api/Visitor": {
      get: {
        tags: ["Visitor"],
        summary: "Get all data from Visitor Table",
        parameters: [
          {
            in: "query",
            name: "id",
            required: false,
            schema: {
              type: "integer",
              example: 1,
            },
            description: "id",
          },

          {
            in: "query",
            name: "Country",
            required: false,
            schema: {
              type: "string",
              example: "UK",
            },
            description: "Country name",
          },

          {
            in: "query",
            name: "ip",
            required: false,
            schema: {
              type: "string",
              format: "date",
              example: "127.0.0.1",
            },
            description: "ip address",
          },

          {
            in: "query",
            name: "Visit_Date",
            required: false,
            schema: {
              type: "string",
              example: "2020-10-09",
            },
            description: "Last visit date from this country",
          },

          {
            in: "query",
            name: "Count",
            required: false,
            schema: {
              type: "integer",
              example: 1,
            },
            description: "Number of visits from this the Country",
          },
        ],

        responses: {
          200: {
            description: "OK",
            schema: {
              $ref: "#/definitions/Visitor",
            },
          },

          400: {
            description: "Bad Request",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Invalid request message parameters",
                },
              },
            },
          },

          500: {
            description: "Internal Server Error",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Field 'Count' doesn't have a default value",
                },
              },
            },
          },
        },
      },

      post: {
        tags: ["Visitor"],
        summary: "Add new data to Visitor",
        parameters: [
          {
            in: "body",
            name: "body",
            description: "Updating data",
            required: true,
            schema: {
              $ref: "#/definitions/Visitor",
            },
          },
        ],

        responses: {
          201: {
            description: "OK",
            schema: {
              $ref: "#/definitions/Visitor",
            },
          },

          400: {
            description: "Bad Request",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Invalid request message parameters",
                },
              },
            },
          },

          500: {
            description: "Internal Server Error",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Field 'Count' doesn't have a default value",
                },
              },
            },
          },
        },
      },

      put: {
        tags: ["Visitor"],
        summary: "Update data in Visitor",
        parameters: [
          {
            in: "query",
            name: "id",
            required: false,
            schema: {
              type: "integer",
              example: 1,
            },
            description: "id",
          },

          {
            in: "query",
            name: "Country",
            required: false,
            schema: {
              type: "string",
              example: "UK",
            },
            description: "Country name",
          },

          {
            in: "query",
            name: "ip",
            required: false,
            schema: {
              type: "string",
              format: "date",
              example: "127.0.0.1",
            },
            description: "ip address",
          },

          {
            in: "query",
            name: "Visit_Date",
            required: false,
            schema: {
              type: "string",
              example: "2020-10-09",
            },
            description: "Last visit date from this country",
          },

          {
            in: "query",
            name: "Count",
            required: false,
            schema: {
              type: "integer",
              example: 1,
            },
            description: "Number of visits from this the Country",
          },
        ],

        responses: {
          204: {
            description: "No Content",
          },

          304: {
            description: "Not Modified",
          },

          400: {
            description: "Bad Request",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Invalid request message parameters",
                },
              },
            },
          },

          500: {
            description: "Internal Server Error",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Field 'Count' doesn't have a default value",
                },
              },
            },
          },
        },
      },

      delete: {
        tags: ["Visitor"],
        summary: "Delete data from Visitor",
        parameters: [
          {
            in: "query",
            name: "id",
            required: false,
            schema: {
              type: "integer",
              example: 1,
            },
            description: "id",
          },

          {
            in: "query",
            name: "Country",
            required: false,
            schema: {
              type: "string",
              example: "UK",
            },
            description: "Country name",
          },

          {
            in: "query",
            name: "ip",
            required: false,
            schema: {
              type: "string",
              format: "date",
              example: "127.0.0.1",
            },
            description: "ip address",
          },

          {
            in: "query",
            name: "Visit_Date",
            required: false,
            schema: {
              type: "string",
              example: "2020-10-09",
            },
            description: "Last visit date from this country",
          },

          {
            in: "query",
            name: "Count",
            required: false,
            schema: {
              type: "integer",
              example: 1,
            },
            description: "Number of visits from this the Country",
          },
        ],

        responses: {
          204: {
            description: "No Content",
          },

          304: {
            description: "Not Modified",
          },

          400: {
            description: "Bad Request",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Invalid request message parameters",
                },
              },
            },
          },

          500: {
            description: "Internal Server Error",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Field 'Count' doesn't have a default value",
                },
              },
            },
          },
        },
      },
    },

    "/api/Visitor/{id}": {
      parameters: [
        {
          in: "path",
          name: "id",
          required: true,
          schema: {
            type: "integer",
          },
          description: "id",
        },
      ],

      get: {
        tags: ["Visitor"],
        summary: "Get data by id from Visitor",

        responses: {
          200: {
            description: "OK",
            schema: {
              $ref: "#/definitions/Visitor",
            },
          },

          400: {
            description: "Bad Request",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Invalid request message parameters",
                },
              },
            },
          },

          500: {
            description: "Internal Server Error",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Field 'Count' doesn't have a default value",
                },
              },
            },
          },
        },
      },

      put: {
        tags: ["Visitor"],
        summary: "Update data by id in Visitor",

        parameters: [
          {
            in: "body",
            name: "body",
            description: "Updating data",
            required: true,
            schema: {
              $ref: "#/definitions/Visitor",
            },
          },
        ],

        responses: {
          204: {
            description: "No Content",
          },

          304: {
            description: "Not Modified",
          },

          400: {
            description: "Bad Request",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Invalid request message parameters",
                },
              },
            },
          },

          500: {
            description: "Internal Server Error",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Field 'Count' doesn't have a default value",
                },
              },
            },
          },
        },
      },

      delete: {
        tags: ["Visitor"],
        summary: "Delete data by id from Visitor",
        responses: {
          204: {
            description: "No Content",
          },

          304: {
            description: "Not Modified",
          },

          400: {
            description: "Bad Request",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Invalid request message parameters",
                },
              },
            },
          },

          500: {
            description: "Internal Server Error",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Field 'Count' doesn't have a default value",
                },
              },
            },
          },
        },
      },
    },

    "/api/Views": {
      get: {
        tags: ["Views"],
        summary: "Get all data from Views Table",
        parameters: [
          {
            in: "query",
            name: "id",
            required: false,
            schema: {
              type: "integer",
              example: 1,
            },
            description: "id",
          },

          {
            in: "query",
            name: "Curr_Date",
            required: false,
            schema: {
              type: "string",
              example: "2020-10-09",
            },
            description: "Current Date",
          },

          {
            in: "query",
            name: "Count",
            required: false,
            schema: {
              type: "integer",
              example: 1,
            },
            description: "Number of visits from this the Country",
          },
        ],

        responses: {
          200: {
            description: "OK",
            schema: {
              $ref: "#/definitions/Views",
            },
          },

          400: {
            description: "Bad Request",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Invalid request message parameters",
                },
              },
            },
          },

          500: {
            description: "Internal Server Error",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Field 'Count' doesn't have a default value",
                },
              },
            },
          },
        },
      },

      post: {
        tags: ["Views"],
        summary: "Add new data to Views",
        parameters: [
          {
            in: "body",
            name: "body",
            description: "Updating data",
            required: true,
            schema: {
              $ref: "#/definitions/Views",
            },
          },
        ],

        responses: {
          201: {
            description: "OK",
            schema: {
              $ref: "#/definitions/Views",
            },
          },

          400: {
            description: "Bad Request",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Invalid request message parameters",
                },
              },
            },
          },

          500: {
            description: "Internal Server Error",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Field 'Count' doesn't have a default value",
                },
              },
            },
          },
        },
      },

      put: {
        tags: ["Views"],
        summary: "Update data in Views",
        parameters: [
          {
            in: "query",
            name: "id",
            required: false,
            schema: {
              type: "integer",
              example: 1,
            },
            description: "id",
          },

          {
            in: "query",
            name: "Curr_Date",
            required: false,
            schema: {
              type: "string",
              example: "2020-10-09",
            },
            description: "Current Date",
          },

          {
            in: "query",
            name: "Count",
            required: false,
            schema: {
              type: "integer",
              example: 1,
            },
            description: "Number of visits from this the Country",
          },
        ],

        responses: {
          204: {
            description: "No Content",
          },

          304: {
            description: "Not Modified",
          },

          400: {
            description: "Bad Request",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Invalid request message parameters",
                },
              },
            },
          },

          500: {
            description: "Internal Server Error",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Field 'Count' doesn't have a default value",
                },
              },
            },
          },
        },
      },

      delete: {
        tags: ["Views"],
        summary: "Delete data from Views",
        parameters: [
          {
            in: "query",
            name: "id",
            required: false,
            schema: {
              type: "integer",
              example: 1,
            },
            description: "id",
          },

          {
            in: "query",
            name: "Curr_Date",
            required: false,
            schema: {
              type: "string",
              example: "2020-10-09",
            },
            description: "Current Date",
          },

          {
            in: "query",
            name: "Count",
            required: false,
            schema: {
              type: "integer",
              example: 1,
            },
            description: "Number of visits from this the Country",
          },
        ],

        responses: {
          204: {
            description: "No Content",
          },

          304: {
            description: "Not Modified",
          },

          400: {
            description: "Bad Request",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Invalid request message parameters",
                },
              },
            },
          },

          500: {
            description: "Internal Server Error",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Field 'Count' doesn't have a default value",
                },
              },
            },
          },
        },
      },
    },

    "/api/Views/{id}": {
      parameters: [
        {
          in: "path",
          name: "id",
          required: true,
          schema: {
            type: "integer",
          },
          description: "id",
        },
      ],

      get: {
        tags: ["Views"],
        summary: "Get data by id from Views",

        responses: {
          200: {
            description: "OK",
            schema: {
              $ref: "#/definitions/Views",
            },
          },

          400: {
            description: "Bad Request",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Invalid request message parameters",
                },
              },
            },
          },

          500: {
            description: "Internal Server Error",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Field 'Count' doesn't have a default value",
                },
              },
            },
          },
        },
      },

      put: {
        tags: ["Views"],
        summary: "Update data by id in Views",

        parameters: [
          {
            in: "body",
            name: "body",
            description: "Updating data",
            required: true,
            schema: {
              $ref: "#/definitions/Views",
            },
          },
        ],

        responses: {
          204: {
            description: "No Content",
          },

          304: {
            description: "Not Modified",
          },

          400: {
            description: "Bad Request",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Invalid request message parameters",
                },
              },
            },
          },

          500: {
            description: "Internal Server Error",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Field 'Count' doesn't have a default value",
                },
              },
            },
          },
        },
      },

      delete: {
        tags: ["Views"],
        summary: "Delete data by id from Views",
        responses: {
          204: {
            description: "No Content",
          },

          304: {
            description: "Not Modified",
          },

          400: {
            description: "Bad Request",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Invalid request message parameters",
                },
              },
            },
          },

          500: {
            description: "Internal Server Error",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Field 'Count' doesn't have a default value",
                },
              },
            },
          },
        },
      },
    },

    "/api/Github": {
      get: {
        tags: ["Github"],
        summary: "Get all data from Github Table",
        parameters: [
          {
            in: "query",
            name: "id",
            required: false,
            schema: {
              type: "integer",
              example: 1,
            },
            description: "id",
          },

          {
            in: "query",
            name: "Name",
            required: false,
            schema: {
              type: "string",
              example: "template",
            },
            description: "repo's name",
          },

          {
            in: "query",
            name: "UpdateAt",
            required: false,
            schema: {
              type: "string",
              example: "2020-10-09",
            },
            description: "Updated Date",
          },

          {
            in: "query",
            name: "Context",
            required: false,
            schema: {
              type: "string",
              maximum: 255,
            },
            description: "Context that taken from README file",
          },
        ],

        responses: {
          200: {
            description: "OK",
            schema: {
              $ref: "#/definitions/Github",
            },
          },

          400: {
            description: "Bad Request",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Invalid request message parameters",
                },
              },
            },
          },

          500: {
            description: "Internal Server Error",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Field 'Count' doesn't have a default value",
                },
              },
            },
          },
        },
      },

      post: {
        tags: ["Github"],
        summary: "Add new data to Github",
        parameters: [
          {
            in: "body",
            name: "body",
            description: "Updating data",
            required: true,
            schema: {
              $ref: "#/definitions/Github",
            },
          },
        ],

        responses: {
          201: {
            description: "OK",
            schema: {
              $ref: "#/definitions/Github",
            },
          },

          400: {
            description: "Bad Request",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Invalid request message parameters",
                },
              },
            },
          },

          500: {
            description: "Internal Server Error",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Field 'Count' doesn't have a default value",
                },
              },
            },
          },
        },
      },

      put: {
        tags: ["Github"],
        summary: "Update data in Github",
        parameters: [
          {
            in: "query",
            name: "id",
            required: false,
            schema: {
              type: "integer",
              example: 1,
            },
            description: "id",
          },

          {
            in: "query",
            name: "Name",
            required: false,
            schema: {
              type: "string",
              example: "template",
            },
            description: "repo's name",
          },

          {
            in: "query",
            name: "UpdateAt",
            required: false,
            schema: {
              type: "string",
              example: "2020-10-09",
            },
            description: "Updated Date",
          },

          {
            in: "query",
            name: "Context",
            required: false,
            schema: {
              type: "string",
              maximum: 255,
            },
            description: "Context that taken from README file",
          },
        ],

        responses: {
          204: {
            description: "No Content",
          },

          304: {
            description: "Not Modified",
          },

          400: {
            description: "Bad Request",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Invalid request message parameters",
                },
              },
            },
          },

          500: {
            description: "Internal Server Error",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Field 'Count' doesn't have a default value",
                },
              },
            },
          },
        },
      },

      delete: {
        tags: ["Github"],
        summary: "Delete data from Github",
        parameters: [
          {
            in: "query",
            name: "id",
            required: false,
            schema: {
              type: "integer",
              example: 1,
            },
            description: "id",
          },

          {
            in: "query",
            name: "Name",
            required: false,
            schema: {
              type: "string",
              example: "template",
            },
            description: "repo's name",
          },

          {
            in: "query",
            name: "UpdateAt",
            required: false,
            schema: {
              type: "string",
              example: "2020-10-09",
            },
            description: "Updated Date",
          },

          {
            in: "query",
            name: "Context",
            required: false,
            schema: {
              type: "string",
              maximum: 255,
            },
            description: "Context that taken from README file",
          },
        ],

        responses: {
          204: {
            description: "No Content",
          },

          304: {
            description: "Not Modified",
          },

          400: {
            description: "Bad Request",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Invalid request message parameters",
                },
              },
            },
          },

          500: {
            description: "Internal Server Error",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Field 'Count' doesn't have a default value",
                },
              },
            },
          },
        },
      },
    },

    "/api/Github/{id}": {
      parameters: [
        {
          in: "path",
          name: "id",
          required: true,
          schema: {
            type: "integer",
          },
          description: "id",
        },
      ],

      get: {
        tags: ["Github"],
        summary: "Get data by id from Github",

        responses: {
          200: {
            description: "OK",
            schema: {
              $ref: "#/definitions/Github",
            },
          },

          400: {
            description: "Bad Request",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Invalid request message parameters",
                },
              },
            },
          },

          500: {
            description: "Internal Server Error",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Field 'Count' doesn't have a default value",
                },
              },
            },
          },
        },
      },

      put: {
        tags: ["Github"],
        summary: "Update data by id in Github",

        parameters: [
          {
            in: "body",
            name: "body",
            description: "Updating data",
            required: true,
            schema: {
              $ref: "#/definitions/Github",
            },
          },
        ],

        responses: {
          204: {
            description: "No Content",
          },

          304: {
            description: "Not Modified",
          },

          400: {
            description: "Bad Request",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Invalid request message parameters",
                },
              },
            },
          },

          500: {
            description: "Internal Server Error",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Field 'Count' doesn't have a default value",
                },
              },
            },
          },
        },
      },

      delete: {
        tags: ["Github"],
        summary: "Delete data by id from Github",
        responses: {
          204: {
            description: "No Content",
          },

          304: {
            description: "Not Modified",
          },

          400: {
            description: "Bad Request",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Invalid request message parameters",
                },
              },
            },
          },

          500: {
            description: "Internal Server Error",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Field 'Count' doesn't have a default value",
                },
              },
            },
          },
        },
      },
    },

    "/api/UsedLanguages": {
      get: {
        tags: ["UsedLanguages"],
        summary: "Get all data from UsedLanguages Table",
        parameters: [
          {
            in: "query",
            name: "id",
            required: false,
            schema: {
              type: "integer",
              example: 1,
            },
            description: "id",
          },

          {
            in: "query",
            name: "Name",
            required: false,
            schema: {
              type: "string",
              example: "JavaScript",
            },
            description: "Programming language name",
          },

          {
            in: "query",
            name: "Size",
            required: false,
            schema: {
              type: "integer",
              example: 55,
            },
            description: "Files size",
          },

          {
            in: "query",
            name: "GithubId",
            required: false,
            schema: {
              type: "integer",
              example: 1,
            },
            description: "Key for Github primaryID (id)",
          },
        ],

        responses: {
          200: {
            description: "OK",
            schema: {
              $ref: "#/definitions/UsedLanguages",
            },
          },

          400: {
            description: "Bad Request",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Invalid request message parameters",
                },
              },
            },
          },

          500: {
            description: "Internal Server Error",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Field 'Count' doesn't have a default value",
                },
              },
            },
          },
        },
      },

      post: {
        tags: ["UsedLanguages"],
        summary: "Add new data to UsedLanguages",
        parameters: [
          {
            in: "body",
            name: "body",
            description: "Updating data",
            required: true,
            schema: {
              $ref: "#/definitions/UsedLanguages",
            },
          },
        ],

        responses: {
          201: {
            description: "OK",
            schema: {
              $ref: "#/definitions/UsedLanguages",
            },
          },

          400: {
            description: "Bad Request",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Invalid request message parameters",
                },
              },
            },
          },

          500: {
            description: "Internal Server Error",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Field 'Count' doesn't have a default value",
                },
              },
            },
          },
        },
      },

      put: {
        tags: ["UsedLanguages"],
        summary: "Update data in UsedLanguages",
        parameters: [
          {
            in: "query",
            name: "id",
            required: false,
            schema: {
              type: "integer",
              example: 1,
            },
            description: "id",
          },

          {
            in: "query",
            name: "Name",
            required: false,
            schema: {
              type: "string",
              example: "JavaScript",
            },
            description: "Programming language name",
          },

          {
            in: "query",
            name: "Size",
            required: false,
            schema: {
              type: "integer",
              example: 55,
            },
            description: "Files size",
          },

          {
            in: "query",
            name: "GithubId",
            required: false,
            schema: {
              type: "integer",
              example: 1,
            },
            description: "Key for Github primaryID (id)",
          },
        ],

        responses: {
          204: {
            description: "No Content",
          },

          304: {
            description: "Not Modified",
          },

          400: {
            description: "Bad Request",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Invalid request message parameters",
                },
              },
            },
          },

          500: {
            description: "Internal Server Error",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Field 'Count' doesn't have a default value",
                },
              },
            },
          },
        },
      },

      delete: {
        tags: ["UsedLanguages"],
        summary: "Delete data from UsedLanguages",
        parameters: [
          {
            in: "query",
            name: "id",
            required: false,
            schema: {
              type: "integer",
              example: 1,
            },
            description: "id",
          },

          {
            in: "query",
            name: "Name",
            required: false,
            schema: {
              type: "string",
              example: "JavaScript",
            },
            description: "Programming language name",
          },

          {
            in: "query",
            name: "Size",
            required: false,
            schema: {
              type: "integer",
              example: 55,
            },
            description: "Files size",
          },

          {
            in: "query",
            name: "GithubId",
            required: false,
            schema: {
              type: "integer",
              example: 1,
            },
            description: "Key for Github primaryID (id)",
          },
        ],

        responses: {
          204: {
            description: "No Content",
          },

          304: {
            description: "Not Modified",
          },

          400: {
            description: "Bad Request",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Invalid request message parameters",
                },
              },
            },
          },

          500: {
            description: "Internal Server Error",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Field 'Count' doesn't have a default value",
                },
              },
            },
          },
        },
      },
    },

    "/api/UsedLanguages/{id}": {
      parameters: [
        {
          in: "path",
          name: "id",
          required: true,
          schema: {
            type: "integer",
          },
          description: "id",
        },
      ],

      get: {
        tags: ["UsedLanguages"],
        summary: "Get data by id from UsedLanguages",

        responses: {
          200: {
            description: "OK",
            schema: {
              $ref: "#/definitions/UsedLanguages",
            },
          },

          400: {
            description: "Bad Request",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Invalid request message parameters",
                },
              },
            },
          },

          500: {
            description: "Internal Server Error",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Field 'Count' doesn't have a default value",
                },
              },
            },
          },
        },
      },

      put: {
        tags: ["UsedLanguages"],
        summary: "Update data by id in UsedLanguages",

        parameters: [
          {
            in: "body",
            name: "body",
            description: "Updating data",
            required: true,
            schema: {
              $ref: "#/definitions/UsedLanguages",
            },
          },
        ],

        responses: {
          204: {
            description: "No Content",
          },

          304: {
            description: "Not Modified",
          },

          400: {
            description: "Bad Request",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Invalid request message parameters",
                },
              },
            },
          },

          500: {
            description: "Internal Server Error",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Field 'Count' doesn't have a default value",
                },
              },
            },
          },
        },
      },

      delete: {
        tags: ["UsedLanguages"],
        summary: "Delete data by id from UsedLanguages",
        responses: {
          204: {
            description: "No Content",
          },

          304: {
            description: "Not Modified",
          },

          400: {
            description: "Bad Request",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Invalid request message parameters",
                },
              },
            },
          },

          500: {
            description: "Internal Server Error",
            schema: {
              properties: {
                message: {
                  type: "string",
                  example: "Field 'Count' doesn't have a default value",
                },
              },
            },
          },
        },
      },
    },
  },

  definitions: {
    Visitor: {
      required: ["Country", "ip", "Visit_Date", "Count"],
      properties: {
        id: {
          type: "integer",
          required: false,
          example: 1,
        },

        Country: {
          type: "string",
          maxLength: 2,
          example: "US",
        },

        ip: {
          type: "string",
          example: "127.0.0.1",
        },

        Visit_Date: {
          type: "string",
          format: "date",
          example: "2020-10-09",
        },

        Count: {
          type: "integer",
          example: 1,
        },
      },
    },

    Views: {
      required: ["Curr_Date", "Count"],
      properties: {
        id: {
          type: "integer",
          required: false,
          example: 1,
        },

        Curr_Date: {
          type: "string",
          format: "date",
          example: "2020-10-09",
        },

        Count: {
          type: "integer",
          example: 1,
        },
      },
    },

    Github: {
      required: ["Name", "UpdateAt", "Context"],
      properties: {
        id: {
          type: "integer",
          required: false,
          example: 1,
        },

        Name: {
          type: "string",
          example: "template",
        },

        UpdateAt: {
          type: "string",
          format: "date",
          example: "2020-10-09",
        },

        Context: {
          type: "string",
          maximum: 255,
        },
      },
    },

    UsedLanguages: {
      required: ["Name", "Size", "GithubId"],
      properties: {
        id: {
          type: "integer",
          required: false,
          example: 1,
        },

        Name: {
          type: "string",
          example: "JavaScript",
        },

        Size: {
          type: "integer",
          example: 55,
        },

        GithubId: {
          type: "integer",
        },
      },
    },
  },
};
