#!/bin/bash

# ScanDiff Quick Start Script
# This script helps you get started with ScanDiff quickly

set -e

echo "=========================================="
echo "  ScanDiff - Host Snapshot Diff Tool"
echo "=========================================="
echo ""

# Check for Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    echo "   Visit: https://docs.docker.com/get-docker/"
    exit 1
fi

# Check for Docker Compose
if ! command -v docker compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    echo "   Visit: https://docs.docker.com/compose/install/"
    exit 1
fi

echo "✅ Docker is installed"
echo "✅ Docker Compose is installed"
echo ""

# Build and start the application
echo "🚀 Building and starting ScanDiff..."
echo ""

docker compose up --build -d

echo ""
echo "⏳ Waiting for services to be ready..."
sleep 5

# Check if services are running
if docker compose ps | grep -q "scandiff-backend.*Up"; then
    echo "✅ Backend is running"
else
    echo "❌ Backend failed to start"
    docker compose logs backend
    exit 1
fi

if docker compose ps | grep -q "scandiff-frontend.*Up"; then
    echo "✅ Frontend is running"
else
    echo "❌ Frontend failed to start"
    docker compose logs frontend
    exit 1
fi

echo ""
echo "=========================================="
echo "  🎉 ScanDiff is now running!"
echo "=========================================="
echo ""
echo "📍 Access the application:"
echo "   Frontend:  http://localhost:3000"
echo "   Backend:   http://localhost:8080/api"
echo "   Health:    http://localhost:8080/health"
echo ""
echo "📁 Sample data is available in: ./sample_data/"
echo ""
echo "📖 Quick Usage:"
echo "   1. Go to http://localhost:3000"
echo "   2. Upload a snapshot from ./sample_data/"
echo "   3. Click 'Hosts' to view uploaded snapshots"
echo "   4. Select a host and choose 2 snapshots to compare"
echo ""
echo "🛑 To stop the application:"
echo "   docker compose down"
echo ""
echo "📚 For more information, see README.md"
echo "=========================================="
