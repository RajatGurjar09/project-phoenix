type StatusCardProps = {
  label: string;
  value: string;
  detail?: string;
  tone?: "success" | "neutral";
  badge?: string;
};

export function StatusCard({
  label,
  value,
  detail,
  tone = "neutral",
  badge,
}: StatusCardProps) {
  const badgeClass =
    tone === "success"
      ? "bg-emerald-50 text-emerald-700 ring-emerald-600/20"
      : "bg-slate-100 text-slate-600 ring-slate-500/10";

  return (
    <article className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
      <div className="flex items-start justify-between gap-4">
        <p className="text-sm font-medium text-slate-500">{label}</p>
        {badge ? (
          <span
            className={`inline-flex items-center gap-1.5 rounded-full px-2.5 py-1 text-xs font-semibold ring-1 ring-inset ${badgeClass}`}
          >
            <span className="h-1.5 w-1.5 rounded-full bg-current" />
            {badge}
          </span>
        ) : null}
      </div>
      <p className="mt-4 text-2xl font-semibold tracking-tight text-slate-950">
        {value}
      </p>
      {detail ? <p className="mt-2 text-sm text-slate-500">{detail}</p> : null}
    </article>
  );
}
