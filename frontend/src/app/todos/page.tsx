"use client"; 

import { Plus } from "lucide-react";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";

import { Navbar } from "@/components/layout/Navbar";
import { MonthCalendar } from "@/components/todo/MonthCalendar";
import { TodayHeader } from "@/components/todo/TodayHeader";
import { TodoForm } from "@/components/todo/TodoForm";
import { TodoList } from "@/components/todo/TodoList";
import { ViewSwitch } from "@/components/todo/ViewSwitch";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { ApiError, createTodo, deleteTodo, getMe, listTodos, logout, updateTodo } from "@/lib/api";
import { addMonths, currentMonth, formatMonthTitle, monthPrefix, todayString } from "@/lib/date";
import { EMPTY_TODO_FORM, formatDueDate, isDueOn, isForToday, toTodoInput } from "@/lib/todo";
import type { User } from "@/types/auth";
import type { CalendarMonth, TodoView } from "@/types/calendar";
import type { Todo, TodoFormValues } from "@/types/todo";

// /todos: the user's todo list, as "Today" or "Month" (calendar)
export default function TodosPage() {
  const router = useRouter();

  const [user, setUser] = useState<User | null>(null);
  const [todos, setTodos] = useState<Todo[]>([]);
  const [loadError, setLoadError] = useState<string | null>(null);
  const [creating, setCreating] = useState(false); // is the "New task" form open?

  // Month tab
  const [view, setView] = useState<TodoView>("today");
  // How many months away from the current month (0 = this month, -1 = last month).
  // Not the month itself: reading the clock in useState runs during Next's prerender, which Next rejects
  const [monthOffset, setMonthOffset] = useState(0);
  const [selectedDate, setSelectedDate] = useState<string | null>(null); // null = whole month

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

  // ---------- Which todos to show ----------
  // Safe to read the clock here: this code only runs in the browser after the user has loaded
  const today = todayString();
  const month: CalendarMonth = addMonths(currentMonth(), monthOffset);

  // Month tab: todos due in the shown month, then (optionally) on the clicked day
  const monthTodos = todos.filter((todo) => isDueOn(todo, monthPrefix(month)));

  let visibleTodos: Todo[];

  if (view === "today") {
    visibleTodos = todos.filter((todo) => isForToday(todo, today));
  } else if (selectedDate !== null) {
    visibleTodos = monthTodos.filter((todo) => isDueOn(todo, selectedDate));
  } else {
    visibleTodos = monthTodos;
  }

  const openTodos = visibleTodos.filter((todo) => !todo.done);
  const doneTodos = visibleTodos.filter((todo) => todo.done);

  // Month tab: list title = clicked day or the month name
  const monthListTitle = selectedDate !== null ? formatDueDate(selectedDate) : formatMonthTitle(month);

  // Change month with ‹ ›; clear the clicked day (it belongs to the old month)
  function handleChangeMonth(amount: number) {
    setMonthOffset((previous) => previous + amount);
    setSelectedDate(null);
  }

  return (
    <>
      <Navbar user={user} onLogout={handleLogout} />

      <main className="mx-auto flex w-full max-w-3xl flex-col gap-8 p-6">
        <ViewSwitch view={view} onChange={setView} />

        {/* Header: date + greeting (Today) or month title (Month), "New task" button on the right */}
        <div className="flex items-center justify-between gap-4">
          {view === "today" ? (
            <TodayHeader now={new Date()} />
          ) : (
            <div>
              <h1 className="text-2xl font-semibold">{monthListTitle}</h1>

              {/* Only when a day is clicked: go back to the whole month */}
              {selectedDate !== null && (
                <button
                  type="button"
                  onClick={() => setSelectedDate(null)}
                  className="mt-1 text-sm font-semibold text-accent"
                >
                  Show whole month
                </button>
              )}
            </div>
          )}

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

        {/* Month tab: calendar above the lists */}
        {view === "month" && (
          <MonthCalendar
            month={month}
            todos={monthTodos}
            today={today}
            selectedDate={selectedDate}
            onSelectDate={setSelectedDate}
            onChangeMonth={handleChangeMonth}
          />
        )}

        {/* Lists */}
        <TodoList
          title="To do"
          todos={openTodos}
          emptyText={view === "today" ? "Nothing to do today. Add a task above." : "No tasks due in this period."}
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
