#!/bin/bash

echo "⏪ Reverting last migration..."

goose -dir "db" -env ".env" down