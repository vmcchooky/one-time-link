export type HealthResponse = {
  service: string;
  status: string;
  time: string;
  store: string;
  mode: string;
};

export type CreateSecretRequest = {
  ciphertext: string;
  nonce: string;
  algorithm: string;
  ttlSeconds: number;
};

export type CreateSecretResponse = {
  secretId: string;
  expiresAt: string;
};

export type SecretState = "pending" | "already_used" | "expired" | "not_found";

export type SecretStatusResponse = {
  secretId: string;
  state: SecretState;
};

export type CreateRevealSessionRequest = {
  secretId: string;
};

export type CreateRevealSessionResponse = {
  sessionId: string;
  expiresAt: string;
};

export type ConsumeSecretRequest = {
  sessionId: string;
};

export type ConsumeSecretResponse = {
  secretId: string;
  ciphertext: string;
  nonce: string;
  algorithm: string;
};
