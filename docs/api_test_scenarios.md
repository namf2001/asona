# 📋 Asona API - Test Scenarios

> **Base URL**: `http://localhost:8080/api/v1`
> **Default Credentials**: `email: user1@example.com` / `password: password123`
> **Status Legend**: 🟢 Pass | 🔴 Fail | ⚪ Not Tested

---

## 1. Auth - Register

| # | Test Case | Method | Endpoint | Auth | Request Input | Expected Status | Expected Output | Status |
|---|-----------|--------|----------|------|---------------|-----------------|-----------------|--------|
| 1.1 | Register thành công | POST | `/register` | ❌ | `{"name":"Test User","username":"testuser","email":"test@test.com","password":"123456"}` | `201` | `{code:"INF001", data:{id, name, username, email}}` | ⚪ |
| 1.2 | Register - thiếu name | POST | `/register` | ❌ | `{"username":"testuser","email":"test@test.com","password":"123456"}` | `400` | `{code:"ERR001", message:"Invalid request parameters"}` | ⚪ |
| 1.3 | Register - thiếu email | POST | `/register` | ❌ | `{"name":"Test","username":"testuser","password":"123456"}` | `400` | `{code:"ERR001"}` | ⚪ |
| 1.4 | Register - email sai format | POST | `/register` | ❌ | `{"name":"Test","username":"testuser","email":"invalid","password":"123456"}` | `400` | `{code:"ERR001"}` | ⚪ |
| 1.5 | Register - password < 6 ký tự | POST | `/register` | ❌ | `{"name":"Test","username":"testuser","email":"test@test.com","password":"123"}` | `400` | `{code:"ERR001"}` | ⚪ |
| 1.6 | Register - email đã tồn tại | POST | `/register` | ❌ | `{"name":"Test","username":"newuser","email":"user1@example.com","password":"123456"}` | `409` | `{code:"ERR008", message:"...already registered"}` | ⚪ |
| 1.7 | Register - body rỗng | POST | `/register` | ❌ | `{}` | `400` | `{code:"ERR001"}` | ⚪ |

---

## 2. Auth - Login

| # | Test Case | Method | Endpoint | Auth | Request Input | Expected Status | Expected Output | Status |
|---|-----------|--------|----------|------|---------------|-----------------|-----------------|--------|
| 2.1 | Login thành công | POST | `/login` | ❌ | `{"email":"user1@example.com","password":"password123"}` | `200` | `{code:"INF003", data:{user:{id,name,username,email}, session_token}}` | ⚪ |
| 2.2 | Login - sai password | POST | `/login` | ❌ | `{"email":"user1@example.com","password":"wrongpass"}` | `401` | `{code:"ERR006"}` | ⚪ |
| 2.3 | Login - email không tồn tại | POST | `/login` | ❌ | `{"email":"noone@test.com","password":"123456"}` | `401` | `{code:"ERR006"}` | ⚪ |
| 2.4 | Login - thiếu email | POST | `/login` | ❌ | `{"password":"123456"}` | `400` | `{code:"ERR001"}` | ⚪ |
| 2.5 | Login - thiếu password | POST | `/login` | ❌ | `{"email":"user1@example.com"}` | `400` | `{code:"ERR001"}` | ⚪ |
| 2.6 | Login - body rỗng | POST | `/login` | ❌ | `{}` | `400` | `{code:"ERR001"}` | ⚪ |

---

## 3. Auth - Profile

| # | Test Case | Method | Endpoint | Auth | Request Input | Expected Status | Expected Output | Status |
|---|-----------|--------|----------|------|---------------|-----------------|-----------------|--------|
| 3.1 | Lấy profile thành công | GET | `/profile` | ✅ | Bearer token | `200` | `{code:"00", data:{id, name, username, email, image}}` | ⚪ |
| 3.2 | Profile - không có token | GET | `/profile` | ❌ | (none) | `401` | `{code:"ERR012"}` | ⚪ |
| 3.3 | Profile - token hết hạn | GET | `/profile` | ⚠️ | Expired token | `401` | `{code:"ERR011"}` | ⚪ |
| 3.4 | Profile - token sai format | GET | `/profile` | ⚠️ | `Bearer invalid_token` | `401` | `{code:"ERR009"}` | ⚪ |

---

## 4. Auth - Logout

| # | Test Case | Method | Endpoint | Auth | Request Input | Expected Status | Expected Output | Status |
|---|-----------|--------|----------|------|---------------|-----------------|-----------------|--------|
| 4.1 | Logout thành công | POST | `/logout` | ✅ | Bearer token | `200` | `{code:"INF004", message:"Logout success"}` | ⚪ |
| 4.2 | Logout - không có token | POST | `/logout` | ❌ | (none) | `401` | `{code:"ERR012"}` | ⚪ |

---

## 5. Auth - Google OAuth

| # | Test Case | Method | Endpoint | Auth | Request Input | Expected Status | Expected Output | Status |
|---|-----------|--------|----------|------|---------------|-----------------|-----------------|--------|
| 5.1 | Redirect tới Google login | GET | `/auth/google` | ❌ | (none) | `302` | Redirect tới Google OAuth consent screen | ⚪ |
| 5.2 | Google callback thành công | GET | `/auth/google/callback` | ❌ | `?code=valid_code&state=...` | `200` | `{code:"00", data:{user, session_token}}` | ⚪ |
| 5.3 | Google callback - code không hợp lệ | GET | `/auth/google/callback` | ❌ | `?code=invalid` | `401/500` | Error response | ⚪ |

---

## 6. Organizations - Create

| # | Test Case | Method | Endpoint | Auth | Request Input | Expected Status | Expected Output | Status |
|---|-----------|--------|----------|------|---------------|-----------------|-----------------|--------|
| 6.1 | Tạo organization thành công | POST | `/organizations` | ✅ | `{"name":"My Org","description":"Desc","logo_url":"https://..."}` | `201` | `{code:"00", data:{id, name, description, logo_url, created_at, updated_at}}` | ⚪ |
| 6.2 | Tạo org - thiếu name | POST | `/organizations` | ✅ | `{"description":"Desc"}` | `400` | `{code:"ERR001"}` | ⚪ |
| 6.3 | Tạo org - không có token | POST | `/organizations` | ❌ | `{"name":"My Org"}` | `401` | `{code:"ERR012"}` | ⚪ |
| 6.4 | Tạo org - body rỗng | POST | `/organizations` | ✅ | `{}` | `400` | `{code:"ERR001"}` | ⚪ |

---

## 7. Organizations - Get

| # | Test Case | Method | Endpoint | Auth | Request Input | Expected Status | Expected Output | Status |
|---|-----------|--------|----------|------|---------------|-----------------|-----------------|--------|
| 7.1 | Lấy org thành công | GET | `/organizations/1` | ✅ | Path: `id=1` | `200` | `{code:"00", data:{id, name, description, logo_url, created_at, updated_at}}` | ⚪ |
| 7.2 | Lấy org - ID không tồn tại | GET | `/organizations/9999` | ✅ | Path: `id=9999` | `404` | `{code:"ERR500"}` | ⚪ |
| 7.3 | Lấy org - ID sai format | GET | `/organizations/abc` | ✅ | Path: `id=abc` | `400` | `{code:"ERR001"}` | ⚪ |
| 7.4 | Lấy org - ID = 0 | GET | `/organizations/0` | ✅ | Path: `id=0` | `400` | `{code:"ERR001"}` | ⚪ |

---

## 8. Projects - Create

| # | Test Case | Method | Endpoint | Auth | Request Input | Expected Status | Expected Output | Status |
|---|-----------|--------|----------|------|---------------|-----------------|-----------------|--------|
| 8.1 | Tạo project thành công | POST | `/projects` | ✅ | `{"workplace_id":1,"name":"New Project","description":"Desc"}` | `201` | `{code:"00", data:{id, workplace_id, name, description, created_by, created_at, updated_at}}` | ⚪ |
| 8.2 | Tạo project - thiếu workplace_id | POST | `/projects` | ✅ | `{"name":"New Project"}` | `400` | `{code:"ERR001"}` | ⚪ |
| 8.3 | Tạo project - thiếu name | POST | `/projects` | ✅ | `{"workplace_id":1}` | `400` | `{code:"ERR001"}` | ⚪ |
| 8.4 | Tạo project - không có token | POST | `/projects` | ❌ | `{"workplace_id":1,"name":"Proj"}` | `401` | `{code:"ERR012"}` | ⚪ |
| 8.5 | Tạo project - workplace_id không tồn tại | POST | `/projects` | ✅ | `{"workplace_id":9999,"name":"Proj"}` | `500` | `{code:"ERR500"}` | ⚪ |

---

## 9. Projects - Get

| # | Test Case | Method | Endpoint | Auth | Request Input | Expected Status | Expected Output | Status |
|---|-----------|--------|----------|------|---------------|-----------------|-----------------|--------|
| 9.1 | Lấy project thành công | GET | `/projects/1` | ✅ | Path: `id=1` | `200` | `{code:"00", data:{id, workplace_id, name, description, created_by, created_at, updated_at}}` | ⚪ |
| 9.2 | Lấy project - ID không tồn tại | GET | `/projects/9999` | ✅ | Path: `id=9999` | `404` | `{code:"ERR500"}` | ⚪ |
| 9.3 | Lấy project - ID sai format | GET | `/projects/abc` | ✅ | Path: `id=abc` | `400` | `{code:"ERR001"}` | ⚪ |

---

## 10. Projects - List by Workplace

| # | Test Case | Method | Endpoint | Auth | Request Input | Expected Status | Expected Output | Status |
|---|-----------|--------|----------|------|---------------|-----------------|-----------------|--------|
| 10.1 | List projects thành công | GET | `/workplaces/1/projects` | ✅ | Path: `id=1` | `200` | `{code:"00", data:[{id, workplace_id, name, ...}]}` | ⚪ |
| 10.2 | List projects - workplace rỗng | GET | `/workplaces/9999/projects` | ✅ | Path: `id=9999` | `200` | `{code:"00", data:[]}` | ⚪ |
| 10.3 | List projects - ID sai format | GET | `/workplaces/abc/projects` | ✅ | Path: `id=abc` | `400` | `{code:"ERR001"}` | ⚪ |

---

## 11. Tasks - Create

| # | Test Case | Method | Endpoint | Auth | Request Input | Expected Status | Expected Output | Status |
|---|-----------|--------|----------|------|---------------|-----------------|-----------------|--------|
| 11.1 | Tạo task thành công | POST | `/tasks` | ✅ | `{"project_id":1,"title":"New Task","description":"Desc","priority":"high"}` | `201` | `{code:"00", data:{id, project_id, title, description, status, priority, created_by, ...}}` | ⚪ |
| 11.2 | Tạo task - thiếu project_id | POST | `/tasks` | ✅ | `{"title":"New Task"}` | `400` | `{code:"ERR001"}` | ⚪ |
| 11.3 | Tạo task - thiếu title | POST | `/tasks` | ✅ | `{"project_id":1}` | `400` | `{code:"ERR001"}` | ⚪ |
| 11.4 | Tạo task - không có token | POST | `/tasks` | ❌ | `{"project_id":1,"title":"Task"}` | `401` | `{code:"ERR012"}` | ⚪ |
| 11.5 | Tạo task - project_id không tồn tại | POST | `/tasks` | ✅ | `{"project_id":9999,"title":"Task"}` | `500` | `{code:"ERR500"}` | ⚪ |
| 11.6 | Tạo task - với due_date | POST | `/tasks` | ✅ | `{"project_id":1,"title":"Task","due_date":"2026-05-01T00:00:00Z"}` | `201` | `{code:"00", data:{..., due_date:"2026-05-01T..."}}` | ⚪ |
| 11.7 | Tạo task - với assignee_id | POST | `/tasks` | ✅ | `{"project_id":1,"title":"Task","assignee_id":2}` | `201` | `{code:"00", data:{..., assignee_id:2}}` | ⚪ |

---

## 12. Tasks - Get

| # | Test Case | Method | Endpoint | Auth | Request Input | Expected Status | Expected Output | Status |
|---|-----------|--------|----------|------|---------------|-----------------|-----------------|--------|
| 12.1 | Lấy task thành công | GET | `/tasks/1` | ✅ | Path: `id=1` | `200` | `{code:"00", data:{id, project_id, title, description, status, priority, ...}}` | ⚪ |
| 12.2 | Lấy task - ID không tồn tại | GET | `/tasks/9999` | ✅ | Path: `id=9999` | `404` | `{message:"Task not found"}` | ⚪ |
| 12.3 | Lấy task - ID sai format | GET | `/tasks/abc` | ✅ | Path: `id=abc` | `400` | `{code:"ERR001"}` | ⚪ |
| 12.4 | Lấy task - ID = 0 | GET | `/tasks/0` | ✅ | Path: `id=0` | `400` | `{code:"ERR001"}` | ⚪ |

---

## 13. Tasks - Update

| # | Test Case | Method | Endpoint | Auth | Request Input | Expected Status | Expected Output | Status |
|---|-----------|--------|----------|------|---------------|-----------------|-----------------|--------|
| 13.1 | Update task title thành công | PUT | `/tasks/1` | ✅ | `{"title":"Updated Title"}` | `200` | `{code:"00", message:"Task updated", data:{message:"Task updated successfully"}}` | ⚪ |
| 13.2 | Update task priority | PUT | `/tasks/1` | ✅ | `{"priority":"highest"}` | `200` | `{code:"00"}` | ⚪ |
| 13.3 | Update task status | PUT | `/tasks/1` | ✅ | `{"status":"in_progress"}` | `200` | `{code:"00"}` | ⚪ |
| 13.4 | Update task - ID không tồn tại | PUT | `/tasks/9999` | ✅ | `{"title":"Test"}` | `500` | `{code:"ERR500"}` | ⚪ |
| 13.5 | Update task - ID sai format | PUT | `/tasks/abc` | ✅ | `{"title":"Test"}` | `400` | `{code:"ERR001"}` | ⚪ |
| 13.6 | Update task - không có token | PUT | `/tasks/1` | ❌ | `{"title":"Test"}` | `401` | `{code:"ERR012"}` | ⚪ |
| 13.7 | Update task - thay đổi due_date | PUT | `/tasks/1` | ✅ | `{"due_date":"2026-06-01T00:00:00Z"}` | `200` | `{code:"00"}` | ⚪ |

---

## 14. Tasks - List by Project

| # | Test Case | Method | Endpoint | Auth | Request Input | Expected Status | Expected Output | Status |
|---|-----------|--------|----------|------|---------------|-----------------|-----------------|--------|
| 14.1 | List tasks thành công | GET | `/projects/1/tasks` | ✅ | Path: `id=1` | `200` | `{code:"00", data:[{id, project_id, title, ...}]}` | ⚪ |
| 14.2 | List tasks - project rỗng | GET | `/projects/9999/tasks` | ✅ | Path: `id=9999` | `200` | `{code:"00", data:[]}` | ⚪ |
| 14.3 | List tasks - ID sai format | GET | `/projects/abc/tasks` | ✅ | Path: `id=abc` | `400` | `{code:"ERR001"}` | ⚪ |
| 14.4 | List tasks - không có token | GET | `/projects/1/tasks` | ❌ | Path: `id=1` | `401` | `{code:"ERR012"}` | ⚪ |

---

## 15. Chat - Create Channel

| # | Test Case | Method | Endpoint | Auth | Request Input | Expected Status | Expected Output | Status |
|---|-----------|--------|----------|------|---------------|-----------------|-----------------|--------|
| 15.1 | Tạo channel thành công | POST | `/channels` | ✅ | `{"workplace_id":1,"name":"new-channel","type":"global"}` | `201` | `{code:"00", data:{id, workplace_id, name, type, created_by}}` | ⚪ |
| 15.2 | Tạo channel - thiếu workplace_id | POST | `/channels` | ✅ | `{"name":"ch","type":"global"}` | `400` | `{code:"ERR001"}` | ⚪ |
| 15.3 | Tạo channel - thiếu name | POST | `/channels` | ✅ | `{"workplace_id":1,"type":"global"}` | `400` | `{code:"ERR001"}` | ⚪ |
| 15.4 | Tạo channel - thiếu type | POST | `/channels` | ✅ | `{"workplace_id":1,"name":"ch"}` | `400` | `{code:"ERR001"}` | ⚪ |
| 15.5 | Tạo channel - không có token | POST | `/channels` | ❌ | `{"workplace_id":1,"name":"ch","type":"global"}` | `401` | `{code:"ERR012"}` | ⚪ |
| 15.6 | Tạo project channel | POST | `/channels` | ✅ | `{"workplace_id":1,"project_id":1,"name":"proj-ch","type":"project"}` | `201` | `{code:"00", data:{..., project_id:1}}` | ⚪ |

---

## 16. Chat - Get Channel

| # | Test Case | Method | Endpoint | Auth | Request Input | Expected Status | Expected Output | Status |
|---|-----------|--------|----------|------|---------------|-----------------|-----------------|--------|
| 16.1 | Lấy channel thành công | GET | `/channels/1` | ✅ | Path: `id=1` | `200` | `{code:"00", data:{id, workplace_id, name, type, created_by}}` | ⚪ |
| 16.2 | Lấy channel - ID không tồn tại | GET | `/channels/9999` | ✅ | Path: `id=9999` | `404` | `{code:"ERR045"}` | ⚪ |
| 16.3 | Lấy channel - ID sai format | GET | `/channels/abc` | ✅ | Path: `id=abc` | `400` | `{code:"ERR001"}` | ⚪ |

---

## 17. Chat - Send Message

| # | Test Case | Method | Endpoint | Auth | Request Input | Expected Status | Expected Output | Status |
|---|-----------|--------|----------|------|---------------|-----------------|-----------------|--------|
| 17.1 | Gửi message thành công | POST | `/messages` | ✅ | `{"channel_id":1,"content":"Hello world"}` | `201` | `{code:"00", data:{id, channel_id, sender_id, content, is_edited, created_at, updated_at}}` | ⚪ |
| 17.2 | Gửi message - reply (thread) | POST | `/messages` | ✅ | `{"channel_id":1,"parent_id":1,"content":"Reply msg"}` | `201` | `{code:"00", data:{..., parent_id:1}}` | ⚪ |
| 17.3 | Gửi message - thiếu channel_id | POST | `/messages` | ✅ | `{"content":"Hello"}` | `400` | `{code:"ERR001"}` | ⚪ |
| 17.4 | Gửi message - thiếu content | POST | `/messages` | ✅ | `{"channel_id":1}` | `400` | `{code:"ERR001"}` | ⚪ |
| 17.5 | Gửi message - không có token | POST | `/messages` | ❌ | `{"channel_id":1,"content":"Hi"}` | `401` | `{code:"ERR012"}` | ⚪ |
| 17.6 | Gửi message - channel_id không tồn tại | POST | `/messages` | ✅ | `{"channel_id":9999,"content":"Hi"}` | `500` | `{code:"ERR051"}` | ⚪ |

---

## 18. Chat - List Messages

| # | Test Case | Method | Endpoint | Auth | Request Input | Expected Status | Expected Output | Status |
|---|-----------|--------|----------|------|---------------|-----------------|-----------------|--------|
| 18.1 | List messages thành công | GET | `/channels/1/messages` | ✅ | Path: `id=1` | `200` | `{code:"00", data:[{id, channel_id, sender_id, content, ...}]}` | ⚪ |
| 18.2 | List messages - với limit/offset | GET | `/channels/1/messages?limit=2&offset=1` | ✅ | Path: `id=1`, Query: `limit=2&offset=1` | `200` | `{code:"00", data:[...]}` (max 2 items) | ⚪ |
| 18.3 | List messages - channel rỗng | GET | `/channels/9999/messages` | ✅ | Path: `id=9999` | `200`/`500` | `{data:[]}` hoặc error | ⚪ |
| 18.4 | List messages - ID sai format | GET | `/channels/abc/messages` | ✅ | Path: `id=abc` | `400` | `{code:"ERR001"}` | ⚪ |
| 18.5 | List messages - limit âm | GET | `/channels/1/messages?limit=-1` | ✅ | Path: `id=1`, Query: `limit=-1` | `200` | Default limit=50 applied | ⚪ |

---

## 📊 Tổng kết

| API Group | Tổng Test Cases | Happy Path | Error Cases |
|---|---|---|---|
| Auth - Register | 7 | 1 | 6 |
| Auth - Login | 6 | 1 | 5 |
| Auth - Profile | 4 | 1 | 3 |
| Auth - Logout | 2 | 1 | 1 |
| Auth - Google | 3 | 2 | 1 |
| Organizations - Create | 4 | 1 | 3 |
| Organizations - Get | 4 | 1 | 3 |
| Projects - Create | 5 | 1 | 4 |
| Projects - Get | 3 | 1 | 2 |
| Projects - List | 3 | 1 | 2 |
| Tasks - Create | 7 | 3 | 4 |
| Tasks - Get | 4 | 1 | 3 |
| Tasks - Update | 7 | 4 | 3 |
| Tasks - List | 4 | 1 | 3 |
| Chat - Create Channel | 6 | 2 | 4 |
| Chat - Get Channel | 3 | 1 | 2 |
| Chat - Send Message | 6 | 2 | 4 |
| Chat - List Messages | 5 | 1 | 4 |
| **TỔNG** | **83** | **27** | **56** |

---

## 📝 Ghi chú

1. **RSA Middleware**: Khi `APP_ENV != local`, request body phải được encrypt bằng RSA. Các test case ở trên giả định `APP_ENV=local` (skip RSA).
2. **Seed Data**: Chạy `make seed` trước khi test để có data phù hợp.
3. **Token**: Sử dụng session token từ Login response cho các API authenticated.
4. **WebSocket**: Không bao gồm ở đây vì cần test riêng (không phải REST API).
