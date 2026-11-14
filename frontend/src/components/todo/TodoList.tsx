import { TodoItem, type TodoItemProps } from "@/components/todo/TodoItem";
import type { Todo } from "@/types/todo";

type TodoListProps = Omit<TodoItemProps, "todo"> & {
  title: string;
  todos: Todo[];
  emptyText: string;
};

// TodoList: a titled group of todos ("To do" / "Completed")
export function TodoList({ title, todos, emptyText, onToggle, onSave, onDelete }: TodoListProps) {
  return (
    <section>
      <h2 className="mb-3 text-lg font-semibold">
        {title} <span className="text-muted">({todos.length})</span>
      </h2>

      {todos.length === 0 ? (
        <p className="text-sm text-muted">{emptyText}</p>
      ) : (
        <ul className="flex flex-col gap-3">
          {todos.map((todo) => (
            <TodoItem key={todo.id} todo={todo} onToggle={onToggle} onSave={onSave} onDelete={onDelete} />
          ))}
        </ul>
      )}
    </section>
  );
}
