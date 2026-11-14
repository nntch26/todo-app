import { z } from "zod";


// ---------- API responses ----------

// User: response of /api/auth/register, /login, /me
export const userSchema = z.object({
  id: z.string(),
  email: z.string(),
});

// Error body: every failed API call returns {"error": "..."}
export const apiErrorSchema = z.object({
  error: z.string(),
});

// ---------- Forms ----------

// Login / register form; same rules as the Go API (authRequest)
export const credentialsSchema = z.object({
  email: z.email("Please enter a valid email"),
  password: z
    .string()
    .min(8, "Password must be at least 8 characters")
    .max(72, "Password must be at most 72 characters"),
});
