"use client";

import { useState } from "react";
import { createDeployment } from "@/lib/backend";

type DeployContainerFormProps = {
  projectId: string;
};

export function DeployContainerForm({
  projectId,
}: DeployContainerFormProps) {
  const [image, setImage] = useState("nginx:alpine");
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleSubmit(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();

    const trimmedImage = image.trim();

    if (!trimmedImage) {
      setError("Docker image is required.");
      return;
    }

    setLoading(true);
    setError("");

    try {
      await createDeployment(projectId, trimmedImage);
      window.location.reload();
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Deployment failed.",
      );
      setLoading(false);
    }
  }

  if (!open) {
    return (
      <button
        type="button"
        onClick={() => setOpen(true)}
        className="rounded-xl bg-indigo-600 px-4 py-2.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500"
      >
        + Deploy container
      </button>
    );
  }

  return (
    <form
      onSubmit={handleSubmit}
      className="w-full rounded-2xl border border-slate-200 bg-white p-5 shadow-sm sm:w-auto"
    >
      <label
        htmlFor="docker-image"
        className="block text-sm font-medium text-slate-700"
      >
        Docker image
      </label>

      <div className="mt-2 flex flex-col gap-2 sm:flex-row">
        <input
          id="docker-image"
          value={image}
          onChange={(event) => setImage(event.target.value)}
          placeholder="nginx:alpine"
          disabled={loading}
          className="min-w-0 rounded-lg border border-slate-300 px-3 py-2 text-sm outline-none focus:border-indigo-500 focus:ring-2 focus:ring-indigo-500/20 sm:w-72"
        />

        <button
          type="submit"
          disabled={loading}
          className="rounded-lg bg-indigo-600 px-4 py-2 text-sm font-semibold text-white hover:bg-indigo-500 disabled:cursor-not-allowed disabled:opacity-50"
        >
          {loading ? "Deploying..." : "Deploy"}
        </button>

        <button
          type="button"
          onClick={() => {
            setOpen(false);
            setError("");
          }}
          disabled={loading}
          className="rounded-lg border border-slate-200 px-4 py-2 text-sm font-medium text-slate-700 hover:bg-slate-50 disabled:opacity-50"
        >
          Cancel
        </button>
      </div>

      {error ? (
        <p className="mt-3 text-sm text-red-600" role="alert">
          {error}
        </p>
      ) : null}
    </form>
  );
}
