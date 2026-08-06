import type { Project } from "@/lib/backend";

type ProjectCardProps = {
  project: Project;
};

const createdDateFormatter = new Intl.DateTimeFormat("en-US", {
  dateStyle: "medium",
  timeZone: "UTC",
});

export function ProjectCard({ project }: ProjectCardProps) {
  return (
    <article className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm">
      <h3 className="font-semibold text-slate-950">{project.name}</h3>
      <p className="mt-3 text-sm leading-6 text-slate-600">
        {project.description ?? "No description provided."}
      </p>
      <p className="mt-6 border-t border-slate-100 pt-4 text-sm text-slate-500">
        Created {createdDateFormatter.format(new Date(project.created_at))}
      </p>
    </article>
  );
}
