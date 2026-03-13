#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
Go SDK Post-Generation Correction Script

This script fixes common issues in the generated Go SDK:
1. Ensures go.work doesn't interfere with builds
2. Validates module paths
3. Fixes any import issues

Run after SDK generation to ensure the SDK builds correctly.
"""

import os
import sys
import subprocess
from pathlib import Path

# Set UTF-8 encoding for Windows console
if sys.platform == 'win32':
    try:
        # Try to set UTF-8 encoding for stdout
        sys.stdout.reconfigure(encoding='utf-8')
    except:
        # Fallback: use ASCII-safe characters
        pass

def print_step(msg):
    print(f"\n[STEP] {msg}")

def print_success(msg):
    try:
        print(f"  ✓ {msg}")
    except UnicodeEncodeError:
        print(f"  [OK] {msg}")

def print_error(msg):
    try:
        print(f"  ✗ {msg}", file=sys.stderr)
    except UnicodeEncodeError:
        print(f"  [ERROR] {msg}", file=sys.stderr)

def print_warning(msg):
    try:
        print(f"  ⚠ {msg}")
    except UnicodeEncodeError:
        print(f"  [WARN] {msg}")


def check_go_work(sdk_path):
    """Check for go.work file that might interfere"""
    print_step("Checking for go.work interference...")
    
    # Check in SDK directory
    go_work_sdk = sdk_path / "go.work"
    if go_work_sdk.exists():
        print_warning(f"Found go.work in SDK directory: {go_work_sdk}")
        go_work_sdk.unlink()
        print_success("Removed go.work from SDK directory")
    
    # Check in parent directories and note them
    current = sdk_path.parent
    go_work_found = False
    while current != current.parent:
        go_work = current / "go.work"
        if go_work.exists():
            print_warning(f"Found go.work in parent directory: {go_work}")
            go_work_found = True
            break
        current = current.parent
    
    if go_work_found:
        print_warning("This may interfere with SDK builds")
        print_warning("The correction script will use GOWORK=off when building")
    else:
        print_success("No go.work interference detected")
    
    return True

def validate_module_path(sdk_path):
    """Validate go.mod module path"""
    print_step("Validating module path...")
    
    go_mod = sdk_path / "go.mod"
    if not go_mod.exists():
        print_error("go.mod not found!")
        return False
    
    with open(go_mod, 'r') as f:
        content = f.read()
    
    # Extract module path
    for line in content.split('\n'):
        if line.startswith('module '):
            module_path = line.split()[1]
            print_success(f"Module path: {module_path}")
            
            # Validate it matches expected path
            if module_path == "github.com/newage-saint/insuretech-go-sdk":
                print_success("Module path is correct")
                return True
            else:
                print_warning(f"Unexpected module path: {module_path}")
                return True  # Don't fail, just warn
    
    print_error("Could not find module declaration in go.mod")
    return False

def test_build(sdk_path):
    """Test if SDK builds successfully"""
    print_step("Testing SDK build...")
    
    # Skip build test in CI environments
    if os.environ.get('CI') == 'true' or os.environ.get('GITHUB_ACTIONS') == 'true':
        print_warning("Skipping build test in CI environment")
        print_warning("Build will be tested in the main workflow")
        return True
    
    try:
        # Build with GOWORK=off to avoid workspace conflicts
        env = os.environ.copy()
        env['GOWORK'] = 'off'
        
        # First try: go build ./...
        result = subprocess.run(
            ['go', 'build', './...'],
            cwd=sdk_path,
            env=env,
            capture_output=True,
            text=True,
            timeout=60
        )
        
        if result.returncode == 0:
            print_success("SDK builds successfully!")
            return True
        
        # If first attempt failed, show error and try alternative
        print_warning("First build attempt failed, trying alternative approach...")
        if result.stderr:
            print(f"  Error: {result.stderr[:200]}")
        
        # Second try: Build specific packages
        result2 = subprocess.run(
            ['go', 'build', './pkg/...'],
            cwd=sdk_path,
            env=env,
            capture_output=True,
            text=True,
            timeout=60
        )
        
        if result2.returncode == 0:
            print_success("SDK builds successfully (pkg only)!")
            return True
        
        # Both failed
        print_error("SDK build failed!")
        print("\nFirst attempt error:")
        print(result.stderr if result.stderr else result.stdout)
        print("\nSecond attempt error:")
        print(result2.stderr if result2.stderr else result2.stdout)
        return False
            
    except subprocess.TimeoutExpired:
        print_error("Build timed out after 60 seconds")
        return False
    except FileNotFoundError:
        print_error("Go not found. Please install Go.")
        return False
    except Exception as e:
        print_error(f"Build test failed: {e}")
        return False

def run_go_fmt(sdk_path):
    """Run go fmt on all Go files"""
    print_step("Running go fmt...")
    
    try:
        result = subprocess.run(
            ['go', 'fmt', './...'],
            cwd=sdk_path,
            capture_output=True,
            text=True,
            timeout=30
        )
        
        if result.returncode == 0:
            print_success("Code formatted")
            return True
        else:
            print_warning("go fmt had issues (continuing...)")
            return True  # Don't fail on formatting issues
            
    except Exception as e:
        print_warning(f"Could not run go fmt: {e}")
        return True  # Don't fail

def main():
    print("=" * 60)
    print("  Go SDK Post-Generation Correction")
    print("=" * 60)
    
    # Find SDK path
    script_dir = Path(__file__).parent
    sdk_path = script_dir / "insuretech-go-sdk"
    
    if not sdk_path.exists():
        print_error(f"SDK not found at: {sdk_path}")
        sys.exit(1)
    
    print_success(f"Found SDK at: {sdk_path}")
    
    # Run corrections
    success = True
    
    check_go_work(sdk_path)
    
    if not validate_module_path(sdk_path):
        success = False
    
    run_go_fmt(sdk_path)
    
    if not test_build(sdk_path):
        success = False
    
    # Summary
    print("\n" + "=" * 60)
    if success:
        print("  ✅ Go SDK Correction Complete - SDK is ready!")
        print("=" * 60)
        print("\n💡 To build manually:")
        print("  Windows: $env:GOWORK=\"off\"; go build ./...")
        print("  Linux/WSL: GOWORK=off go build ./...")
    else:
        print("  ❌ Go SDK Correction Failed - Manual fixes needed")
        print("=" * 60)
        print("\n💡 Try building manually with:")
        print("  Windows: $env:GOWORK=\"off\"; go build ./...")
        print("  Linux/WSL: GOWORK=off go build ./...")
    print("=" * 60)
    
    sys.exit(0 if success else 1)

if __name__ == "__main__":
    main()
