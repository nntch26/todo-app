import type { ComponentProps } from "react";

type TextareaProps = ComponentProps<"textarea"> & {
  label: string;
  error?: string; // shown in red under the box
};

// Textarea: label + multi-line text box + optional error message (same look as Input)
export function Textarea({ label, error, id, className = "", ...props }: TextareaProps) {
  // Link the label and error to the textarea for screen readers
  const textareaId = id ?? props.name;
  const errorId = `${textareaId}-error`;

  return (
    <div className="flex flex-col gap-2">
      <label htmlFor={textareaId} className="text-sm text-muted">
        {label}
      </label>

      <textarea
        id={textareaId}
        aria-invalid={error ? true : undefined}
        aria-describedby={error ? errorId : undefined}
        className={`resize-y rounded-2xl border bg-surface-raised px-4 py-3 outline-none focus:border-accent ${
          error ? "border-red-400" : "border-transparent"
        } ${className}`}
        {...props}
      />

      {error && (
        <p id={errorId} className="text-sm text-red-400">
          {error}
        </p>
      )}
    </div>
  );
}
