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

export type SecretStatus = {
  secretId: string;
  status: "pending" | "not_found";
  createdAt?: string;
  expiresAt?: string;
  message?: string;
};

export type ConsumeSecretResponse = {
  secretId: string;
  ciphertext: string;
  nonce: string;
  algorithm: string;
  consumedAt: string;
};
