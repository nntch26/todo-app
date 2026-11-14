import { z } from "zod";


// ---------- API responses ----------

// User: response of /api/auth/register, /login, /me
export const userSchema = z.object({
  id: z.string(),
  email: z.string(),
});

// Todo: one item from /api/todos (backend/internal/handler/todo.go todoResponse)
export const todoSchema = z.object({
  id: z.string(),
  title: z.string(),
  description: z.string(),
  due_date: z.iso.date().nullable(), // "2026-12-31" or null
  done: z.boolean(),
  created_at: z.string(),
  updated_at: z.string(),
});

// GET /api/todos returns a list
export const todoListSchema = z.array(todoSchema);

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

// Create / edit todo form; same limits as the Go API (createTodoRequest)
export const todoFormSchema = z.object({
  title: z.string().trim().min(1, "Title is required").max(200, "Title must be at most 200 characters"),
  description: z.string().trim().max(2000, "Description must be at most 2000 characters"),
  dueDate: z.iso.date("Please pick a valid date").or(z.literal("")),
});
