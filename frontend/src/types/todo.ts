import type { z } from "zod";

import type { todoFormSchema, todoSchema } from "@/lib/schemas";

// Todo: one item from the API
export type Todo = z.infer<typeof todoSchema>;

// TodoFormValues: what the create / edit form holds (dueDate "" = none)
export type TodoFormValues = z.infer<typeof todoFormSchema>;


export type TodoFormErrors = Partial<Record<keyof TodoFormValues, string>>;


// TodoInput: body of POST /api/todos
export type TodoInput = {
  title: string;
  description: string;
  due_date: string | null;
};

// TodoPatch: body of PATCH /api/todos/:id; missing field = keep, due_date null = remove
export type TodoPatch = Partial<TodoInput & { done: boolean }>;
