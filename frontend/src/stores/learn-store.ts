import { create } from "zustand";

interface LearnState {
  theme: "light" | "dark";
  toggleTheme: () => void;
}

export const useLearnStore = create<LearnState>((set) => ({
  theme: "light",
  toggleTheme: () =>
    set((state: any) => ({ theme: state.theme === "light" ? "dark" : "light" })),
}));  

// Auth Store
interface AuthState {
  user: any;
  login: (user: any) => void;
  logout: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  login: (user: any) => set({ user }),
  logout: () => set({ user: null }),
}));