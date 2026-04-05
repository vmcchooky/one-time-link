import { useEffect, useState } from "react";
import { apiBaseUrl, fetchHealth } from "./lib/api";
import type { HealthResponse } from "./lib/types";

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
        <h1>Secure one-time secret sharing, scaffolded for real deployment.</h1>
        <p className="lede">
          This frontend is the starting shell for the create flow, reveal flow,
          and client-side encryption work described in the product docs.
        </p>
      </section>

      <section className="panel">
        <h2>Current scaffold status</h2>
        <ul>
          <li>Frontend stack: React + TypeScript + Vite</li>
          <li>Backend target: Go API behind a reverse proxy</li>
          <li>Storage target: Redis for TTL and atomic consume</li>
          <li>Deployment target: Vercel + VPS primary + Oracle standby</li>
        </ul>
      </section>

      <section className="grid">
        <article className="card">
          <h3>API Base URL</h3>
          <code>{apiBaseUrl}</code>
        </article>

        <article className="card">
          <h3>Backend health</h3>
          {health ? (
            <div>
              <p className="status ok">Connected</p>
              <p>{health.service}</p>
              <p>{health.time}</p>
            </div>
          ) : error ? (
            <div>
              <p className="status warn">Unavailable</p>
              <p>{error}</p>
            </div>
          ) : (
            <p className="status wait">Checking...</p>
          )}
        </article>
      </section>

      <section className="panel">
        <h2>Next implementation steps</h2>
        <ol>
          <li>Build the create-secret form and TTL selection.</li>
          <li>Add Web Crypto helpers for encrypt and decrypt flows.</li>
          <li>Implement Go endpoints for create, status, reveal session, and consume.</li>
          <li>Wire Redis for TTL storage and atomic one-time consume.</li>
        </ol>
      </section>
    </main>
  );
}
