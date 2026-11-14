import type { ComponentProps } from "react";

export function Card({ className = "", ...props }: ComponentProps<"div">) {
  return <div className={`rounded-3xl bg-surface p-6 ${className}`} {...props} />;
}
