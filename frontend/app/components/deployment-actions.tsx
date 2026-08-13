"use client";

import { useState } from "react";
import {
  removeDeployment,
  restartDeployment,
  stopDeployment,
} from "@/lib/backend";

type DeploymentActionsProps = {
  deploymentId: string;
  status: string;
};

export function DeploymentActions({
  deploymentId,
  status,
}: DeploymentActionsProps) {
  const [loading, setLoading] = useState<string | null>(null);
  const [error, setError] = useState("");

  async function runAction(
    action: "stop" | "restart" | "remove",
  ) {
    setLoading(action);
    setError("");

    try {
      if (action === "stop") {
        await stopDeployment(deploymentId);
      }

      if (action === "restart") {
        await restartDeployment(deploymentId);
      }

      if (action === "remove") {
        await removeDeployment(deploymentId);
      }

      window.location.reload();
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Deployment action failed",
      );
      setLoading(null);
    }
  }

  return (
    <div>
      <div className="flex flex-wrap gap-2">
        {status === "running" ? (
          <button
            type="button"
            disabled={loading !== null}
            onClick={() => runAction("stop")}
            className="rounded-lg border border-slate-200 px-3 py-2 text-sm font-medium hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {loading === "stop" ? "Stopping..." : "Stop"}
          </button>
        ) : null}

        {status === "stopped" ? (
          <button
            type="button"
            disabled={loading !== null}
            onClick={() => runAction("restart")}
            className="rounded-lg border border-slate-200 px-3 py-2 text-sm font-medium hover:bg-slate-50 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {loading === "restart" ? "Restarting..." : "Restart"}
          </button>
        ) : null}

        {status !== "removed" && status !== "pending" ? (
          <button
            type="button"
            disabled={loading !== null}
            onClick={() => runAction("remove")}
            className="rounded-lg border border-red-200 px-3 py-2 text-sm font-medium text-red-600 hover:bg-red-50 disabled:cursor-not-allowed disabled:opacity-50"
          >
            {loading === "remove" ? "Removing..." : "Remove"}
          </button>
        ) : null}
      </div>
      
      {status === "pending" ? (
        <span className="text-sm text-slate-500">Deploying...</span>
      ) : null}

      {error ? (
        <p className="mt-2 text-sm text-red-600" role="alert">
          {error}
        </p>
      ) : null}
    </div>
  );
}
