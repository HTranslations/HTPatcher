<script lang="ts">
  import { domain } from "../../wailsjs/go/models.js";

  export let show: boolean;
  export let game: domain.LocatedGame | null;
  export let onClose: () => void;
  export let onSave: (friendlyName: string, tags: string[], playStatus: string) => void;
  export let onDelete: (game: domain.LocatedGame) => void;
  export let onExport: (game: domain.LocatedGame) => void;

  let friendlyName = "";
  let tagsInput = "";
  let playStatus = "unplayed";

  // Update local state when game prop changes
  $: if (game && show) {
    friendlyName = game.friendlyName || "";
    tagsInput = game.tags?.join(", ") || "";
    playStatus = (game as any).playStatus || "unplayed";
  }

  function handleSave() {
    const tagArray = tagsInput
      .split(",")
      .map((tag) => tag.trim())
      .filter((tag) => tag.length > 0);

    onSave(friendlyName, tagArray, playStatus);
  }

  function handleCancel() {
    onClose();
  }

  function handleDelete() {
    if (game) {
      onDelete(game);
      onClose();
    }
  }

  function handleExport() {
    if (game) {
      onExport(game);
    }
  }

  $: canSave = friendlyName.trim().length > 0;
</script>

{#if show && game}
  <div class="fixed inset-0 z-50 flex text-left">
    <div
      class="flex-1 bg-zinc-950/75"
      on:click={handleCancel}
      role="button"
      tabindex="0"
      aria-label="Close drawer"
    ></div>
    <div
      class="w-full max-w-2xl bg-zinc-900 border-l border-zinc-800 h-full overflow-hidden flex flex-col"
    >
      <!-- Header -->
      <div
        class="bg-zinc-900 border-b border-zinc-800 px-6 py-4 flex items-center justify-between"
      >
        <h2 class="text-xl font-semibold">Edit Game Details</h2>
        <button
          on:click={handleCancel}
          class="text-zinc-400 hover:text-zinc-300"
          aria-label="Close drawer"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            stroke="currentColor"
            class="w-6 h-6"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M6 18L18 6M6 6l12 12"
            />
          </svg>
        </button>
      </div>

      <!-- Main Content -->
      <div class="flex-1 overflow-y-auto p-6">
        <div class="flex flex-col gap-6">
          <!-- Game Info Display -->
          <div class="flex flex-col gap-3">
            <span
              class="text-sm font-medium text-zinc-400 uppercase tracking-wide"
            >
              Game Information
            </span>
            <div class="bg-zinc-800 border border-zinc-700 px-4 py-3">
              <div class="flex flex-col gap-2">
                <div class="flex items-center gap-2">
                  <span class="text-xs text-zinc-500 w-20">RJ Code:</span>
                  <span class="text-sm text-zinc-300 font-mono"
                    >{game.rjCode || "N/A"}</span
                  >
                </div>
                <div class="flex items-center gap-2">
                  <span class="text-xs text-zinc-500 w-20">Path:</span>
                  <span class="text-xs text-zinc-400 font-mono truncate"
                    >{game.exePath || game.gameDir}</span
                  >
                </div>
              </div>
            </div>
          </div>

          <!-- Friendly Name Input -->
          <div class="flex flex-col gap-3">
            <label
              for="friendlyName"
              class="text-sm font-medium text-zinc-400 uppercase tracking-wide"
            >
              Friendly Name *
            </label>
            <input
              id="friendlyName"
              type="text"
              bind:value={friendlyName}
              placeholder="Enter a friendly name for the game"
              class="bg-zinc-800 border border-zinc-700 px-4 py-3 text-sm text-zinc-300 focus:outline-none focus:border-zinc-600 transition-colors"
            />
            <p class="text-xs text-zinc-500">
              The display name shown in your games collection.
            </p>
          </div>

          <!-- Tags Input -->
          <div class="flex flex-col gap-3">
            <label
              for="tags"
              class="text-sm font-medium text-zinc-400 uppercase tracking-wide"
            >
              Tags
            </label>
            <input
              id="tags"
              type="text"
              bind:value={tagsInput}
              placeholder="e.g., RPG, Translated, Favorite"
              class="bg-zinc-800 border border-zinc-700 px-4 py-3 text-sm text-zinc-300 focus:outline-none focus:border-zinc-600 transition-colors"
            />
            <p class="text-xs text-zinc-500">
              Comma-separated tags to help organize and search for games.
            </p>

            <!-- Tags Preview -->
            {#if tagsInput.trim()}
              <div class="flex flex-wrap gap-2 mt-2">
                {#each tagsInput
                  .split(",")
                  .map((t) => t.trim())
                  .filter((t) => t) as tag}
                  <span
                    class="text-xs bg-zinc-800 border border-zinc-700 text-zinc-400 px-2 py-1"
                  >
                    {tag}
                  </span>
                {/each}
              </div>
            {/if}
          </div>

          <!-- Play Status Dropdown -->
          <div class="flex flex-col gap-3">
            <label
              for="playStatus"
              class="text-sm font-medium text-zinc-400 uppercase tracking-wide"
            >
              Play Status
            </label>
            <select
              id="playStatus"
              bind:value={playStatus}
              class="bg-zinc-800 border border-zinc-700 px-4 py-3 text-sm text-zinc-300 focus:outline-none focus:border-zinc-600 transition-colors"
            >
              <option value="unplayed">Unplayed</option>
              <option value="playing">Playing</option>
              <option value="on-hold">On Hold</option>
              <option value="finished">Finished</option>
              <option value="given-up">Given Up</option>
            </select>
            <p class="text-xs text-zinc-500">
              Track your progress with this game.
            </p>
          </div>

          <!-- Export Patched Files (only shown when game is patched) -->
          {#if game.translated}
            <div class="flex flex-col gap-3">
              <span
                class="text-sm font-medium text-zinc-400 uppercase tracking-wide"
              >
                Export
              </span>
              <button
                on:click={handleExport}
                class="flex items-center justify-center gap-2 bg-zinc-800 hover:bg-zinc-700 border border-zinc-700 text-zinc-300 px-4 py-3 text-sm font-semibold uppercase tracking-wide transition-colors"
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke-width="1.5"
                  stroke="currentColor"
                  class="w-5 h-5"
                >
                  <path
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    d="M3 16.5v2.25A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75V16.5m-13.5-9L12 3m0 0l4.5 4.5M12 3v13.5"
                  />
                </svg>
                Export Patched Files as ZIP
              </button>
              <p class="text-xs text-zinc-500">
                Export only the patched files (translations) as a ZIP archive.
              </p>
            </div>
          {/if}
        </div>
      </div>

      <!-- Footer Actions -->
      <div class="border-t border-zinc-800 p-6 flex gap-3 flex-shrink-0">
        <button
          on:click={handleDelete}
          class="bg-red-600 hover:bg-red-500 border border-red-600 text-white px-4 py-3 text-sm font-semibold uppercase tracking-wide transition-colors"
          title="Delete game"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            stroke-width="1.5"
            stroke="currentColor"
            class="w-5 h-5"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0"
            />
          </svg>
        </button>
        <button
          on:click={handleCancel}
          class="flex-1 bg-zinc-800 hover:bg-zinc-700 border border-zinc-700 text-zinc-300 px-4 py-3 text-sm font-semibold uppercase tracking-wide transition-colors"
        >
          Cancel
        </button>
        <button
          on:click={handleSave}
          disabled={!canSave}
          class="flex-1 bg-emerald-600 hover:bg-emerald-500 disabled:bg-zinc-700 disabled:cursor-not-allowed text-white px-4 py-3 text-sm font-semibold uppercase tracking-wide transition-colors"
        >
          Save Changes
        </button>
      </div>
    </div>
  </div>
{/if}


