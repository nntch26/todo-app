"use client"; // runs in the browser: form state and events

import { useState, type SubmitEvent } from "react";

import { Button } from "@/components/ui/Button";
import { Input } from "@/components/ui/Input";
import { Textarea } from "@/components/ui/Textarea";
import { errorMessage } from "@/lib/api";
import { todoFormSchema } from "@/lib/schemas";
import type { TodoFormErrors, TodoFormValues } from "@/types/todo";

type TodoFormProps = {
  initialValues: TodoFormValues;
  submitLabel: string;
  submittingLabel: string;
  onSubmit: (values: TodoFormValues) => Promise<void>; // throws on API error
  onCancel?: () => void; // shows a Cancel button when given (edit mode)
};

// TodoForm: title + description + due date; used to create and to edit a todo
export function TodoForm({ initialValues, submitLabel, submittingLabel, onSubmit, onCancel }: TodoFormProps) {
  const [values, setValues] = useState<TodoFormValues>(initialValues);
  const [fieldErrors, setFieldErrors] = useState<TodoFormErrors>({});
  const [formError, setFormError] = useState<string | null>(null); // error from the API
  const [submitting, setSubmitting] = useState(false);

  // Update one field, keep the others
  function setField(field: keyof TodoFormValues, value: string) {
    setValues((previous) => ({ ...previous, [field]: value }));
  }

  async function handleSubmit(event: SubmitEvent<HTMLFormElement>) {
    event.preventDefault(); // stop the browser from reloading the page
    setFormError(null);

    // 1. Validate in the browser first (also trims spaces)
    const parsed = todoFormSchema.safeParse(values);

    if (!parsed.success) {
      const errors: TodoFormErrors = {};

      for (const issue of parsed.error.issues) {
        const field = issue.path[0];

        if (field === "title" || field === "description" || field === "dueDate") {
          errors[field] ??= issue.message; // keep the first message per field
        }
      }

      setFieldErrors(errors);

      return;
    }

    setFieldErrors({});

    // 2. Let the parent call the API
    setSubmitting(true);

    try {
      await onSubmit(parsed.data);

      // 3. Saved: clear the form (create mode); in edit mode the parent closes it
      setValues(initialValues);
    } catch (error) {
      setFormError(errorMessage(error));
    } finally {
      setSubmitting(false);
    }
  }

  return (
    // noValidate: we show our own messages instead of the browser's popups
    <form onSubmit={handleSubmit} noValidate className="flex flex-col gap-3">
      <Input
        label="Title"
        name="title"
        autoFocus // the form opens on a click, so jump straight to typing
        value={values.title}
        onChange={(event) => setField("title", event.target.value)}
        error={fieldErrors.title}
      />

      <Textarea
        label="Description"
        name="description"
        rows={2}
        value={values.description}
        onChange={(event) => setField("description", event.target.value)}
        error={fieldErrors.description}
      />

      {/* Bottom row: due date on the left, buttons on the right (wraps on small screens) */}
      <div className="flex flex-wrap items-end justify-between gap-3">
        <Input
          label="Due date"
          name="dueDate"
          type="date"
          value={values.dueDate}
          onChange={(event) => setField("dueDate", event.target.value)}
          error={fieldErrors.dueDate}
          className="[color-scheme:dark]" // dark calendar popup
        />

        {/* Main action on the right, as most forms do */}
        <div className="flex gap-2">
          {onCancel && (
            <Button variant="secondary" size="sm" onClick={onCancel} disabled={submitting}>
              Cancel
            </Button>
          )}

          <Button type="submit" size="sm" disabled={submitting}>
            {submitting ? submittingLabel : submitLabel}
          </Button>
        </div>
      </div>

      {/* Error from the API */}
      {formError && (
        <p role="alert" className="text-sm text-red-400">
          {formError}
        </p>
      )}
    </form>
  );
}
