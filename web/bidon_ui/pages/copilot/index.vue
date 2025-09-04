<template>
  <PageContainer>
    <Card>
      <template #header>
        <div class="flex items-center justify-between px-3 sm:px-4 py-3">
          <div class="flex items-center gap-2">
            <Button
              label="New Chat"
              icon="pi pi-plus"
              text
              @click="onNewChat"
            />
          </div>
        </div>
      </template>

      <template #content>
        <div
          class="px-3 sm:px-4 space-y-3 sm:space-y-4 max-h-[60vh] sm:max-h-[65vh] overflow-y-auto"
          data-testid="copilot-messages"
        >
          <div v-if="historyLoading" class="text-sm text-gray-400">
            Loading conversation…
          </div>
          <div
            v-for="(msg, idx) in messages"
            :key="idx"
            :class="[
              'flex',
              msg.role === 'user' ? 'justify-end' : 'justify-start',
            ]"
          >
            <div
              :class="[
                'px-3 py-2 sm:px-4 sm:py-3 rounded-lg max-w-[85%] sm:max-w-[80%] whitespace-pre-wrap',
                msg.role === 'user'
                  ? 'bg-blue-500 text-white'
                  : 'bg-gray-100 text-gray-800',
              ]"
              :data-testid="
                msg.role === 'user' ? 'user-message' : 'assistant-message'
              "
            >
              {{ msg.content }}
            </div>
          </div>
          <div v-if="loading" class="text-sm text-gray-400">Thinking…</div>
        </div>
      </template>

      <template #footer>
        <div class="flex items-center gap-2 px-3 sm:px-4 py-3 pt-0">
          <InputText
            v-model="input"
            data-testid="copilot-input"
            placeholder="Type a message..."
            class="flex-1 min-w-0"
            @keydown.enter.prevent="onSend"
          />
          <Button
            :disabled="!canSend"
            data-testid="copilot-send"
            label="Send"
            icon="pi pi-send"
            class="shrink-0"
            @click="onSend"
          />
        </div>
      </template>
    </Card>
  </PageContainer>
</template>

<script setup lang="ts">
import { onMounted } from "vue";
import { Client } from "@langgraph-js/sdk";
import { useToast } from "primevue/usetoast";

definePageMeta({
  middleware: ["admin-only"],
});

interface ChatMessage {
  role: "user" | "assistant";
  content: string;
}

// Types to represent remote/copilot state and streaming payloads without using 'any'
type RemoteMessagePart = string | { text?: string };
interface RemoteMessage {
  type?: string;
  role?: string;
  content?: string | RemoteMessagePart[];
}
type CopilotValues = { messages?: RemoteMessage[] };

function getThreadId(t: unknown): string | null {
  if (t && typeof t === "object") {
    const obj = t as Record<string, unknown>;
    const raw = (obj["thread_id"] ?? obj["threadId"]) as unknown;
    return typeof raw === "string" ? raw : null;
  }
  return null;
}

function getErrorMessage(e: unknown): string {
  if (e instanceof Error) return e.message;
  if (e && typeof e === "object" && "message" in e) {
    const m = (e as { message?: unknown }).message;
    return typeof m === "string" ? m : "";
  }
  return typeof e === "string" ? e : "";
}

const config = useRuntimeConfig();
const copilotBase = config.public.copilotBase || "/api/copilot";

const messages = ref<ChatMessage[]>([]);
const input = ref("");
const loading = ref(false);
const historyLoading = ref(false);
const assistantId = ref<string | null>(null);
const copilotStore = useCopilotStore();
const threadId = computed(() => copilotStore.threadId as string | null);
const toast = useToast();

let client: Client;

onMounted(async () => {
  const apiUrl = new URL(copilotBase, window.location.origin).toString();
  client = new Client({
    apiUrl,
    defaultHeaders: {
      "X-Bidon-App": "web",
    },
  });
  // Use the conventional default assistant id; avoids extra calls and works for most setups
  assistantId.value = "agent";

  // Restore persisted thread or create a new one
  if (threadId.value) {
    try {
      historyLoading.value = true;
      // Use getState to fetch the latest state values for the thread
      const state = await client.threads.getState<CopilotValues>(
        threadId.value,
      );
      const raw: RemoteMessage[] =
        state?.values?.messages ?? ([] as RemoteMessage[]);

      const filtered = raw.filter((m: RemoteMessage) => m.type !== "tool");
      const normalized = filtered.map((m: RemoteMessage): ChatMessage => {
        const role = (m.role || m.type || "").toLowerCase();
        // Coerce role names to our UI roles
        const uiRole: "user" | "assistant" =
          role.includes("human") || role === "user" ? "user" : "assistant";
        let content = "";
        const c = m.content;
        if (typeof c === "string") content = c;
        else if (Array.isArray(c)) {
          for (const part of c) {
            if (typeof part === "string") content += part;
            else if (
              typeof part !== "string" &&
              part &&
              typeof part.text === "string"
            )
              content += part.text;
          }
        }
        return { role: uiRole, content };
      });
      messages.value = normalized;
    } catch (e: unknown) {
      console.error("[Copilot] failed to load history", e);
      toast.add({
        severity: "warn",
        summary: "History unavailable",
        detail: getErrorMessage(e) || "Could not fetch previous messages",
        life: 2500,
      });
    } finally {
      historyLoading.value = false;
    }
  } else {
    try {
      const thread = await client.threads.create();
      const id = getThreadId(thread);
      copilotStore.setThreadId(id);
    } catch (e) {
      console.error("[Copilot] failed to create thread", e);
    }
  }
});

const canSend = computed(
  () => !!input.value.trim() && !!assistantId.value && !loading.value,
);

async function onNewChat() {
  try {
    // Clear UI
    messages.value = [];
    input.value = "";

    // Create new thread
    const thread = await client.threads.create();
    const id = getThreadId(thread);
    copilotStore.setThreadId(id);

    // Visual feedback
    toast.add({
      severity: "success",
      summary: "New chat",
      detail: "Started a new conversation",
      life: 2000,
    });
  } catch (e: unknown) {
    console.error("[Copilot] failed to start new chat", e);
    toast.add({
      severity: "error",
      summary: "New chat failed",
      detail: getErrorMessage(e) || "Could not create a new thread",
      life: 3000,
    });
  }
}

async function onSend() {
  if (!canSend.value || !assistantId.value) return;
  const text = input.value.trim();
  input.value = "";

  messages.value.push({ role: "user", content: text });
  const assistantIndex =
    messages.value.push({ role: "assistant", content: "" }) - 1;
  loading.value = true;

  try {
    // Threadless streaming run for a simple chat
    const payload = {
      input: { messages: [{ role: "user", content: text }] },
      // messages-tuple returns [messageChunk, metadata]
      streamMode: "messages-tuple" as const,
    };

    let stream;
    if (threadId.value) {
      stream = client.runs.stream(threadId.value, assistantId.value, payload);
    } else {
      stream = client.runs.stream(null, assistantId.value, payload);
    }

    for await (const chunk of stream as AsyncIterable<{
      event: string;
      data: unknown;
    }>) {
      if (!chunk || chunk.event !== "messages") continue;

      const data = chunk.data as unknown;

      // messages-tuple -> [messageChunk, metadata]; plain 'messages' -> messageChunk
      const tuple = Array.isArray(data) ? data : null;
      const messageChunk = (tuple ? tuple[0] : data) as {
        type?: string;
        content?: unknown;
      };
      if (!messageChunk || messageChunk.type !== "AIMessageChunk") continue;

      let delta = "";
      const content = messageChunk.content as unknown;

      // Content can be a string or array of parts like [{ type: 'text', text: '...' }]
      if (typeof content === "string") {
        delta = content;
      } else if (Array.isArray(content)) {
        for (const part of content) {
          if (typeof part === "string") {
            delta += part;
          } else if (part && typeof part.text === "string") {
            delta += part.text;
          }
        }
      }

      if (delta) {
        messages.value[assistantIndex].content += delta;
      }
    }
  } catch (err: unknown) {
    console.error("[Copilot stream error]", err);
    const detail = getErrorMessage(err) || "Failed to send message";
    messages.value[assistantIndex].content = `Error: ${detail}`;
  } finally {
    loading.value = false;
  }
}
</script>
