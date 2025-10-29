#!/bin/bash
sed -i "s|\$SERVICE|$SERVICE|g" /reflex.conf
reflex -c /reflex.conf
