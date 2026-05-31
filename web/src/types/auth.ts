export interface UmiUser {
    username: string;
    token: string;
}

export interface UmiAuthState {
    user: UmiUser | null;
    loading: boolean;
    error: string | null;
}

export interface UmiLoginPayload {
    username: string;
    password: string;
}
