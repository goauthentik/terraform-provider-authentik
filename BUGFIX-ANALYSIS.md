# Bug Analysis: URL Path Component Stripped in Terraform Provider

**Issue:** XPO-772 - Provider strips path component when Authentik is deployed at a subpath

**Date:** 2025-11-13

**Severity:** High - Blocks usage with path-based deployments

## Problem Summary

When the Terraform provider is configured with a URL containing a path component (e.g., `https://api.example.com/sso/`), the provider strips the path and makes API requests to the wrong endpoint.

### Expected Behavior
```
Provider URL: https://api.example.com/sso/
API Request:  https://api.example.com/sso/api/v3/core/applications/
```

### Actual Behavior
```
Provider URL: https://api.example.com/sso/
API Request:  https://api.example.com/api/v3/core/applications/  ❌ Missing /sso/
```

### Result
- 404 errors from ingress controller
- `Provider produced inconsistent result after apply` errors
- Unable to manage Authentik resources when deployed at subpath

## Root Cause Analysis

### Location
`pkg/provider/provider.go`, lines 254-260 in function `providerConfigure()`

### Code Issue
```go
akURL, err := url.Parse(apiURL)
if err != nil {
    return nil, diag.FromErr(err)
}

config := api.NewConfiguration()
config.Debug = true
config.UserAgent = fmt.Sprintf("authentik-terraform@%s", version)
config.Host = akURL.Host      // ✓ Sets hostname: "api.example.com"
config.Scheme = akURL.Scheme  // ✓ Sets scheme: "https"
// ❌ MISSING: Path component "/sso/" is never set
```

### Why This Happens

1. **URL Parsing:** The code correctly parses the full URL including the path
2. **Configuration:** Only `Host` and `Scheme` are extracted and set
3. **Path Ignored:** `akURL.Path` (containing `/sso/`) is completely ignored
4. **API Client:** The OpenAPI-generated client constructs URLs using only Host + Scheme + API paths
5. **Result:** Base path `/sso/` is lost, requests go to wrong endpoint

## Impact Assessment

### Affected Deployments
- Any Authentik instance hosted at a subpath (e.g., `/sso/`, `/auth/`, etc.)
- Reverse proxy configurations with path-based routing
- Multi-service deployments sharing a domain via path segments
- Kubernetes ingress configurations with path prefixes

### Common Use Cases
1. **Shared Domain Pattern:**
   - `https://api.example.com/sso/` → Authentik
   - `https://api.example.com/app/` → Application
   - `https://api.example.com/api/` → Other API

2. **Ingress Path Routing:**
   ```yaml
   - path: /sso
     pathType: Prefix
     backend:
       service:
         name: authentik
         port: 9000
   ```

3. **Reverse Proxy Configuration:**
   ```nginx
   location /sso/ {
       proxy_pass http://authentik:9000/;
   }
   ```

## Technical Solution

### Required Changes

The provider must preserve and use the URL path component. The OpenAPI-generated API client typically provides a field for setting the base path.

### Implementation
```go
config := api.NewConfiguration()
config.Debug = true
config.UserAgent = fmt.Sprintf("authentik-terraform@%s", version)
config.Host = akURL.Host
config.Scheme = akURL.Scheme

// ADD: Preserve path component
if akURL.Path != "" && akURL.Path != "/" {
    config.BasePath = strings.TrimSuffix(akURL.Path, "/")
}
```

### Verification Steps

1. **Check API Client Structure:**
   - Examine `goauthentik.io/api/v3` Configuration struct
   - Identify correct field name (BasePath, Servers, BaseURL, etc.)
   - Verify how path is used in URL construction

2. **Implementation:**
   - Add path preservation logic
   - Handle edge cases (trailing slashes, empty paths, root path)
   - Maintain backward compatibility for root-path deployments

3. **Testing:**
   - Test with path-based URL: `https://api.example.com/sso/`
   - Test with root URL: `https://api.example.com/`
   - Test with trailing slash variations
   - Verify existing tests still pass

## Improvement Suggestions

### 1. Enhanced URL Validation
Add validation to ensure the provider URL is properly formatted:
```go
// Validate URL format
if !strings.HasSuffix(apiURL, "/") {
    return nil, diag.Errorf("URL must end with trailing slash: %s", apiURL)
}
```

### 2. Better Error Messages
When 404 errors occur, provide diagnostic information:
```go
// Add context to HTTP errors
if resp.StatusCode == 404 {
    return fmt.Errorf("API endpoint not found. If Authentik is deployed at a subpath, ensure the provider URL includes the full path (e.g., https://api.example.com/sso/)")
}
```

### 3. Configuration Documentation
Update provider documentation to explicitly mention path support:
```hcl
provider "authentik" {
  # For root-path deployment
  url   = "https://sso.example.com/"
  
  # For subpath deployment (fixed in v2025.x.x+)
  url   = "https://api.example.com/sso/"
  
  token = var.authentik_token
}
```

### 4. Debug Logging
Add debug output showing the constructed API base URL:
```go
if config.Debug {
    log.Printf("[DEBUG] Authentik API Base URL: %s://%s%s", 
               config.Scheme, config.Host, config.BasePath)
}
```

### 5. Integration Tests
Add test cases covering path-based deployments:
```go
func TestProvider_PathBasedURL(t *testing.T) {
    testCases := []struct {
        name     string
        url      string
        expected string
    }{
        {
            name:     "Root path",
            url:      "https://api.example.com/",
            expected: "https://api.example.com",
        },
        {
            name:     "Single segment path",
            url:      "https://api.example.com/sso/",
            expected: "https://api.example.com/sso",
        },
        {
            name:     "Multi-segment path",
            url:      "https://api.example.com/auth/sso/",
            expected: "https://api.example.com/auth/sso",
        },
    }
    // ... test implementation
}
```

## Migration Considerations

### Backward Compatibility
- Root-path deployments must continue to work
- No breaking changes to existing configurations
- Provider version should support both old and new URL formats

### Deployment Scenarios
1. **Existing root-path users:** No changes needed, continues working
2. **New subpath users:** Can now use path-based URLs
3. **Migrating deployments:** Can update URL configuration as needed

## Related Issues

### Upstream OpenAPI Client
The root cause may also exist in the OpenAPI-generated client library. Consider:
- Reporting issue to `goauthentik.io/api` repository
- Checking if client properly supports BasePath configuration
- Verifying URL construction logic in generated code

### Similar Terraform Providers
Other providers may have encountered similar issues:
- Check how other providers handle URL paths
- Review Terraform provider best practices
- Consider standardized approach across providers

## References

- **Issue:** XPO-772
- **Provider Repository:** terraform-provider-authentik
- **API Client:** goauthentik.io/api/v3 v3.2025100.4
- **Terraform SDK:** hashicorp/terraform-plugin-sdk/v2 v2.37.0

## Next Steps

1. ✅ Confirm bug existence (COMPLETED)
2. ✅ Document findings (COMPLETED)
3. ⏳ Examine API client structure
4. ⏳ Implement fix
5. ⏳ Add tests
6. ⏳ Submit upstream PR
7. ⏳ Update documentation

---

**Status:** Analysis complete, ready for implementation
**Assignee:** smarkmann
**Priority:** High
