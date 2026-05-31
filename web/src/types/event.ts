export type UmiEventLevel = "info" | "warning" | "error" | "critical";

export interface UmiEvent {
    id: string;
    title: string;
    body: string | null;
    level: UmiEventLevel | null;
    action_url: string | null;
    icon_url: string | null;
    tags: string[];
    service: string;
    metadata: Record<string, string> | null;
    created_at: string; // ISO 8601
}

export interface UmiEventPayload {
    title: string;
    body?: string;
    level?: UmiEventLevel;
    action_url?: string;
    icon_url?: string;
    tags: string[];
    metadata?: Record<string, string>;
}
