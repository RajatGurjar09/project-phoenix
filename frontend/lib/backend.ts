export type BackendHealth = {
  status: string;
  version: string;
  uptime: string;
  hostname: string;
  timestamp: string;
};

export type BackendStatus = {
  health: BackendHealth;
  responseTimeMs: number;
};

export type Project = {
  id: string;
  name: string;
  description: string | null;
  created_at: string;
  updated_at: string;
};

const backendUrl = process.env.BACKEND_URL ?? "http://backend:8080";

async function fetchBackend<T>(path: string): Promise<T> {
  const response = await fetch(`${backendUrl}${path}`, { cache: "no-store" });

  if (!response.ok) {
    throw new Error(`Backend request failed with status ${response.status}`);
  }

  return response.json() as Promise<T>;
}

export async function getBackendStatus(): Promise<BackendStatus> {
  const startedAt = Date.now();
  const health = await fetchBackend<BackendHealth>("/health");

  return {
    health,
    responseTimeMs: Date.now() - startedAt,
  };
}

export async function getProjects(): Promise<Project[]> {
  return fetchBackend<Project[]>("/projects");
}
