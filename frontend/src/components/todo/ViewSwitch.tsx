import type { TodoView } from "@/types/calendar";

const OPTIONS: { value: TodoView; label: string }[] = [
  { value: "today", label: "Today" },
  { value: "month", label: "Month" },
];

type ViewSwitchProps = {
  view: TodoView;
  onChange: (view: TodoView) => void;
};

// ViewSwitch: pill-shaped toggle between the Today and Month tabs
export function ViewSwitch({ view, onChange }: ViewSwitchProps) {
  return (
    // role="tablist" + aria-selected: screen readers announce it as tabs
    <div role="tablist" className="flex rounded-full bg-surface p-1">
      {OPTIONS.map((option) => {
        const selected = option.value === view;

        return (
          <button
            key={option.value}
            type="button"
            role="tab"
            aria-selected={selected}
            onClick={() => onChange(option.value)}
            className={`flex-1 rounded-full py-2 text-sm font-semibold transition-colors ${
              selected ? "bg-accent text-white" : "text-muted hover:text-foreground"
            }`}
          >
            {option.label}
          </button>
        );
      })}
    </div>
  );
}
