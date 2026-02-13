## 📦 Project Structure

.
cmd/ # Entry point for the application
internal/
├─ config/ # Configuration (Viper, Gorm, Logger, etc)
├─ delivery/
│ ├─ http/ # HTTP handlers / controllers
│ │ ├─ auth_controller.go
│ │ ├─ category_controller.go
│ │ ├─ customer_controller.go
│ │ ├─ order_controller.go
│ │ ├─ product_controller.go
│ │ ├─ user_controller.go
│ │ └─ route/ # API route definitions
├─ middleware/ # Custom middlewares (e.g., auth)
├─ entity/ # Core entities / domain models
├─ migration/ # Database migrations
├─ model/ # Database models
│ └─ converter/ # Model converters
├─ repository/ # Database repository layer
├─ service/ # Business logic / service layer
├─ usecase/ # Application use cases
└─ utils/ # Utilities (e.g., validator, cloudinary helpers)

---

## ⚡ Installation

1. Clone the repository

```bash
git clone https://github.com/ojihalaw/ojihalaw-daily-coffee-api.git
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
