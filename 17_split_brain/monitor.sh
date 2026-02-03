#!/bin/bash
while true; do
    clear
    echo "=== Cluster Status ==="
    date
    echo ""
    for i in {1..5}; do
        echo -n "Node $i: "
        docker logs split-brian-node${i}-1 --tail 1 2>/dev/null | grep "Status:" || echo "no status yet"
    done
    sleep 2
done