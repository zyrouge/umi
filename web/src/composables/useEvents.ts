import { ref, shallowRef } from "vue";
import type { UmiEvent } from "../types/event";

const RECONNECT_BASE_MS = 1000;
const RECONNECT_MAX_MS = 30000;
const HISTORY_LIMIT = 200;

export function useEvents() {
    const events = shallowRef<UmiEvent[]>([]);
    const connected = ref(false);
    const error = ref<string | null>(null);

    let ws: WebSocket | null = null;
    let reconnectTimer: ReturnType<typeof setTimeout> | null = null;
    let reconnectDelay = RECONNECT_BASE_MS;
    let currentTags: string[] = [];
    let currentToken = "";
    let intentionalClose = false;

    function buildUrl(tags: string[], token: string): string {
        const proto = location.protocol === "https:" ? "wss" : "ws";
        const tagParam = encodeURIComponent(tags.join(","));
        return `${proto}://${location.host}/api/ws?tags=${tagParam}&token=${encodeURIComponent(token)}`;
    }

    function connect(tags: string[], token: string): void {
        currentTags = tags;
        currentToken = token;
        intentionalClose = false;
        openSocket();
    }

    function openSocket(): void {
        if (ws) {
            ws.onclose = null;
            ws.close();
        }
        ws = new WebSocket(buildUrl(currentTags, currentToken));

        ws.onopen = () => {
            connected.value = true;
            error.value = null;
            reconnectDelay = RECONNECT_BASE_MS;
        };

        ws.onmessage = (e: MessageEvent) => {
            try {
                const event: UmiEvent = JSON.parse(e.data);
                // Prepend newest first, cap at HISTORY_LIMIT
                const updated = [event, ...events.value];
                events.value = updated.slice(0, HISTORY_LIMIT);
            } catch {
                // ignore malformed messages
            }
        };

        ws.onerror = () => {
            error.value = "WebSocket error";
        };

        ws.onclose = () => {
            connected.value = false;
            ws = null;
            if (!intentionalClose) {
                scheduleReconnect();
            }
        };
    }

    function scheduleReconnect(): void {
        reconnectTimer = setTimeout(() => {
            reconnectDelay = Math.min(reconnectDelay * 2, RECONNECT_MAX_MS);
            openSocket();
        }, reconnectDelay);
    }

    function disconnect(): void {
        intentionalClose = true;
        if (reconnectTimer) {
            clearTimeout(reconnectTimer);
            reconnectTimer = null;
        }
        if (ws) {
            ws.close();
            ws = null;
        }
        connected.value = false;
        events.value = [];
    }

    return { events, connected, error, connect, disconnect };
}
