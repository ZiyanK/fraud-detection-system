# Fraud Detection System

(WIP)

## Description
This repository implements a transactional fraud detection service taking transactions via a message broker as well as an REST API. The idea is to develop this to further add metrics and proper logging.

## Tools/Technologies Used
1. [Gin](https://github.com/gin-gonic/gin) web framework
2. PostgreSQL (database)
3. Apache Kafka (message broker)
4. [sqlx](https://github.com/jmoiron/sqlx) is used to interact with the database
5. [zap](https://github.com/uber-go/zap) for logging
6. [viper](https://github.com/spf13/viper) to get the env variables
7. Docker and docker-compose for easy deployment