import { ChevronLeft, ChevronRight } from "lucide-react";

import { Card } from "@/components/ui/Card";
import { formatMonthTitle, monthGrid } from "@/lib/date";
import type { CalendarMonth } from "@/types/calendar";
import type { Todo } from "@/types/todo";

const WEEKDAYS = ["Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"];

type MonthCalendarProps = {
  month: CalendarMonth;
  todos: Todo[]; // todos due in this month (for the dots)
  today: string; // "2026-02-04"
  selectedDate: string | null; // null = whole month
  onSelectDate: (date: string | null) => void;
  onChangeMonth: (amount: number) => void; // -1 = previous, +1 = next
};

// MonthCalendar: month grid with a dot on days that have tasks
export function MonthCalendar({ month, todos, today, selectedDate, onSelectDate, onChangeMonth }: MonthCalendarProps) {
  // Which days have open / only finished tasks -> orange / gray dot
  const openDays = new Set<string>();
  const doneDays = new Set<string>();

  for (const todo of todos) {
    if (todo.due_date === null) {
      continue;
    }

    if (todo.done) {
      doneDays.add(todo.due_date);
    } else {
      openDays.add(todo.due_date);
    }
  }

  return (
    <Card>
      {/* Header: ‹ February 2026 › */}
      <div className="mb-4 flex items-center justify-between">
        <button
          type="button"
          onClick={() => onChangeMonth(-1)}
          aria-label="Previous month"
          className="rounded-full p-2 text-muted hover:bg-surface-raised hover:text-foreground"
        >
          <ChevronLeft className="h-5 w-5" aria-hidden="true" />
        </button>

        <h2 className="text-lg font-semibold">{formatMonthTitle(month)}</h2>

        <button
          type="button"
          onClick={() => onChangeMonth(1)}
          aria-label="Next month"
          className="rounded-full p-2 text-muted hover:bg-surface-raised hover:text-foreground"
        >
          <ChevronRight className="h-5 w-5" aria-hidden="true" />
        </button>
      </div>

      {/* Weekday names */}
      <div className="grid grid-cols-7 text-center text-xs text-muted">
        {WEEKDAYS.map((weekday) => (
          <span key={weekday} className="py-2">
            {weekday}
          </span>
        ))}
      </div>

      {/* Days */}
      <div className="grid grid-cols-7 gap-y-1 text-center">
        {monthGrid(month).map((cell, index) => {
          // Empty cell before the 1st
          if (cell === null) {
            return <span key={`empty-${index}`} />;
          }

          const selected = cell.date === selectedDate;
          const isToday = cell.date === today;
          const hasOpen = openDays.has(cell.date);
          const hasDoneOnly = !hasOpen && doneDays.has(cell.date);

          // Day style: selected (orange fill) > today (orange ring) > normal
          let dayClass = "hover:bg-surface-raised";

          if (selected) {
            dayClass = "bg-accent text-white";
          } else if (isToday) {
            dayClass = "ring-1 ring-accent text-accent";
          }

          return (
            <button
              key={cell.date}
              type="button"
              // Click the selected day again = back to the whole month
              onClick={() => onSelectDate(selected ? null : cell.date)}
              aria-pressed={selected}
              aria-label={`${cell.date}${hasOpen ? ", has tasks" : ""}`}
              className="flex flex-col items-center py-1"
            >
              <span className={`flex h-9 w-9 items-center justify-center rounded-full text-sm ${dayClass}`}>
                {cell.day}
              </span>

              {/* Dot: orange = open tasks, gray = all done; empty space keeps rows aligned */}
              <span
                className={`mt-0.5 h-1.5 w-1.5 rounded-full ${
                  hasOpen ? "bg-accent" : hasDoneOnly ? "bg-muted" : "bg-transparent"
                }`}
                aria-hidden="true"
              />
            </button>
          );
        })}
      </div>
    </Card>
  );
}
