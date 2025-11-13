# Better Auth Integration Guide

This document provides a complete overview of how to integrate **Better Auth** into a full-stack architecture combining **SvelteKit (Bun)**, **Go (Echo)**, and **PostgreSQL**.

It explains:
- What Better Auth is
- How the authentication flow works
- How to set up the SvelteKit auth server
- How to protect pages in SvelteKit
- How to validate JWT in your Go API using JWKS
- How to handle social login (Google)
- Useful development notes

Code samples are intentionally excluded — instead, this guide links directly to the corresponding files inside your repository.

---

## 🔐 What Better Auth Is

**Better Auth** is an all‑in‑one authentication engine designed to run directly inside your SvelteKit application. It provides:

- Email/password authentication
- Social login (Google, GitHub, Discord, etc.)
- Database persistence (PostgreSQL in your case)
- Secure session handling
- Built‑in JWT support for external APIs
- Automatic user/account linking (e.g., Google login on an email/password account)

Better Auth exposes its features through:
- A server‑side API (`auth.api.*`)
- A client‑side SDK (`authClient.*`)
- A generated router (`/api/auth/...`)
- A JWKS endpoint (`/api/auth/jwks`) used by external services

---

## 🧩 How the Architecture Works

Your stack contains 3 components:

### **1. SvelteKit App (Bun) — Authentication Layer**

Responsible for:
- Serving Better Auth
- Generating sessions
- Generating JWTs for your Go API
- Protecting routes (`/me`, `/todos`)
- Rendering UI for login/register/todos

See:
- [`src/lib/auth.ts`](./frontend/src/lib/auth.ts)
- [`src/hooks.server.ts`](./frontend/src/hooks.server.ts)
- [`src/lib/auth-client.ts`](./frontend/src/lib/auth-client.ts)

### **2. Go API — Business Logic**

Responsible for:
- Verifying JWT tokens issued by Better Auth
- Extracting the user ID from claims
- Performing business logic (todos CRUD)

See:
- [`middleware/jwt.go`](./backend/middleware/auth.go)
- [`handlers/todos.go`](./backend/handlers/todos.go)
- [`main.go`](./backend/main.go)

### **3. PostgreSQL — Persistent Storage**

Used by Better Auth and your Go API.

Better Auth automatically manages:
- `user` table
- `account` table (for social providers)
- `session` table

Your API adds its own tables (e.g., `todos`).
See:

- [`migrations`](./frontend/better-auth_migrations/2025-11-12T23-02-36.934Z.sql)

---

## 🔄 End‑to‑End Authentication Flow

The full flow looks like this:
1. **User logs in** (email/password or Google) from SvelteKit
2. **Better Auth creates a session** stored in PostgreSQL
3. The frontend requests a **JWT** using `authClient.token()`
4. The JWT is sent to the Go API via `Authorization: Bearer <token>`
5. The Go API verifies the token using:
   - Better Auth JWKS endpoint
   - Ed25519 signature
   - `sub` claim (the user ID)
6. The Go API returns protected data
This creates a secure boundary:

- **SvelteKit handles auth**

- **Go handles protected business logic**

---

## 🏗️ Step 1 — Install Better Auth (SvelteKit + Bun)

Better Auth runs fully inside your frontend server.

You installed it via:
```shell
bun add better-auth better-auth/svelte
```

Then configure the server instance:
- Database connection
- Secret
- Base URL
- Enabled auth methods (email/password + Google)
- JWT plugin (for Go)

See:
- [`src/lib/auth.ts`](./frontend/src/lib/auth.ts)

---

## 🗂️ Step 2 — Install the Client SDK

Installs the Svelte client:
```shell
bun add better-auth/client
```

Configured through:

- [`src/lib/auth-client.ts`](./frontend/src/lib/auth-client.ts)
This file also includes the **JWT client plugin**, which is required so your Go API can accept tokens.

---

## 🛡️ Step 3 — Protect Server Routes in SvelteKit

Your server hook:
- Fetches the current session
- Populates `event.locals`
- Protects `/me` and `/todos`
- Hands routing control to Better Auth

See:
- [`src/hooks.server.ts`](./frontend/src/hooks.server.ts)

This ensures that:
- Any unauthenticated request to protected routes results in a redirect to `/login`
- Client scripts don't even run if the user is not authenticated

---

## 🧪 Step 4 — Frontend Pages

### Implemented pages:

- `/login`
- `/signup`
- `/me`
- `/todos`

These pages consume the Better Auth client:

- `signIn.email`
- `signUp.email`
- `signIn.social` (Google)
- `useSession()` / `getSession()`
- `token()` (for Go API calls)

See:
- [`src/routes/login/+page.svelte`](./frontend/src/routes/login/+page.svelte)
- [`src/routes/signup/+page.svelte`](./frontend/src/routes/signup/+page.svelte)
- [`src/routes/me/+page.svelte`](./frontend/src/routes/me/+page.svelte)
- [`src/routes/todos/+page.svelte`](./frontend/src/routes/todos/+page.svelte)

---

## 🌐 Step 5 — Adding Google OAuth Login

Better Auth supports social login via providers. In this project, Google is used as an example, but the same pattern applies to GitHub, Discord, etc.

The idea is:
- Google authenticates the user and returns an authorization code
- Better Auth exchanges this code for an access token and profile
- Better Auth either creates a new user or links the Google account to an existing user (same verified email)
- The user ends up with a normal Better Auth session, just like with email/password

You do not talk directly to Google from the frontend or the Go API. Only Better Auth does.

### 5.1 Creating Google OAuth Credentials

You first need an OAuth 2.0 Client ID from Google:

1. Go to Google Cloud Console: APIs & Services → Credentials.
2. If needed, create or select a project.
3. Click “Create Credentials” → “OAuth client ID”.
4. Choose Application type: Web application.
5. Set a descriptive name (e.g. better-auth-sveltekit-local).
6. In Authorized redirect URIs, add the URL handled by Better Auth:
  - For local development with your current setup:
    - http://localhost:5173/api/auth/callback/google
  - In production, this must match your real domain, for example:
    - https://your-domain.com/api/auth/callback/google
7. Validate and copy the generated Client ID and Client Secret.

You then store them in your SvelteKit .env (never commit these):

 - ``GOOGLE_CLIENT_ID``
 - ``GOOGLE_CLIENT_SECRET``

These values are read in your Better Auth configuration.

### 5.2 Wiring Google into Better Auth

On the SvelteKit side, Google is plugged into Better Auth via the ``socialProviders.google`` section. This tells Better Auth:
 - which clientId and clientSecret to use
 - for which provider name (google) requests should be handled

See:
- [`src/lib/auth.ts`](./frontend/src/lib/auth.ts)

When a user authenticates with Google for the first time:
 - Better Auth fetches the user profile
 - If the email is verified and matches an existing user, Better Auth links the Google account to that user instead of creating a duplicate
 - Otherwise, a new user row is created in the user table, with a corresponding entry in the account table for the google provider

This is why, in your tests, signing in with Google on an email that already existed did not create a second account.

### 5.3 Adding the “Continue with Google” Button (Frontend)

On the login page, you expose Google as a one-click option next to the email/password form.

The button:
 - calls the Better Auth client (authClient.signIn.social)
 - passes provider: "google"
 - optionally passes callbackURL (where to go after a successful login) and errorCallbackURL (where to go on failure)

Better Auth then:
 - redirects the browser to Google’s consent screen
 - receives the callback on /api/auth/callback/google
 - creates/links the user and session
 - redirects back to your callbackURL, where the user is now considered logged in

See:
 - [`src/routes/login/+page.svelte`](./frontend/src/routes/login/+page.svelte)

From this point on, there is no difference for the rest of the stack:
 - The session object looks the same to your SvelteKit pages
 - The JWTs issued to the Go API work exactly the same
 - The Go backend never needs to know whether the user logged in via email/password or Google.

## 🔑 Step 6 — Token Validation in Go (JWKS)

Better Auth exposes a JWKS endpoint used by your Go API to verify JWT signatures.
JWKS URL:

```
/api/auth/jwks
```
Your Go server:
- Fetches JWKS at startup
- Caches + auto-refreshes keys
- Validates EdDSA tokens
- Extracts the `sub` claim (user ID)

See:
- [`middleware/jwt.go`](./backend/middleware/auth.go)
- [`main.go`](./backend/main.go)

---

## 🧱 Step 7 — Protecting Go API Routes

You protect routes using Echo middleware.

It:
- Reads `Authorization: Bearer` header
- Verifies the JWT signature via JWKS
- Checks standard claims (expiry, etc.)
- Injects claims into the Echo context

Protected routes include:
- `GET /api/todos`
- `POST /api/todos`
- `PATCH /api/todos/:id/toggle`
- `DELETE /api/todos/:id`

Handlers simply extract the user ID from claims.

See:
- [`middleware/jwt.go`](./backend/middleware/auth.go)
- [`handlers/todos.go`](./backend/handlers/todos.go)

---

## 🗃️ Step 8 — Database Schema (PostgreSQL)

### Better Auth tables (auto‑managed):

- `user`
- `account`
- `session`

### Your tables:

- `todos` with user foreign key

See migrations:
- [`better auth migrations`](./frontend/better-auth_migrations/2025-11-12T23-02-36.934Z.sql)
- [`todos table migration`](./backend/database/20251113001752_todos_table.sql)

---

## 🎉 Summary

This stack achieves:
- Full session-based authentication in SvelteKit
- Automatic JWT issuance for external APIs
- Secure user account linking across providers
- Clean Go API authentication via JWKS
- A minimal and modern development setup powered by Bun

With this architecture, you get a foundation that is:
- Secure
- Scalable
- Extensible
- Pleasant to work with

Perfect for both personal projects and production systems.
