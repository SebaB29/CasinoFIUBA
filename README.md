# 🎰 CasinoFIUBA

Este es un proyecto universitario que simula un sitio de apuestas tipo casino. Incluye:

- Backend en **Go**
- Frontend en **React**
- Base de datos **PostgreSQL**
- Orquestado con **Docker Compose**

## 🚀 Requisitos

Antes de empezar, asegurate de tener instalado:

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)


## ⚙️ Instalación y ejecución

### 1. Clonar el repositorio

```bash
git clone git@github.com:SebaB29/CasinoFIUBA.git
cd CasinoFIUBA
```

### 2. Crear el archivo .env

```bash
cp .env.example .env
```

⚠️ Este archivo contiene variables de entorno necesarias para la base de datos y el backend.

### 3. Ejecutar el entorno

```bash
sh start.sh
```

Este script compila y levanta:
* PostgreSQL
* Backend en Go
* Frontend en React

## 📁 Estructura del proyecto

```bash
.
├── backend/          # Backend en Go
├── frontend/         # Frontend en React
├── .env.example      # Variables de entorno de ejemplo
├── docker-compose.yml
├── start.sh          # Script para levantar el entorno
└── README.md
```

## 🧼 Apagar el entorno

```bash
docker-compose down
```

## 🌐 URLs importantes

| Servicio   | URL                                   |
|------------|----------------------------------------|
| Frontend   | [http://localhost:3000](http://localhost:3000) |
| Backend    | [http://localhost:8080](http://localhost:8080) |
| PostgreSQL | `localhost:5432` (internamente)        |

