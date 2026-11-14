import type { ComponentProps } from "react";

type InputProps = ComponentProps<"input"> & {
  label: string;
  error?: string; // shown in red under the input
};

// Input: label + text box + optional error message
export function Input({ label, error, id, className = "", ...props }: InputProps) {
  // Link the label and error to the input, so screen readers read them together
  const inputId = id ?? props.name;
  const errorId = `${inputId}-error`;

  return (
    <div className="flex flex-col gap-2">
      <label htmlFor={inputId} className="text-sm text-muted">
        {label}
      </label>

      <input
        id={inputId}
        aria-invalid={error ? true : undefined}
        aria-describedby={error ? errorId : undefined}
        className={`rounded-2xl border bg-surface-raised px-4 py-3 outline-none focus:border-accent ${
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
