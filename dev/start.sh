#!/bin/bash

set -e

APP_NAME="raus-damit"

# Load environment variables from .env file
if [ -f .env ]; then
  echo "🔧 Loading environment variables from .env..."
  export $(grep -v '^#' .env | xargs)
else
  echo "⚠️ No .env file found. Please ensure variables are exported manually."
fi

# Always build the app
echo "🔨 Building $APP_NAME..."
go build -o $APP_NAME ../

# Run the app
echo "🚀 Running $APP_NAME..."
./$APP_NAME
