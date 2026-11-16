import type { CalendarDay, CalendarMonth } from "@/types/calendar";

// All dates here use the user's LOCAL time.
// (toISOString() is UTC: in Thailand between 00:00 and 07:00 it would still say "yesterday")

// toDateString: Date -> "2026-02-04" (same format as todo.due_date)
export function toDateString(date: Date): string {
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0"); // getMonth() is 0–11
  const day = String(date.getDate()).padStart(2, "0");

  return `${year}-${month}-${day}`;
}

// todayString: today's date as "2026-02-04"
export function todayString(): string {
  return toDateString(new Date());
}

// currentMonth: the month we are in now
export function currentMonth(): CalendarMonth {
  const now = new Date();

  return { year: now.getFullYear(), month: now.getMonth() };
}

// addMonths: { 2026, 0 } + -1 -> { 2025, 11 } (Date handles the year change)
export function addMonths(current: CalendarMonth, amount: number): CalendarMonth {
  const moved = new Date(current.year, current.month + amount, 1);

  return { year: moved.getFullYear(), month: moved.getMonth() };
}

// monthPrefix: "2026-02-" — every due_date in that month starts with this
export function monthPrefix(current: CalendarMonth): string {
  return `${current.year}-${String(current.month + 1).padStart(2, "0")}-`;
}

// formatMonthTitle: { 2026, 1 } -> "February 2026"
export function formatMonthTitle(current: CalendarMonth): string {
  return new Date(current.year, current.month, 1).toLocaleDateString("en-GB", {
    month: "long",
    year: "numeric",
  });
}

// monthGrid: the cells of a month calendar, weeks start on Monday.
// Empty cells (null) come first so the 1st lands under the right weekday.
export function monthGrid(current: CalendarMonth): CalendarDay[] {
  const firstDay = new Date(current.year, current.month, 1);
  const daysInMonth = new Date(current.year, current.month + 1, 0).getDate(); // day 0 of next month = last day of this month

  // getDay(): 0 = Sunday ... 6 = Saturday -> shift so Monday = 0
  const emptyCells = (firstDay.getDay() + 6) % 7;

  const cells: CalendarDay[] = [];

  for (let i = 0; i < emptyCells; i++) {
    cells.push(null);
  }

  for (let day = 1; day <= daysInMonth; day++) {
    const date = toDateString(new Date(current.year, current.month, day));

    cells.push({ date, day });
  }

  return cells;
}

// greeting: text for the time of day
export function greeting(now: Date): string {
  const hour = now.getHours();

  if (hour < 12) {
    return "Good Morning.";
  }

  if (hour < 18) {
    return "Good Afternoon.";
  }

  return "Good Evening.";
}
