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

1. **Server**:
   - **Language**: Go (Golang)
   - **Functionality**:
     - Reads RSSI data from `/proc/net/wireless` at regular intervals.
     - Broadcasts the collected RSSI data to all connected WebSocket clients.
     - Serves the client web interface (`/monitor`) for visualization.
   - **Key Features**:
     - Supports dynamic configuration via command-line flags.
     - Handles multiple concurrent WebSocket connections efficiently.
     - Provides real-time updates with minimal latency.

2. **Client**:
   - **Technology**: HTML, CSS, JavaScript
   - **Functionality**:
     - Connects to the server via WebSockets to receive live RSSI data.
     - Displays the data in a structured table format for easy monitoring.
     - Offers a clean and responsive user interface for optimal user experience.

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

3. **Move the Binary to a Suitable Location** (Optional):

   You can move the binary to `/usr/local/bin` for easier access.

   ```bash
   sudo mv rssi_reader_armv7 /usr/local/bin/rssi_reader
   ```

### Building from Source

If you prefer to build the binary yourself or want to modify the source code, follow these steps:

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/abdullahwaqar/rssi_reader.git
   cd rssi_reader
   ```

2. **Build the Binary**:

   Ensure you have [Go](https://golang.org/dl/) installed (version 1.16 or later recommended).

   ```bash
   go build -o rssi_reader
   ```

   This command compiles the source code and generates an executable named `rssi_reader` in the current directory.

3. **Move the Binary to a Suitable Location** (Optional):

   ```bash
   sudo mv rssi_reader /usr/local/bin/
   ```

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
