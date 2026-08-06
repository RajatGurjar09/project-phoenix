"use client";

export default function Error({ reset }: { reset: () => void }) {
  return (
    <main className="flex min-h-screen items-center justify-center bg-slate-50 px-6">
      <div className="max-w-md rounded-2xl border border-rose-200 bg-white p-8 text-center shadow-sm">
        <p className="text-sm font-semibold uppercase tracking-[0.2em] text-rose-600">
          Connection error
        </p>
        <h1 className="mt-3 text-2xl font-semibold text-slate-950">
          Phoenix Console is unavailable
        </h1>
        <p className="mt-3 text-sm leading-6 text-slate-600">
          We couldn’t retrieve dashboard data. Check the backend service and try again.
        </p>
        <button
          type="button"
          onClick={() => reset()}
          className="mt-6 rounded-lg bg-indigo-600 px-4 py-2 text-sm font-semibold text-white transition hover:bg-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-600 focus:ring-offset-2"
        >
          Try again
        </button>
      </div>
    </main>
  );
}
