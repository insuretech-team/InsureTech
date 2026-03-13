#!/bin/bash
# Check all restarting container logs on remote
REMOTE_HOST="insureadmin@146.190.97.242"
REMOTE_DIR="/home/insureadmin/insuretech"

SSH_SOCKET="/tmp/insuretech_svc_$$"
SSH_OPTS="-o ControlMaster=auto -o ControlPath=$SSH_SOCKET -o ControlPersist=60 -o StrictHostKeyChecking=no"

trap "ssh -O exit -o ControlPath=$SSH_SOCKET $REMOTE_HOST 2>/dev/null || true; rm -f $SSH_SOCKET" EXIT

echo ">> Connecting to $REMOTE_HOST ..."
ssh $SSH_OPTS -N -f "$REMOTE_HOST"
echo ">> Connected."
echo ""

ssh $SSH_OPTS "$REMOTE_HOST" bash << 'REMOTE'
cd /home/insureadmin/insuretech

for svc in authn authz b2b; do
    echo "============================================================"
    echo " LOGS: insuretech-$svc (last 50 lines)"
    echo "============================================================"
    docker compose --profile full logs --no-log-prefix --tail=50 $svc 2>&1
    echo ""
done

echo "============================================================"
echo " NGINX error"
echo "============================================================"
cat /var/log/nginx/error.log 2>/dev/null | tail -20 || echo 'no nginx error log'
REMOTE