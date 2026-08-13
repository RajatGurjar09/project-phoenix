import Link from "next/link";
import { notFound } from "next/navigation";
import { getDeployments, getProject } from "@/lib/backend";
import { DeployContainerForm } from "@/app/components/deploy-container-form";
import { DeploymentActions } from "@/app/components/deployment-actions";

type ProjectPageProps = {
  params: Promise<{
    id: string;
  }>;
};

function statusClass(status: string) {
  switch (status.toLowerCase()) {
    case "running":
      return "bg-emerald-50 text-emerald-700 ring-emerald-600/20";
    case "stopped":
      return "bg-amber-50 text-amber-700 ring-amber-600/20";
    case "failed":
      return "bg-red-50 text-red-700 ring-red-600/20";
    case "removed":
      return "bg-slate-100 text-slate-600 ring-slate-500/10";
    default:
      return "bg-blue-50 text-blue-700 ring-blue-600/20";
  }
}

export default async function ProjectPage({ params }: ProjectPageProps) {
  const { id } = await params;

  let project;

  try {
    project = await getProject(id);
  } catch {
    notFound();
  }

  const deployments = await getDeployments(id);

  return (
    <main className="min-h-screen bg-slate-50 px-6 py-10 text-slate-950 sm:px-10 lg:px-16">
      <div className="mx-auto max-w-6xl">
        <Link
          href="/"
          className="text-sm font-medium text-indigo-600 hover:text-indigo-500"
        >
          ← Back to projects
        </Link>

        <header className="mt-8 border-b border-slate-200 pb-8">
          <p className="text-sm font-semibold uppercase tracking-[0.2em] text-indigo-600">
            Project
          </p>

          <div className="mt-3 flex flex-col justify-between gap-5 sm:flex-row sm:items-end">
            <div>
              <h1 className="text-4xl font-semibold tracking-tight">
                {project.name}
              </h1>

              <p className="mt-3 max-w-2xl text-lg leading-8 text-slate-600">
                {project.description ?? "No description provided."}
              </p>
            </div>

            <span className="rounded-full bg-emerald-50 px-3 py-1.5 text-sm font-semibold text-emerald-700">
              Active
            </span>
          </div>
        </header>

        <section className="pt-10">
          <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-end">
            <div>
              <h2 className="text-xl font-semibold">Deployments</h2>
              <p className="mt-1 text-sm text-slate-500">
                Docker containers deployed for this project.
              </p>
            </div>

            <DeployContainerForm projectId={id} />
          </div>

          {deployments.length === 0 ? (
            <div className="mt-6 rounded-2xl border border-dashed border-slate-300 bg-white p-10 text-center">
              <h3 className="font-semibold">No deployments yet</h3>
              <p className="mt-2 text-sm text-slate-500">
                Deploy a Docker image to get started.
              </p>
            </div>
          ) : (
            <div className="mt-6 space-y-4">
              {deployments.map((deployment) => (
                <article
                  key={deployment.id}
                  className="rounded-2xl border border-slate-200 bg-white p-6 shadow-sm"
                >
                  <div className="flex flex-col justify-between gap-4 sm:flex-row sm:items-start">
                    <div className="min-w-0">
                      <div className="flex flex-wrap items-center gap-3">
                        <h3 className="text-lg font-semibold">
                          {deployment.image}
                        </h3>

                        <span
                          className={`rounded-full px-2.5 py-1 text-xs font-semibold capitalize ring-1 ring-inset ${statusClass(
                            deployment.status,
                          )}`}
                        >
                          {deployment.status}
                        </span>
                      </div>

                      <p className="mt-3 break-all text-sm text-slate-500">
                        Deployment ID: {deployment.id}
                      </p>

                      {deployment.container_id ? (
                        <p className="mt-1 break-all text-sm text-slate-500">
                          Container: {deployment.container_id}
                        </p>
                      ) : null}
                    </div>

                    <DeploymentActions
                      deploymentId={deployment.id}
                      status={deployment.status}
                    />
                  </div>

                  <div className="mt-5 border-t border-slate-100 pt-4 text-sm text-slate-500">
                    Created{" "}
                    {new Intl.DateTimeFormat("en-US", {
                      dateStyle: "medium",
                      timeStyle: "short",
                      timeZone: "UTC",
                    }).format(new Date(deployment.created_at))}
                  </div>
                </article>
              ))}
            </div>
          )}
        </section>
      </div>
    </main>
  );
}
