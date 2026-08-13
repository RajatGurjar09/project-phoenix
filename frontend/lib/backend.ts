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

export type Deployment = {
  id: string;
  project_id: string;
  image: string;
  status: string;
  container_id?: string;
  created_at: string;
  updated_at: string;
};

type FetchBackendOptions = {
  method?: "GET" | "POST" | "DELETE";
  body?: unknown;
};

const backendUrl =
  typeof window === "undefined"
    ? process.env.BACKEND_URL ?? "http://backend:8080"
    : process.env.NEXT_PUBLIC_BACKEND_URL ?? "http://localhost:8080";
    
async function fetchBackend<T>(
  path: string,
  options: FetchBackendOptions = {},
): Promise<T> {
  const response = await fetch(`${backendUrl}${path}`, {
    method: options.method ?? "GET",
    cache: "no-store",
    headers: options.body
      ? {
          "Content-Type": "application/json",
        }
      : undefined,
    body: options.body ? JSON.stringify(options.body) : undefined,
  });

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

export async function getProject(id: string): Promise<Project> {
  return fetchBackend<Project>(`/projects/${id}`);
}

export async function getDeployments(
  projectId: string,
): Promise<Deployment[]> {
  return fetchBackend<Deployment[]>(
    `/projects/${projectId}/deployments`,
  );
}

export async function createDeployment(
  projectId: string,
  image: string,
): Promise<Deployment> {
  return fetchBackend<Deployment>(`/projects/${projectId}/deployments`, {
    method: "POST",
    body: {
      image,
      status: "pending",
    },
  });
}

export async function stopDeployment(id: string): Promise<Deployment> {
  return fetchBackend<Deployment>(`/deployments/${id}/stop`, {
    method: "POST",
  });
}

export async function restartDeployment(id: string): Promise<Deployment> {
  return fetchBackend<Deployment>(`/deployments/${id}/restart`, {
    method: "POST",
  });
}

export async function removeDeployment(id: string): Promise<Deployment> {
  return fetchBackend<Deployment>(`/deployments/${id}`, {
    method: "DELETE",
  });
}
