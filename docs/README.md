## The Anima Engineering Roadmap: From Core Logic to Cloud-Native

This roadmap is designed to build a robust, testable, and scalable application by mastering one architectural layer at a time. Each phase produces a working, valuable piece of software and directly maps to a critical engineering skill.

### Phase 0: The Core Engine (CLI MVP)

**Objective:** Build a fully-functional, locally-runnable command-line journal. We will focus exclusively on the core business logic, perfectly tested, with zero network or infrastructure overhead. This forces a clean separation between your application's "what" and its "how."

- [ ] **Project Scaffolding:**
    - [ ] Initialize a Go module for the project.
    - [ ] Set up a standard Go project layout (e.g., `/cmd`, `/internal`).
    - [ ] Initialize Git for version control.
- [ ] **TDD First:**
    - [ ] Set up Go's built-in testing framework.
    - [ ] Write your first failing test for creating a journal entry.
- [ ] **Core Journaling Logic (Go):**
    - [ ] Implement the "Create Journal Entry" function. An entry will have an ID, content, and timestamps.
    - [ ] Implement "Get All Entries" and "Get Entry by ID" functions.
- [ ] **Local Persistence:**
    - [ ] Implement a storage mechanism that reads from and writes the journal to a local JSON file (`journal.json`). We will use an interface so we can easily swap this for a real database later.
- [ ] **CLI Interface:**
    - [ ] Use a simple CLI library (like `cobra`) to create commands:
        - [ ] `anima add "My first entry"`
        - [ ] `anima list`
        - [ ] `anima show <entry_id>`

**🔑 Skills Mastered:** **Idiomatic Go**, **Test-Driven Development (TDD)** as a non-negotiable habit, **SOLID Principles** (specifically the 'D' with our storage interface), and fundamental application design.

---

### Phase 1: The Networked Service (API & Database)

**Objective:** Evolve the core logic into a networked API service with a professional database. The CLI from Phase 0 will become the first "client" of this new API.

- [ ] **Database Integration:**
    - [ ] Write the database schema for `users` and `entries` in PostgreSQL.
    - [ ] Swap the local JSON file storage with a PostgreSQL implementation, satisfying the storage interface from Phase 0.
- [ ] **API Layer (Go):**
    - [ ] Build a REST API server using the standard `net/http` package.
    - [ ] Create API endpoints for all core journaling functions.
- [ ] **Authentication Service (Go):**
    - [ ] Implement "Sign Up with Email" and "Login" endpoints.
    - [ ] Issue JWTs upon successful login.
    - [ ] Implement middleware to protect the journaling endpoints, requiring a valid JWT.
- [ ] **Refactor CLI:**
    - [ ] Modify the CLI commands to make HTTP requests to your new local API instead of accessing the file system directly.

**🔑 Skills Mastered:** **REST API Design**, **PostgreSQL** schema design and interaction, **The Testing Pyramid** (unit tests for logic, integration tests for the database), and security best practices (password hashing, JWTs).

---

### Phase 2: The Production-Ready System (Containerization & Deployment)

**Objective:** Take the working API and prepare it for real-world deployment. This phase is about making the system scalable, manageable, and automated—the core of cloud-native thinking.

- [ ] **Containerization (Docker):**
    - [ ] Write a `Dockerfile` for the Go API service.
    - [ ] Use `docker-compose.yml` to orchestrate the Go API and the PostgreSQL database for a one-command local setup.
- [ ] **Infrastructure as Code (Terraform):**
    - [ ] Write Terraform scripts to provision a managed PostgreSQL database on AWS (RDS).
    - [ ] Write Terraform scripts to provision a managed Kubernetes cluster on AWS (EKS).
- [ ] **Deployment (Kubernetes):**
    - [ ] Write Kubernetes manifest files (`Deployment`, `Service`, `Ingress`) to deploy the containerized Go application to your EKS cluster.
- [ ] **CI/CD Automation (GitHub Actions):**
    - [ ] Create a GitHub Actions pipeline that on every push to `main`:
        1. [ ] Runs all tests.
        2. [ ] Builds the Go binary.
        3. [ ] Builds and pushes the Docker image to a registry (AWS ECR).
        4. [ ] Deploys the new image to the Kubernetes cluster.

**🔑 Skills Mastered:** **Docker**, **Kubernetes**, **Terraform (IaC)**, **AWS** (RDS, EKS, ECR), and **CI/CD** best practices. This phase is your ticket to a DevOps and systems architecture mindset.

---

### Phase 3: The User Experience (Frontend MVP)

**Objective:** With a robust, tested, and deployed backend API in place, we now build the user-facing application. The backend is our source of truth.

- [ ] **Frontend Scaffolding (Next.js):**
    - [ ] Set up the Next.js application (now in a monorepo alongside the Go backend).
- [ ] **Connect to the API:**
    - [ ] Implement API client logic to communicate with your deployed Go backend.
- [ ] **Build Core UI:**
    - [ ] Create the Sign Up and Login pages.
    - [ ] Implement state management for handling authentication and the JWT.
    - [ ] Create a simple dashboard to display journal entries.
    - [ ] Create a form/editor to create and update entries.

**🔑 Skills Mastered:** Frontend development with **Next.js**, state management, and interacting with a secure, production API.

---

### Phase 4 & Beyond: The Scalable Vision

Now, with a full-stack, production-grade application running, we can strategically implement the advanced features from your original roadmap. They are no longer just features; they are extensions to a stable platform.

- **Enhancements:** OAuth, media uploads (with S3), tagging, search, etc.
- **Microservices:** As complexity grows, we can break out functionality (e.g., an `Insights Service`) following the patterns we've already established.
- **AI Integration:** Finally, we can build the AI layer. It can be a new Go service that analyzes entries stored in our robust database, providing the insights that fulfill your ultimate vision for *Anima*.
