<script lang="ts">
  import { main } from "../../wailsjs/go/models.js";
  
  export let show: boolean;
  export let gameInfo: main.GameInfo | null;
  export let logs: Array<{ message: string; type: "info" | "success" | "error" }>;
  export let isRestoring: boolean;
  
  export let onClose: () => void;
  export let onRestoreBackup: () => void;
  
  let logContainer: HTMLDivElement;
  let previousLogsLength = 0;
  
  $: {
    if (logContainer && logs.length > 0) {
      if (logs.length !== previousLogsLength) {
        // Only auto-scroll if user is near the bottom (within 100px) or it's the first log
        const isNearBottom = logContainer.scrollHeight - logContainer.scrollTop - logContainer.clientHeight < 100;
        if (isNearBottom || previousLogsLength === 0) {
          setTimeout(() => {
            logContainer.scrollTop = logContainer.scrollHeight;
          }, 0);
        }
        previousLogsLength = logs.length;
      }
    }
  }
</script>

{#if show}
  <div class="fixed inset-0 z-50 flex">
    <div 
      class="flex-1 bg-zinc-950/75" 
      onclick={onClose}
      role="button"
      tabindex="0"
      aria-label="Close drawer"
    ></div>
    <div class="w-full max-w-2xl bg-zinc-900 border-l border-zinc-800 h-full overflow-hidden flex flex-col">
      <!-- Header -->
      <div class="bg-zinc-900 border-b border-zinc-800 px-6 py-4 flex items-center justify-between">
        <h2 class="text-xl font-semibold">Remove Translation</h2>
        <button
          onclick={onClose}
          class="text-zinc-400 hover:text-zinc-300"
          aria-label="Close drawer"
        >
          <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Main Content -->
      <div class="flex-1 flex flex-col overflow-hidden">
        <!-- Scrollable Content Area -->
        <div class="flex flex-col gap-6 p-6 overflow-y-auto flex-shrink-0">
          <!-- Game Selection Section -->
          <div class="flex flex-col gap-3">
            <div class="flex items-center justify-between">
              <span class="text-sm font-medium text-zinc-400 uppercase tracking-wide">
                Game Executable
              </span>
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
          </div>

          <!-- Info Section -->
          <div class="bg-blue-900/20 border border-blue-700/50 px-4 py-3 text-left">
            <div class="flex items-start gap-3">
              <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5 text-blue-400 flex-shrink-0 mt-0.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M11.25 11.25l.041-.02a.75.75 0 011.063.852l-.708 2.836a.75.75 0 001.063.853l.041-.021M21 12a9 9 0 11-18 0 9 9 0 0118 0zm-9-3.75h.008v.008H12V8.25z" />
              </svg>
              <div class="flex-1">
                <p class="text-sm text-blue-200 font-medium mb-1">Restore Backup</p>
                <p class="text-sm text-blue-300/90">
                  This will restore the game to its original state before translation. 
                  All translated files will be replaced with the backed-up originals.
                </p>
              </div>
            </div>
          </div>
        </div>

        <!-- Logs Section - Fills available space -->
        <div class="flex-1 border-t border-zinc-800 p-6 flex flex-col min-h-0">
          <div
            bind:this={logContainer}
            class="flex-1 bg-black px-5 py-4 overflow-y-scroll font-mono text-xs text-left"
          >
            {#each logs as log}
              <div class="py-0.5 {log.type === 'error' ? 'text-red-400' : log.type === 'success' ? 'text-emerald-400' : 'text-zinc-300'}">
                <span class="text-zinc-600">›</span>
                {log.message}
              </div>
            {/each}
          </div>
        </div>

        <!-- Footer Action - Stuck to bottom -->
        <div class="border-t border-zinc-800 p-6 flex flex-col gap-3 flex-shrink-0">
          <button
            onclick={onRestoreBackup}
            disabled={isRestoring || !gameInfo}
            class="w-full bg-blue-600 hover:bg-blue-500 disabled:bg-zinc-700 disabled:cursor-not-allowed px-4 py-3 text-sm font-semibold uppercase tracking-wide transition-colors"
          >
            {isRestoring ? "Restoring Backup..." : gameInfo ? "Restore Backup" : "No Game Selected"}
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}

