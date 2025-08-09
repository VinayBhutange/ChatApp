# Postman Beginner's Guide for Chat Application

## Table of Contents
1. [What is Postman?](#what-is-postman)
2. [Installing Postman](#installing-postman)
3. [Postman Interface Overview](#postman-interface-overview)
4. [Setting up Environment Variables](#setting-up-environment-variables)
5. [Manual API Testing Step-by-Step](#manual-api-testing-step-by-step)
6. [Understanding HTTP Methods](#understanding-http-methods)
7. [Working with Headers](#working-with-headers)
8. [Request Body Types](#request-body-types)
9. [Reading Responses](#reading-responses)
10. [Creating Collections](#creating-collections)
11. [Writing Tests](#writing-tests)
12. [Importing Collections](#importing-collections)
13. [Advanced Features](#advanced-features)

## What is Postman?

Postman is a popular API testing tool that allows you to:
- Send HTTP requests to APIs
- Test API endpoints
- Organize requests into collections
- Set up environment variables
- Write automated tests
- Generate API documentation

Think of it as a user-friendly interface to interact with your backend APIs without needing a frontend application.

## Installing Postman

1. **Download Postman:**
   - Go to [https://www.postman.com/downloads/](https://www.postman.com/downloads/)
   - Download the desktop app for Windows
   - Install and create a free account (recommended for saving your work)

2. **Alternative - Web Version:**
   - You can also use Postman in your browser at [https://web.postman.co/](https://web.postman.co/)
   - Requires account creation

## Postman Interface Overview

When you open Postman, you'll see:

```
┌─────────────────────────────────────────────────────────────┐
│ File Edit View Help                                         │
├─────────────────────────────────────────────────────────────┤
│ [Workspaces] [Collections] [APIs] [Environments]           │
├─────────────────────────────────────────────────────────────┤
│ Sidebar          │ Main Request Area                        │
│ - Collections    │ ┌─────────────────────────────────────┐   │
│ - Environments   │ │ GET  [URL Input Box]        [Send] │   │
│ - History        │ ├─────────────────────────────────────┤   │
│                  │ │ Params | Authorization | Headers   │   │
│                  │ │ Body | Pre-request | Tests         │   │
│                  │ ├─────────────────────────────────────┤   │
│                  │ │ Response Area                       │   │
│                  │ │ - Status, Time, Size               │   │
│                  │ │ - Body, Headers, Test Results      │   │
│                  │ └─────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

**Key Components:**
- **Request Builder:** Where you build your HTTP requests
- **Response Viewer:** Shows the API response
- **Collections:** Organize related requests
- **Environments:** Store variables like URLs and tokens
- **History:** Shows previously sent requests

## Setting up Environment Variables

Environment variables help you avoid repeating values like URLs and tokens.

### Step 1: Create an Environment
1. Click the **"Environments"** tab in the sidebar
2. Click **"Create Environment"** (+ icon)
3. Name it: `Chat App Local`

### Step 2: Add Variables
Add these variables to your environment:

| Variable Name | Initial Value | Current Value |
|---------------|---------------|---------------|
| `base_url` | `http://localhost:8082` | `http://localhost:8082` |
| `auth_token` | (leave empty) | (leave empty) |
| `user_id` | (leave empty) | (leave empty) |
| `room_id` | (leave empty) | (leave empty) |

### Step 3: Select Environment
1. Click the environment dropdown (top right)
2. Select **"Chat App Local"**

**Why use variables?**
- Easy to switch between development/production
- Automatically update tokens after login
- Reuse values across multiple requests

## Manual API Testing Step-by-Step

Let's test your chat application APIs manually. **Make sure your backend server is running first!**

### Step 1: Start Your Backend Server

Open terminal in your backend folder and run:
```bash
cd C:\CodeBase\PlayGround\ChatApp\backend
go run cmd/server/main.go
```

You should see:
```
Starting Chat Application Server...
Server starting on port 8082...
```

### Test 1: Health Check

**Purpose:** Verify server is running

1. **Create New Request:**
   - Click **"New"** → **"HTTP Request"**
   - Or use **Ctrl+N**

2. **Set Request Details:**
   - **Method:** `GET` (default)
   - **URL:** `{{base_url}}/`
   - **Headers:** None needed
   - **Body:** None needed

3. **Send Request:**
   - Click **"Send"** button
   - **Expected Response:** `200 OK` with text "WebSocket Chat Server"

4. **What to Look For:**
   - **Status:** `200 OK` (green)
   - **Response Body:** "WebSocket Chat Server"
   - **Time:** Should be under 100ms

**If it fails:**
- Check if server is running
- Verify URL is correct
- Check if port 8082 is blocked

### Test 2: Register New User

**Purpose:** Create a new user account

1. **Create New Request:**
   - Click **"New"** → **"HTTP Request"**

2. **Set Request Details:**
   - **Method:** Change to `POST`
   - **URL:** `{{base_url}}/api/register`

3. **Add Headers:**
   - Click **"Headers"** tab
   - Add: `Content-Type` = `application/json`

4. **Add Request Body:**
   - Click **"Body"** tab
   - Select **"raw"**
   - Select **"JSON"** from dropdown
   - Enter:
   ```json
   {
     "username": "testuser1",
     "password": "password123"
   }
   ```

5. **Send Request:**
   - Click **"Send"**
   - **Expected Response:** `200 OK` with user details

6. **Response Should Look Like:**
   ```json
   {
     "message": "User registered successfully",
     "user": {
       "id": "some-uuid-here",
       "username": "testuser1"
     }
   }
   ```

7. **Save User ID:**
   - Copy the `id` from response
   - Go to **Environments** → **Chat App Local**
   - Set `user_id` = copied ID value

**Common Issues:**
- **400 Bad Request:** Check JSON format
- **500 Internal Server Error:** Check server logs
- **Username already exists:** Use different username

### Test 3: Login User

**Purpose:** Authenticate and get JWT token

1. **Create New Request:**
   - **Method:** `POST`
   - **URL:** `{{base_url}}/api/login`

2. **Add Headers:**
   - `Content-Type` = `application/json`

3. **Add Body:**
   ```json
   {
     "username": "testuser1",
     "password": "password123"
   }
   ```

4. **Send Request:**
   - **Expected Response:** `200 OK` with token

5. **Response Should Include:**
   ```json
   {
     "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
     "user": {
       "id": "user-uuid",
       "username": "testuser1"
     }
   }
   ```

6. **Save Token:**
   - Copy the entire `token` value
   - Go to **Environments** → **Chat App Local**
   - Set `auth_token` = copied token

**Important:** This token is required for all protected endpoints!

### Test 4: Create Chat Room

**Purpose:** Create a new chat room (requires authentication)

1. **Create New Request:**
   - **Method:** `POST`
   - **URL:** `{{base_url}}/api/rooms/create`

2. **Add Headers:**
   - `Content-Type` = `application/json`
   - `Authorization` = `Bearer {{auth_token}}`

   **Note:** The `{{auth_token}}` will automatically use the token from your environment!

3. **Add Body:**
   ```json
   {
     "name": "My Test Room",
     "roomType": "public"
   }
   ```

4. **Send Request:**
   - **Expected Response:** `200 OK` with room details

5. **Response Should Include:**
   ```json
   {
     "message": "Room created successfully",
     "room": {
       "id": "room-uuid",
       "name": "My Test Room",
       "ownerId": "your-user-id",
       "roomType": "public"
     }
   }
   ```

6. **Save Room ID:**
   - Copy the room `id`
   - Set `room_id` in environment

**Common Issues:**
- **401 Unauthorized:** Check Authorization header format
- **Token expired:** Login again to get new token

### Test 5: List All Rooms

**Purpose:** Get all available chat rooms

1. **Create New Request:**
   - **Method:** `GET`
   - **URL:** `{{base_url}}/api/rooms`

2. **Add Headers:**
   - `Authorization` = `Bearer {{auth_token}}`

3. **Send Request:**
   - **Expected Response:** `200 OK` with rooms array

4. **Response Should Look Like:**
   ```json
   {
     "rooms": [
       {
         "id": "room-uuid",
         "name": "My Test Room",
         "ownerId": "your-user-id",
         "roomType": "public"
       }
     ]
   }
   ```

## Understanding HTTP Methods

Your chat app uses these HTTP methods:

| Method | Purpose | Example |
|--------|---------|---------|
| `GET` | Retrieve data | Get list of rooms |
| `POST` | Create new data | Register user, Create room |
| `PUT` | Update existing data | Update room settings |
| `DELETE` | Remove data | Delete room |

**In Postman:**
- Click the dropdown next to URL input
- Select appropriate method
- Different methods may require different body content

## Working with Headers

Headers provide metadata about your request:

### Common Headers in Chat App:

1. **Content-Type: application/json**
   - Tells server you're sending JSON data
   - Required for POST requests with JSON body

2. **Authorization: Bearer {token}**
   - Provides authentication token
   - Required for protected endpoints
   - Format: `Bearer ` + your JWT token

### Adding Headers in Postman:
1. Click **"Headers"** tab below URL
2. Click **"Key"** field, type header name
3. Click **"Value"** field, type header value
4. Use `{{variable}}` for environment variables

## Request Body Types

Different endpoints require different body types:

### 1. No Body (GET requests)
- Used for: Health check, List rooms
- Body tab: Select **"none"**

### 2. JSON Body (POST requests)
- Used for: Register, Login, Create room
- Body tab: Select **"raw"** → **"JSON"**
- Content must be valid JSON format

### 3. Form Data
- Used for: File uploads (if implemented)
- Body tab: Select **"form-data"**

**JSON Format Tips:**
- Use double quotes for strings: `"username"`
- No trailing commas
- Validate JSON online if unsure

## Reading Responses

Postman shows response information in several tabs:

### 1. Status Information
```
Status: 200 OK    Time: 45ms    Size: 156 B
```

**Status Codes:**
- `200 OK` - Success
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Authentication required/failed
- `404 Not Found` - Endpoint doesn't exist
- `500 Internal Server Error` - Server error

### 2. Response Body
Shows the actual data returned by API:
- **Pretty:** Formatted JSON (easiest to read)
- **Raw:** Unformatted text
- **Preview:** HTML preview (for web pages)

### 3. Headers Tab
Shows response headers from server:
```
Content-Type: application/json
Content-Length: 156
Date: Thu, 08 Aug 2024 14:49:50 GMT
```

### 4. Test Results
Shows results of automated tests (if written)

## Creating Collections

Collections organize related requests:

### Step 1: Create Collection
1. Click **"Collections"** in sidebar
2. Click **"Create Collection"**
3. Name: `Chat Application API`
4. Add description: `API endpoints for chat application`

### Step 2: Add Requests to Collection
1. Click **"Save"** after creating a request
2. Choose your collection
3. Give request a descriptive name

### Step 3: Organize with Folders
1. Right-click collection → **"Add Folder"**
2. Create folders like:
   - `Authentication` (Register, Login)
   - `Room Management` (Create, List rooms)
   - `Error Testing` (Invalid requests)

**Benefits:**
- Keep related requests together
- Share collections with team
- Run all requests in sequence

## Writing Tests

Tests automatically verify API responses:

### Basic Test Structure
Click **"Tests"** tab and add:

```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Response has message", function () {
    const response = pm.response.json();
    pm.expect(response.message).to.exist;
});
```

### Practical Tests for Chat App

#### For Registration:
```javascript
pm.test("User registration successful", function () {
    pm.response.to.have.status(200);
    const response = pm.response.json();
    pm.expect(response.message).to.include("successfully");
    
    // Save user ID to environment
    if (response.user && response.user.id) {
        pm.environment.set("user_id", response.user.id);
    }
});
```

#### For Login:
```javascript
pm.test("Login successful", function () {
    pm.response.to.have.status(200);
    const response = pm.response.json();
    pm.expect(response.token).to.exist;
    
    // Save token to environment
    pm.environment.set("auth_token", response.token);
});
```

#### For Room Creation:
```javascript
pm.test("Room created successfully", function () {
    pm.response.to.have.status(200);
    const response = pm.response.json();
    pm.expect(response.room.id).to.exist;
    
    // Save room ID
    pm.environment.set("room_id", response.room.id);
});
```

**Test Benefits:**
- Automatic validation
- Environment variable updates
- Quick error detection

## Importing Collections

To use the provided collection file:

### Step 1: Import Collection
1. Click **"Import"** button (top left)
2. Click **"Upload Files"**
3. Select `ChatApp-Postman-Collection.json`
4. Click **"Import"**

### Step 2: Set Environment
1. Make sure your environment is selected
2. Update `base_url` if needed

### Step 3: Run Collection
1. Right-click collection → **"Run collection"**
2. Select requests to run
3. Click **"Run Chat Application API"**

**Collection Benefits:**
- Pre-configured requests
- Automated tests included
- Environment variables setup
- Error cases included

## Advanced Features

### 1. Pre-request Scripts
Run code before sending request:

```javascript
// Generate random username
const randomUser = "user" + Math.floor(Math.random() * 1000);
pm.environment.set("random_username", randomUser);
```

### 2. Collection Variables
Variables specific to a collection:
- Right-click collection → **"Edit"**
- Go to **"Variables"** tab
- Add collection-specific values

### 3. Data Files
Test with multiple data sets:
- Create CSV/JSON file with test data
- Use in Collection Runner
- Test multiple scenarios automatically

### 4. Mock Servers
Create fake API responses:
- Useful when backend isn't ready
- Test frontend with mock data

### 5. Documentation
Generate API documentation:
- Add descriptions to requests
- Include examples
- Share with team

## Troubleshooting Common Issues

### 1. "Could not get any response"
**Causes:**
- Server not running
- Wrong URL
- Firewall blocking

**Solutions:**
- Check server is running on correct port
- Verify URL spelling
- Try `127.0.0.1` instead of `localhost`

### 2. "401 Unauthorized"
**Causes:**
- Missing Authorization header
- Wrong token format
- Expired token

**Solutions:**
- Check Authorization header: `Bearer {token}`
- Login again to get fresh token
- Verify token is copied completely

### 3. "400 Bad Request"
**Causes:**
- Invalid JSON format
- Missing required fields
- Wrong Content-Type header

**Solutions:**
- Validate JSON syntax
- Check API documentation for required fields
- Ensure Content-Type is `application/json`

### 4. Environment Variables Not Working
**Causes:**
- Wrong environment selected
- Variable name mismatch
- Typo in variable usage

**Solutions:**
- Check environment dropdown (top right)
- Verify variable names match exactly
- Use `{{variable_name}}` format

## Practice Exercises

Try these exercises to master Postman:

### Exercise 1: Complete User Flow
1. Register a new user
2. Login with that user
3. Create a room
4. List all rooms
5. Verify your room appears in the list

### Exercise 2: Error Testing
1. Try registering with same username twice
2. Try logging in with wrong password
3. Try creating room without authentication
4. Observe different error responses

### Exercise 3: Multiple Users
1. Register two different users
2. Login as first user, create a room
3. Login as second user, list rooms
4. Verify both users can see the room

### Exercise 4: Environment Management
1. Create "Development" environment
2. Create "Production" environment (when available)
3. Switch between environments
4. Test same requests on different environments

## Next Steps

After mastering these basics:

1. **Learn Collection Runner**
   - Run multiple requests in sequence
   - Use data files for testing

2. **Explore Newman**
   - Command-line collection runner
   - Integrate with CI/CD pipelines

3. **API Documentation**
   - Generate docs from collections
   - Share with team members

4. **Advanced Scripting**
   - Complex pre-request scripts
   - Dynamic data generation
   - Custom test assertions

## Quick Reference

### Essential Keyboard Shortcuts
- `Ctrl+N` - New request
- `Ctrl+S` - Save request
- `Ctrl+Enter` - Send request
- `Ctrl+/` - Search collections

### Common Variable Usage
- `{{base_url}}` - Base API URL
- `{{auth_token}}` - JWT authentication token
- `{{user_id}}` - Current user ID
- `{{room_id}}` - Current room ID

### HTTP Status Codes
- `200` - OK (Success)
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `404` - Not Found
- `500` - Server Error

This guide should give you a solid foundation for using Postman to test your chat application APIs. Start with the manual testing steps, then gradually explore the advanced features as you become more comfortable with the tool!
