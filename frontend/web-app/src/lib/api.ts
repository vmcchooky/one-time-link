import type { HealthResponse } from "./types";

const defaultApiBaseUrl = "http://localhost:8080";

export const apiBaseUrl =
  import.meta.env.VITE_API_BASE_URL?.trim() || defaultApiBaseUrl;

export async function fetchHealth(): Promise<HealthResponse> {
  const response = await fetch(`${apiBaseUrl}/healthz`, {
    headers: {
      Accept: "application/json",
    },
  });

  if (!response.ok) {
    throw new Error(`Health check failed with status ${response.status}`);
  }

  return (await response.json()) as HealthResponse;
}
