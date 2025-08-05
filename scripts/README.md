# Scripts

This directory contains utility scripts for the API Gateway project.

## Available Scripts

### `test_rate_limit.sh`
A bash script to test the token bucket rate limiting functionality of the API Gateway.

**Usage:**
```bash
./scripts/test_rate_limit.sh
```

**What it does:**
- Tests basic rate limiting by making multiple requests
- Shows token consumption and remaining tokens
- Tests rapid requests to see rate limiting in action
- Tests token refill after waiting periods

**Prerequisites:**
- API Gateway must be running on `localhost:8080`
- Redis should be enabled for rate limiting to work

**Example output:**
```
🧪 Testing Token Bucket Rate Limiting
=====================================

📊 Testing token consumption:
Request 1:
  Remaining tokens: 99/100
  Refill rate: 1.67 tokens/second

⚡ Testing rapid requests (should consume tokens quickly):
Request 1: 98 tokens remaining | Request 2: 97 tokens remaining | ...

⏰ Testing token refill (waiting 2 seconds):
Before wait:
  Remaining tokens: 95
After wait:
  Remaining tokens: 98

✅ Token bucket rate limiting test completed!
``` 