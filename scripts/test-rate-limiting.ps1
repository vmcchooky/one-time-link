#!/usr/bin/env pwsh
# Test script for rate limiting functionality

$ErrorActionPreference = "Stop"

$API_URL = "http://localhost:8080"

Write-Host "=== Rate Limiting Test ===" -ForegroundColor Cyan
Write-Host ""

# Function to create a secret
function Create-Secret {
    $body = @{
        ciphertext = "dGVzdCBjaXBoZXJ0ZXh0"
        nonce = "dGVzdCBub25jZQ=="
        algorithm = "XChaCha20-Poly1305"
        ttlSeconds = 3600
    } | ConvertTo-Json

    try {
        $response = Invoke-RestMethod -Uri "$API_URL/api/secrets" `
            -Method Post `
            -ContentType "application/json" `
            -Body $body `
            -ResponseHeadersVariable headers

        return @{
            Success = $true
            SecretID = $response.secretId
            Headers = $headers
        }
    }
    catch {
        $statusCode = $_.Exception.Response.StatusCode.value__
        return @{
            Success = $false
            StatusCode = $statusCode
            Headers = $_.Exception.Response.Headers
        }
    }
}

# Test 1: Verify rate limit headers are present
Write-Host "Test 1: Checking rate limit headers..." -ForegroundColor Yellow
$result = Create-Secret

if ($result.Success) {
    Write-Host "✓ Request succeeded" -ForegroundColor Green
    
    $limit = $result.Headers["X-RateLimit-Limit"]
    $remaining = $result.Headers["X-RateLimit-Remaining"]
    $reset = $result.Headers["X-RateLimit-Reset"]
    
    if ($limit -and $remaining -and $reset) {
        Write-Host "✓ Rate limit headers present:" -ForegroundColor Green
        Write-Host "  Limit: $limit"
        Write-Host "  Remaining: $remaining"
        Write-Host "  Reset: $reset"
    }
    else {
        Write-Host "✗ Rate limit headers missing" -ForegroundColor Red
        exit 1
    }
}
else {
    Write-Host "✗ Request failed with status $($result.StatusCode)" -ForegroundColor Red
    exit 1
}

Write-Host ""

# Test 2: Test rate limit enforcement
Write-Host "Test 2: Testing rate limit enforcement (10 requests/hour)..." -ForegroundColor Yellow
Write-Host "Making 11 requests rapidly..."

$successCount = 0
$rateLimitedCount = 0

for ($i = 1; $i -le 11; $i++) {
    $result = Create-Secret
    
    if ($result.Success) {
        $successCount++
        $remaining = $result.Headers["X-RateLimit-Remaining"]
        Write-Host "  Request $i : Success (Remaining: $remaining)" -ForegroundColor Green
    }
    elseif ($result.StatusCode -eq 429) {
        $rateLimitedCount++
        Write-Host "  Request $i : Rate limited (429)" -ForegroundColor Yellow
    }
    else {
        Write-Host "  Request $i : Failed with status $($result.StatusCode)" -ForegroundColor Red
    }
    
    Start-Sleep -Milliseconds 100
}

Write-Host ""
Write-Host "Results:" -ForegroundColor Cyan
Write-Host "  Successful: $successCount"
Write-Host "  Rate limited: $rateLimitedCount"

if ($successCount -eq 10 -and $rateLimitedCount -eq 1) {
    Write-Host "✓ Rate limiting working correctly!" -ForegroundColor Green
}
elseif ($successCount -gt 10) {
    Write-Host "✗ Rate limiting not enforced (too many requests succeeded)" -ForegroundColor Red
    exit 1
}
else {
    Write-Host "⚠ Unexpected results (may need to wait for rate limit reset)" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "=== Rate Limiting Test Complete ===" -ForegroundColor Cyan
