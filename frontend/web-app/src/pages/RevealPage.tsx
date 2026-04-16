import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { getSecretStatus, consumeSecret } from "../lib/api";
import { decryptSecret, decodeBase64Url, importKeyFromBase64Url } from "../lib/crypto";

type RevealState = 
  | { status: "loading" }
  | { status: "ready"; expiresAt: string }
  | { status: "revealing" }
  | { status: "revealed"; plaintext: string }
  | { status: "error"; message: string };

export function RevealPage() {
  const { secretId } = useParams<{ secretId: string }>();
  const [state, setState] = useState<RevealState>({ status: "loading" });

  // Extract key from URL fragment
  const fragmentKey = window.location.hash.slice(1); // Remove '#' prefix

  useEffect(() => {
    if (!secretId) {
      setState({
        status: "error",
        message: "Liên kết không hợp lệ. Thiếu ID bí mật.",
      });
      return;
    }

    if (!fragmentKey) {
      setState({
        status: "error",
        message: "Liên kết không hợp lệ. Thiếu khóa giải mã.",
      });
      return;
    }

    // Check secret status without consuming
    getSecretStatus(secretId)
      .then((status) => {
        if (status.status === "pending") {
          setState({
            status: "ready",
            expiresAt: status.expiresAt || "",
          });
        } else {
          setState({
            status: "error",
            message: "Bí mật không tồn tại hoặc đã hết hạn.",
          });
        }
      })
      .catch((error) => {
        setState({
          status: "error",
          message: error.message || "Không thể kiểm tra trạng thái bí mật.",
        });
      });
  }, [secretId, fragmentKey]);

  const handleReveal = async () => {
    if (!secretId || !fragmentKey) {
      return;
    }

    setState({ status: "revealing" });

    try {
      // Consume the secret
      const response = await consumeSecret(secretId);

      // Import the key from URL fragment
      const key = await importKeyFromBase64Url(fragmentKey);

      // Decode ciphertext and nonce
      const ciphertext = decodeBase64Url(response.ciphertext);
      const nonce = decodeBase64Url(response.nonce);

      // Decrypt the secret
      const plaintext = await decryptSecret(ciphertext, key, nonce);

      setState({
        status: "revealed",
        plaintext,
      });
    } catch (error: any) {
      if (error.message === "already_consumed") {
        setState({
          status: "error",
          message: "Bí mật này đã được xem trước đó.",
        });
      } else if (error.message === "not_found") {
        setState({
          status: "error",
          message: "Bí mật không tồn tại hoặc đã hết hạn.",
        });
      } else {
        setState({
          status: "error",
          message: "Liên kết không hợp lệ hoặc dữ liệu đã bị hỏng.",
        });
      }
    }
  };

  const handleCopy = () => {
    if (state.status === "revealed") {
      navigator.clipboard.writeText(state.plaintext);
    }
  };

  return (
    <main className="shell">
      <section className="hero">
        <p className="eyebrow">quorix.io.vn / one-time-link</p>
        <h1>Xem bí mật</h1>
      </section>

      <section className="panel">
        {state.status === "loading" && (
          <div className="reveal-state">
            <p className="status wait">Đang kiểm tra...</p>
          </div>
        )}

        {state.status === "ready" && (
          <div className="reveal-gate">
            <div className="reveal-info">
              <p className="reveal-message">
                Bí mật này chỉ có thể xem <strong>một lần duy nhất</strong>.
              </p>
              <p className="reveal-warning">
                Sau khi bạn nhấn nút bên dưới, bí mật sẽ bị xóa vĩnh viễn.
              </p>
              {state.expiresAt && (
                <p className="reveal-expiry">
                  Hết hạn lúc: {new Date(state.expiresAt).toLocaleString("vi-VN")}
                </p>
              )}
            </div>
            <button
              type="button"
              className="btn-primary btn-large"
              onClick={handleReveal}
            >
              Nhấn để xem bí mật
            </button>
          </div>
        )}

        {state.status === "revealing" && (
          <div className="reveal-state">
            <p className="status wait">Đang giải mã...</p>
          </div>
        )}

        {state.status === "revealed" && (
          <div className="reveal-success">
            <p className="reveal-label">Nội dung bí mật:</p>
            <div className="secret-content">
              <pre>{state.plaintext}</pre>
            </div>
            <button
              type="button"
              className="btn-secondary"
              onClick={handleCopy}
            >
              Sao chép
            </button>
            <p className="reveal-notice">
              Bí mật này đã bị xóa và không thể xem lại.
            </p>
          </div>
        )}

        {state.status === "error" && (
          <div className="reveal-error">
            <p className="status warn">Lỗi</p>
            <p className="error-message">{state.message}</p>
            <a href="/" className="btn-secondary">
              Tạo bí mật mới
            </a>
          </div>
        )}
      </section>
    </main>
  );
}
