import type { 
  HealthResponse, 
  CreateSecretRequest, 
  CreateSecretResponse,
  SecretStatus,
  ConsumeSecretResponse
} from "./types";

const defaultApiBaseUrl = "http://localhost:8080";

export const apiBaseUrl =
  import.meta.env.VITE_API_BASE_URL?.trim() || defaultApiBaseUrl;

/**
 * Generate a UUID v4 for request tracking
 */
function generateRequestId(): string {
  return crypto.randomUUID();
}

export async function fetchHealth(): Promise<HealthResponse> {
  const response = await fetch(`${apiBaseUrl}/healthz`, {
    headers: {
      Accept: "application/json",
      "X-Request-ID": generateRequestId(),
    },
  });

  if (!response.ok) {
    throw new Error(`Health check failed with status ${response.status}`);
  }

  return (await response.json()) as HealthResponse;
}

export async function createSecret(
  request: CreateSecretRequest
): Promise<CreateSecretResponse> {
  const requestId = generateRequestId();
  
  const response = await fetch(`${apiBaseUrl}/api/secrets`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
      "X-Request-ID": requestId,
    },
    body: JSON.stringify(request),
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new Error(
      errorData.message || `Failed to create secret: ${response.status}`
    );
  }

  return (await response.json()) as CreateSecretResponse;
}

export async function getSecretStatus(
  secretId: string
): Promise<SecretStatus> {
  const requestId = generateRequestId();
  
  const response = await fetch(`${apiBaseUrl}/api/secrets/${secretId}/status`, {
    method: "GET",
    headers: {
      Accept: "application/json",
      "X-Request-ID": requestId,
    },
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new Error(
      errorData.message || `Failed to get secret status: ${response.status}`
    );
  }

  return (await response.json()) as SecretStatus;
}

export async function consumeSecret(
  secretId: string
): Promise<ConsumeSecretResponse> {
  const requestId = generateRequestId();
  
  const response = await fetch(`${apiBaseUrl}/api/secrets/${secretId}/consume`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
      "X-Request-ID": requestId,
    },
    body: JSON.stringify({}),
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    
    // Handle specific error cases
    if (response.status === 410) {
      throw new Error("already_consumed");
    }
    if (response.status === 404) {
      throw new Error("not_found");
    }
    
    throw new Error(
      errorData.message || `Failed to consume secret: ${response.status}`
    );
  }

  return (await response.json()) as ConsumeSecretResponse;
}
