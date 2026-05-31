<script setup lang="ts">
import { reactive } from "vue";
import { useAuth } from "./composables/useAuth";
import { useEvents } from "./composables/useEvents";

const auth = useAuth();
const eventsState = useEvents();

const loginForm = reactive({ username: "", password: "" });
const tagInput = reactive({ value: "" });

const handleLogin = async () => {
    await auth.login({
        username: loginForm.username,
        password: loginForm.password,
    });
};

const handleLogout = () => {
    eventsState.disconnect();
    auth.logout();
};

const handleConnect = () => {
    const tags = tagInput.value
        .split(",")
        .map((t) => t.trim())
        .filter(Boolean);
    if (!tags.length) return;
    const token = auth.state.user!.token;
    eventsState.connect(tags, token);
};

const handleDisconnect = () => {
    eventsState.disconnect();
};

const formatTime = (iso: string): string => {
    return new Date(iso).toLocaleString();
};
</script>

<template>
    <div v-if="!auth.state.user" data-view="login">
        <form @submit.prevent="handleLogin">
            <input
                v-model="loginForm.username"
                type="text"
                name="username"
                autocomplete="username"
                placeholder="Username"
                required
            />
            <input
                v-model="loginForm.password"
                type="password"
                name="password"
                autocomplete="current-password"
                placeholder="Password"
                required
            />
            <button type="submit" :disabled="auth.state.loading">
                {{ auth.state.loading ? "Signing in…" : "Sign in" }}
            </button>
            <p v-if="auth.state.error" data-role="error">
                {{ auth.state.error }}
            </p>
        </form>
    </div>

    <div v-else data-view="main">
        <header>
            <span data-role="username">{{ auth.state.user.username }}</span>
            <button @click="handleLogout">Sign out</button>
        </header>

        <section data-role="subscribe">
            <input
                v-model="tagInput"
                type="text"
                placeholder="Tags (comma-separated, e.g. alerts,production)"
                :disabled="eventsState.connected.value"
            />
            <button v-if="!eventsState.connected.value" @click="handleConnect">
                Connect
            </button>
            <button v-else @click="handleDisconnect">Disconnect</button>
            <span
                data-role="status"
                :data-connected="eventsState.connected.value"
            >
                {{ eventsState.connected.value ? "Connected" : "Disconnected" }}
            </span>
            <p v-if="eventsState.error.value" data-role="error">
                {{ eventsState.error.value }}
            </p>
        </section>

        <section data-role="feed">
            <p v-if="eventsState.events.value.length === 0" data-role="empty">
                No events yet.
            </p>
            <article
                v-for="event in eventsState.events.value"
                :key="event.id"
                data-role="event"
                :data-level="event.level ?? 'none'"
            >
                <div data-role="event-header">
                    <img
                        v-if="event.icon_url"
                        :src="event.icon_url"
                        :alt="event.service"
                        data-role="event-icon"
                    />
                    <strong data-role="event-title">{{ event.title }}</strong>
                    <span v-if="event.level" data-role="event-level">{{
                        event.level
                    }}</span>
                    <span data-role="event-service">{{ event.service }}</span>
                    <time data-role="event-time" :datetime="event.created_at">{{
                        formatTime(event.created_at)
                    }}</time>
                </div>
                <p v-if="event.body" data-role="event-body">{{ event.body }}</p>
                <div v-if="event.tags.length" data-role="event-tags">
                    <span
                        v-for="tag in event.tags"
                        :key="tag"
                        data-role="event-tag"
                        >{{ tag }}</span
                    >
                </div>
                <a
                    v-if="event.action_url"
                    :href="event.action_url"
                    target="_blank"
                    rel="noopener"
                    data-role="event-action"
                >
                    View details
                </a>
                <details
                    v-if="
                        event.metadata && Object.keys(event.metadata).length > 0
                    "
                    data-role="event-metadata"
                >
                    <summary>Metadata</summary>
                    <dl>
                        <template
                            v-for="(val, key) in event.metadata"
                            :key="key"
                        >
                            <dt>{{ key }}</dt>
                            <dd>{{ val }}</dd>
                        </template>
                    </dl>
                </details>
            </article>
        </section>
    </div>
</template>
