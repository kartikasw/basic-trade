# Basic Trade

## Overview

This project is a RESTful API built with GoLang. It is designed to showcase a backend system for managing products and their variants, with authentication. It supports essential CRUD operations (Create, Read, Update, Delete) to efficiently handle product data.

## Tech Stack

- Go
- Gin
- PostgreSQL with pgx driver
- Migrate
- SQLC
- Go JWT
- Cloudinary
- Testing purpose: testify, gomock, dockertest

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Cloudinary
- Docker
- PostgreSQL

### Installation

1. **Clone the repository:**

    ```sh
    git clone https://github.com/kartikasw/basic-trade
    cd your-repo-name
    ```

2. **Set up the environment variables:**

    Create a `.env` file in the root directory of the project. Refer to `.example.env` for the content. 

    > [NOTES]
    >
    > - Ensure that your private key and public key are in the Base64 string format of RSA PEM keys. See [OpenSSL-RSA](https://www.openssl.org/docs/manmaster/man1/openssl-rsa.html) or any RSA Key Generator for generating these keys (for simple/personal project purposes only).
    > - The Migration URL should point to the location of your migration file. In this case, it is located in the migration folder: `file://migration`.

## API Routes

Here is a list of existing routes available in the API:

| Method   | URL              | Description                   |
| :------- |:---------------- | :---------------------------- |
| POST     | /auth/register   | Register new account          |
| POST     | /auth/login      | Login with registered account |
| POST     | /products        | Create new product            |
| PUT      | /products/:uuid  | Update existing product       |
| GET      | /products/:uuid  | Get product's detail          |
| GET      | /products        | Get all or search products    |
| DELETE   | /products/:uuid  | Delete existing product       |
| POST     | /variants        | Create new variant            |
| PUT      | /variants/:uuid  | Update existing variant       |
| GET      | /variants/:uuid  | Get variant's detail          |
| GET      | /variants        | Get all or search variants    |
| DELETE   | /variants/:uuid  | Delete existing variant       |
