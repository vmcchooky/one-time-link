import { BrowserRouter, Routes, Route } from "react-router-dom";
import { HomePage } from "./pages/HomePage";
import { RevealPage } from "./pages/RevealPage";

export function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/reveal/:secretId" element={<RevealPage />} />
      </Routes>
    </BrowserRouter>
  );
}
