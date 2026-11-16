"use client"; // runs in the browser: edit / delete state and events

import { useState } from "react";

import { TodoForm } from "@/components/todo/TodoForm";
import { errorMessage } from "@/lib/api";
import { todayString } from "@/lib/date";
import { formatDueDate, isOverdue, toFormValues } from "@/lib/todo";
import type { Todo, TodoFormValues } from "@/types/todo";

export type TodoItemProps = {
  todo: Todo;
  onToggle: (todo: Todo) => Promise<void>;
  onSave: (id: string, values: TodoFormValues) => Promise<void>;
  onDelete: (id: string) => Promise<void>;
};

// Which view the item shows
type ItemMode = "view" | "edit" | "confirmDelete";

// TodoItem: one todo card with checkbox, edit and delete
export function TodoItem({ todo, onToggle, onSave, onDelete }: TodoItemProps) {
  const [mode, setMode] = useState<ItemMode>("view");
  const [busy, setBusy] = useState(false); // waiting for toggle / delete
  const [error, setError] = useState<string | null>(null);

  const overdue = isOverdue(todo, todayString());

  // runAction: call toggle / delete, show the error on this card if it fails
  async function runAction(action: () => Promise<void>) {
    setBusy(true);
    setError(null);

    try {
      await action();
    } catch (err) {
      setError(errorMessage(err));
    } finally {
      setBusy(false);
    }
  }

  // Save from the edit form; errors are shown inside the form (TodoForm catches them)
  async function handleSave(values: TodoFormValues) {
    await onSave(todo.id, values);

    setMode("view");
  }

  // ---------- Edit mode: the form replaces the card content ----------
  if (mode === "edit") {
    return (
      <li className="rounded-2xl bg-surface p-4">
        <TodoForm
          initialValues={toFormValues(todo)}
          submitLabel="Save"
          submittingLabel="Saving..."
          onSubmit={handleSave}
          onCancel={() => setMode("view")}
        />
      </li>
    );
  }

  // ---------- View mode ----------
  return (
    <li className="rounded-2xl bg-surface-raised p-4">
      <div className="flex items-start gap-3">
        {/* Done checkbox */}
        <input
          type="checkbox"
          checked={todo.done}
          disabled={busy}
          onChange={() => runAction(() => onToggle(todo))}
          aria-label={todo.done ? `Mark "${todo.title}" as not done` : `Mark "${todo.title}" as done`}
          className="mt-1 h-5 w-5 shrink-0 cursor-pointer accent-accent"
        />

        {/* Text */}
        <div className="min-w-0 flex-1">
          <p className={`font-medium break-words ${todo.done ? "text-muted line-through" : ""}`}>{todo.title}</p>

          {todo.description && (
            <p className="mt-1 text-sm break-words whitespace-pre-line text-muted">{todo.description}</p>
          )}

          {/* Due date: red + "Overdue" when it has passed and the task isn't done */}
          {todo.due_date && (
            <p className={`mt-2 text-xs font-semibold ${overdue ? "text-red-400" : "text-accent"}`}>
              {overdue ? "Overdue · " : "Due "}
              {formatDueDate(todo.due_date)}
            </p>
          )}
        </div>
      </div>

      {/* Actions */}
      {mode === "confirmDelete" ? (
        <div className="mt-4 flex items-center justify-end gap-3 text-sm">
          <span className="text-muted">Delete this task?</span>

          <button type="button" onClick={() => setMode("view")} disabled={busy} className="font-semibold">
            Cancel
          </button>

          <button
            type="button"
            onClick={() => runAction(() => onDelete(todo.id))}
            disabled={busy}
            className="font-semibold text-red-400"
          >
            {busy ? "Deleting..." : "Delete"}
          </button>
        </div>
      ) : (
        <div className="mt-4 flex justify-end gap-4 text-sm">
          {/* Completed todos can't be edited (uncheck first to edit) */}
          {!todo.done && (
            <button
              type="button"
              onClick={() => setMode("edit")}
              disabled={busy}
              className="font-semibold text-muted hover:text-foreground"
            >
              Edit
            </button>
          )}

          <button
            type="button"
            onClick={() => setMode("confirmDelete")}
            disabled={busy}
            className="font-semibold text-muted hover:text-red-400"
          >
            Delete
          </button>
        </div>
      )}

      {/* Error from toggle / delete */}
      {error && (
        <p role="alert" className="mt-2 text-sm text-red-400">
          {error}
        </p>
      )}
    </li>
  );
}
