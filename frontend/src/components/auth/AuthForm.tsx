"use client"; // runs in the browser: uses state, events and the router

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState, type SubmitEvent } from "react";

import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Input } from "@/components/ui/Input";
import { errorMessage, login, register } from "@/lib/api";
import { credentialsSchema } from "@/lib/schemas";
import type { AuthFieldErrors, AuthMode } from "@/types/auth";

// Text that differs between the two pages
const TEXT = {
  login: {
    title: "Welcome back",
    subtitle: "Log in to see your tasks",
    submit: "Log in",
    submitting: "Logging in...",
    switchText: "Don't have an account?",
    switchLink: "Sign up",
    switchHref: "/register",
  },
  register: {
    title: "Create account",
    subtitle: "Start organizing your tasks",
    submit: "Sign up",
    submitting: "Creating account...",
    switchText: "Already have an account?",
    switchLink: "Log in",
    switchHref: "/login",
  },
} as const;

// AuthForm: the login and register pages share this form
export function AuthForm({ mode }: { mode: AuthMode }) {
  const router = useRouter();
  const text = TEXT[mode];

  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [fieldErrors, setFieldErrors] = useState<AuthFieldErrors>({});
  const [formError, setFormError] = useState<string | null>(null); // error from the API
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(event: SubmitEvent<HTMLFormElement>) {
    event.preventDefault(); // stop the browser from reloading the page
    setFormError(null);

    // 1. Validate in the browser first (same rules as the API)
    const parsed = credentialsSchema.safeParse({ email, password });

    if (!parsed.success) {
      const errors: AuthFieldErrors = {};

      for (const issue of parsed.error.issues) {
        const field = issue.path[0];

        if (field === "email" || field === "password") {
          errors[field] ??= issue.message; // keep the first message per field
        }
      }

      setFieldErrors(errors);

      return;
    }

    setFieldErrors({});

    // 2. Call the API
    setSubmitting(true);

    try {
      if (mode === "register") {
        await register(parsed.data);
      }

      // Log in right after register too, so the user doesn't type it twice
      await login(parsed.data);

      // 3. Logged in: go to the todo list
      router.replace("/todos");
    } catch (error) {
      setFormError(errorMessage(error));
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <Card className="w-full max-w-sm">
      <h1 className="text-2xl font-semibold">{text.title}</h1>
      <p className="mt-1 text-muted">{text.subtitle}</p>

      {/* noValidate: we show our own messages instead of the browser's popups */}
      <form onSubmit={handleSubmit} noValidate className="mt-6 flex flex-col gap-4">
        <Input
          label="Email"
          name="email"
          type="email"
          autoComplete="email"
          value={email}
          onChange={(event) => setEmail(event.target.value)}
          error={fieldErrors.email}
        />

        <Input
          label="Password"
          name="password"
          type="password"
          autoComplete={mode === "login" ? "current-password" : "new-password"}
          value={password}
          onChange={(event) => setPassword(event.target.value)}
          error={fieldErrors.password}
        />

        {/* Error from the API, e.g. wrong password / email taken */}
        {formError && (
          <p role="alert" className="text-sm text-red-400">
            {formError}
          </p>
        )}

        <Button type="submit" disabled={submitting} className="mt-2">
          {submitting ? text.submitting : text.submit}
        </Button>
      </form>

      {/* Link to the other page */}
      <p className="mt-6 text-center text-sm text-muted">
        {text.switchText}{" "}
        <Link href={text.switchHref} className="font-semibold text-accent">
          {text.switchLink}
        </Link>
      </p>
    </Card>
  );
}
