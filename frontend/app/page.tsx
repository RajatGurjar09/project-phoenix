import { getBackendStatus, getProjects } from "@/lib/backend";
import { ProjectCard } from "./components/project-card";
import { StatusCard } from "./components/status-card";

export const dynamic = "force-dynamic";

export default async function Home() {
  const [{ health, responseTimeMs }, projects] = await Promise.all([
    getBackendStatus(),
    getProjects(),
  ]);
  const isHealthy = health.status.toLowerCase() === "ok";

  return (
    <main className="min-h-screen bg-slate-50 px-6 py-10 text-slate-950 sm:px-10 lg:px-16">
      <div className="mx-auto max-w-6xl">
        <header className="border-b border-slate-200 pb-10">
          <p className="text-sm font-semibold uppercase tracking-[0.2em] text-indigo-600">
            Operations dashboard
          </p>
          <div className="mt-4 flex flex-col justify-between gap-5 sm:flex-row sm:items-end">
            <div>
              <h1 className="text-4xl font-semibold tracking-tight sm:text-5xl">
                Phoenix Console
              </h1>
              <p className="mt-4 max-w-2xl text-lg leading-8 text-slate-600">
                Monitor the health and runtime details of your Phoenix platform.
              </p>
            </div>
            <span className="text-sm text-slate-400">Live on load</span>
          </div>
        </header>

        <section className="pt-10" aria-labelledby="backend-heading">
          <div>
            <h2 id="backend-heading" className="text-xl font-semibold">
              Backend overview
            </h2>
            <p className="mt-1 text-sm text-slate-500">
              Runtime information reported by the Phoenix API.
            </p>
          </div>

          <div className="mt-6 grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
            <StatusCard
              label="Status"
              value={isHealthy ? "Operational" : health.status}
              detail={`Health endpoint: ${health.status}`}
              tone={isHealthy ? "success" : "neutral"}
              badge={isHealthy ? "Healthy" : "Attention"}
            />
            <StatusCard
              label="Version"
              value={health.version}
              detail="Current API release"
            />
            <StatusCard
              label="Hostname"
              value={health.hostname}
              detail="Backend host"
            />
            <StatusCard label="Uptime" value={health.uptime} detail="Since startup" />
            <StatusCard
              label="Timestamp"
              value={health.timestamp}
              detail="Latest health check (UTC)"
            />
            <StatusCard
              label="API response time"
              value={`${responseTimeMs} ms`}
              detail="Measured from the console"
            />
          </div>
        </section>

        <section className="pt-12" aria-labelledby="projects-heading">
          <div>
            <h2 id="projects-heading" className="text-xl font-semibold">
              Projects
            </h2>
            <p className="mt-1 text-sm text-slate-500">
              Projects managed by the Phoenix platform.
            </p>
          </div>

          {projects.length === 0 ? (
            <div className="mt-6 rounded-2xl border border-dashed border-slate-300 bg-white p-8 text-center">
              <h3 className="font-semibold text-slate-950">No projects yet</h3>
              <p className="mt-2 text-sm text-slate-500">
                Create a project through the Phoenix API to see it here.
              </p>
            </div>
          ) : (
            <div className="mt-6 grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
              {projects.map((project) => (
                <ProjectCard key={project.id} project={project} />
              ))}
            </div>
          )}
        </section>
      </div>
    </main>
  );
}
