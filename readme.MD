## ⚡ Installation

1. Clone the repository

```bash
git clone https://github.com/ojihalaw/merapi_test.git
cd project-name
```

2. Install dependencies

```bash
go mod tidy
```

3. Setup environment variables

```bash
# Copy the example environment file
cp .env_example .env

# Edit the .env file
```

4. Run the project

```bash
go run cmd/main.go
```

## 📄 Logging

- All logs use **Logrus**
- Supports log levels: `debug`, `info`, `warn`, `error`, `fatal`

---

## 🔧 Configuration

- **Viper** handles configuration from:
  - `.env` file
  - Environment variables
  - Default values

---
