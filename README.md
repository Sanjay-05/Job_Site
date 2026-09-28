# 🧑‍💼 JobSite API

A RESTful backend service for a job portal built using Go (Golang) with Go-Kit and Gorilla Mux. This API supports user authentication, job listings, and admin job management.

<p align="center">
  <img src="https://github.com/user-attachments/assets/f9f2b3bc-8b81-43b5-ab0c-2c23bf865535" width="100" />
  <img width="100" height="70" alt="image" src="https://github.com/user-attachments/assets/be507447-e8f6-4ac2-bad6-3c1cc06f4dd5" />

  <img width="100" height="70" alt="image" src="https://github.com/user-attachments/assets/ad488acd-4fdd-418a-9ac2-1b76a4c89bc1" />
  <img width="100" height="70" alt="image" src="https://github.com/user-attachments/assets/00606063-9d2f-479e-bcc5-0a2ba8b61267" />
</p>




## 🎯 Features
- User Registration & Login
- JWT-based Authentication
- Forgot & Reset Password
- Job Listings (Public)
- Lookup APIs (Countries, Job Types)
- Admin Job Management (Create, Update, Delete)
- Middleware-based Role Authorization

## 🚀 Installation
1. Clone the repo: `git clone https://github.com/NishanthSaravanan/Senthil-Nov2025-RTE.git`
2. Install dependencies: `go mod tidy`
3. Run the app: `go run main.go`

## 📁 Project Structure

```
RTE_jobsite/
│
├── handler/        # HTTP request/response encoding & decoding
├── service/        # Business logic
├── middleware/     # JWT & role-based middleware
├── router/         # API route definitions
├── models/         # Data models (if applicable)
└── main.go         # Application entry point
```

### ➡️ Control Flow

```mermaid
flowchart LR
   Request e1@==> Server
   e1@{ animate: true }
   Server e2@==> Router
   e2@{ animate: true }
   M@{ shape: processes, label: "Middlewares" }
   Router e3@==> M
   e3@{ animate: true }
   C@{ shape: processes, label: "Handlers" }
   M e4@==> C
   e4@{ animate: true }
   R@{ shape: processes, label: "Repository" }
   C e5@==> R
   e5@{ animate: true }
   id1[(Database)]
   R e6@==> id1
   e6@{ animate: true }
```

1. **Request**: Server receives the incoming request.
2. **Server**: Server processes the request and forwards it to the router.
3. **Router**: Router directs the request to the appropriate middleware(s).
4. **Middlewares**: The middlewares handle various tasks such as logging, authentication, security headers, tracing etc.,
5. **Handlers**: The request is passed to the appropriate handler, which validates the request and forwards it to the repository layer.
6. **Repository**: The repository layer communicates with the database to perform CRUD operations.


## 🔐 Authentication
- Uses JWT for secure authentication
- Role-based access (admin protected routes)

## 📡 API Endpoints
### 🔓 Public Routes
| Method | Endpoint               | Description            |
| ------ | ---------------------- | ---------------------- |
| POST   | `/api/register`        | Register new user      |
| POST   | `/api/login`           | User login             |
| POST   | `/api/forgot-password` | Request password reset |
| POST   | `/api/reset-password`  | Reset password         |
| GET    | `/api/jobs`            | List all jobs          |
| GET    | `/api/countries`       | List countries         |
| GET    | `/api/job-types`       | List job types         |

### 🔒 Admin Routes (JWT Required)
```Base Path: /api/admin ```
| Method | Endpoint     | Description |
| ------ | ------------ | ----------- |
| POST   | `/jobs`      | Create job  |
| PUT    | `/jobs`      | Update job  |
| DELETE | `/jobs/{id}` | Delete job  |


## 🛠️ Tech Stack
- Golang
- Go-Kit (Transport Layer)
- Gorilla Mux (Router)
- JWT Authentication
- Clean Architecture (Handler → Service → Repository)


## Contributing
Feel free to fork the repo and submit pull requests.

## Author
Senthilkanagaraj

GitHub: https://github.com/Senthilkanagaraj