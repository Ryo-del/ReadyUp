export type AuthMode = "login" | "register";

export type AuthRequest = {
  email: string;
  password: string;
  username?: string;
};

export type AuthResponse = {
  token: string;
  token_type: string;
  expires_in: number;
};

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:3001/auth";

export async function submitAuth(mode: AuthMode, payload: AuthRequest) {
  const endpoint = mode === "login" ? "/login" : "/register";
  const response = await fetch(`${API_URL}${endpoint}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  });

  const contentType = response.headers.get("content-type");
  const data = contentType?.includes("application/json")
    ? await response.json()
    : await response.text();

  if (!response.ok) {
    const message =
      typeof data === "string"
        ? data.trim()
        : data?.message ?? "Что-то пошло не так. Попробуйте снова.";
    throw new Error(message || "Что-то пошло не так. Попробуйте снова.");
  }

  return data as AuthResponse;
}
