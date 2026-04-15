import { useEffect, useState } from "react";
import { apiBaseUrl, fetchHealth } from "./lib/api";
import type { HealthResponse } from "./lib/types";
import { CreateSecretForm } from "./components/CreateSecretForm";

export function App() {
  const [health, setHealth] = useState<HealthResponse | null>(null);
  const [error, setError] = useState<string>("");

  useEffect(() => {
    let cancelled = false;

    fetchHealth()
      .then((nextHealth) => {
        if (!cancelled) {
          setHealth(nextHealth);
        }
      })
      .catch((nextError: Error) => {
        if (!cancelled) {
          setError(nextError.message);
        }
      });

    return () => {
      cancelled = true;
    };
  }, []);

  return (
    <main className="shell">
      <section className="hero">
        <p className="eyebrow">quorix.io.vn / one-time-link</p>
        <h1>Chia sẻ bí mật an toàn một lần</h1>
        <p className="lede">
          Mã hóa phía client, chỉ xem được một lần, tự động hết hạn.
        </p>
      </section>

      <section className="panel">
        <CreateSecretForm />
      </section>

      <section className="grid">
        <article className="card">
          <h3>API Base URL</h3>
          <code>{apiBaseUrl}</code>
        </article>

        <article className="card">
          <h3>Trạng thái Backend</h3>
          {health ? (
            <div>
              <p className="status ok">Đã kết nối</p>
              <p>{health.service}</p>
              <p>{health.timestamp}</p>
              <p>Version: {health.version}</p>
            </div>
          ) : error ? (
            <div>
              <p className="status warn">Không khả dụng</p>
              <p>{error}</p>
            </div>
          ) : (
            <p className="status wait">Đang kiểm tra...</p>
          )}
        </article>
      </section>
    </main>
  );
}
