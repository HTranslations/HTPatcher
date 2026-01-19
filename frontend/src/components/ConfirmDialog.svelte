<script lang="ts">
  export let show: boolean;
  export let title: string;
  export let message: string;
  export let confirmText: string = "Confirm";
  export let cancelText: string = "Cancel";
  export let onConfirm: () => void;
  export let onCancel: () => void;
</script>

{#if show}
  <div class="fixed inset-0 z-50 flex items-center justify-center">
    <div 
      class="absolute inset-0 bg-zinc-950/75" 
      onclick={onCancel}
      role="button"
      tabindex="0"
      aria-label="Close dialog"
    ></div>
    <div class="relative bg-zinc-900 border border-zinc-800 p-6 max-w-md w-full mx-4">
      <h3 class="text-lg font-semibold text-zinc-100 mb-2">{title}</h3>
      <p class="text-sm text-zinc-400 mb-6">{message}</p>
      <div class="flex items-center justify-between gap-3">
        <button
          onclick={onCancel}
          class="px-4 py-2 text-sm w-full bg-zinc-800 border border-zinc-700 text-zinc-300 hover:bg-zinc-700 transition-colors"
        >
          {cancelText}
        </button>
        <button
          onclick={onConfirm}
          class="px-4 py-2 text-sm w-full bg-red-600 border border-red-600 text-white hover:bg-red-700 transition-colors"
        >
          {confirmText}
        </button>
      </div>
    </div>
  </div>
{/if}





