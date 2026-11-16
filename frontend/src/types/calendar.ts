// TodoView: which tab of the todos page is open
export type TodoView = "today" | "month";

// CalendarMonth: one month; month is 0–11 like JavaScript's Date (0 = January)
export type CalendarMonth = {
  year: number;
  month: number;
};

// CalendarDay: one cell in the month grid; null = empty cell before the 1st
export type CalendarDay = {
  date: string; // "2026-02-04"
  day: number; // 4
} | null;
