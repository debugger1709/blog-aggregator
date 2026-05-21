# Blog Aggregator (Gator)

A command-line RSS feed aggregator written in Go. This tool allows users to manage accounts, subscribe to RSS feeds, automatically fetch blog posts via a background worker, and browse content directly from the terminal.

Built as part of the Boot.dev backend career path to practice SQL, database migrations, and advanced CLI command patterns in Go.

## Features

- **User Management**: Register, login, and list users.
- **Feed Management**: Add RSS feeds and track which user added them.
- **Follow System**: Follow, unfollow, and view followed feeds.
- **Aggregation Worker**: Background daemon to periodically scrape and save blog posts.
- **Post Browser**: View aggregated blog posts directly in the CLI.

---

## Prerequisites

Ensure you have the following installed:
- **Go** (version 1.20+)
- **PostgreSQL**
- **Goose** (for database migrations)
- **SQLC** (for generating type-safe SQL code)

---

## Setup & Installation

### 1. Configuration File
The application reads configuration from your home directory. Create a file named `.gatorconfig.json` in your home folder:

```json
{
  "db_url": "postgres://username:password@localhost:5432/your_database_name?sslmode=disable",
  "current_user_name": ""
}
```

### 2. Database Migrations
Navigate to your database schema directory and run your SQL migrations using Goose:
```bash
cd sql/schema
goose postgres "postgres://username:password@localhost:5432/your_database_name?sslmode=disable" up
```

### 3. Generate SQL Code
Generate the type-safe Go queries using SQLC from the project root:
```bash
sqlc generate
```

### 4. Build the Project
Compile the binary:
```bash
go build -o gator
```

---

## Command Reference

Run commands using the compiled binary: `./gator <command> [args...]`.

### Authentication & Users
- **`./gator register <username>`**  
  Creates a new user and logs them in.
- **`./gator login <username>`**  
  Switches the active user in the configuration file.
- **`./gator users`**  
  Lists all registered users in the system.
- **`./gator reset`**  
  Clears database contents (useful for development).

### Feed Management
- **`./gator addfeed <name> <url>`**  
  Adds a new RSS feed and automatically follows it.
- **`./gator feeds`**  
  Lists all feeds stored in the system.

### Follow System
- **`./gator follow <url>`**  
  Follows an existing RSS feed.
- **`./gator following`**  
  Lists all feeds followed by the current user.
- **`./gator unfollow <url>`**  
  Unfollows a specific RSS feed.

### Scraping & Browsing
- **`./gator agg <time_duration>`**  
  Starts the scraping daemon. Example: `./gator agg 1m` or `./gator agg 30s`.
- **`./gator browse [limit]`**  
  Displays stored posts. Optional `limit` parameter sets how many posts to show.
