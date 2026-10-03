# AI-Powered Lead Qualification & CRM Automation

> A production-style business automation system that combines **Go, PostgreSQL, AI, and n8n** to automate lead qualification, sales follow-up, and CRM workflows.

---

## 📸 Workflow Preview

![n8n Workflow](docs/n8n/workflow.png)

The system connects a custom Go backend with AI and n8n to transform incoming leads into qualified sales opportunities and automated follow-up actions.

---

## 🎯 Problem

Businesses receive leads from websites, forms, and other channels, but the process of reviewing, qualifying, and following up with those leads is often manual.

This can result in:

* Slow response times
* Missed leads
* Inconsistent qualification
* Repetitive sales work
* Manual CRM updates
* Time spent writing repetitive emails

This project demonstrates how those processes can be automated while keeping important business decisions under human control.

---

# 💡 Solution

The system provides an automated pipeline:

```text
Lead Submission
      ↓
Go REST API
      ↓
PostgreSQL
      ↓
AI Lead Qualification
      ↓
HOT / WARM / COLD
      ↓
n8n Webhook
      ↓
Automation Workflow
      ↓
AI Email Draft
      ↓
Human Review
      ↓
CRM / Email / Notification
```

The backend controls the core business logic, while n8n handles external workflow orchestration and integrations.

---

# 🏗️ Architecture

```text
                         ┌──────────────────────┐
                         │   Customer / Lead    │
                         │     Website Form     │
                         └──────────┬───────────┘
                                    │
                                    ▼
                         ┌──────────────────────┐
                         │       Go API         │
                         │       :8081          │
                         └──────────┬───────────┘
                                    │
                         ┌──────────┴───────────┐
                         │                      │
                         ▼                      ▼
                ┌─────────────────┐    ┌─────────────────┐
                │   PostgreSQL    │    │  AI Qualifier   │
                │                 │    │     Ollama      │
                │ Source of Truth │    └────────┬────────┘
                └─────────────────┘             │
                                                ▼
                                     ┌────────────────────┐
                                     │  Qualification     │
                                     │ HOT / WARM / COLD  │
                                     └─────────┬──────────┘
                                               │
                                               │ HTTP Webhook
                                               ▼
                                     ┌────────────────────┐
                                     │        n8n         │
                                     │ Workflow Engine    │
                                     └─────────┬──────────┘
                                               │
                              ┌────────────────┼────────────────┐
                              ▼                ▼                ▼
                         AI Email Draft       CRM         Notifications
                              │
                              ▼
                        Human Review
                              │
                              ▼
                           Send
```

---

# 🔄 End-to-End Workflow

## 1. Lead Capture

A customer submits a lead through a website or frontend application.

The Go backend exposes:

```http
POST /api/v1/leads
```

Example:

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "company": "Example Inc",
  "message": "We need AI automation for our sales process",
  "budget": 5000,
  "timeline": "1 month"
}
```

---

## 2. Go API

The Go application:

* Receives the request
* Validates the input
* Applies business logic
* Persists the lead
* Calls the AI qualification service
* Validates the AI result
* Sends the qualified lead to n8n

The main application flow is:

```text
HTTP Handler
     ↓
Lead Service
     ↓
Repository
     ↓
PostgreSQL
     ↓
AI Qualifier
     ↓
n8n Integration
```

---

# 🗄️ PostgreSQL

PostgreSQL acts as the application's **source of truth**.

The lead data includes:

```text
id
name
email
company
message
budget
timeline
created_at
```

Database initialization is located in:

```text
migrations/001_init.sql
```

---

# 🤖 AI Lead Qualification

The system uses **Ollama** to run a local LLM for lead qualification.

The AI analyzes:

* Business/company
* Lead message
* Budget
* Timeline
* Business requirements

It returns structured information:

```json
{
  "score": 85,
  "category": "HOT",
  "summary": "The company has a clear automation need, sufficient budget, and a short implementation timeline.",
  "recommended_action": "Contact the lead within 24 hours."
}
```

### Qualification Categories

```text
HOT
 └── High-priority opportunity

WARM
 └── Potential opportunity requiring follow-up

COLD
 └── Low-priority or long-term opportunity
```

---

# 🛡️ AI Output Validation

The LLM output is treated as **untrusted input**.

Before the result is passed to the automation layer, Go validates:

* JSON format
* Required fields
* Score range: `0–100`
* Category: `HOT`, `WARM`, or `COLD`

The architecture follows:

```text
AI generates
     ↓
Go validates
     ↓
Automation executes
```

This prevents malformed AI responses from directly controlling downstream business workflows.

---

# 🔌 Go → n8n Integration

After qualification, Go sends the lead and AI result to an n8n webhook.

Example payload:

```json
{
  "lead": {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com",
    "company": "Example Inc",
    "message": "We need AI automation for our sales process",
    "budget": 5000,
    "timeline": "1 month"
  },
  "qualification": {
    "score": 85,
    "category": "HOT",
    "summary": "Strong business need with sufficient budget and short timeline.",
    "recommended_action": "Contact the lead within 24 hours."
  }
}
```

The integration is implemented in:

```text
internal/automation/n8n.go
```

---

# 🔗 Webhook Integration

![n8n Webhook Receiver](docs/n8n/webhook-receiver.png)

The webhook provides the connection between the custom Go backend and the n8n automation layer.

```text
Go Application
      │
      │ HTTP POST
      ▼
n8n Webhook
      │
      ▼
Automation Workflow
```

---

# ⚙️ n8n Automation

![n8n Workflow](docs/n8n/workflow.png)

n8n is responsible for coordinating external business workflows.

The workflow can:

* Receive qualified leads
* Route leads based on category
* Generate personalized email drafts
* Send notifications
* Update CRM systems
* Trigger follow-up workflows
* Connect to other business applications

The exported n8n workflow is included in this repository:

```text
n8n/
└── lead-qualification-crm-automation.json
```

The workflow can be imported into another n8n instance.

---

# 📩 AI Sales Email Draft

![AI Email Draft Configuration](docs/n8n/email-draft.png)

For qualified leads, n8n can use an LLM to generate a personalized email draft.

The AI receives:

```text
Lead Information
+
Qualification
+
Recommended Action
```

And generates:

```json
{
  "subject": "Following up on your automation requirements",
  "body": "Hi John,\n\nThanks for reaching out..."
}
```

The system follows a **human-in-the-loop** approach:

```text
AI generates draft
       ↓
Salesperson reviews
       ↓
Salesperson edits if necessary
       ↓
Salesperson sends
```

The AI does not blindly send communication without human review.

---

# 🧠 Why Go + n8n?

The architecture intentionally separates responsibilities.

### Go

Handles:

* Core business logic
* REST APIs
* Validation
* Authentication
* Database operations
* AI orchestration
* AI output validation
* Webhook communication

### PostgreSQL

Handles:

* Persistent business data
* Lead records
* Application state

### Ollama

Handles:

* Local LLM inference
* Lead qualification
* Structured AI generation

### n8n

Handles:

* Workflow orchestration
* Email integrations
* CRM integrations
* Notifications
* External SaaS APIs

This allows the backend to remain independent from individual third-party integrations.

---

# 🛠️ Technology Stack

| Technology        | Purpose                        |
| ----------------- | ------------------------------ |
| **Go**            | Backend API and business logic |
| **PostgreSQL**    | Database                       |
| **Ollama**        | Local AI/LLM inference         |
| **n8n**           | Workflow automation            |
| **Docker**        | Containerization               |
| **REST API**      | Backend communication          |
| **HTTP Webhooks** | Go ↔ n8n integration           |

---

# 📁 Project Structure

```text
.
├── cmd
│   └── api
│       └── main.go
│
├── internal
│   ├── ai
│   │   └── qualifier.go
│   │
│   ├── automation
│   │   └── n8n.go
│   │
│   └── lead
│       ├── handler.go
│       ├── model
│       │   └── model.go
│       ├── repository.go
│       └── service.go
│
├── migrations
│   └── 001_init.sql
│
├── n8n
│   └── lead-qualification-crm-automation.json
│
├── docs
│   └── n8n
│       ├── workflow.png
│       ├── webhook-receiver.png
│       └── email-draft.png
│
├── docker-compose.yml
├── Dockerfile
├── go.mod
├── go.sum
└── README.md
```

---

# 🐳 Deployment

The system is designed to be containerized and can be deployed in environments such as:

```text
Cloud Server
     │
     ├── Go API
     ├── PostgreSQL
     └── n8n
```

It can also be adapted for a company's private infrastructure:

```text
Company Infrastructure
        │
        ├── Go Application
        ├── PostgreSQL
        └── n8n
```

Docker provides a repeatable deployment environment for the application components.

---

# 🚀 Running Locally

## Requirements

* Go
* Docker
* Docker Compose
* Ollama
* An Ollama model
* n8n

---

## 1. Clone the repository

```bash
git clone <your-repository-url>
cd lead-automation
```

---

## 2. Start PostgreSQL

```bash
docker compose up -d postgres
```

Check the container:

```bash
docker compose ps
```

---

## 3. Configure Environment Variables

Create a `.env` file.

Example:

```env
PORT=8081

DATABASE_URL=postgres://postgres:postgres@localhost:5432/leads?sslmode=disable

OLLAMA_URL=http://localhost:11434/api/generate
OLLAMA_MODEL=qwen2.5-coder:1.5b

N8N_WEBHOOK_URL=<your-n8n-webhook-url>
```

Do not commit production credentials or secrets to GitHub.

---

## 4. Run Database Migration

Apply:

```text
migrations/001_init.sql
```

to the PostgreSQL database.

---

## 5. Start the Go API

```bash
go run ./cmd/api
```

The API runs on:

```text
http://localhost:8081
```

---

## 6. Import the n8n Workflow

Open your n8n instance and import:

```text
n8n/lead-qualification-crm-automation.json
```

Configure the required credentials and webhook URL.

---

# 🧪 API Example

Create a lead:

```bash
curl -X POST http://localhost:8081/api/v1/leads \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "company": "Example Inc",
    "message": "We need AI automation for our sales process",
    "budget": 5000,
    "timeline": "1 month"
  }'
```

Expected processing:

```text
Request
  ↓
Go Handler
  ↓
Lead Service
  ↓
PostgreSQL
  ↓
Ollama
  ↓
Qualification
  ↓
n8n Webhook
  ↓
Automation Workflow
```

---

# 🔐 Security Considerations

This project is designed with several security and reliability principles in mind.

### Secrets

Sensitive credentials should be provided through environment variables rather than committed to source control.

### AI Validation

AI responses are validated before being used by downstream automation.

### Data Ownership

PostgreSQL remains the source of truth for application data.

### Human Approval

AI-generated customer communication can require human review before sending.

### Service Separation

Core business logic is kept separate from external workflow integrations.

---

# 📈 Business Use Cases

The architecture can be adapted to many business workflows, including:

* Lead qualification
* Sales follow-up
* CRM synchronization
* Customer support automation
* Email classification
* Lead routing
* Meeting scheduling
* Notifications
* Task creation
* Automated reports
* Document processing
* Internal workflow automation

The underlying pattern is:

```text
Business Event
      ↓
Backend Logic
      ↓
AI / Rules
      ↓
Automation
      ↓
External Business Systems
```

---

# 🔮 Future Improvements

The current architecture can be extended with:

* [ ] Sales dashboard
* [ ] Authentication and authorization
* [ ] Role-based access control
* [ ] Email draft management
* [ ] Human approval workflow
* [ ] Email sending from dashboard
* [ ] CRM synchronization
* [ ] Conversation history
* [ ] Automated follow-ups
* [ ] Calendar integration
* [ ] Meeting scheduling
* [ ] Lead activity timeline
* [ ] Retry and failure handling
* [ ] Structured application logging
* [ ] Monitoring and observability
* [ ] Automated deployment
* [ ] Multi-tenant architecture

---

# 🎯 Engineering Highlights

This project demonstrates the ability to build an end-to-end business automation system rather than an isolated AI demo.

### Backend Engineering

* Go REST API
* Layered service architecture
* PostgreSQL persistence
* Repository pattern
* Request validation
* HTTP integrations

### AI Engineering

* Local LLM integration
* Structured AI output
* Prompt design
* AI response validation
* Human-in-the-loop workflows

### Automation Engineering

* n8n workflow design
* Webhook-driven architecture
* External API integration
* Conditional workflow routing
* Automated sales workflows

### Infrastructure

* Docker
* Docker Compose
* Environment-based configuration
* Deployable service architecture

---

# 📌 Key Design Principle

The system does not treat AI as the entire application.

Instead:

```text
                 AI
                 │
                 ▼
         ┌───────────────┐
         │ Business Logic│
         │      Go       │
         └───────┬───────┘
                 │
                 ▼
          Automation
              n8n
                 │
        ┌────────┼────────┐
        ▼        ▼        ▼
       CRM     Email   Calendar
```

**AI provides intelligence.**

**Go provides control and business logic.**

**PostgreSQL provides persistence.**

**n8n provides workflow orchestration.**

**External systems provide business capabilities.**

---

# 👨‍💻 Portfolio Objective

This project was built to explore how modern backend engineering, AI, and workflow automation can be combined to solve practical business problems.

The focus is on building systems that are:

* Reliable
* Integratable
* Deployable
* Maintainable
* Business-oriented

Rather than replacing existing business software, the architecture is designed to **connect existing systems and automate the processes between them**.

---

## 📄 License

Add your preferred license here.
