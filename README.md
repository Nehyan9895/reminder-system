# 📅 Dynamic Reminder System

[![Go](https://img.shields.io/badge/Go-1.21-blue?logo=go)](https://golang.org/) 
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

A **Dynamic Reminder System** built in **Go (Golang)** with configurable rules, background scheduling, and an audit trail.  
This system allows you to create tasks, set reminders (before due or at intervals), and logs all activity for auditing purposes.

---

## 🚀 Features

- **Task Management**
  - Create, update, delete, and list tasks
  - Track due dates and status (pending/done)
- **Reminder Rules**
  - **Before Due:** Remind X minutes before task is due
  - **Interval:** Repeat reminders every Y minutes until task is done
  - Activate/deactivate rules dynamically
- **Scheduler**
  - Runs periodically (every minute)
  - Applies active reminder rules to tasks
  - Simulates sending reminders via console logs
- **Audit Trail**
  - Logs all rule changes (create/update/delete/activate/deactivate)
  - Logs every triggered reminder with rule and task details
- **Sample Data**
  - Pre-seeded tasks for demonstration

---

## 🛠 Technologies

- **Language:** Go (Golang)
- **Database:** PostgreSQL
- **Web Framework:** [Chi](https://github.com/go-chi/chi)
- **ORM:** [GORM](https://gorm.io/)
- **Logging:** Logrus

---

## ⚙️ Setup

1. **Clone the repository:**

```bash
git clone https://github.com/Nehyan9895/reminder-system.git

cd reminder-system
```

2. **Install dependencies:**

```bash
go mod tidy
```

3. **Configure the database using environment variables or config file.**


4. **Start the server:**

```bash
 go run ./cmd
```

- The API runs at: http://localhost:8080 at default
- The simple UI is available at: http://localhost:8080/
- Hosted system available at https://reminder-system-4v3s.onrender.com/
---

## 📦 API Endpoints

**Tasks**

| Method | Endpoint      | Description       |
| ------ | ------------- | ----------------- |
| GET    | `/tasks`      | List all tasks    |
| POST   | `/tasks`      | Create a new task |
| GET    | `/tasks/{id}` | Get task by ID    |
| PUT    | `/tasks/{id}` | Update task by ID |
| DELETE | `/tasks/{id}` | Delete task       |


**Reminder Rules**

| Method | Endpoint                 | Description       |
| ------ | ------------------------ | ----------------- |
| GET    | `/rules`                 | List all rules    |
| POST   | `/rules`                 | Create a new rule |
| GET    | `/rules/{id}`            | Get rule by ID    |
| PUT    | `/rules/{id}`            | Update rule by ID |
| DELETE | `/rules/{id}`            | Delete rule by ID |
| POST   | `/rules/{id}/activate`   | Activate a rule   |
| POST   | `/rules/{id}/deactivate` | Deactivate a rule |

**Audit**

| Method | Endpoint | Description                |
| ------ | -------- | -------------------------- |
| GET    | `/audit` | Retrieve audit logs/events |
 
---

## 🕒 Scheduler Logic

- Before Due: triggers a reminder X minutes before the task’s due date

- Interval: triggers a reminder every Y minutes until task is marked as done

- Scheduler runs every minute by default (configurable for demo purposes)

- Each reminder execution is logged in the audit trail

---

## 🛡️ Validation Rules for Reminder Rules

**To prevent duplicate or conflicting reminders, the system enforces the following constraints:**

1. **at_due Rule**
    
    - Only one at_due rule is allowed in the system.
    - Attempting to create another at_due rule will return an error.

2. **before_due Rule**

    - Cannot create a before_due rule with the same minutes_before value as an existing before_due rule.
    - Ensures that reminders scheduled for the same time before a task do not conflict.

3. **interval Rule**

    - Cannot create an interval rule with the same interval_min value as an existing interval rule.
    - Prevents duplicate interval reminders for the same task timing.

4. **Update Behavior**

    - The same validations are applied when updating a rule.
    - The rule being updated is ignored during conflict checks, so it can keep its own value if unchanged.

**Example Error Response:**

```bash
{
  "error": "before_due rule with 5 minutes already exists"
}
```

## 📋 Sample Tasks

Pre-seeded tasks for demonstration:

1. Submit assignment
2. Daily workout
3. Project meeting
4. Pay bills
5. Grocery shopping

---

## 💻 Demo / Console Output

**Example of scheduler logs and audit entries:**

- 2025-10-04T12:12:27+05:30 [INFO] [scheduler] run pass
- 2025-10-04T12:12:27+05:30 [INFO] IntervalReminder(rule:every 2 min) -> Task:7 Submit assignment (past due: 2025-10-03T22:39:37+05:30)
- 2025-10-04T12:12:27+05:30 [INFO] IntervalReminder(rule:every 2 min) -> Task:8 Daily workout (past due: 2025-10-03T23:22:00+05:30)

**Audit Log:**

- Oct 4, 2025, 05:54 PM — [reminder.trigger] Reminder triggered [Rule #2: every 2 min] -> [Task #4: Call supplier]
- Oct 4, 2025, 05:54 PM — [reminder.trigger] Reminder triggered [Rule #2: every 2 min] -> [Task #2: Submit assignment]
- Oct 4, 2025, 05:54 PM — [reminder.trigger] Reminder triggered [Rule #2: every 2 min] -> [Task #1: Pay electricity bill]
- Oct 4, 2025, 05:42 PM — [reminder.trigger] Reminder triggered [Rule #2: every 2 min] -> [Task #4: Call supplier]


✅ The console output shows reminders for tasks and their respective audit log entries.

## 🔧 Usage

- Load tasks via /tasks endpoint or UI.

- Create reminder rules (before_due or interval).

- Activate rules and let the scheduler trigger reminders.

- Check the audit logs for every reminder executed and rule change.

---


## 💡 Notes

- Reminders are simulated via console logs (no actual emails).

- Scheduler is idempotent, ensuring reminders are not duplicated for the same task in the same time window.

- Easily extendable for email/SMS notifications in the future.
