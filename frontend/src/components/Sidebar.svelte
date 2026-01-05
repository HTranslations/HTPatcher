<script lang="ts">
  export let currentPage: "games" | "patches";
  export let selectedCategory: string;
  export let categories: Array<{ id: string; name: string; count: number }>;
</script>

<div class="w-64 bg-zinc-900 border-r border-zinc-800 overflow-y-auto text-left">
  <div class="p-4">
    <div class="mb-6">
      <h2 class="text-xs font-semibold text-zinc-500 uppercase tracking-wide mb-3">
        Navigation
      </h2>
      <div class="space-y-1">
        <button
          onclick={() => currentPage = "games"}
          class="w-full text-left px-3 py-2 text-sm border border-zinc-800 {currentPage === 'games' ? 'bg-zinc-800 text-zinc-100' : 'bg-zinc-900 text-zinc-400'}"
        >
          Games
        </button>
        <button
          onclick={() => currentPage = "patches"}
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
            onclick={() => selectedCategory = category.id}
            class="w-full flex items-center justify-between px-3 py-2 text-sm border border-zinc-800 {selectedCategory === category.id ? 'bg-zinc-800 text-zinc-100' : 'bg-zinc-900 text-zinc-400'}"
          >
            <span>{category.name}</span>
            <span class="text-xs text-zinc-500">{category.count}</span>
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
        <button class="w-full text-left px-3 py-2 text-sm bg-zinc-900 border border-zinc-800 text-zinc-400">
          Settings
        </button>
        <button class="w-full text-left px-3 py-2 text-sm bg-zinc-900 border border-zinc-800 text-zinc-400">
          Request Translation
        </button>
      </div>
    </div>
  </div>
</div>

