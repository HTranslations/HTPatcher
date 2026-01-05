<script lang="ts">
  import { main } from "../../wailsjs/go/models.js";
  import PageHeader from "./PageHeader.svelte";
  import GameCard from "./GameCard.svelte";
  
  export let games: main.LocatedGame[];
  export let onAddGame: () => void;
  export let onOpenFolder: (game: main.LocatedGame) => void;
  export let onLaunchGame: (game: main.LocatedGame) => void;
  export let onLocateGame: (game: main.LocatedGame) => void;
  export let onDeleteGame: (game: main.LocatedGame) => void;
  export let onTranslateGame: (game: main.LocatedGame) => void;
</script>

<PageHeader
  title="Your Games"
  description="Select a game to view available translations and apply patches"
/>

<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
  <!-- Add Game Card -->
  <button
    onclick={onAddGame}
    class="bg-zinc-900 border border-zinc-800 text-left"
  >
    <div class="aspect-4/3 bg-zinc-800 flex items-center justify-center">
      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-12 h-12 text-zinc-500">
        <path stroke-linecap="round" stroke-linejoin="round" d="M12 4.5v15m7.5-7.5h-15" />
      </svg>
    </div>
    <div class="p-4">
      <h3 class="font-medium text-sm mb-1">Add Game</h3>
      <p class="text-xs text-zinc-500 mb-3">Locate game from disk</p>
      <div class="w-full bg-zinc-800 border border-zinc-700 text-zinc-300 text-sm px-3 py-2 text-center">
        Select Game to Add
      </div>
    </div>
  </button>

  {#each games as game}
    <GameCard
      {game}
      onOpenGame={onTranslateGame}
      {onOpenFolder}
      {onLaunchGame}
      {onLocateGame}
      {onDeleteGame}
    />
  {/each}
</div>

