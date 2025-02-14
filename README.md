# TweetVault

TweetVault is a full-stack web application designed to help you manage, tag, search, and archive your Twitter bookmarks. With a sleek and responsive interface built using React and Next.js and a powerful backend powered by Go (using Gin and Gorm), TweetVault makes it easy to organize your saved tweets and keep track of your reading or to-do items.

## Demo

### Part 1: Basic Features
![TweetVault Demo Part 1](./docs/assets/part1.gif)

### Part 2: Advanced Features
![TweetVault Demo Part 2](./docs/assets/part2.gif)

---

## Features

- **Bookmark Management:**  
  View detailed bookmarks including full text (with clickable links), images, and media.
  
- **Tagging System:**  
  - Create, rename, and delete custom tags.
  - Includes fixed *standard tags* ("To do", "To read") that cannot be renamed or deleted.
  - Toggle tag completion status on individual bookmarks.

- **Search and Filter:**  
  Quickly search through your bookmarks and filter them by tag.
  
- **Pagination:**  
  Adjustable page size and navigation through bookmark pages.
  
- **Batch Selection:**  
  Easily select multiple bookmarks and perform batch operations (e.g., delete, archive).

- **Statistics Dashboard:**  
  Get a summary of your bookmarks and tag usage including total counts, active and archived bookmarks, and popular tags with completion stats.

- **Upload Handling:**  
  Import Twitter bookmarks from a JSON file (within a ZIP archive) via a simple upload interface.

---

## Architecture Overview

### Frontend
- **Framework:** Next.js with React and TypeScript  
- **Components:**  
  - **BookmarkCard:** Displays individual bookmark details, with options for tagging, archiving, and deletion.
  - **TagMenu & SearchAndFilter:** Provides an interactive UI for managing and filtering tags.
  - **Pagination & Button:** Custom UI components for navigation and interaction.
  - **Statistics:** A dashboard component to display bookmark and tag statistics.
  
The frontend is organized within the `frontend/src/components` directory where you'll find subdirectories for bookmarks, UI elements, theme management, and more.

### Backend
- **Language & Framework:** Written in Go using the Gin framework and Gorm ORM.
- **Key Modules:**  
  - **Handlers:** Route handlers for bookmarks, tags, uploads, and statistics (located in `backend/internal/api/handlers`).
  - **Services:** Business logic encapsulated in services such as `BookmarkService` (located in `backend/internal/services`).
  - **Database:** Postgres is used as the primary data store. Auto-migrations and initial seed data (e.g., standard tags) are managed on startup.
  
The backend's entry point is located at `backend/cmd/server/main.go`, and routes are defined in `backend/internal/api/routes/routes.go`.

---

## Installation

### Prerequisites

- **Node.js** (v14 or higher) & **npm** (or yarn) – for frontend dependencies.
- **Go** (v1.18 or higher) – for backend.
- **PostgreSQL** – Database server.

### Setup Instructions

1. **Clone the repository:**

   ```bash
   git clone https://github.com/yourusername/tweetvault.git
   cd tweetvault
   ```

2. **Frontend Setup:**

   ```bash
   cd frontend
   npm install
   npm run dev
   ```

   This will start the frontend development server (typically on [http://localhost:3000](http://localhost:3000)).

3. **Backend Setup:**

   - Set up your database and ensure PostgreSQL is running.
   - Configure your environment variables. You might create a `.env` file in the backend root with settings such as:
     
     ```env
     DBHost=localhost
     DBUser=your_db_user
     DBPassword=your_db_password
     DBName=tweetvault
     DBPort=5432
     ServerPort=8080
     ```
     
   - Install Go dependencies and run the server:

     ```bash
     cd backend
     go mod tidy
     go run cmd/server/main.go
     ```

     The backend server will start (typically on [http://localhost:8080](http://localhost:8080)) and automatically run migrations as well as insert standard tags if missing.

---

## API Endpoints

The backend API is organized under an `/api` route and provides endpoints for core functionalities:

- **Bookmarks:**
  - `GET /api/bookmarks` – List bookmarks with optional filtering by tag or search query.
  - `GET /api/bookmarks/:id` – Retrieve details of a single bookmark.
  - `PUT /api/bookmarks/:id` – Update bookmark details and tags.
  - `DELETE /api/bookmarks/:id` – Delete a bookmark.
  - `POST /api/bookmarks/:id/toggle-archive` – Archive/unarchive a bookmark.

- **Tags:**
  - `GET /api/tags` – Retrieve all tags.
  - `POST /api/tags` – Create a new tag.
  - `PUT /api/tags/:id` – Update a tag.
  - `DELETE /api/tags/:id` – Delete a tag (except standard tags).
  - `GET /api/tags/:id/count` – Get the count of bookmarks using a specific tag.

- **Statistics:**
  - `GET /api/statistics` – Retrieve summary statistics for bookmarks and tags.

- **Uploads:**
  - `POST /api/upload` – Process and import Twitter bookmark data from a ZIP file.
