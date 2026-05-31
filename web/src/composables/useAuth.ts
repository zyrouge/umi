import { reactive } from "vue";
import type { UmiAuthState, UmiLoginPayload, UmiUser } from "../types/auth";

const TOKEN_KEY = "umi_token";
const USERNAME_KEY = "umi_username";

const state = reactive<UmiAuthState>({
    user: loadStoredUser(),
    loading: false,
    error: null,
});

function loadStoredUser(): UmiUser | null {
    const token = localStorage.getItem(TOKEN_KEY);
    const username = localStorage.getItem(USERNAME_KEY);
    if (token && username) {
        return { token, username };
    }
    return null;
}

async function login(payload: UmiLoginPayload): Promise<void> {
    state.loading = true;
    state.error = null;
    try {
        const res = await fetch("/api/auth/login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload),
        });
        const json = await res.json();
        if (!res.ok || !json.success) {
            state.error = json.data ?? "Login failed";
            return;
        }
        const user: UmiUser = {
            token: json.data.token,
            username: json.data.username,
        };
        localStorage.setItem(TOKEN_KEY, user.token);
        localStorage.setItem(USERNAME_KEY, user.username);
        state.user = user;
    } catch (e) {
        state.error = "Network error";
    } finally {
        state.loading = false;
    }
}

function logout(): void {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USERNAME_KEY);
    state.user = null;
    state.error = null;
}

export function useAuth() {
    return { state, login, logout };
}
