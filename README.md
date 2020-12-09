# NodeJS RESTful API

## Features

- [x] Framework: Express
- [x] API Routes: Express
- [x] HTTP Status Codes
- [x] ORM: Sequelize
- [x] Request handler:
  - [x] JSON
  - [x] XML
- [x] CRUD
  - [x] CRUD: Controllers
  - [x] CRUD: Services
  - [x] CRUD: Repositories
- [x] Basic npm scripts
- [x] Design / Documentation
  - [x] Files Structure
  - [x] Unit test: Mocha/Chai
  - [x] Logging
  - [x] Error handler
  - [x] Swagger
- [x] Docker Container description
- [x] Kubernetes description

## How to use this API

- Copy .env_template to .env
- Fill .env file with your configurations
- Change the /modules files if needed
- The documentation for the CRUD Requests located in localhost:31337/api/doc
- Start api

  ```
    npm run
  ```

- Debug api

  ```
  npm run dev
  ```

- Run the test

  ```
  npm test
  ```

- Create and run a docker container
  ```
    docker build -t api_server .
    docker run -p 31337:31337 api_server
    # Now you can your server on: localhost:31337/api
  ```
- Run Kubernetes
  ```
    kubectl apply -f kubernetes/set
    kubectl apply -f kubernetes
  ```
- For getting a simple navigation just send: GET localhost:31337/api
