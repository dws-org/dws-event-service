#!/bin/bash

echo "Regenerating Prisma client..."

# Generate Prisma client
go run github.com/steebchen/prisma-client-go generate

echo "Prisma client regenerated successfully!"
echo "You may need to restart your Go application for changes to take effect."
