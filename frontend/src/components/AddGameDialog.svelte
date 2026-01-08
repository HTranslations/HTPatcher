<script lang="ts">
  import { domain } from "../../wailsjs/go/models.js";
  import { getDlsiteImageUrl } from "../lib/utils.js";
  
  export let show: boolean;
  export let locatedGame: domain.LocatedGame | null;
  export let rjCode: string;
  export let friendlyName: string;
  export let tags: string;
  export let previewImageUrl: string;
  export let previewImageLoaded: boolean;
  
  export let onClose: () => void;
  export let onRjCodeInput: (event: Event) => void;
  export let onFriendlyNameInput: (event: Event) => void;
  export let onTagsInput: (event: Event) => void;
  export let onPreviewImageLoad: () => void;
  export let onPreviewImageError: () => void;
  export let onAddGame: () => void;
</script>

{#if show && locatedGame}
  <div class="fixed inset-0 z-50 flex">
    <div 
      class="flex-1 bg-zinc-950/75" 
      onclick={onClose}
      onkeydown={(e) => e.key === 'Escape' && onClose()}
      role="button"
      tabindex="0"
      aria-label="Close drawer"
    ></div>
    <div class="w-full max-w-lg bg-zinc-900 border-l border-zinc-800 h-full overflow-y-auto">
      <div class="p-6">
        <div class="flex items-center justify-between mb-6">
          <h2 class="text-xl font-semibold">Add Game to Collection</h2>
          <button
            onclick={onClose}
            class="text-zinc-400"
            aria-label="Close drawer"
          >
            <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <div class="space-y-4 text-left">
          <!-- Game Info -->
          <div class="space-y-4">
            <div>
              <label for="game-dir-input" class="block text-sm font-medium text-zinc-400 mb-2">
                Game Directory
              </label>
              <input
                id="game-dir-input"
                type="text"
                value={locatedGame.gameDir}
                readonly
                class="w-full bg-zinc-800 border border-zinc-700 px-4 py-2 text-sm text-zinc-300 font-mono focus:outline-none focus:border-zinc-600"
              />
            </div>
            <div>
              <label for="exe-path-input" class="block text-sm font-medium text-zinc-400 mb-2">
                Executable Path
              </label>
              <input
                id="exe-path-input"
                type="text"
                value={locatedGame.exePath}
                readonly
                class="w-full bg-zinc-800 border border-zinc-700 px-4 py-2 text-sm text-zinc-300 font-mono focus:outline-none focus:border-zinc-600"
              />
            </div>
          </div>

          <!-- Friendly Name Input -->
          <div>
            <label for="friendly-name-input" class="block text-sm font-medium text-zinc-400 mb-2">
              Friendly Name <span class="text-red-400">*</span>
            </label>
            <input
              id="friendly-name-input"
              type="text"
              value={friendlyName}
              oninput={onFriendlyNameInput}
              placeholder="My Favorite Game"
              class="w-full bg-zinc-800 border border-zinc-700 px-4 py-2 text-sm text-zinc-300 focus:outline-none focus:border-zinc-600"
            />
          </div>

          <!-- RJ Code Input -->
          <div>
            <label for="rj-code-input" class="block text-sm font-medium text-zinc-400 mb-2">
              RJ Code <span class="text-red-400">*</span>
            </label>
            <input
              id="rj-code-input"
              type="text"
              value={rjCode}
              oninput={onRjCodeInput}
              placeholder="RJ00000000"
              class="w-full bg-zinc-800 border border-zinc-700 px-4 py-2 text-sm text-zinc-300 font-mono focus:outline-none focus:border-zinc-600"
            />
          </div>

          <!-- Tags Input -->
          <div>
            <label for="tags-input" class="block text-sm font-medium text-zinc-400 mb-2">
              Tags <span class="text-zinc-500">(optional)</span>
            </label>
            <input
              id="tags-input"
              type="text"
              value={tags}
              oninput={onTagsInput}
              placeholder="RPG, Fantasy, Action (comma-separated)"
              class="w-full bg-zinc-800 border border-zinc-700 px-4 py-2 text-sm text-zinc-300 focus:outline-none focus:border-zinc-600"
            />
            <p class="text-xs text-zinc-500 mt-1">Separate tags with commas</p>
          </div>

          <!-- Preview Image -->
          <div>
            <span class="block text-sm font-medium text-zinc-400 mb-2">Preview</span>
            <div class="bg-zinc-800 border border-zinc-700 aspect-[4/3] relative">
              {#if previewImageUrl && previewImageLoaded}
                <img
                  src={previewImageUrl}
                  alt="Game preview"
                  class="w-full h-full object-cover"
                  onload={onPreviewImageLoad}
                  onerror={onPreviewImageError}
                />
              {:else}
                <div class="w-full h-full flex items-center justify-center">
                  <div class="text-center">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-12 h-12 text-zinc-600 mx-auto mb-2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M2.25 15.75l5.159-5.159a2.25 2.25 0 013.182 0l5.159 5.159m-1.5-1.5l1.409-1.409a2.25 2.25 0 013.182 0l2.909 2.909m-18 3.75h16.5a1.5 1.5 0 001.5-1.5V6a1.5 1.5 0 00-1.5-1.5H3.75A1.5 1.5 0 002.25 6v12a1.5 1.5 0 001.5 1.5zm10.5-11.25h.008v.008h-.008V8.25zm.375 0a.375.375 0 11-.75 0 .375.375 0 01.75 0z" />
                    </svg>
                    <p class="text-sm text-zinc-500">No preview available</p>
                  </div>
                </div>
              {/if}
              {#if previewImageUrl && !previewImageLoaded}
                <img
                  src={previewImageUrl}
                  alt=""
                  class="hidden"
                  onload={onPreviewImageLoad}
                  onerror={onPreviewImageError}
                />
              {/if}
            </div>
          </div>

          <!-- Actions -->
          <div class="flex items-center justify-end gap-3 pt-4 border-t border-zinc-800">
            <button
              onclick={onClose}
              class="px-4 py-2 text-sm bg-zinc-800 border border-zinc-700 text-zinc-300"
            >
              Cancel
            </button>
            <button
              onclick={onAddGame}
              disabled={!rjCode || rjCode.length < 3 || !friendlyName || friendlyName.trim().length === 0}
              class="px-4 py-2 text-sm bg-emerald-600 text-white disabled:bg-zinc-700 disabled:text-zinc-500 disabled:cursor-not-allowed"
            >
              Add Game
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
{/if}

