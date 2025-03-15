#!/bin/bash

set -e

echo "building kuma-waybar..."
go build

if [ ! -d ~/.config/waybar ]; then
    echo "waybar not found"
    exit 1
fi

echo "copying kuma-waybar to ~/.local/bin/"
mkdir -p ~/.local/bin
cp ./kuma-waybar ~/.local/bin/kuma-waybar
echo  "Installing at /usr/local/bin/kuma-waybar (CTRL+C to cancel system install)"
sudo cp ./kuma-waybar /usr/local/bin/kuma-waybar

echo "done."
echo ""
echo "Please add the following to your waybar config:"
echo "    \"custom/kuma-waybar\": {"
echo "        \"exec\": \"kuma-waybar --format=waybar --env=\$HOME/.config/waybar/kuma.env\","
echo "        \"interval\": 60",
echo "        \"on-click\": \"kuma-waybar open --env=\$HOME/.config/waybar/kuma.env\"",
echo "        \"format\": \"Kuma {}\","
echo "    },"
echo ""
echo "Please add your UPTIME_KUMA_API_KEY & UPTIME_KUMA_BASE_URL to ~/.config/waybar/scripts/kuma-waybar.env"
echo "Optionally use --env=./second.env to allow for multiple Uptime Kuma"
