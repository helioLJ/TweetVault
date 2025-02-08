# TweetVault

TweetVault is a web application that helps you organize, tag, and categorize your Twitter bookmarks along with their associated media. It provides an intuitive interface for managing your Twitter content collection.

## Features

- 📤 Upload Twitter bookmarks via JSON file and media ZIP
- 🏷️ Tag and categorize bookmarks
- 🔍 Advanced search and filtering capabilities
- 📱 Media preview and management
- 🔄 Duplicate detection and handling
- 📑 "Read Later" and "To Do" organization

## Tech Stack

- **Frontend:** React + Next.js
- **Backend:** Golang
- **Database:** PostgreSQL
- **Styling:** Tailwind CSS

## Getting Started

### Prerequisites

- Node.js (v16+)
- Go (v1.19+)
- PostgreSQL (v14+)
- Docker (optional)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/tweetvault.git
cd tweetvault
```

2. Set up the backend:
```bash
cd backend
go mod download
```

3. Set up the frontend:
```bash
cd frontend
npm install
```

4. Configure your environment variables:
```bash
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env.local
```

### Development

1. Start the backend server:
```bash
cd backend
go run main.go
```

2. Start the frontend development server:
```bash
cd frontend
npm run dev
```

3. Visit `http://localhost:3000` in your browser

## Project Structure

```
TweetVault/
├── backend/         # Golang API server
├── frontend/        # Next.js web application
├── data/           # Sample data and schemas
└── docs/           # Documentation
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with ❤️ for Twitter power users
- Inspired by the need for better bookmark organization