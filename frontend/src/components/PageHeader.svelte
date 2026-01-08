<script lang="ts">
  export let title: string;
  export let description: string;
</script>

<div class="mb-8">
  <div class="flex items-center gap-6 mb-3 text-left">
    <div class="h-20 w-1 bg-emerald-500"></div>
    <div>
      <h2 class="text-3xl font-bold mb-1">{title}</h2>
      <p class="text-sm text-zinc-400">
        {description}
      </p>
    </div>
  </div>
</div>

