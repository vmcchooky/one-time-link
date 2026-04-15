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
  ttlSeconds: number;
};

export type CreateSecretResponse = {
  secretId: string;
  expiresAt: string;
};

export type SecretState = "pending" | "alreadyUsed" | "expired" | "notFound";

export type SecretStatusResponse = {
  secretId: string;
  status: SecretState;
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
