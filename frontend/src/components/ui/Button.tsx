import type { ComponentProps } from "react";

// Button: the main orange rounded button.
export function Button({ className = "", type = "button", ...props }: ComponentProps<"button">) {
  return (
    <button
      type={type} 
      className={`w-full rounded-full bg-accent py-3 font-semibold text-white transition-colors hover:bg-accent-hover disabled:opacity-50 ${className}`}
      {...props}
    />
  );
}
