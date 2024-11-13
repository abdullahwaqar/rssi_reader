#!/bin/bash

# install.sh - Installation script for rssi_reader

set -e

# Variables
BINARY_NAME="rssi_reader_armv7"
INSTALL_PATH="/usr/local/bin/rssi_reader"
SERVICE_NAME="rssi_reader.service"
SERVICE_PATH="/etc/systemd/system/$SERVICE_NAME"
PORT="8723"

echo_info() {
    echo -e "\e[32m[INFO]\e[0m $1"
}

echo_warn() {
    echo -e "\e[33m[WARN]\e[0m $1"
}

echo_error() {
    echo -e "\e[31m[ERROR]\e[0m $1"
}

# Check if the script is run as root
if [[ "$EUID" -ne 0 ]]; then
    echo_error "Please run as root or use sudo."
    exit 1
fi

# Determine the service user
if [[ -n "$SUDO_USER" && "$SUDO_USER" != "root" ]]; then
    SERVICE_USER="$SUDO_USER"
else
    echo_warn "SUDO_USER not set or running as root directly. The service will run as root."
    SERVICE_USER="root"
fi

echo_info "Service will run as user: $SERVICE_USER"

# Check if the binary exists in the current directory
if [[ ! -f "$BINARY_NAME" ]]; then
    echo_error "Binary '$BINARY_NAME' not found in the current directory."
    exit 1
fi

# Move the binary to /usr/local/bin
echo_info "Moving binary to $INSTALL_PATH."
mv "$BINARY_NAME" "$INSTALL_PATH"

# Ensure the binary has execute permissions
chmod +x "$INSTALL_PATH"

# Change ownership to the service user
chown "$SERVICE_USER":"$SERVICE_USER" "$INSTALL_PATH"

# Create the systemd service file if it doesn't exist
if [[ -f "$SERVICE_PATH" ]]; then
    echo_info "Service file '$SERVICE_NAME' already exists."
else
    echo_info "Creating systemd service file at '$SERVICE_PATH'."
    cat <<EOF >"$SERVICE_PATH"
[Unit]
Description=RSSI Reader WebSocket Server
After=network.target

[Service]
Type=simple
User=$SERVICE_USER
Group=$SERVICE_USER
ExecStart=$INSTALL_PATH --port=$PORT
Restart=on-failure
RestartSec=5s
Environment=PATH=/usr/local/bin:/usr/bin:/bin
Environment=PORT=$PORT

[Install]
WantedBy=multi-user.target
EOF
    echo_info "Service file created successfully."
fi

# Reload systemd daemon to recognize the new service
echo_info "Reloading systemd daemon."
systemctl daemon-reload

# Enable the service to start on boot
echo_info "Enabling the '$SERVICE_NAME' service to start on boot."
systemctl enable "$SERVICE_NAME"

# Start or restart the service
if systemctl is-active --quiet "$SERVICE_NAME"; then
    echo_info "Restarting the '$SERVICE_NAME' service."
    systemctl restart "$SERVICE_NAME"
else
    echo_info "Starting the '$SERVICE_NAME' service."
    systemctl start "$SERVICE_NAME"
fi

# Check the status of the service
echo_info "Checking the status of the '$SERVICE_NAME' service."
systemctl status "$SERVICE_NAME" --no-pager

# Display access instructions
echo_info "Installation completed successfully!"
echo_info "Access the client interface at: http://<Raspberry_Pi_IP>:$PORT/monitor"

exit 0
