#!/bin/env bash

echo "GET http://localhost:8080/api/v1/asdasdasd" | \
  vegeta attack -rate 50 -duration 10s -header "Authorization: ${TOKEN}" | \
  vegeta report
