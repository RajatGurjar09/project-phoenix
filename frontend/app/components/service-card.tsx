type ServiceCardProps = {
  name: string;
};

export function ServiceCard({ name }: ServiceCardProps) {
  return (
    <article className="rounded-2xl border border-slate-200 bg-white p-5 shadow-sm">
      <div className="flex items-center justify-between gap-4">
        <h3 className="font-semibold text-slate-950">{name}</h3>
        <span className="rounded-full bg-slate-100 px-2.5 py-1 text-xs font-semibold text-slate-500">
          Planned
        </span>
      </div>
      <div className="mt-6 flex items-center justify-between border-t border-slate-100 pt-4 text-sm">
        <span className="text-slate-500">Status</span>
        <span className="font-medium text-slate-700">Coming Soon</span>
      </div>
    </article>
  );
}
