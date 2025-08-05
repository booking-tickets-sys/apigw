#!/bin/bash

echo "🧪 Testing Token Bucket Rate Limiting"
echo "====================================="
echo

# Test basic rate limiting
echo "📊 Testing token consumption:"
for i in {1..5}; do
    echo "Request $i:"
    response=$(curl -s -I http://localhost:8080/health)
    remaining=$(echo "$response" | grep "X-Ratelimit-Remaining" | cut -d' ' -f2)
    limit=$(echo "$response" | grep "X-Ratelimit-Limit" | cut -d' ' -f2)
    refill_rate=$(echo "$response" | grep "X-Ratelimit-Refillrate" | cut -d' ' -f2)
    echo "  Remaining tokens: $remaining/$limit"
    echo "  Refill rate: $refill_rate tokens/second"
    echo
done

# Test rapid requests to see rate limiting
echo "⚡ Testing rapid requests (should consume tokens quickly):"
for i in {1..10}; do
    response=$(curl -s -I http://localhost:8080/health)
    remaining=$(echo "$response" | grep "X-Ratelimit-Remaining" | cut -d' ' -f2)
    echo -n "Request $i: $remaining tokens remaining | "
    if [ "$i" -eq 10 ]; then
        echo ""
    fi
done
echo

# Test waiting for token refill
echo "⏰ Testing token refill (waiting 2 seconds):"
echo "Before wait:"
response=$(curl -s -I http://localhost:8080/health)
remaining=$(echo "$response" | grep "X-Ratelimit-Remaining" | cut -d' ' -f2)
echo "  Remaining tokens: $remaining"

echo "Waiting 2 seconds for token refill..."
sleep 2

echo "After wait:"
response=$(curl -s -I http://localhost:8080/health)
remaining=$(echo "$response" | grep "X-Ratelimit-Remaining" | cut -d' ' -f2)
echo "  Remaining tokens: $remaining"
echo

echo "✅ Token bucket rate limiting test completed!" 