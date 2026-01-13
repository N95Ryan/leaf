#!/bin/bash
# Test script for Leaf project using gotestsum

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}🌱 Leaf Test Suite${NC}"
echo ""

# Check if gotestsum is installed
if ! command -v gotestsum &> /dev/null; then
    echo -e "${YELLOW}⚠️  gotestsum is not installed${NC}"
    echo "Installing gotestsum..."
    go install gotest.tools/gotestsum@latest
    echo -e "${GREEN}✓ gotestsum installed${NC}"
    echo ""
fi

# Run tests with gotestsum
echo -e "${GREEN}Running tests...${NC}"
gotestsum --format pkgname-and-test-fails -- -cover -v ./...

echo ""
echo -e "${GREEN}✓ Tests completed${NC}"
