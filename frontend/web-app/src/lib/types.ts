export type HealthResponse = {
  service: string;
  status: string;
  timestamp: string;
  version: string;
};

export type CreateSecretRequest = {
  ciphertext: string;
  nonce: string;
  algorithm: string;
  ttl_seconds: number;
};

export type CreateSecretResponse = {
  secret_id: string;
  expires_at: string;
};

export type SecretState = "pending" | "already_used" | "expired" | "not_found";

export type SecretStatusResponse = {
  secret_id: string;
  status: SecretState;
};

export type CreateRevealSessionRequest = {
  secret_id: string;
};

export type CreateRevealSessionResponse = {
  session_id: string;
  expires_at: string;
};

export type ConsumeSecretRequest = {
  session_id: string;
};

export type ConsumeSecretResponse = {
  secret_id: string;
  ciphertext: string;
  nonce: string;
  algorithm: string;
};
