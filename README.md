---

# 🚀 API Documentation for Duck Duck Go Chat Application

Welcome to the API documentation for the **Duck Duck Go Chat Application**! This guide will walk you through how to use the endpoints with **Postman**. The application supports creating chat sessions, sending messages, deleting chat sessions, and checking the health of the service.

---

## **📂 Endpoints**

The application exposes the following endpoints:

1. **POST `/chat/:model`**: Create or interact with a chat session. 💬
2. **DELETE `/chat/:model`**: Delete an existing chat session. 🗑️
3. **GET `/health`**: Check the health/status of the service. 🩺

---

## **1. POST `/chat/:model`** 💬

This endpoint is used to create a new chat session or send a message to an existing chat session.

### **📝 Request Details**
- **Method**: `POST`
- **URL**: `http://localhost:8080/chat/{model}`
  - Replace `{model}` with one of the following model aliases:
    - `gpt-4o-mini` 🤖
    - `claude-3-haiku` 🐧
    - `llama` 🦙
    - `mixtral` 🌪️
- **Headers**:
  - `User-ID`: A unique identifier for the user (e.g., `12345`). 🆔
  - `Content-Type`: `application/json`. 📄
- **Body** (JSON):
  ```json
  {
    "message": "Your message here"
  }
  ```

### **🌐 Example Request**
- **URL**: `http://localhost:8080/chat/gpt-4o-mini`
- **Headers**:
  ```
  User-ID: 12345
  Content-Type: application/json
  ```
- **Body**:
  ```json
  {
    "message": "Hello, how are you?"
  }
  ```

### **📤 Example Response**
```json
{
  "response": "Hello! I'm just a program, so I don't have feelings, but I'm here to help. How can I assist you today?"
}
```

---

## **2. DELETE `/chat/:model`** 🗑️

This endpoint is used to delete an existing chat session for a specific user.

### **📝 Request Details**
- **Method**: `DELETE`
- **URL**: `http://localhost:8080/chat/{model}`
  - Replace `{model}` with the model alias used to create the chat session.
- **Headers**:
  - `User-ID`: The unique identifier for the user whose session is being deleted. 🆔

### **🌐 Example Request**
- **URL**: `http://localhost:8080/chat/gpt-4o-mini`
- **Headers**:
  ```
  User-ID: 12345
  ```

### **📤 Example Response**
```json
{
  "message": "Chat session deleted"
}
```

---

## **3. GET `/health`** 🩺

This endpoint is used to check the health/status of the service.

### **📝 Request Details**
- **Method**: `GET`
- **URL**: `http://localhost:8080/health`

### **🌐 Example Request**
- **URL**: `http://localhost:8080/health`

### **📤 Example Response**
```json
{
  "status": "ok"
}
```

---

## **📊 Summary of Endpoints**

| **Endpoint**            | **Method** | **URL Example**                          | **Headers**                          | **Body**                              |
|--------------------------|------------|------------------------------------------|--------------------------------------|---------------------------------------|
| Create/Send Chat Message | `POST`     | `http://localhost:8080/chat/gpt-4o-mini` | `User-ID: 12345`, `Content-Type: application/json` | `{"message": "Hello, how are you?"}` |
| Delete Chat Session      | `DELETE`   | `http://localhost:8080/chat/gpt-4o-mini` | `User-ID: 12345`                     | None                                  |
| Health Check             | `GET`      | `http://localhost:8080/health`           | None                                 | None                                  |

---

## **📌 Notes**
1. **User-ID**: Ensure you use the same `User-ID` for creating and deleting chat sessions. This is how the application tracks sessions for each user. 🆔
2. **Streaming Responses**: The `POST /chat/:model` endpoint streams responses in chunks. To observe the streaming behavior, use the **Postman Console** (`View > Show Postman Console`). 📡
3. **Errors**: If you encounter errors (e.g., invalid model, missing `User-ID`), Postman will display the error message in the response body. ❌

---

## **🔧 Example Workflow in Postman**
1. **Create a Chat Session**:
   - Send a `POST` request to `http://localhost:8080/chat/gpt-4o-mini` with a `User-ID` and a message.
   - Receive the assistant's response. 💬

2. **Delete the Chat Session**:
   - Send a `DELETE` request to `http://localhost:8080/chat/gpt-4o-mini` with the same `User-ID`.
   - Receive a confirmation that the session was deleted. 🗑️

3. **Health Check**:
   - Send a `GET` request to `http://localhost:8080/health` to ensure the service is running. 🩺

---
