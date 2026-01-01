<script lang="ts">
  import "./app.css";
  import { main } from "../wailsjs/go/models.js";
  import { EventsOn, BrowserOpenURL } from "../wailsjs/runtime/runtime.js";
  import { onMount } from "svelte";

  import {
    ApplyPatch,
    SelectGameExeFile,
    SelectPatchFile,
    FetchAllPatches,
    DownloadPatch,
  } from "../wailsjs/go/main/App.js";

  interface LogMessage {
    message: string;
    type: "info" | "success" | "error";
  }

  let gameInfo: main.GameInfo | null = null;
  let patchInfo: main.PatchInfo | null = null;
  let logs: LogMessage[] = [];
  let isPatching = false;
  let launchAfterPatch = true;
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

  let patches: main.PatchEntry[] = [];
  let patchSearchQuery = "";
  let selectedPatch: main.PatchEntry | null = null;

  onMount(async () => {
    patches = await FetchAllPatches();
  });

  $: filteredPatches = (() => {
    const query = patchSearchQuery.toLowerCase();
    return patches.filter(
      (patch) =>
        patch.title.toLowerCase().includes(query) ||
        patch.systemGameTitle.toLowerCase().includes(query)
    )
    // Sort patches by game title match
    .sort((a, b) => {
      if (a.systemGameTitle === gameInfo?.gameTitle && b.systemGameTitle !== gameInfo?.gameTitle) return -1;
      if (a.systemGameTitle !== gameInfo?.gameTitle && b.systemGameTitle === gameInfo?.gameTitle) return 1;
      return 0;
    });
  })();

  async function selectGameExeFile(): Promise<void> {
    gameInfo = await SelectGameExeFile();
    selectedPatch ??= patches.find(patch => patch.systemGameTitle === gameInfo.gameTitle);
  }

  function togglePatch(patch: main.PatchEntry): void {
    if (selectedPatch === patch) {
      selectedPatch = null;
    } else {
      selectedPatch = patch;
      patchInfo = null;
    }
  }

  async function selectPatchFile(): Promise<void> {
    patchInfo = await SelectPatchFile();
    selectedPatch = null;
  }

  function clearCustomPatch(): void {
    patchInfo = null;
  }

  async function applyPatch(): Promise<void> {
    logs = [];
    isPatching = true;
    try {
      if (selectedPatch) {
        patchInfo = await DownloadPatch(selectedPatch.patchDownloadId);
      }
      await ApplyPatch(gameInfo, patchInfo, launchAfterPatch);
    } catch (error) {
      logs = [...logs, { message: `Error: ${error}`, type: "error" }];
    } finally {
      isPatching = false;
    }
  }

  function openLink(url: string): void {
    BrowserOpenURL(url);
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
        onclick={() => openLink('https://ko-fi.com/htranslations')}
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
          <div class="flex flex-col gap-3 relative">
            <div class="flex items-center justify-between">
              <span
                class="text-sm font-medium text-zinc-400 uppercase tracking-wide"
                >Patch File</span
              >
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
                  onclick={clearCustomPatch}
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
              <div class="relative">
                <div class="flex items-center gap-2">
                  <input
                    type="text"
                    bind:value={patchSearchQuery}
                    placeholder="Search patches..."
                    class="flex-1 bg-zinc-800 border border-zinc-700 px-4 py-3 text-sm text-zinc-300 focus:outline-none focus:ring-0 focus-visible:outline-none transition-colors"
                  />
                  <button
                    type="button"
                    onclick={selectPatchFile}
                    class="bg-zinc-800 hover:bg-zinc-700 active:bg-zinc-600 border border-zinc-700 px-3 py-3 text-zinc-300 transition-colors duration-150 flex items-center justify-center"
                    title="Select patch file"
                  >
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-5 h-5">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5m-13.5-9L12 3m0 0l4.5 4.5M12 3v13.5" />
                    </svg>
                  </button>
                </div>
                {#if filteredPatches.length > 0}
                  <div class="absolute z-50 w-full mt-3 bg-zinc-800 border border-zinc-700 max-h-38 overflow-y-auto">
                    {#each filteredPatches as patch}
                      <button
                        type="button"
                        onmousedown={() => togglePatch(patch)}
                        class="w-full text-left px-4 py-2 text-sm text-zinc-300 hover:bg-zinc-700 transition-colors flex items-center justify-between relative {selectedPatch === patch ? 'border-2 border-emerald-500' : 'border-b border-zinc-700 last:border-b-0'}"
                      >
                        {#if selectedPatch === patch}
                          <span class="absolute top-0 right-0 bg-emerald-500 text-white text-xs font-bold px-1.5 py-0.5">✓</span>
                        {/if}
                        <div class="flex-1">
                          <div class="font-medium">{patch.title}</div>
                          {#if patch.systemGameTitle}
                            <div class="text-xs text-zinc-500 mt-1">{patch.systemGameTitle}</div>
                          {/if}
                        </div>
                      </button>
                    {/each}
                  </div>
                {:else}
                  <div class="absolute z-50 w-full mt-3 bg-zinc-800 border border-zinc-700 px-4 py-2 text-sm text-zinc-500">
                    No patches found
                  </div>
                {/if}
              </div>
            {/if}
          </div>
      </div>

      <!-- Footer Action -->
        <div class="border-t border-zinc-800 p-6 flex flex-col gap-3">
          <label class="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              bind:checked={launchAfterPatch}
              class="w-4 h-4 bg-zinc-800 border-zinc-700 focus:ring-0 focus:ring-offset-0"
            />
            <span class="text-sm text-zinc-400">
              Launch game after patching
            </span>
          </label>
          <button
            onclick={applyPatch}
            disabled={isPatching || !gameInfo || !(patchInfo || selectedPatch)}
            class="w-full bg-emerald-600 hover:bg-emerald-500 active:bg-emerald-700 disabled:bg-zinc-700 disabled:cursor-not-allowed px-4 py-4 text-sm font-semibold uppercase tracking-wide transition-colors duration-150"
          >
            {isPatching ? "Applying Patch..." : gameInfo && (patchInfo || selectedPatch) ? "Apply Patch" : "Select Game.exe and Patch File"}
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
