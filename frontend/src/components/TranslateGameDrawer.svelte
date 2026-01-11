<script lang="ts">
  import { domain } from "../../wailsjs/go/models.js";
  
  export let show: boolean;
  export let gameInfo: domain.GameInfo | null;
  export let patches: domain.PatchEntry[];
  export let logs: Array<{ message: string; type: "info" | "success" | "error" | "warning" }>;
  export let isPatching: boolean;
  export let launchAfterPatch: boolean;
  export let selectedPatch: domain.PatchEntry | null;
  export let patchInfo: domain.PatchInfo | null;
  export let patchSearchQuery: string;
  
  export let onClose: () => void;
  export let onSelectPatchFile: () => void;
  export let onTogglePatch: (patch: domain.PatchEntry) => void;
  export let onClearCustomPatch: () => void;
  export let onApplyPatch: () => void;
  export let onLaunchAfterPatchChange: (value: boolean) => void;
  export let onPatchSearchQueryChange: (value: string) => void;
  
  let logContainer: HTMLDivElement;
  let previousLogsLength = 0;
  let patchSuccess = false;
  
  $: filteredPatches = (() => {
    const query = patchSearchQuery.toLowerCase();
    return patches.filter(
      (patch) =>
        patch.title.toLowerCase().includes(query) ||
        patch.systemGameTitle?.toLowerCase().includes(query)
    )
    .sort((a, b) => {
      if (a.systemGameTitle === gameInfo?.gameTitle && b.systemGameTitle !== gameInfo?.gameTitle) return -1;
      if (a.systemGameTitle !== gameInfo?.gameTitle && b.systemGameTitle === gameInfo?.gameTitle) return 1;
      return 0;
    });
  })();
  
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
        
        // Check if patching was successful
        const lastLog = logs[logs.length - 1];
        if (!isPatching && lastLog.type === "success" && lastLog.message.includes("✓")) {
          patchSuccess = true;
        }
      }
    }
  }
  
  // Reset success state when drawer is shown
  $: if (show) {
    patchSuccess = false;
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
        <h2 class="text-xl font-semibold">Translate Game</h2>
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

      <!-- Main Content - Single Column Layout -->
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

          <!-- Patch Selection Section -->
          <div class="flex flex-col gap-3">
            <div class="flex items-center justify-between">
              <span class="text-sm font-medium text-zinc-400 uppercase tracking-wide">
                Patch File
              </span>
              <div class="flex items-center gap-2">
                {#if patchInfo}
                  <span class="text-xs text-emerald-400">✓ Selected</span>
                {/if}
              </div>
            </div>
            
            {#if patchInfo && !selectedPatch}
              <!-- Custom File Selected -->
              <div class="bg-zinc-800 border border-zinc-700 px-4 py-3 flex items-center justify-between">
                <p class="text-sm text-zinc-300 font-mono truncate flex-1">
                  {patchInfo.patchPath}
                </p>
                <button
                  type="button"
                  onclick={onClearCustomPatch}
                  class="ml-3 text-zinc-400 hover:text-zinc-300 transition-colors"
                  title="Clear custom patch file"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
            {:else}
              <!-- Patches List -->
              <div class="text-sm text-zinc-500 text-left">
                Search for patches or select a .htpatch file manually.
              </div>
              
              <!-- Autocomplete Input -->
              <div class="flex items-center gap-2">
                <input
                  type="text"
                  value={patchSearchQuery}
                  oninput={(e) => onPatchSearchQueryChange((e.target as HTMLInputElement).value)}
                  placeholder="Search patches..."
                  class="flex-1 bg-zinc-800 border border-zinc-700 px-4 py-3 text-sm text-zinc-300 focus:outline-none focus:ring-0 focus-visible:outline-none transition-colors"
                />
                <button
                  type="button"
                  onclick={onSelectPatchFile}
                  class="bg-zinc-800 hover:bg-zinc-700 border border-zinc-700 px-3 py-3 text-zinc-300 transition-colors flex items-center justify-center"
                  title="Select patch file"
                >
                  <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5m-13.5-9L12 3m0 0l4.5 4.5M12 3v13.5" />
                  </svg>
                </button>
              </div>

              <!-- Autocomplete Dropdown -->
              <div class="bg-zinc-800 border border-zinc-700 max-h-34 overflow-y-auto shadow-lg">
                {#if filteredPatches.length > 0}
                  {#each filteredPatches as patch}
                    <button
                      type="button"
                      onmousedown={() => onTogglePatch(patch)}
                      class="w-full text-left px-4 py-2 text-sm text-zinc-300 hover:bg-zinc-700 transition-colors flex items-center justify-between relative border-b border-zinc-700 last:border-b-0 {selectedPatch === patch ? 'bg-zinc-700 border-l-2 border-l-emerald-500' : ''}"
                    >
                      {#if selectedPatch === patch}
                        <span class="absolute top-1/2 -translate-y-1/2 right-4 bg-emerald-500 text-white text-xs font-bold px-1.5 py-0.5">✓</span>
                      {/if}
                      <div class="flex-1 pr-8">
                        <div class="font-medium">{patch.title}</div>
                        {#if patch.systemGameTitle}
                          <div class="text-xs text-zinc-500 mt-1">{patch.systemGameTitle}</div>
                        {/if}
                      </div>
                    </button>
                  {/each}
                {:else}
                  <div class="px-4 py-2 text-sm text-zinc-500">
                    No patches found
                  </div>
                {/if}
              </div>
            {/if}
          </div>
        </div>

        <!-- Logs Section - Fills available space -->
        <div class="flex-1 border-t border-zinc-800 p-6 flex flex-col min-h-0">
          <div
            bind:this={logContainer}
            class="flex-1 bg-black px-5 py-4 overflow-y-scroll font-mono text-xs text-left"
          >
            {#each logs as log}
              <div class="py-0.5 {log.type === 'error' ? 'text-red-400' : log.type === 'success' ? 'text-emerald-400' : log.type === 'warning' ? 'text-amber-400' : 'text-zinc-300'}">
                <span class="text-zinc-600">›</span>
                {log.message}
              </div>
            {/each}
          </div>
        </div>

        <!-- Footer Action - Stuck to bottom -->
        <div class="border-t border-zinc-800 p-6 flex flex-col gap-3 flex-shrink-0">
          <label class="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              checked={launchAfterPatch}
              onchange={(e) => onLaunchAfterPatchChange((e.target as HTMLInputElement).checked)}
              class="w-4 h-4 bg-zinc-800 border-zinc-700 focus:ring-0 focus:ring-offset-0"
            />
            <span class="text-sm text-zinc-400">
              Launch game after patching
            </span>
          </label>
          <button
            onclick={onApplyPatch}
            disabled={isPatching || !gameInfo || !(patchInfo || selectedPatch) || patchSuccess}
            class="w-full bg-emerald-600 hover:bg-emerald-500 disabled:bg-zinc-700 disabled:cursor-not-allowed px-4 py-3 text-sm font-semibold uppercase tracking-wide transition-colors"
          >
            {patchSuccess ? "Patch Applied Successfully" : isPatching ? "Applying Patch..." : gameInfo && (patchInfo || selectedPatch) ? "Apply Patch" : "Select Patch File"}
          </button>
        </div>
      </div>
    </div>
  </div>
{/if}

