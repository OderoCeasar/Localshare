# LocalShare

A cross-platform file-sharing tool for local networks. Transfer files between devices on your LAN through a simple web interface.

## Features

- **Easy Setup**: Single command to start the server
- **Web Interface**: Access from any browser on your network
- **PIN Protection**: Optional PIN to secure file access
- **Admin Authentication**: Control who can upload files
- **Cross-Platform**: Works on Windows, macOS, and Linux
- **No Cloud Required**: Everything stays on your local network

## Tech Stack

- **Backend**: Go + Gin framework
- **Frontend**: React + Tailwind CSS
- **CLI**: Cobra for professional command-line interface
- **Session Management**: Cookie-based sessions

## Installation

### Prerequisites

- Go 1.21 or higher
- Node.js 18+ (for building the frontend)

### Setup

1. **Clone the repository**
```bash
git clone https://github.com/OderoCeasar/localshare.git
cd localshare
```

2. **Install Go dependencies**
```bash
go mod download
```

3. **Build the binary**
```bash
go build -o localshare
```

## Usage

### Basic Usage

Start a server on port 8080:
```bash
./localshare
```

### With PIN Protection

Require a PIN to access files:
```bash
./localshare --pin 1234
```

### With Admin Authentication

Require authentication for uploads:
```bash
./localshare --admin --admin-pass mypassword
```

### Custom Port and Directory

```bash
./localshare --port 3000 --dir /path/to/uploads
```

### All Options

```bash
./localshare --help
```

**Available Flags:**
- `-p, --port` - Port to run server on (default: 8080)
- `-d, --dir` - Directory to store files (default: ./uploads)
- `--pin` - Optional PIN for file access (4-6 digits)
- `--admin` - Enable admin authentication
- `--admin-user` - Admin username (default: admin)
- `--admin-pass` - Admin password
- `--max-size` - Maximum file size in MB (default: 500)

## Project Structure

```
localshare/
├── main.go           # Entry point and CLI setup
├── server.go         # HTTP server and routes
├── go.mod            # Go dependencies
├── frontend/         # React frontend
│   ├── src/
│   │   └── App.jsx   # Main React component
│   └── package.json
└── uploads/          # Default upload directory
```

## How It Works

### Backend (Go + Gin)

The backend is built with the Gin web framework and provides:

1. **RESTful API** for file operations
2. **Session Management** using cookie-based sessions
3. **Security Features**:
   - PIN verification
   - Admin authentication
   - File size limits
   - Path traversal prevention

### Frontend (React)

The frontend is a single-page React application with:

1. **File Upload/Download** interface
2. **Authentication Forms** for PIN and admin login
3. **Real-time Updates** of file list
4. **Responsive Design** with Tailwind CSS

### CLI (Cobra)

The CLI is built with Cobra, providing:

1. **Professional Help Text**
2. **Flag Parsing** for configuration
3. **User-Friendly Messages**

## Security Considerations

- **Local Network Only**: The server binds to all interfaces but is designed for LAN use
- **PIN Protection**: Uses constant-time comparison to prevent timing attacks
- **Admin Auth**: Credentials are hashed and verified securely
- **Path Traversal**: File paths are sanitized to prevent directory traversal
- **File Size Limits**: Configurable maximum file size

## Development

### Running in Development

1. **Start the Go server** (Terminal 1):
```bash
cd /home/ceasar/cza/Projects/localshare
go run ./cmd/localshare
```

You'll see output like:
```
╔════════════════════════════════════════════════════════════╗
║              LocalShare Server Started                      ║
╠════════════════════════════════════════════════════════════╣
║  Local:    http://localhost:8080                          ║
║  Network:  http://192.168.100.13:8080                      ║
╚════════════════════════════════════════════════════════════╝
```

2. **Start the Vite frontend dev server** (Terminal 2):
```bash
cd /home/ceasar/cza/Projects/localshare/web
pnpm run dev
```

Vite will start on `http://localhost:5173` and proxy all `/api` requests to the Go backend.

### Testing Locally

#### Test on Your Laptop Browser

1. Open `http://localhost:5173` or `http://192.168.100.13:5173` in your browser
2. You should see the LocalShare interface
3. Try uploading a file, downloading it, and deleting it

#### Test on Your Phone (Same Network)

1. Make sure your phone is on the same Wi-Fi network as your laptop
2. Get your laptop's IP address from the Go server startup banner (e.g., `192.168.100.13`)
3. On your phone browser, go to: `http://192.168.100.13:5173`
4. Test the same features:
   - **Upload a file** from your phone
   - **Check the laptop** to see if the file appears
   - **Download** the file on your phone
   - **Delete** files

#### Feature Testing Checklist

- [ ] **Upload** - Upload files from both laptop and phone
- [ ] **Download** - Download files from both devices
- [ ] **Delete** - Delete files and verify they disappear
- [ ] **File List** - Refresh and verify files appear correctly
- [ ] **File Size Display** - Check files show correct sizes (B, KB, MB, GB)
- [ ] **Modification Time** - Verify timestamps are accurate

#### Testing with Security Features

**Test with PIN Protection:**
```bash
go run ./cmd/localshare --pin 1234
```
Then try accessing the app - you should see a PIN entry screen before accessing files.

**Test with Admin Authentication:**
```bash
go run ./cmd/localshare --admin --admin-pass secret123
```
Then try uploading a file - you should be prompted to login with admin credentials.

**Test Both PIN and Admin:**
```bash
go run ./cmd/localshare --pin 1234 --admin --admin-pass secret123
```

#### Cross-Device Sync Test

1. Upload a file from your **phone** 
2. Check your **laptop** - it should appear immediately
3. Download/delete the file from your **laptop**
4. Refresh on your **phone** - changes should be reflected

### Building for Production

1. **Build the frontend**:
```bash
cd /home/ceasar/cza/Projects/localshare/web
npm run build
```

2. **Embed the frontend in the Go binary** (optional):
   - Use `go:embed` to include the built frontend
   - Or serve the build directory statically

3. **Build the final binary**:
```bash
go build -o localShare
```

## Example Scenarios

### Scenario 1: Quick File Transfer
```bash
# On your laptop
./localshare

# On your phone, open browser and go to:
# http://192.168.1.100:8080
```

### Scenario 2: Secure Team Share
```bash
# Start with PIN and admin auth
./localShare --pin 5678 --admin --admin-pass teampass123

# Team members can download files with PIN
# Only admins can upload new files
```

### Scenario 3: Large File Transfer
```bash
# Increase max file size to 2GB
./localShare --max-size 2000 --dir ~/large-files
```

## What You'll Learn

Building this project teaches:

1. **Go Web Development**: HTTP servers, routing, middleware
2. **Session Management**: Cookie-based authentication
3. **File Handling**: Upload, download, and storage
4. **CLI Development**: Building professional command-line tools with Cobra
5. **Security**: Authentication, authorization, and input validation
6. **Full-Stack Integration**: Connecting React frontend with Go backend
7. **Network Programming**: Local network discovery and binding

## Troubleshooting

**Server won't start**
- Check if the port is already in use
- Try a different port with `--port 3000`

**Can't access from other devices**
- Ensure devices are on the same network
- Check firewall settings
- Use the network IP shown in the startup message

**Upload fails**
- Check admin authentication if enabled
- Verify file size is within limits
- Ensure upload directory has write permissions

## Future Enhancements

- QR code generation for easy mobile access
- Drag-and-drop file upload
- Multiple file upload
- File search and filtering
- Transfer history
- Progressive web app (PWA) support
- HTTPS support with self-signed certificates
- WebSocket for real-time updates

## License

MIT License - feel free to use this for learning and personal projects!

## Contributing

This is a learning project, but contributions are welcome! Feel free to:
- Report bugs
- Suggest features
- Submit pull requests
- Share your modifications

## Acknowledgments

- Built with [Gin](https://github.com/gin-gonic/gin) web framework
- CLI powered by [Cobra](https://github.com/spf13/cobra)
- UI components from [Lucide React](https://lucide.dev/)
- Styling with [Tailwind CSS](https://tailwindcss.com/)

---

