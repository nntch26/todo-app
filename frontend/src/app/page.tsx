export default function Home() {
  return (
    <main className="flex flex-1 items-center justify-center p-6">
      {/* Card */}
      <div className="w-full max-w-sm rounded-3xl bg-surface p-6">
        <h1 className="text-2xl font-semibold">Todo</h1>
        <p className="mt-1 text-muted">Theme preview</p>

        {/* Item inside the card */}
        <div className="mt-6 rounded-2xl bg-surface-raised p-4">
          <p className="font-medium">Buy milk</p>
          <p className="text-sm text-muted">Due 31 Dec 2026</p>
        </div>

        {/* Main orange button */}
        <button
          type="button"
          className="mt-6 w-full rounded-full bg-accent py-3 font-semibold text-white transition-colors hover:bg-accent-hover"
        >
          Create task
        </button>
      </div>
    </main>
  );
}
