import type { Todo, TodoFormValues, TodoInput } from "@/types/todo";

// Empty create form
export const EMPTY_TODO_FORM: TodoFormValues = {
  title: "",
  description: "",
  dueDate: "",
};

// toFormValues: todo from the API -> values for the edit form
export function toFormValues(todo: Todo): TodoFormValues {
  return {
    title: todo.title,
    description: todo.description,
    dueDate: todo.due_date ?? "", // null -> "" (empty date input)
  };
}

// toTodoInput: form values -> body for the API
export function toTodoInput(values: TodoFormValues): TodoInput {
  return {
    title: values.title,
    description: values.description,
    due_date: values.dueDate === "" ? null : values.dueDate, // "" -> null (no due date)
  };
}

// formatDueDate: "2026-12-31" -> "31 Dec 2026"
const dueDateFormat = new Intl.DateTimeFormat("en-GB", {
  day: "numeric",
  month: "short",
  year: "numeric",
  timeZone: "UTC", // the date has no time; UTC on both sides so it never shifts a day
});

export function formatDueDate(date: string): string {
  return dueDateFormat.format(new Date(`${date}T00:00:00Z`));
}
