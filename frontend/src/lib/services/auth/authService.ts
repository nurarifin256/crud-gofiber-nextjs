import { create } from "zustand";

export interface ILogin {
    email: string;
    password: string;
}

export const initLogin: ILogin = {
    email: "",
    password: "",
}

interface AuthService {
    error: any;
    resetError: () => void;
}

export const useAuthService = create<AuthService>((set) => ({
    error: null,
    resetError: () => set({ error: null})
}))