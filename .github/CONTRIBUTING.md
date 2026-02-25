# Contributing to CasinoFIUBA 🎰

Thank you for your interest in contributing! This project is a simulated casino environment designed to test backend logic, concurrency in Go, and database integrity. Whether you are fixing a bug or adding a new game, your help is appreciated.

## How to Contribute

1. **Fork** the repository.
2. Create your branch: `git checkout -b feature/new-game-mechanic`.
3. Make your changes and commit them: `git commit -m "Add new game: Blackjack"`.
4. Push your branch: `git push origin feature/new-game-mechanic`.
5. Open a **Pull Request**.

## 🛠️ Development Workflow

- **Dockerized Environment:** Always use `sh start.sh` to test your changes. This ensures your local Go environment matches the containerized environment.
- **Database Schema:** If you modify the database structure, ensure you update the initialization scripts to maintain consistency.
- **Environment Variables:** If your feature requires a new secret or config, add it to `.env.example` with a dummy value.

## 💡 Ideas for Contribution

- **🎨 Create a Frontend:** Since the project is currently a **Backend API only**, building a web interface (using React, Vue, or Next.js) would be a massive contribution.
- **Enhanced Security:** Improve user balance updates using SQL transactions to prevent race conditions.
- **Admin Dashboard:** Create endpoints to monitor total bets, house edge, and active users.
- **Unit Testing:** Increase coverage in the `test/` directory, especially for betting odds math.
- **Performance:** Optimize Docker images for faster builds and smaller footprints.

## 📝 Coding Standards (Go)

- Follow standard [Go Formatting](https://go.dev/doc/effective_go#formatting) (`go fmt`).
- Ensure all database connections are properly managed or pooled.
- Write descriptive error messages for API responses.

Thank you for helping us build a better **CasinoFIUBA**! 🚀
