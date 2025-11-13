# Upstream Submission Guide

This document provides instructions for submitting the bugfix to the upstream terraform-provider-authentik repository.

## Summary

**Issue:** Terraform provider strips URL path component when Authentik is deployed at a subpath  
**Fix:** Modified provider configuration to use `Servers` field instead of deprecated `Host`/`Scheme` pattern  
**Commit:** 1b03d1c - "Fix: Preserve URL path component for subpath deployments"

## Files Changed

1. **pkg/provider/provider.go**
   - Modified `providerConfigure()` function
   - Now properly constructs full server URL including path component
   - Uses `config.Servers` field for proper URL handling

2. **pkg/provider/provider_pathurl_test.go** (NEW)
   - Comprehensive test suite for path-based URLs
   - Tests root paths, single/multi-segment paths, with/without trailing slashes
   - All tests pass ✅

3. **BUGFIX-ANALYSIS.md** (NEW)
   - Detailed technical analysis of the bug
   - Root cause explanation
   - Implementation details
   - Improvement suggestions

## Before Submitting

### 1. Fork the Repository
```bash
# On GitHub: Fork goauthentik/terraform-provider-authentik
# Then add your fork as a remote
cd /home/smarkmann/Developer/Git/skyla/terraform-provider-authentik
git remote add fork git@github.com:YOUR_USERNAME/terraform-provider-authentik.git
```

### 2. Create Feature Branch
```bash
git checkout -b fix/url-path-component
git push fork fix/url-path-component
```

### 3. Run All Tests
```bash
# Run all provider tests to ensure nothing broke
go test ./pkg/provider/...

# Run specific test
go test -v ./pkg/provider -run TestProviderConfigure_PathBasedURL
```

### 4. Check Code Style
```bash
# Format code
go fmt ./...

# Run linter if available
golangci-lint run ./pkg/provider/
```

## Pull Request Template

### Title
```
Fix: Preserve URL path component for subpath deployments
```

### Description
```markdown
## Problem

When the Terraform provider is configured with a URL containing a path component (e.g., `https://api.example.com/sso/`), the provider strips the path and makes API requests to the wrong endpoint.

**Expected:**
- Provider URL: `https://api.example.com/sso/`
- API Request: `https://api.example.com/sso/api/v3/core/applications/`

**Actual:**
- Provider URL: `https://api.example.com/sso/`
- API Request: `https://api.example.com/api/v3/core/applications/` ❌ (Missing `/sso/`)

This blocks usage with path-based deployments, which are common in:
- Multi-service reverse proxy configurations
- Kubernetes ingress with path prefixes
- Shared domain deployments

## Root Cause

In `pkg/provider/provider.go`, the configuration only set `config.Host` and `config.Scheme`, completely ignoring `akURL.Path`:

```go
config.Host = akURL.Host      // Only hostname
config.Scheme = akURL.Scheme  // Only scheme
// Path component (/sso/) is lost!
```

## Solution

Modified the provider to use the `config.Servers` field, which properly supports full URLs including path components:

```go
// Construct full server URL including path component
serverURL := fmt.Sprintf("%s://%s", akURL.Scheme, akURL.Host)
if akURL.Path != "" && akURL.Path != "/" {
    serverURL += strings.TrimSuffix(akURL.Path, "/")
}

config.Servers = api.ServerConfigurations{
    {
        URL:         serverURL,
        Description: "Authentik API Server",
    },
}
```

## Testing

- Added comprehensive test suite (`pkg/provider/provider_pathurl_test.go`)
- Tests cover:
  - Root paths with/without trailing slashes
  - Single-segment paths (e.g., `/sso`)
  - Multi-segment paths (e.g., `/auth/sso`)
  - HTTP and HTTPS schemes
- All tests pass ✅

## Backward Compatibility

✅ Root-path deployments continue to work unchanged  
✅ No breaking changes to existing configurations  
✅ Maintains all existing functionality

## Example Usage

After this fix, users can configure the provider with path-based URLs:

```hcl
provider "authentik" {
  url   = "https://api.example.com/sso/"  # Now works correctly!
  token = var.authentik_token
}
```
```

### Additional Notes for PR
- Reference any related GitHub issues if they exist
- Mention this fixes a common deployment scenario
- Highlight backward compatibility
- Include test results

## Submission Checklist

- [ ] Fork repository and add as remote
- [ ] Create feature branch from main
- [ ] Run all tests and ensure they pass
- [ ] Run code formatter (`go fmt`)
- [ ] Run linter if available
- [ ] Push to your fork
- [ ] Create Pull Request with template above
- [ ] Link any related GitHub issues
- [ ] Respond to reviewer feedback promptly

## After Submission

1. **Monitor the PR** for comments and reviews
2. **Respond to feedback** within 24-48 hours
3. **Make requested changes** if any
4. **Be patient** - review times vary

## Alternative: Create GitHub Issue First

If you prefer, you can create a GitHub issue first to discuss the approach:

### Issue Template
```markdown
**Title:** Provider strips URL path component for subpath deployments

**Description:**
When configuring the provider with a URL like `https://api.example.com/sso/`, the path component `/sso/` is stripped from API requests, causing 404 errors.

**Expected behavior:**
API requests should preserve the path: `https://api.example.com/sso/api/v3/...`

**Actual behavior:**
API requests strip the path: `https://api.example.com/api/v3/...`

**Impact:**
Blocks usage with:
- Reverse proxy path-based routing
- Kubernetes ingress path prefixes
- Multi-service shared domain deployments

**Proposed solution:**
Use `config.Servers` instead of `config.Host`/`config.Scheme` to preserve full URL including path.

I have a working fix with tests ready to submit as PR if this approach is acceptable.
```

## Contact

For questions about this fix:
- **Developer:** smarkmann
- **Issue Tracker:** XPO-772 (internal)
- **Repository:** /home/smarkmann/Developer/Git/skyla/terraform-provider-authentik

---

**Status:** Ready for upstream submission  
**Date:** 2025-11-13  
**Commit:** 1b03d1c
