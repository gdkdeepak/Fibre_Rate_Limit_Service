Fibre Rate-Limit Service: Testing Checklist
1. Start the server
cd C:\Users\deepakkumar\go\projects\Sample\Fibre_Rate_Limit_Service
go run ./cmd/server


Server should start on http://localhost:8080.

2. Add a Token Bucket Limiter

Request:

$body = @{
    name        = "test_bucket"
    type        = "token-bucket"
    capacity    = 5
    refill_rate = 1
    refill_every = 10  # seconds
    ttl         = 60   # seconds
} | ConvertTo-Json

Invoke-RestMethod -Uri http://localhost:8080/admin/limiters -Method POST -Body $body -ContentType "application/json"


Expected Result:

{"message":"limiter added/updated successfully"}

3. Add a Policy Rule

Request:

$body = @{
    route = "/some-route"
    header = "X-Secret"
    value  = "123"
} | ConvertTo-Json

Invoke-RestMethod -Uri http://localhost:8080/admin/policies -Method POST -Body $body -ContentType "application/json"


Expected Result:

{"message":"policy added/updated successfully"}

4. Test /check endpoint with correct header

Request:

$body = @{
    client_id = "client1"
    route     = "/some-route"
} | ConvertTo-Json

$headers = @{ "X-Secret" = "123" }

$response = Invoke-RestMethod -Uri http://localhost:8080/check -Method POST -Body $body -ContentType "application/json" -Headers $headers
$response


Expected Result:

{
  "allowed": true,
  "remaining": 4,
  "reset_at": "2025-11-24T15:51:23+05:30",
  "reason": ""
}

5. Test /check endpoint with missing/wrong header

Request: Use header X-Secret: 000 or remove header.
Expected Result:

{
  "allowed": false,
  "reason": "Header X-Secret must equal 123"
}

6. Test Rate Limiting

Run multiple requests quickly to exceed the token bucket capacity:

for ($i = 1; $i -le 10; $i++) {
    $response = Invoke-RestMethod -Uri http://localhost:8080/check -Method POST -Body $body -ContentType "application/json" -Headers $headers
    Write-Output "Request $i: Allowed=$($response.allowed), Remaining=$($response.remaining), ResetAt=$($response.reset_at), Reason=$($response.reason)"
    Start-Sleep -Seconds 1
}


Expected Behavior:

First 5 requests: Allowed = true.

Next requests: Allowed = false, Reason = "rate limit exceeded".

remaining decrements correctly, reset_at shows correct refill time.

7. Verify /admin/snapshot endpoint
$response = Invoke-RestMethod -Uri http://localhost:8080/admin/snapshot -Method GET
$response | ConvertTo-Json -Depth 5


Expected Result:

Shows all limiters with token counts.

Reflects the current state of /some-route.

8. Test live update of limiter

Update limiter capacity or refill rate:

$body = @{
    name        = "test_bucket"
    type        = "token-bucket"
    capacity    = 10
    refill_rate = 2
    refill_every = 5
    ttl         = 60
} | ConvertTo-Json

Invoke-RestMethod -Uri http://localhost:8080/admin/limiters -Method POST -Body $body -ContentType "application/json"


Repeat /check requests and confirm behavior reflects new config immediately.

9. Test live update of policy

Update header requirement:

$body = @{
    route = "/some-route"
    header = "X-Secret"
    value  = "456"
} | ConvertTo-Json

Invoke-RestMethod -Uri http://localhost:8080/admin/policies -Method POST -Body $body -ContentType "application/json"


/check should now require X-Secret=456.

10. Concurrency Test

Use a simple PowerShell loop or multiple terminals to send simultaneous requests to /check.

Confirm:

No race conditions.

Token counts are accurate.

Policies still enforced correctly.

===========================================================================================

✅ Completion Criteria

All endpoints respond as expected.

Token bucket behaves correctly under load.

Policy rules are enforced.

Snapshot accurately reflects limiter states.

Live updates to limiters and policies work without server restart.

No compilation or runtime errors.