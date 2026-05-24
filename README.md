# 🎰 Casino FIUBA

A robust university project simulating a casino betting platform. This system features a variety of games like Roulette, Plinko, and Minesweeper, all powered by a high-performance backend. It showcases modern development practices including containerization and relational database management.

# 📍 Table of Contents
- [📝 Description](#-description)
  - [🧩 Key Features](#-key-features)
  - [🏗️ System Architecture](#️-system-architecture)
  - [🛠️ Technologies](#️-technologies)
- [🚀 Getting Started](#-getting-started)
  - [📋 Prerequisites](#-prerequisites)
  - [⚙️ Setup & Execution](#️-setup--execution)
- [🌐 Important URLs](#-important-urls)
- [🤝 Contributing](#-contributing)
- [👥 Team](#-team)
- [📄 License](#-license)

---

# 📝 Description
**CasinoFIUBA** is a backend-focused application designed to handle betting logic, user balances, and game results securely and efficiently. By using **Go**, the project ensures high concurrency, while **PostgreSQL** provides data integrity for financial transactions (bets and wins).

## 🧩 Key Features
- **Game Variety:** Implementation of classic casino games (Roulette, Plinko, Mines).
- **Relational Storage:** Persistent management of users and game history.
- **Environment Isolation:** Fully dockerized setup for easy deployment and consistent development environments.
- **Backend API:** Built with Go for optimal performance and type safety.

## 🏗️ System Architecture
```text
.
├── backend/          # Go source code (API, Logic, DB connection)
├── test/             # Automated and manual test suites
├── .env.example      # Example environment variables
├── docker-compose.yml # Docker orchestration
└── start.sh          # Automation script for environment setup
```

## 🛠️ Technologies
* **Language**: Go (Golang)
* **Database**: PostgreSQL
* **DevOps**: Docker & Docker Compose
* **Scripting**: Shell Script (sh)

# 🚀 Getting Started
## 📋 Prerequisites
Ensure you have the following installed:
* Docker
* Docker Compose

## ⚙️ Setup & Execution
1. Clone the repository:
   ```bash
   git clone git@github.com:SebaB29/casino-FIUBA.git
   cd casino-FIUBA
   ```

2. Configure environment variables:
   ```bash
   cp .env.example .env
   ```
   ⚠️ Note: Open .env and fill in the required database credentials.

3. Launch the environment:
   ```bash
   sh start.sh
   ```
   This script will build the Go binary and spin up the PostgreSQL container.

4. Shutdown:
   To stop the services, simply press `Ctrl+C` or run `docker-compose down`.

# 🌐 Important URLs
| Service     | URL                              |
| ----------- | -------------------------------- |
| Backend API | http://localhost:8080            |
| PostgreSQL  | localhost:5432 (Internal access) |

# 🤝 Contributing
1. Fork the project.
2. Create your Feature Branch (git checkout -b feature/AmazingFeature).
3. Commit your changes (git commit -m 'Add some AmazingFeature').
4. Push to the Branch (git push origin feature/AmazingFeature).
5. Open a Pull Request.

# 👥 Team
| Nombre             | GitHub                                             |
| ------------------ | -------------------------------------------------- |
| Sebastián Brizuela | [@SebaB29](https://github.com/SebaB29)             |
| Mauri Laganga      | [@Mauri-laganga](https://github.com/Mauri-laganga) |

# 📄 License
This project is licensed under the MIT License - see the LICENSE file for details.
