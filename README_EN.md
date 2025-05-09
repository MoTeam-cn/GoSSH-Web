# GoSSH-Web

<div align="center">

![GoSSH-Web Logo](assets/logo.svg)

[![Go Version](https://img.shields.io/github/go-mod/go-version/MoTeam-cn/GoSSH-Web)](https://golang.org/doc/devel/release.html)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

*🚀 A Powerful Modern Web Terminal Solution*

English | [简体中文](README.md)

</div>

## 📖 Introduction

GoSSH-Web is a web-based SSH terminal tool developed in Go. It provides a modern web interface that allows users to securely connect to and manage SSH servers through their browsers. The project focuses on providing a smooth terminal experience and robust session management capabilities.

## ✨ Features

- 🌈 Modern web terminal interface based on xterm.js
- ⚡️ Real-time command execution and response with WebSocket support
- 🔐 SSH password authentication and session encryption
- 🔄 Smart session persistence and auto-reconnection
- 📦 Multi-session management and status monitoring
- ⌨️ Complete terminal shortcut support
- 🔍 Powerful terminal output search
- 📱 Responsive design for mobile devices

## 🛠 Tech Stack

- **Backend Framework:** Go + Gin
- **Frontend Technologies:** 
  - HTML5 + JavaScript
  - xterm.js Terminal Emulation
  - WebSocket Real-time Communication
  - Bootstrap 5 Responsive Interface
- **Core Features:**
  - SSH Protocol Support
  - PTY (Pseudo Terminal) Support
  - Session Persistence
  - Real-time Data Transmission

## 🚀 Quick Start

### Requirements

- Go 1.21+
- Modern Browser (Chrome, Firefox, Edge, etc.)

### Installation

1. Clone the repository
   ```bash
   git clone https://github.com/MoTeam-cn/GoSSH-Web.git
   cd GoSSH-Web
   ```

2. Install dependencies
   ```bash
   go mod tidy
   ```

3. Run the server
   ```bash
   go run main.go
   ```

4. Access the application
   Open your browser and visit http://localhost:8080

## 💡 Usage Guide

### Basic Usage

1. **Connect to Server**
   - Enter server information (host, port, username, password)
   - Click "Connect" button

2. **Terminal Operations**
   - Supports all standard terminal operations
   - Common shortcut key support
   - Terminal output search
   - One-click screen clear

3. **Session Management**
   - Auto-save connection configuration
   - Smart reconnection on disconnection
   - Parallel multi-session processing
   - Real-time connection status display

### Advanced Features

- **Terminal Customization**
  - Font and color settings
  - Terminal size adjustment
  - Shortcut key configuration

- **Security Features**
  - HTTPS/WSS support
  - Data encryption in transit
  - Session timeout handling

## 🔒 Security Notes

- HTTPS is strongly recommended in production
- All sensitive information is encrypted using TLS
- Automatic session timeout disconnect
- Access control and authentication mechanism recommended

## 🗺 Roadmap

### Short-term Plans
- [ ] SSH key authentication support
- [ ] Terminal split-screen functionality
- [ ] File transfer system
- [ ] Session recording and playback

### Long-term Goals
- [ ] Terminal theme marketplace
- [ ] Plugin system
- [ ] User authentication system
- [ ] Cluster management support

## 🤝 Contributing

We welcome all forms of contributions, including but not limited to:

- Submitting issues and suggestions
- Improving documentation
- Contributing code improvements
- Sharing usage experiences

Please check [CONTRIBUTING.md](CONTRIBUTING.md) for more details.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

Thanks to all the developers who have contributed to this project!

---

<div align="center">

If this project helps you, please consider giving it a star ⭐️

</div> 