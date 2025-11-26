#!/bin/bash

docker build -t dezzwwi/booking-service .

docker push dezzwwi/booking-service

echo "Готово"