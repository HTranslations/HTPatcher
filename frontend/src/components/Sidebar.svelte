<script lang="ts">
  import { onMount } from "svelte";
  import { GetCurrentVersion } from "../../wailsjs/go/main/App.js";

  export let currentPage: "games" | "patches" | "settings" | "request-translation";
  export let selectedCategory: string;
  export let categories: Array<{ id: string; name: string; count: number }>;
  export let playStatusCategories: Array<{ id: string; name: string; count: number }> = [];

  let currentVersion = 0;

  onMount(async () => {
    try {
      currentVersion = await GetCurrentVersion();
    } catch (error) {
      console.error("Failed to get current version:", error);
    }
  });

  function getStatusDotColor(id: string): string {
    switch (id) {
      case "playing":
        return "bg-blue-500";
      case "on-hold":
        return "bg-amber-500";
      case "finished":
        return "bg-emerald-500";
      case "given-up":
        return "bg-red-500";
      case "unplayed":
      default:
        return "bg-zinc-600";
    }
  }
</script>

<div class="w-64 bg-zinc-900 border-r border-zinc-800 overflow-y-auto text-left flex flex-col">
  <div class="p-4 flex-1">
    <div class="mb-6">
      <h2 class="text-xs font-semibold text-zinc-500 uppercase tracking-wide mb-3">
        Navigation
      </h2>
      <div class="space-y-1">
        <button
          on:click={() => currentPage = "games"}
          class="w-full text-left px-3 py-2 text-sm border border-zinc-800 {currentPage === 'games' ? 'bg-zinc-800 text-zinc-100' : 'bg-zinc-900 text-zinc-400'}"
        >
          Games
        </button>
        <button
          on:click={() => currentPage = "patches"}
          class="w-full text-left px-3 py-2 text-sm border border-zinc-800 {currentPage === 'patches' ? 'bg-zinc-800 text-zinc-100' : 'bg-zinc-900 text-zinc-400'}"
        >
          Translation Patches
        </button>
      </div>
    </div>

    {#if currentPage === "games"}
    <div class="mb-6">
      <h2 class="text-xs font-semibold text-zinc-500 uppercase tracking-wide mb-3">
        Library
      </h2>
      <div class="space-y-1">
        {#each categories as category}
          <button
            on:click={() => selectedCategory = category.id}
            class="w-full flex items-center justify-between px-3 py-2 text-sm border border-zinc-800 {selectedCategory === category.id ? 'bg-zinc-800 text-zinc-100' : 'bg-zinc-900 text-zinc-400'}"
          >
            <span>{category.name}</span>
            <span class="text-xs text-zinc-500">{category.count}</span>
          </button>
        {/each}
      </div>
    </div>

    <div class="mb-6">
      <h2 class="text-xs font-semibold text-zinc-500 uppercase tracking-wide mb-3">
        Play Status
      </h2>
      <div class="space-y-1">
        {#each playStatusCategories as status}
          <button
            on:click={() => selectedCategory = status.id}
            class="w-full flex items-center justify-between px-3 py-2 text-sm border border-zinc-800 {selectedCategory === status.id ? 'bg-zinc-800 text-zinc-100' : 'bg-zinc-900 text-zinc-400'}"
          >
            <span class="flex items-center gap-2">
              <span class="w-2 h-2 rounded-full {getStatusDotColor(status.id)}"></span>
              {status.name}
            </span>
            <span class="text-xs text-zinc-500">{status.count}</span>
          </button>
        {/each}
      </div>
    </div>
    {/if}

    <div>
      <h2 class="text-xs font-semibold text-zinc-500 uppercase tracking-wide mb-3">
        Quick Actions
      </h2>
      <div class="space-y-1">
        <button
          on:click={() => currentPage = "settings"}
          class="w-full text-left px-3 py-2 text-sm border border-zinc-800 {currentPage === 'settings' ? 'bg-zinc-800 text-zinc-100' : 'bg-zinc-900 text-zinc-400'}"
        >
          Settings
        </button>
        <button
          on:click={() => currentPage = "request-translation"}
          class="w-full text-left px-3 py-2 text-sm border border-zinc-800 {currentPage === 'request-translation' ? 'bg-zinc-800 text-zinc-100' : 'bg-zinc-900 text-zinc-400'}"
        >
          Request Translation
        </button>
      </div>
    </div>
  </div>

  <!-- Version at bottom -->
  <div class="p-4 border-t border-zinc-800">
    <span class="text-xs text-zinc-500">v{currentVersion}</span>
  </div>
</div>
