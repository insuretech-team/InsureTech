#!/usr/bin/env python3
"""
Set InsureTech-specific secrets on the FLVE HuggingFace Space.

Usage:
    python scripts/deploy_flve_hf_secrets.py

Required environment variables (set in .env or shell):
    HF_TOKEN          — HuggingFace token with write access to the Space
    HF_SPACE_ID       — e.g. "farukhannan/flve" or "newage-saint/flve-inference-server"
    ACCESS_KEY        — DigitalOcean Spaces access key
    SECRET_KEY        — DigitalOcean Spaces secret key
    FLVE_API_TOKEN    — Bearer token the Space requires on every request

Optional:
    DO_BUCKET         — Spaces bucket (default: lcst)
    DO_REGION         — Spaces region (default: sgp1)
    MAIN_CDN          — CDN base URL (default: https://lcst.sgp1.cdn.digitaloceanspaces.com)
    PROFILE_PATH_PREFIX — Image path prefix (default: insuretech/kyc/profile)
"""

import os
import sys
import requests
import json

def _require(key: str) -> str:
    val = os.environ.get(key, "").strip()
    if not val:
        print(f"ERROR: Required env var '{key}' is not set.", file=sys.stderr)
        sys.exit(1)
    return val

def set_hf_space_secret(token: str, space_id: str, key: str, value: str) -> bool:
    """Set a single secret on a HuggingFace Space via the Hub API."""
    url = f"https://huggingface.co/api/spaces/{space_id}/secrets"
    headers = {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json",
    }
    resp = requests.post(url, headers=headers, json={"key": key, "value": value}, timeout=30)
    if resp.status_code in (200, 201, 204):
        print(f"  ✅ {key} set")
        return True
    else:
        print(f"  ❌ {key} failed: {resp.status_code} {resp.text}", file=sys.stderr)
        return False

def main():
    print("=" * 60)
    print("InsureTech → FLVE HuggingFace Space Secret Deployer")
    print("=" * 60)

    hf_token   = _require("HF_TOKEN")
    space_id   = _require("HF_SPACE_ID")
    access_key = _require("ACCESS_KEY")
    secret_key = _require("SECRET_KEY")
    flve_token = _require("FLVE_API_TOKEN")

    # Required + optional secrets with InsureTech-specific defaults
    secrets = {
        "ACCESS_KEY":           access_key,
        "SECRET_KEY":           secret_key,
        "FLVE_API_TOKEN":       flve_token,
        "DO_BUCKET":            os.environ.get("DO_BUCKET", "lcst"),
        "DO_REGION":            os.environ.get("DO_REGION", "sgp1"),
        "MAIN_CDN":             os.environ.get("MAIN_CDN", "https://lcst.sgp1.cdn.digitaloceanspaces.com"),
        "PROFILE_PATH_PREFIX":  os.environ.get("PROFILE_PATH_PREFIX", "insuretech/kyc/profile"),
        "MODELS_PATH_PREFIX":   os.environ.get("MODELS_PATH_PREFIX", "labaidai/models"),
    }

    # Redis session persistence (optional — skip if not configured)
    redis_url = os.environ.get("REDIS_URL", "").strip()
    if redis_url:
        secrets["REDIS_URL"] = redis_url
        print(f"  REDIS_URL configured — sessions will persist across Space restarts")

    print(f"\nSetting {len(secrets)} secrets on Space: {space_id}\n")
    ok = 0
    for key, value in secrets.items():
        display = "[SET]" if value else "[EMPTY]"
        print(f"  Setting {key} = {display}")
        if set_hf_space_secret(hf_token, space_id, key, value):
            ok += 1

    print(f"\n{'✅ All' if ok == len(secrets) else f'⚠️ {ok}/{len(secrets)}'} secrets deployed to {space_id}")
    print("\nNext steps:")
    print("  1. Go to https://huggingface.co/spaces/{} → Settings → Secrets to verify".format(space_id))
    print("  2. Restart the Space (Factory Reboot) to pick up new secrets")
    print("  3. Set FLVE_API_TOKEN= and FLVE_API_TOKEN= in InsureTech .env")
    print("  4. Confirm /health returns { status: healthy } before running KYC flow")

if __name__ == "__main__":
    main()
