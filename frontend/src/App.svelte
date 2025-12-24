<script lang="ts">
  import "./app.css";
  import { main } from "../wailsjs/go/models.js";
  import { EventsOn } from "../wailsjs/runtime/runtime.js";
  import { onMount } from "svelte";

  import {
    ApplyPatch,
    SelectGameExeFile,
    SelectPatchFile,
  } from "../wailsjs/go/main/App.js";

  interface LogMessage {
    message: string;
    type: "info" | "success" | "error";
  }

  let gameInfo: main.GameInfo | null = null;
  let patchInfo: main.PatchInfo | null = null;
  let logs: LogMessage[] = [];
  let isPatching = false;
  let logContainer: HTMLDivElement;

  onMount(() => {
    EventsOn("log", (logMessage: LogMessage) => {
      logs = [...logs, logMessage];
      setTimeout(() => {
        if (logContainer) {
          logContainer.scrollTop = logContainer.scrollHeight;
        }
      }, 0);
    });
  });

  async function selectGameExeFile(): Promise<void> {
    gameInfo = await SelectGameExeFile();
  }

  async function selectPatchFile(): Promise<void> {
    patchInfo = await SelectPatchFile();
  }

  async function applyPatch(): Promise<void> {
    logs = [];
    isPatching = true;
    try {
      await ApplyPatch(gameInfo, patchInfo);
    } catch (error) {
      logs = [...logs, { message: `Error: ${error}`, type: "error" }];
    } finally {
      isPatching = false;
    }
  }

  function clearLogs(): void {
    logs = [];
  }

  function openLink(url: string): void {
    window.open(url, '_blank');
  }
</script>

<div class="flex flex-col h-screen bg-zinc-900 text-zinc-100">
  <!-- Header -->
  <div class="bg-zinc-950 border-b border-zinc-800 px-6 py-4 flex items-center justify-between">
    <h1 class="text-lg font-semibold tracking-tight">HTranslations Patcher</h1>
    
    <div class="flex items-center gap-4">
      <button
        onclick={() => openLink('https://htranslations.com')}
        class="text-sm text-zinc-400 transition-colors cursor-pointer"
      >
        Website
      </button>
      <button
        onclick={() => openLink('https://discord.gg/sKXZDn72cr')}
        class="text-sm text-zinc-400 transition-colors cursor-pointer"
      >
        Discord
      </button>
      <button
        onclick={() => openLink('https://ko-fi.com/h-translations')}
        class="text-sm text-emerald-400 transition-colors font-medium cursor-pointer"
      >
        Support Us
      </button>
    </div>
  </div>

  <!-- Main Content - Two Column Layout -->
  <div class="flex-1 flex overflow-hidden">
    <!-- Left Column - Controls -->
    <div class="flex flex-col w-1/2 border-r border-zinc-800">
      <div class="flex-1 flex flex-col gap-6 p-6 overflow-y-auto">
        <!-- Game Selection Section -->
        <div class="flex flex-col gap-3">
          <div class="flex items-center justify-between">
            <span class="text-sm font-medium text-zinc-400 uppercase tracking-wide"
              >Game Executable</span
            >
            {#if gameInfo}
              <span class="text-xs text-emerald-400">✓ Selected</span>
            {/if}
          </div>

          {#if gameInfo}
            <div class="bg-zinc-800 border border-zinc-700 px-4 py-3">
              <p class="text-sm text-zinc-300 font-mono truncate">
                {gameInfo.exePath}
              </p>
            </div>
          {:else}
            <div class="bg-zinc-800 border border-zinc-700 px-4 py-3">
              <p class="text-sm text-zinc-500 italic">No file selected</p>
            </div>
          {/if}

          <button
            onclick={selectGameExeFile}
            class="bg-zinc-800 hover:bg-zinc-700 active:bg-zinc-600 border border-zinc-700 px-4 py-3 text-sm font-medium transition-colors duration-150"
          >
            {gameInfo ? "Change Game.exe" : "Select Game.exe"}
          </button>
        </div>

        <!-- Patch Selection Section -->
          <div class="flex flex-col gap-3">
            <div class="flex items-center justify-between">
              <span
                class="text-sm font-medium text-zinc-400 uppercase tracking-wide"
                >Patch File</span
              >
              {#if patchInfo}
                <span class="text-xs text-emerald-400">✓ Selected</span>
              {/if}
            </div>

            {#if patchInfo}
              <div class="bg-zinc-800 border border-zinc-700 px-4 py-3">
                <p class="text-sm text-zinc-300 font-mono truncate">
                  {patchInfo.patchPath}
                </p>
              </div>
            {:else}
              <div class="bg-zinc-800 border border-zinc-700 px-4 py-3">
                <p class="text-sm text-zinc-500 italic">No patch selected</p>
              </div>
            {/if}

            <button
              onclick={selectPatchFile}
              class="bg-zinc-800 hover:bg-zinc-700 active:bg-zinc-600 border border-zinc-700 px-4 py-3 text-sm font-medium transition-colors duration-150"
            >
              {patchInfo ? "Change Patch File" : "Select Patch File"}
            </button>
          </div>
      </div>

      <!-- Footer Action -->
        <div class="border-t border-zinc-800 p-6">
          <button
            onclick={applyPatch}
            disabled={isPatching || !gameInfo || !patchInfo}
            class="w-full bg-emerald-600 hover:bg-emerald-500 active:bg-emerald-700 disabled:bg-zinc-700 disabled:cursor-not-allowed px-4 py-4 text-sm font-semibold uppercase tracking-wide transition-colors duration-150"
          >
            {isPatching ? "Applying Patch..." : gameInfo && patchInfo ? "Apply Patch" : "Select Game.exe and Patch File"}
          </button>
        </div>
    </div>

    <!-- Right Column - Logs -->
    <div class="flex flex-col w-1/2">
      <div
        bind:this={logContainer}
        class="flex-1 bg-black px-5 pb-5 pt-4 overflow-y-auto font-mono text-xs text-left"
      >
          {#each logs as log}
            <div class="py-0.5 {log.type === 'error' ? 'text-red-400' : log.type === 'success' ? 'text-emerald-400' : 'text-zinc-300'}">
              <span class="text-zinc-600">›</span>
              {log.message}
            </div>
          {/each}
      </div>
    </div>
  </div>
</div>
