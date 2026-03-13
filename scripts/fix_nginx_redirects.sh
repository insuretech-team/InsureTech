#!/bin/bash
# Fix ERR_TOO_MANY_REDIRECTS caused by certbot duplicate redirects
# Certbot adds a new server block with `if ($host = ...) return 301`
# on top of our existing return 301 in the HTTP block = double redirect.

REMOTE_HOST="insureadmin@146.190.97.242"
SSH_SOCKET="/tmp/insuretech_fixnginx_$$"
SSH_OPTS="-o ControlMaster=auto -o ControlPath=$SSH_SOCKET -o ControlPersist=60 -o StrictHostKeyChecking=no"
trap "ssh -O exit -o ControlPath=$SSH_SOCKET $REMOTE_HOST 2>/dev/null || true; rm -f $SSH_SOCKET" EXIT

echo ">> Connecting..."
ssh $SSH_OPTS -N -f "$REMOTE_HOST"
read -rsp ">> Sudo password: " SUDO_PASS; echo ""

ssh $SSH_OPTS "$REMOTE_HOST" bash << REMOTE
set -euo pipefail
SUDO_PASS='${SUDO_PASS}'
s() { printf '%s\n' "\$SUDO_PASS" | sudo -S -p '' "\$@"; }

echo "=== Current nginx site configs ==="
for f in /etc/nginx/sites-enabled/*.conf; do
    echo "--- \$f ---"
    cat "\$f"
    echo ""
done

echo "=== Fixing duplicate certbot redirects ==="
# Certbot adds server blocks like:
#   server {
#       if (\$host = domain.com) {
#           return 301 https://\$host\$request_uri;
#       } # managed by Certbot
#       listen 80;
#       server_name domain.com;
#       return 404; # managed by Certbot
#   }
# These need to be removed since our HTTP blocks already have return 301.
# We use python3 to safely parse and remove these certbot-added blocks.
for f in /etc/nginx/sites-enabled/*.conf; do
    echo "Fixing: \$f"
    s python3 - "\$f" << 'PYEOF'
import sys, re
with open(sys.argv[1]) as fh:
    content = fh.read()
# Remove certbot-managed server blocks (the duplicate redirect ones)
# Pattern: server block containing 'managed by Certbot' and 'return 404'
cleaned = re.sub(
    r'server\s*\{[^}]*managed by Certbot[^}]*return 404[^}]*\}\s*',
    '',
    content,
    flags=re.DOTALL
)
if cleaned != content:
    with open(sys.argv[1], 'w') as fh:
        fh.write(cleaned)
    print(f'Fixed: {sys.argv[1]}')
else:
    print(f'No certbot blocks found: {sys.argv[1]}')
PYEOF
done

echo ""
echo "=== Testing nginx config ==="
s nginx -t
s systemctl reload nginx
echo "Done. nginx reloaded."
REMOTE
