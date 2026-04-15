/**
 * Web Crypto API helpers for AES-GCM encryption/decryption
 * Following specs from docs/contracts/crypto-and-api-decisions.md
 */

const ALGORITHM = "AES-GCM";
const KEY_LENGTH = 256;
const NONCE_LENGTH = 12; // 12 bytes = 96 bits

/**
 * Generate a 256-bit AES-GCM key
 */
export async function generateKey(): Promise<CryptoKey> {
  return await crypto.subtle.generateKey(
    {
      name: ALGORITHM,
      length: KEY_LENGTH,
    },
    true, // extractable
    ["encrypt", "decrypt"]
  );
}

/**
 * Generate a 12-byte (96-bit) nonce for AES-GCM
 */
export function generateNonce(): Uint8Array {
  return crypto.getRandomValues(new Uint8Array(NONCE_LENGTH));
}

/**
 * Encrypt plaintext using AES-GCM
 * @param plaintext - The secret text to encrypt
 * @param key - The AES-GCM key
 * @param nonce - The 12-byte nonce
 * @returns The ciphertext as Uint8Array
 */
export async function encryptSecret(
  plaintext: string,
  key: CryptoKey,
  nonce: Uint8Array
): Promise<Uint8Array> {
  const encoder = new TextEncoder();
  const plaintextBytes = encoder.encode(plaintext);

  const ciphertextBuffer = await crypto.subtle.encrypt(
    {
      name: ALGORITHM,
      iv: nonce,
    },
    key,
    plaintextBytes
  );

  return new Uint8Array(ciphertextBuffer);
}

/**
 * Decrypt ciphertext using AES-GCM
 * @param ciphertext - The encrypted data
 * @param key - The AES-GCM key
 * @param nonce - The 12-byte nonce
 * @returns The decrypted plaintext string
 */
export async function decryptSecret(
  ciphertext: Uint8Array,
  key: CryptoKey,
  nonce: Uint8Array
): Promise<string> {
  const plaintextBuffer = await crypto.subtle.decrypt(
    {
      name: ALGORITHM,
      iv: nonce,
    },
    key,
    ciphertext
  );

  const decoder = new TextDecoder();
  return decoder.decode(plaintextBuffer);
}

/**
 * Encode bytes to base64url format (RFC 4648)
 * Safe for URLs and fragments
 */
export function encodeBase64Url(bytes: Uint8Array): string {
  // Convert to base64
  let base64 = btoa(String.fromCharCode(...bytes));
  
  // Convert to base64url: replace +/= with -_
  return base64
    .replace(/\+/g, "-")
    .replace(/\//g, "_")
    .replace(/=/g, "");
}

/**
 * Decode base64url format to bytes
 */
export function decodeBase64Url(base64url: string): Uint8Array {
  // Convert base64url to base64: replace -_ with +/
  let base64 = base64url
    .replace(/-/g, "+")
    .replace(/_/g, "/");
  
  // Add padding if needed
  while (base64.length % 4 !== 0) {
    base64 += "=";
  }
  
  // Decode base64 to bytes
  const binaryString = atob(base64);
  const bytes = new Uint8Array(binaryString.length);
  for (let i = 0; i < binaryString.length; i++) {
    bytes[i] = binaryString.charCodeAt(i);
  }
  
  return bytes;
}

/**
 * Export a CryptoKey to base64url format for URL fragments
 */
export async function exportKeyToBase64Url(key: CryptoKey): Promise<string> {
  const rawKey = await crypto.subtle.exportKey("raw", key);
  return encodeBase64Url(new Uint8Array(rawKey));
}

/**
 * Import a key from base64url format
 */
export async function importKeyFromBase64Url(base64url: string): Promise<CryptoKey> {
  const keyBytes = decodeBase64Url(base64url);
  
  return await crypto.subtle.importKey(
    "raw",
    keyBytes,
    {
      name: ALGORITHM,
      length: KEY_LENGTH,
    },
    true,
    ["encrypt", "decrypt"]
  );
}
