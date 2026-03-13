#!/bin/bash
# Check service status

echo "=== Service Status ==="
for service in gateway authn authz tenant b2b; do
    echo ""
    echo "--- insuretech-$service ---"
    if systemctl is-active --quiet insuretech-$service 2>/dev/null; then
        echo "Status: RUNNING"
    else
        echo "Status: NOT RUNNING"
        systemctl status insuretech-$service --no-pager -l 2>/dev/null | tail -10
    fi
done

echo ""
echo "=== Port Status ==="
sudo ss -tlnp | grep -E '(8080|3000|5005[0-9]|5006[0-9]|5007[0-9])' || echo "No services listening"

echo ""
echo "=== Recent Error Logs ==="
for service in gateway authn authz tenant b2b; do
    if [ -f "/home/insureadmin/insuretech/logs/$service.error.log" ]; then
        ERROR_COUNT=$(wc -l < /home/insureadmin/insuretech/logs/$service.error.log)
        if [ "$ERROR_COUNT" -gt 0 ]; then
            echo ""
            echo "--- $service errors (last 10 lines) ---"
            tail -10 /home/insureadmin/insuretech/logs/$service.error.log
        fi
    fi
done

echo ""
echo "=== Gateway Health Check ==="
curl -s http://localhost:8080/healthz | jq '.' 2>/dev/null || curl -s http://localhost:8080/healthz
