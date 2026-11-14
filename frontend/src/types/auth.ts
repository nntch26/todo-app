import type { z } from "zod";

import type { credentialsSchema, userSchema } from "@/lib/schemas";


// User: { id, email } from the API
export type User = z.infer<typeof userSchema>;

// Credentials: { email, password } from the login / register form
export type Credentials = z.infer<typeof credentialsSchema>;

// AuthMode: which page the AuthForm is on
export type AuthMode = "login" | "register";

export type AuthFieldErrors = {
  email?: string;
  password?: string;
};
