# rssi_reader

## Table of Contents

- [Overview](#overview)
- [Architecture Overview](#architecture-overview)
- [Installation](#installation)
  - [Precompiled Binary](#precompiled-binary)
  - [Building from Source](#building-from-source)
- [Usage](#usage)
  - [Running the Binary](#running-the-binary)
  - [Using the `--port` Flag](#using-the-port-flag)
- [Accessing the Client Interface](#accessing-the-client-interface)

## Overview

**rssi_reader** is a lightweight and efficient tool designed to monitor and broadcast Wi-Fi RSSI (Received Signal Strength Indicator) data in real-time. It leverages WebSockets to provide live updates to connected clients, enabling users to visualize signal strength across different wireless interfaces through a user-friendly web interface.

## Architecture Overview

The architecture of **rssi_reader** comprises two main components:

**Server**:

- **Language**: Go (Golang)
- **Functionality**:
  - Reads RSSI data from `/proc/net/wireless` at regular intervals.
  - Broadcasts the collected RSSI data to all connected WebSocket clients.
  - Serves the client web interface (`/monitor`) for visualization.
- **Key Features**:
  - Handles multiple concurrent WebSocket connections efficiently.
  - Provides real-time updates with minimal latency.

## Installation

### Precompiled Binary

For convenience, a precompiled binary for ARMv7 architecture is provided. This is ideal for devices like Raspberry Pi.

1. **Download the Binary**:

   Navigate to the [Releases](https://github.com/abdullahwaqar/rssi_reader/releases) section of this repository and download the latest `rssi_reader_armv7` binary.

   ```bash
   wget https://github.com/abdullahwaqar/rssi_reader/releases/download/v1.0.0/rssi_reader_armv7
   ```

2. **Make the Binary Executable**:

   After downloading, ensure the binary has execute permissions.

   ```bash
   chmod +x rssi_reader_armv7
   ```

### Building from Source

If you prefer to build the binary yourself or want to modify the source code, follow these steps:

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/abdullahwaqar/rssi_reader.git
   cd rssi_reader
   ```

2. **Build the Binary**:

   Ensure you have [Go](https://golang.org/dl/) installed.

   ```bash
   go build -o rssi_reader
   ```

   This command compiles the source code and generates an executable named `rssi_reader` in the current directory.

3. **Building for arm based arch**:
    This is used to build for arm based processors, in this case for Raspberry Pi.

    ```
    make arm
    ```

    This will output the binary in a `bin` folder.

## Usage

### Running the Binary

To start the **rssi_reader** server, execute the binary. You can specify the port using the `--port` flag or use the default port `8080`.

```bash
./rssi_reader --port=9090
```

**If you've moved the binary to `/usr/local/bin` and renamed it to `rssi_reader`:**

```bash
rssi_reader --port=9090
```

**Default Port (`8080`):**

```bash
rssi_reader
```

### Using the `--port` Flag

The `--port` flag allows you to specify the port on which the server listens for incoming connections. This is useful if the default port `8080` is already in use or if you prefer to use a different port for organizational reasons.

**Syntax:**

```bash
./rssi_reader --port=<PORT_NUMBER>
```

**Examples:**

1. **Run on Port `9090`:**

   ```bash
   ./rssi_reader --port=9090
   ```

2. **Run on Port `3000`:**

   ```bash
   ./rssi_reader --port=3000
   ```

**Note:** Ensure that the chosen port is open and not restricted by your firewall or network policies.

## Accessing the Client Interface

Once the server is running, you can access the web-based client interface to monitor RSSI data.

1. **Open a Web Browser:**

   Navigate to the following URL, replacing `<Raspberry_Pi_IP>` with your Raspberry Pi's IP address and `<PORT>` with the port number you've set (default is `8080`).

   ```
   http://<Raspberry_Pi_IP>:<PORT>/monitor
   ```

   **Example:**

   ```
   http://192.168.1.100:8080/monitor
   ```

2. **Interact with the Interface:**

   - **Status Indicator:** Displays the connection status to the WebSocket server.
   - **RSSI Table:** Dynamically updates with real-time RSSI values from your wireless interfaces.

3. **Developer Tools (Optional):**

   For debugging or advanced usage, you can open the browser's developer tools (`F12` or `Ctrl+Shift+I`) to monitor console logs and network activity.

## Automated Installation Using `install.sh`

To simplify the deployment process of **rssi_reader**, an `install.sh` script is provided. This script automates the following tasks:

1. **Moves the binary to `/usr/local/bin`.**
2. **Creates a systemd service file if it doesn't exist.**
3. **Enables and starts/restarts the service.**

### **Prerequisites**

Before running the installation script, ensure you have the following:

- **Compiled Binary:** Ensure you have the `rssi_reader_armv7` binary ready in your current directory. If you haven't compiled it yet, refer to the [Building from Source](#building-from-source) section above.
- **Root or Sudo Access:** The script requires administrative privileges to move binaries to system directories and create systemd service files.

### **Step-by-Step Installation**

Follow these steps to install and set up **rssi_reader** using the `install.sh` script:

#### **1. Make the `install.sh` Script Executable**

Ensure the installation script has the necessary execute permissions:

```bash
chmod +x install.sh
```

#### **2. Run the Installation Script**

Execute the `install.sh` script with root privileges to perform the installation:

```bash
sudo ./install.sh
```

**Note:** Running the script without `sudo` or as a non-root user will result in permission errors.

#### **3. Installation Script Actions**

The `install.sh` script performs the following actions:

1. **Binary Installation:**
   - Moves the `rssi_reader_armv7` binary to `/usr/local/bin/rssi_reader`.
   - Sets the executable permissions.
   - Changes ownership of the binary to the user executing the script.

2. **Service Setup:**
   - Creates a systemd service file at `/etc/systemd/system/rssi_reader.service` if it doesn't already exist.
   - Configures the service to run under the current user.
   - Sets the service to start on system boot.
   - Starts or restarts the service based on its current state.

3. **Systemd Configuration:**
   - Reloads the systemd daemon to recognize the new service.
   - Enables the service to start automatically at boot.
   - Starts or restarts the service to apply any new configurations.

#### **4. Verify the Installation**

After running the installation script, perform the following checks to ensure everything is set up correctly:

##### **a. Check Service Status**

Verify that the **rssi_reader** service is active and running:

```bash
systemctl status rssi_reader.service
```

**Expected Output:**

```
● rssi_reader.service - RSSI Reader WebSocket Server
     Loaded: loaded (/etc/systemd/system/rssi_reader.service; enabled; vendor preset: enabled)
     Active: active (running) since Thu 2024-04-27 12:00:00 UTC; 5s ago
   Main PID: 1234 (rssi_reader)
      Tasks: 2 (limit: 4915)
     Memory: 1.2M
     CGroup: /system.slice/rssi_reader.service
             └─1234 /usr/local/bin/rssi_reader --port=8723
```

##### **b. Access the Client Interface**

Open your web browser and navigate to the client interface to start monitoring RSSI data:

```
http://<Raspberry_Pi_IP>:8723/monitor
```

**Replace `<Raspberry_Pi_IP>`** with the actual IP address of your Raspberry Pi or server. If you specified a different port during installation, replace `8080` with your chosen port number.

**Example:**

```
http://192.168.1.100:8723/monitor
```

##### **c. Verify Real-Time Data**

- **Status Indicator:** Should display "Connected to WebSocket."
- **RSSI Table:** Should dynamically update with the latest RSSI values from your wireless interfaces.

#### **5. Customizing the Installation**

By default, the installation script sets the server to listen on port `8080`. To customize the port, you can modify the `PORT` variable within the `install.sh` script before running it.

**Example: Setting Port to `8723`**

1. **Edit the `install.sh` Script:**

   Open the script in your preferred text editor:

   ```bash
   nano install.sh
   ```

2. **Modify the `PORT` Variable:**

   Locate the following line and change `8723` to your desired port number (e.g., `8000`):

   ```bash
   PORT="8723"  # Default port; can be modified as needed
   ```

3. **Save and Exit:**

   - **In `nano`:** Press `Ctrl + O` to save, then `Ctrl + X` to exit.

4. **Run the Installation Script Again:**

   ```bash
   sudo ./install.sh
   ```

5. **Access the Client Interface on the New Port:**

   ```
   http://<Raspberry_Pi_IP>:9090/monitor
   ```

#### **6. Uninstallation (Optional)**

If you need to remove **rssi_reader** from your system, follow these steps:

1. **Stop the Service:**

   ```bash
   sudo systemctl stop rssi_reader.service
   ```

2. **Disable the Service from Starting on Boot:**

   ```bash
   sudo systemctl disable rssi_reader.service
   ```

3. **Remove the Service File:**

   ```bash
   sudo rm /etc/systemd/system/rssi_reader.service
   ```

4. **Reload systemd Daemon:**

   ```bash
   sudo systemctl daemon-reload
   ```

5. **Remove the Binary:**

   ```bash
   sudo rm /usr/local/bin/rssi_reader
   ```

6. **Verify Removal:**

   ```bash
   systemctl status rssi_reader.service
   ```

   **Expected Output:**

   ```
   Unit rssi_reader.service could not be found.
   ```

#### **7. Troubleshooting**

If you encounter any issues during or after installation, consider the following steps:

##### **a. Check Service Logs**

View real-time logs to identify any runtime errors:

```bash
journalctl -u rssi_reader.service -f
```
