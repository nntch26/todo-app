"use client"; // runs in the browser: loads data, handles create / edit / delete

import { Plus } from "lucide-react";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

import { Navbar } from "@/components/layout/Navbar";
import { TodoForm } from "@/components/todo/TodoForm";
import { TodoList } from "@/components/todo/TodoList";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { ApiError, createTodo, deleteTodo, getMe, listTodos, logout, updateTodo } from "@/lib/api";
import { EMPTY_TODO_FORM, toTodoInput } from "@/lib/todo";
import type { User } from "@/types/auth";
import type { Todo, TodoFormValues } from "@/types/todo";

// /todos: the user's todo list
export default function TodosPage() {
  const router = useRouter();

  const [user, setUser] = useState<User | null>(null);
  const [todos, setTodos] = useState<Todo[]>([]);
  const [loadError, setLoadError] = useState<string | null>(null);
  const [creating, setCreating] = useState(false); // is the "New task" form open?

  // ---------- Load the user + todos once; not logged in -> /login ----------
  useEffect(() => {
    async function loadPage() {
      try {
        // Both requests at the same time
        const [me, myTodos] = await Promise.all([getMe(), listTodos()]);

        setUser(me);
        setTodos(myTodos);
      } catch (err) {
        if (err instanceof ApiError && err.status === 401) {
          router.replace("/login");

          return;
        }

        console.error(err);
        setLoadError("Could not load your tasks. Please refresh the page.");
      }
    }

    loadPage();
  }, [router]);

  // ---------- Actions (errors are thrown to the component that called them, which shows them) ----------

  // Create: call API, put the new todo at the top, then close the form
  async function handleCreate(values: TodoFormValues) {
    const created = await createTodo(toTodoInput(values));

    setTodos((previous) => [created, ...previous]);
    setCreating(false);
  }

  // Edit: call API, then replace the old todo with the saved one
  async function handleSave(id: string, values: TodoFormValues) {
    const updated = await updateTodo(id, toTodoInput(values));

    setTodos((previous) => previous.map((todo) => (todo.id === id ? updated : todo)));
  }

  // Check / uncheck done
  async function handleToggle(todo: Todo) {
    const updated = await updateTodo(todo.id, { done: !todo.done });

    setTodos((previous) => previous.map((item) => (item.id === todo.id ? updated : item)));
  }

  // Delete: call API, then remove it from the list
  async function handleDelete(id: string) {
    await deleteTodo(id);

    setTodos((previous) => previous.filter((todo) => todo.id !== id));
  }

  // Logout: clear the cookie on the server, then go to /login
  async function handleLogout() {
    try {
      await logout();
    } finally {
      router.replace("/login"); // leave the page even if the call failed
    }
  }

  // ---------- Render ----------

  // Error / loading: no navbar yet because we don't know the user
  if (loadError) {
    return <p className="p-6 text-sm text-red-400">{loadError}</p>;
  }

  if (!user) {
    return <p className="p-6 text-muted">Loading...</p>;
  }

  const openTodos = todos.filter((todo) => !todo.done);
  const doneTodos = todos.filter((todo) => todo.done);

  return (
    <>
      <Navbar user={user} onLogout={handleLogout} />

      <main className="mx-auto flex w-full max-w-3xl flex-col gap-8 p-6">
        {/* Header: title on the left, "New task" button on the right */}
        <div className="flex items-end justify-between gap-4">
          <div>
            <h1 className="text-2xl font-semibold">My tasks</h1>
            <p className="mt-1 text-muted">
              {openTodos.length} to do · {doneTodos.length} completed
            </p>
          </div>

          {/* Hidden while the form is open */}
          {!creating && (
            <Button size="sm" onClick={() => setCreating(true)}>
              <Plus className="h-4 w-4" aria-hidden="true" />
              New task
            </Button>
          )}
        </div>

        {/* Create form: only after clicking "New task" */}
        {creating && (
          <Card>
            <h2 className="mb-4 text-lg font-semibold">New task</h2>

            <TodoForm
              initialValues={EMPTY_TODO_FORM}
              submitLabel="Add task"
              submittingLabel="Adding..."
              onSubmit={handleCreate}
              onCancel={() => setCreating(false)}
            />
          </Card>
        )}

        {/* Lists */}
        <TodoList
          title="To do"
          todos={openTodos}
          emptyText="Nothing to do. Add a task above."
          onToggle={handleToggle}
          onSave={handleSave}
          onDelete={handleDelete}
        />

        <TodoList
          title="Completed"
          todos={doneTodos}
          emptyText="No completed tasks yet."
          onToggle={handleToggle}
          onSave={handleSave}
          onDelete={handleDelete}
        />
      </main>
    </>
  );
}
