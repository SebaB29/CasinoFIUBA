---
name: 🐛 Bug report
about: Create a report to help us improve CasinoFIUBA
title: '[BUG] '
labels: 'bug'
assignees: ''

---

## 📝 Description

Briefly describe the problem. Is it a Go compilation error, a database connection failure, or an incorrect API response?

## 👣 How to Reproduce

Steps to reproduce the behavior:
1. Run `sh start.sh`
2. Try to access the endpoint: `GET/POST ...`
3. Perform the action: '...' (e.g., placing a bet with zero balance)
4. See error: '...' (e.g., the container crashes, or the database returns a 500 error)

## 🎯 Expected Behavior

A clear and concise description of what you expected to happen (e.g., the API should return a 400 Bad Request with a clear error message).

## 📸 Screenshots (if applicable)

Add screenshots or logs to help explain your problem (especially useful for Postman/Insomnia results or Docker logs).

## 💻 Environment

- **OS:** (e.g. Windows 11, macOS, Ubuntu)
- **Docker Version:** (e.g. 24.x.x)
- **Go Version:** (if running locally without Docker)
- **Other relevant data:** (e.g., if you modified the `.env` file or `docker-compose.yml`)

## 🔍 Additional Context

Add any other context about the problem here. Please paste the output of `docker-compose logs` or the Go Traceback below:
```text
[Paste your error log or container logs here]
