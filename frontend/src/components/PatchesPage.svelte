<script lang="ts">
  import { main } from "../../wailsjs/go/models.js";
  import PageHeader from "./PageHeader.svelte";
  import PatchCard from "./PatchCard.svelte";
  
  export let patches: main.PatchEntry[];
</script>

<PageHeader
  title="New Translations"
  description="Latest translation patches released"
/>

<div class="space-y-4">
  {#each patches as patch, index}
    <PatchCard {patch} {index} />
  {/each}

  {#if patches.length === 0}
    <div class="text-center py-12">
      <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor" class="w-16 h-16 text-zinc-600 mx-auto mb-4">
        <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z" />
      </svg>
      <p class="text-zinc-400 text-lg mb-2">No patches available</p>
      <p class="text-zinc-500 text-sm">Check back later for new translations</p>
    </div>
  {/if}
</div>

