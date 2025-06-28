# 🎰 CasinoFIUBA

Proyecto universitario que simula un sitio de apuestas tipo casino. Incluye juegos como ruleta, plinko, buscaminas, entre otros. 

Backend en **Go**, base de datos **PostgreSQL**, todo orquestado con **Docker Compose**.

---

## 📑 Índice

- [🚀 Requisitos](#-requisitos)
- [⚙️ Instalación y ejecución](#️-instalación-y-ejecución)
- [📁 Estructura del proyecto](#-estructura-del-proyecto)
- [🧼 Apagar el entorno](#-apagar-el-entorno)
- [🌐 URLs importantes](#-urls-importantes)
- [👥 Participantes](#-participantes)

---

## 🚀 Requisitos

Asegurate de tener instalado:

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

---

## ⚙️ Instalación y ejecución

### 1. Clonar el repositorio

```bash
git clone git@github.com:SebaB29/CasinoFIUBA.git
cd CasinoFIUBA
```

### 2. Crear el archivo `.env`

```bash
cp .env.example .env
```
`⚠️ Este archivo contiene variables de entorno necesarias para la base de datos y el backend.`

### 3. Ejecutar el entorno

```bash
sh start.sh
```

Este script compila y levanta:
* PostgreSQL
* Backend en Go
* Frontend en React

## 📁 Estructura del proyecto
```
.
├── backend/          # Backend en Go
├── test/             # Pruebas automáticas o manuales
├── .env.example      # Variables de entorno de ejemplo
├── .gitignore
├── docker-compose.yml
├── start.sh          # Script para levantar el entorno
└── README.md
```

## 🧼 Apagar el entorno
Para detener los servicios, simplemente presioná `Ctrl+C` en la terminal donde ejecutaste el script.

## 🌐 URLs importantes
| Servicio   | URL                                            |
| ---------- | ---------------------------------------------- |
| Backend    | [http://localhost:8080](http://localhost:8080) |
| PostgreSQL | `localhost:5432` (accesible internamente)      |

## 👥 Participantes
| Nombre             | GitHub                                             |
| ------------------ | -------------------------------------------------- |
| Sebastián Brizuela | [@SebaB29](https://github.com/SebaB29)             |
| Mauri Laganga      | [@Mauri-laganga](https://github.com/Mauri-laganga) |
