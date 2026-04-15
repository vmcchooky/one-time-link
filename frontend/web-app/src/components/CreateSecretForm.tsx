import { useState } from "react";
import { createSecret } from "../lib/api";
import {
  generateKey,
  generateNonce,
  encryptSecret,
  encodeBase64Url,
  exportKeyToBase64Url,
} from "../lib/crypto";

const TTL_OPTIONS = [
  { value: 3600, label: "1 giờ" },
  { value: 86400, label: "24 giờ" },
  { value: 604800, label: "7 ngày" },
];

const MAX_PLAINTEXT_SIZE = 10 * 1024; // 10KB

export function CreateSecretForm() {
  const [plaintext, setPlaintext] = useState("");
  const [ttl, setTtl] = useState(3600);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [secretLink, setSecretLink] = useState("");

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setSecretLink("");

    // Validate plaintext size
    const plaintextBytes = new TextEncoder().encode(plaintext);
    if (plaintextBytes.length > MAX_PLAINTEXT_SIZE) {
      setError(`Nội dung vượt quá giới hạn ${MAX_PLAINTEXT_SIZE / 1024}KB`);
      return;
    }

    if (!plaintext.trim()) {
      setError("Vui lòng nhập nội dung bí mật");
      return;
    }

    setLoading(true);

    try {
      // Generate encryption key and nonce
      const key = await generateKey();
      const nonce = generateNonce();

      // Encrypt the plaintext
      const ciphertextBytes = await encryptSecret(plaintext, key, nonce);

      // Encode to base64url
      const ciphertext = encodeBase64Url(ciphertextBytes);
      const nonceB64 = encodeBase64Url(nonce);

      // Create secret via API
      const response = await createSecret({
        ciphertext,
        nonce: nonceB64,
        algorithm: "AES-GCM",
        ttlSeconds: ttl,
      });

      // Export key for URL fragment
      const keyB64 = await exportKeyToBase64Url(key);

      // Build secret link
      const baseUrl = window.location.origin;
      const link = `${baseUrl}/reveal/${response.secretId}#${keyB64}`;

      setSecretLink(link);
      setPlaintext(""); // Clear form
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Không thể tạo liên kết bí mật"
      );
    } finally {
      setLoading(false);
    }
  };

  const copyToClipboard = async () => {
    try {
      await navigator.clipboard.writeText(secretLink);
      alert("Đã sao chép liên kết!");
    } catch (err) {
      alert("Không thể sao chép. Vui lòng sao chép thủ công.");
    }
  };

  return (
    <div className="create-secret-form">
      <h2>Tạo liên kết bí mật một lần</h2>

      {!secretLink ? (
        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label htmlFor="plaintext">Nội dung bí mật:</label>
            <textarea
              id="plaintext"
              value={plaintext}
              onChange={(e) => setPlaintext(e.target.value)}
              placeholder="Nhập mật khẩu, API key, hoặc thông tin bí mật..."
              rows={6}
              disabled={loading}
              maxLength={MAX_PLAINTEXT_SIZE}
            />
            <small>
              {new TextEncoder().encode(plaintext).length} / {MAX_PLAINTEXT_SIZE} bytes
            </small>
          </div>

          <div className="form-group">
            <label htmlFor="ttl">Thời gian hết hạn:</label>
            <select
              id="ttl"
              value={ttl}
              onChange={(e) => setTtl(Number(e.target.value))}
              disabled={loading}
            >
              {TTL_OPTIONS.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          </div>

          {error && <div className="error-message">{error}</div>}

          <button type="submit" disabled={loading || !plaintext.trim()}>
            {loading ? "Đang tạo..." : "Tạo liên kết"}
          </button>
        </form>
      ) : (
        <div className="success-result">
          <h3>✓ Liên kết đã được tạo!</h3>
          <p>Liên kết này chỉ có thể xem được <strong>một lần duy nhất</strong>.</p>

          <div className="secret-link-box">
            <input
              type="text"
              value={secretLink}
              readOnly
              onClick={(e) => e.currentTarget.select()}
            />
            <button onClick={copyToClipboard}>Sao chép</button>
          </div>

          <button
            onClick={() => {
              setSecretLink("");
              setPlaintext("");
            }}
            className="secondary"
          >
            Tạo liên kết mới
          </button>
        </div>
      )}
    </div>
  );
}
