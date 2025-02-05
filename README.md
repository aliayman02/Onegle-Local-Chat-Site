# Onegle Chat

**Onegle Chat** is a simple real-time chat application built with **Go** and **WebSockets**. It allows users to create or join chat rooms and exchange messages or images. This app is designed for **local network use**, making it perfect for small group chats in the same Wi-Fi network.

---

## Features

- **Real-time messaging** using WebSockets.
- **Create or join chat rooms** dynamically.
- **Send images** along with text messages.
- **Responsive design** for both desktop and mobile devices.
- **Automatic room updates** for all connected users.

---

## Project Structure

```
/onegle-chat
├── main.go             # Go server handling WebSocket connections
├── go.mod              # Go module file for dependencies
├── go.sum              # Checksum file for module verification
├── public/
│   ├── index.html      # Frontend of the chat app (HTML structure)
│   └── styles.css      # External CSS for responsive styling
└── .gitignore          # Git ignore file to exclude unnecessary files
```

---

## Prerequisites

- **Go** installed on your system. [Download Go](https://golang.org/dl/)
- A **modern web browser** (Chrome, Firefox, etc.)

---

## How to Run the Server

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/yourusername/onegle-chat.git
   cd onegle-chat
   ```

2. **(Optional) Download Dependencies:**  
   If you encounter missing package errors, run:

   ```bash
   go mod tidy
   ```

3. **Run the Go Server:**

   ```bash
   go run main.go
   ```

4. **Access the Chat App:**

   Open your browser and navigate to:

   ```
   http://localhost:8080
   ```

---

## Connecting from Other Devices (Local Network)

If others on the **same Wi-Fi network** want to join:

1. **Find your local IP address:**

   - On **Windows**:
     ```bash
     ipconfig
     ```
     Look for the **IPv4 Address** under your active network connection.

   - On **macOS/Linux**:
     ```bash
     ifconfig
     ```
     Look for the **inet** address under your active network interface.

2. **Share your IP** with others. They can access the app using:

   ```
   http://YOUR_LOCAL_IP:8080
   ```

   Example:
   ```
   http://192.168.1.5:8080
   ```

---

## How to Stop the Server

- Press **`Ctrl + C`** in the terminal to stop the server gracefully.

---

## Troubleshooting

- **Port already in use?**  
  If you see an error like `address already in use`, change the port in `main.go`:

  ```go
  err := http.ListenAndServe(":8081", nil)
  ```

  Then access the app at `http://localhost:8081`.

- **Can't connect from another device?**  
  - Ensure both devices are on the **same Wi-Fi network**.
  - Disable firewalls or allow Go through your firewall settings.

---

## License

This project is licensed under the **MIT License**.

---

## Acknowledgments

- [Gorilla WebSocket](https://github.com/gorilla/websocket) for WebSocket support in Go.
- Inspired by simple, real-time chat applications for local network use.
