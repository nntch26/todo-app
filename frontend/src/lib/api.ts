import type { z } from "zod";

import { apiErrorSchema, userSchema } from "@/lib/schemas";
import type { Credentials, User } from "@/types/auth";

// ApiError: a failed API call; message is safe to show to the user
export class ApiError extends Error {
  readonly status: number;

  constructor(status: number, message: string) {
    super(message);
    this.status = status;
  }
}

// ---------- Helpers ----------

// send makes the request; throws ApiError if the network fails or the status is not 2xx
async function send(path: string, init: RequestInit = {}): Promise<Response> {
  let response: Response;

  try {
    response = await fetch(path, {
      ...init,
      headers: { "Content-Type": "application/json" },
    });
  } catch {
    // fetch only throws when the server can't be reached at all
    throw new ApiError(0, "Cannot reach the server. Please try again.");
  }

  if (!response.ok) {
    const message = await readErrorMessage(response);

    throw new ApiError(response.status, message);
  }

  return response;
}

// readErrorMessage gets "..." from {"error": "..."}, or a generic message if the body isn't that shape
async function readErrorMessage(response: Response): Promise<string> {
  const body: unknown = await response.json().catch(() => null); // body may not be JSON (e.g. proxy error page)
  const parsed = apiErrorSchema.safeParse(body);

  if (parsed.success) {
    return parsed.data.error;
  }

  return "Something went wrong. Please try again.";
}

// readBody parses the JSON body and checks it matches the schema (validate at the boundary)
async function readBody<T>(response: Response, schema: z.ZodType<T>): Promise<T> {
  const body: unknown = await response.json();
  const parsed = schema.safeParse(body);

  if (!parsed.success) {
    throw new Error(`Unexpected response from ${response.url}: ${parsed.error.message}`);
  }

  return parsed.data;
}

// ---------- Auth ----------

// register POST /api/auth/register (does not log in)
export async function register(credentials: Credentials): Promise<User> {
  const response = await send("/api/auth/register", {
    method: "POST",
    body: JSON.stringify(credentials),
  });

  return readBody(response, userSchema);
}

// login POST /api/auth/login; the server sets the token cookie
export async function login(credentials: Credentials): Promise<User> {
  const response = await send("/api/auth/login", {
    method: "POST",
    body: JSON.stringify(credentials),
  });

  return readBody(response, userSchema);
}

// logout POST /api/auth/logout; the server clears the cookie
export async function logout(): Promise<void> {
  await send("/api/auth/logout", { method: "POST" });
}

// getMe GET /api/auth/me; throws ApiError 401 if not logged in
export async function getMe(): Promise<User> {
  const response = await send("/api/auth/me");

  return readBody(response, userSchema);
}
