# Go JWT system

## What is this about?

This is a small project dedicated to creating a JWT and refresh token pairs based on the user GUID. We have two REST endpoints that are used to get the token pairs and to refresh the JWT token using the access token.

## How to run it?

You can run the application by installing Docker first and navigating in your command line to the project folder. After that run ```docker compose up```
There are going to be two containers, one is the PostgreSQL database and the other is the API server. To run the tests, run ```go test ```.

Enjoy :D