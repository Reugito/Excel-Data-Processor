# Golang CRUD Application with Excel Import, MySQL, and Redis

## Overview

This project is a Golang application that imports data from an Excel file, stores it into a MySQL database, and caches it in Redis. It provides a simple CRUD (Create, Read, Update) system to manage the imported data.

## Features

- Upload and parse Excel (.xlsx) file
- Store parsed data in MySQL
- Cache data in Redis with expiration
- View, edit, and store records

## Prerequisites

- Golang 1.19+
- MySQL
- Redis
- Postman or any other API testing tool

## Setup

### 1. Install Dependencies
```bash
go mod tidy
```

### 2. Add MySQL and Redis Config Details in config.yaml
```yaml
mysql:
  host: "localhost"
  port: "3306"
  username: "root"
  password: "226594"
  database: "test88"

redis:
  address: "localhost:6379"
  password: ""
  db: 0
```

### 3. Provide the yaml file path in main.go 

### 4. Run `go run main.go` to start the server.

### 5. Used Dependencies
```makefile
	github.com/gin-gonic/gin v1.10.0
	github.com/go-playground/validator/v10 v10.20.0
	github.com/go-redis/redis/v8 v8.11.5
	github.com/google/uuid v1.6.0
	github.com/xuri/excelize/v2 v2.8.1
	gopkg.in/yaml.v2 v2.4.0
	gorm.io/driver/mysql v1.5.6
	gorm.io/gorm v1.25.10

```


## Endpoints

### 1. Import Data

**Description:** This endpoint allows you to upload an Excel file for importing data into the MySQL database and caching it in Redis.

**URL:** `/import`

**Method:** `POST`

**Request:**

- **Header:**
    - `Content-Type: multipart/form-data`
- **Body:**
    - `file`: The Excel file to be uploaded with xlsx extension

**Example Curl Command:**
```shell
curl --location 'localhost:8080/v1/contact/upload' \
--form 'file=@"/C:/Users/raosa/Downloads/Sample_Employee_data_xlsx (2).xlsx"'
```

### 2. View Imported Data

**Description:** This endpoint allows you to view the imported data. If the data is available in Redis cache, it retrieves from there; otherwise, it fetches from the MySQL database

**URL:** `/getAll`

**Method:** `GET`

**Example Curl Command:**
```shell
curl --location 'localhost:8080/v1/contact/getAll'
```

### 3. Edit a Record

**Description:** This endpoint allows you to edit a specific record by ID. The changes will be reflected in both MySQL and Redis cache (will reflect in Redis if cache available if not then it will add it from db to cache).

**URL:** `/update/:refId`

**Method:** `PUT`

**Example Curl Command:**

```shell
curl --location --request PUT 'localhost:8080/v1/contact/update/5522cdb6-1921-4e3d-b988-8cb5bfdf73aa' \
--header 'Content-Type: application/json' \
--data-raw '{
            "ref_id": "5522cdb6-1921-4e3d-b988-8cb5bfdf73aa",
            "first_name": "rao",
            "last_name": "jk",
            "company_name": "ABCMfg Corp",
            "address": "529 mtd Rd #3680",
            "city": "Pune Walden",
            "county": "Essex",
            "postal": "CB11 4DJ",
            "phone": "01260-744622",
            "email": "earnestine_casper@hotmail.com",
            "web": "https://www.abc.co.uk"
        }'
```

### 4. Get Contact By RefId

**Description:** This endpoint allows you to view the imported data by RefId 

**URL:** `/getByRefId`

**Method:** `GET`

**Example Curl Command:**
```shell
curl --location 'http://localhost:8080/v1/contact/getByRefId/5522cdb6-1921-4e3d-b988-8cb5bfdf73aa'
```
