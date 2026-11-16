import { greeting } from "@/lib/date";

// TodayHeader: big date on the left, greeting on the right ("Feb 4 | Good Afternoon.")
export function TodayHeader({ now }: { now: Date }) {
  const month = now.toLocaleDateString("en-GB", { month: "short" }); // "Feb"
  const day = now.getDate(); // 4

  return (
    <div className="flex items-center gap-4">
      {/* Big date */}
      <div className="flex flex-col items-center leading-none">
        <span className="text-sm font-semibold text-muted uppercase">{month}</span>
        <span className="text-4xl font-bold text-accent">{day}</span>
      </div>

      {/* Divider */}
      <div className="h-10 w-px bg-border" aria-hidden="true" />

      <div>
        <h1 className="text-2xl font-semibold">{greeting(now)}</h1>
        <p className="text-muted">What&apos;s your plan for today?</p>
      </div>
    </div>
  );
}
