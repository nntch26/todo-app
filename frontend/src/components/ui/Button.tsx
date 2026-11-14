import type { ComponentProps } from "react";

type ButtonProps = ComponentProps<"button"> & {
  variant?: "primary" | "secondary"; // primary = orange (main action), secondary = outline (cancel / other)
  size?: "md" | "sm"; // md = full width (auth forms), sm = compact (inline forms, navbar)
};

const VARIANT_CLASSES = {
  primary: "bg-accent text-white hover:bg-accent-hover",
  secondary: "border border-border text-foreground hover:border-accent hover:text-accent",
};

const SIZE_CLASSES = {
  md: "w-full py-3",
  sm: "px-4 py-2 text-sm",
};

// Button: rounded button in the app's style
export function Button({ variant = "primary", size = "md", className = "", type = "button", ...props }: ButtonProps) {
  return (
    <button
      type={type}
      // inline-flex + gap: an icon inside the button sits next to the text, vertically centered
      className={`inline-flex items-center justify-center gap-1.5 rounded-full font-semibold transition-colors disabled:opacity-50 ${VARIANT_CLASSES[variant]} ${SIZE_CLASSES[size]} ${className}`}
      {...props}
    />
  );
}
