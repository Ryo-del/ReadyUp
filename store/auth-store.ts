import { create } from "zustand";
import type { AuthMode } from "@/lib/auth-api";

type AuthState = {
  mode: AuthMode;
  token: string | null;
  tokenType: string | null;
  setMode: (mode: AuthMode) => void;
  toggleMode: () => void;
  setSession: (token: string, tokenType: string) => void;
  clearSession: () => void;
};

export const useAuthStore = create<AuthState>((set) => ({
  mode: "login",
  token: null,
  tokenType: null,
  setMode: (mode) => set({ mode }),
  toggleMode: () =>
    set((state) => ({ mode: state.mode === "login" ? "register" : "login" })),
  setSession: (token, tokenType) => set({ token, tokenType }),
  clearSession: () => set({ token: null, tokenType: null }),
}));
