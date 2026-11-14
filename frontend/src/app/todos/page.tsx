"use client"; // runs in the browser: loads the user and handles logout

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

import { Navbar } from "@/components/layout/Navbar";
import { Card } from "@/components/ui/Card";
import { ApiError, getMe, logout } from "@/lib/api";
import type { User } from "@/types/auth";

// /todos: for now only shows who is logged in (the todo list comes in Step 8)
export default function TodosPage() {
  const router = useRouter();

  const [user, setUser] = useState<User | null>(null);
  const [error, setError] = useState<string | null>(null);

  // Load the logged-in user once; not logged in -> go to /login
  useEffect(() => {
    async function loadUser() {
      try {
        const me = await getMe();

        setUser(me);
      } catch (err) {
        if (err instanceof ApiError && err.status === 401) {
          router.replace("/login");

          return;
        }

        console.error(err);
        setError("Could not load your account. Please refresh the page.");
      }
    }

    loadUser();
  }, [router]);

  // Logout: clear the cookie on the server, then go to /login
  async function handleLogout() {
    try {
      await logout();
    } finally {
      router.replace("/login"); // leave the page even if the call failed
    }
  }

  // Error / loading: no navbar yet because we don't know the user
  if (error) {
    return <p className="p-6 text-sm text-red-400">{error}</p>;
  }

  if (!user) {
    return <p className="p-6 text-muted">Loading...</p>;
  }

  return (
    <>
      <Navbar user={user} onLogout={handleLogout} />

      <main className="mx-auto w-full max-w-3xl p-6">
        <Card>
          <h1 className="text-2xl font-semibold">My tasks</h1>
          <p className="mt-1 text-muted">Your todo list will appear here.</p>
        </Card>
      </main>
    </>
  );
}
