export default function Loading() {
  return (
    <main className="flex min-h-screen items-center justify-center bg-slate-50 px-6">
      <div className="text-center">
        <div className="mx-auto h-8 w-8 animate-spin rounded-full border-4 border-slate-200 border-t-indigo-600" />
        <p className="mt-4 text-sm font-medium text-slate-600">
          Connecting to Phoenix API…
        </p>
      </div>
    </main>
  );
}
