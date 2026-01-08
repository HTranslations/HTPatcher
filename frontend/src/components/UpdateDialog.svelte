<script lang="ts">
  import { domain } from "../../wailsjs/go/models.js";
  import {
    DownloadUpdate,
    ApplyUpdate,
  } from "../../wailsjs/go/main/App.js";
  import { EventsOn, EventsOff, BrowserOpenURL } from "../../wailsjs/runtime/runtime.js";
  import { onDestroy } from "svelte";

  export let show: boolean;
  export let releaseInfo: domain.ReleaseInfo | null = null;
  export let onClose: () => void;

  // Download states: 'idle' | 'downloading' | 'ready' | 'error'
  let downloadState: 'idle' | 'downloading' | 'ready' | 'error' = 'idle';
  let downloadProgress = 0;
  let downloadedBytes = 0;
  let totalBytes = 0;
  let errorMessage = '';

  $: exeAsset = releaseInfo?.assets?.find((asset) =>
    asset.name.endsWith(".exe")
  );

  // Reset state when dialog is shown
  $: if (show) {
    downloadState = 'idle';
    downloadProgress = 0;
    downloadedBytes = 0;
    totalBytes = 0;
    errorMessage = '';
  }

  function formatFileSize(bytes: number): string {
    if (bytes === 0) return "0 Bytes";
    const k = 1024;
    const sizes = ["Bytes", "KB", "MB", "GB"];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round((bytes / Math.pow(k, i)) * 100) / 100 + " " + sizes[i];
  }

  // Parse progress from log messages
  function parseProgress(message: string): boolean {
    // Format: UPDATE_PROGRESS:downloaded:total:percentage
    if (message.startsWith('UPDATE_PROGRESS:')) {
      const parts = message.split(':');
      if (parts.length >= 4) {
        downloadedBytes = parseInt(parts[1], 10);
        totalBytes = parseInt(parts[2], 10);
        downloadProgress = parseFloat(parts[3]);
        return true;
      }
    }
    return false;
  }

  // Listen to log events for progress updates
  let unsubscribe: (() => void) | null = null;

  function setupProgressListener() {
    unsubscribe = EventsOn('log', (logMessage: { message: string; type: string }) => {
      if (downloadState === 'downloading') {
        parseProgress(logMessage.message);
      }
    });
  }

  function cleanupProgressListener() {
    if (unsubscribe) {
      EventsOff('log');
      unsubscribe = null;
    }
  }

  onDestroy(() => {
    cleanupProgressListener();
  });

  async function downloadAndInstall() {
    if (!releaseInfo) return;

    try {
      downloadState = 'downloading';
      downloadProgress = 0;
      downloadedBytes = 0;
      totalBytes = exeAsset?.size || 0;
      errorMessage = '';

      setupProgressListener();

      await DownloadUpdate(releaseInfo);

      cleanupProgressListener();
      downloadState = 'ready';
    } catch (err: any) {
      cleanupProgressListener();
      downloadState = 'error';
      errorMessage = err?.message || 'Download failed';
    }
  }

  async function installNow() {
    try {
      await ApplyUpdate(); // This will exit the app
    } catch (err: any) {
      errorMessage = err?.message || 'Installation failed';
      downloadState = 'error';
    }
  }

  function handleManualDownload() {
    if (releaseInfo?.html_url) {
      BrowserOpenURL(releaseInfo.html_url);
    }
    onClose();
  }

  function handleClose() {
    if (downloadState === 'downloading') {
      // Don't allow closing during download
      return;
    }
    onClose();
  }
</script>

{#if show && releaseInfo && exeAsset}
  <div class="fixed inset-0 z-50 flex items-center justify-center">
    <div
      class="absolute inset-0 bg-zinc-950/75"
      onclick={handleClose}
      role="button"
      tabindex="0"
      aria-label="Close dialog"
    ></div>
    <div class="relative bg-zinc-900 border border-zinc-800 p-6 max-w-lg w-full mx-4">
      <h3 class="text-lg font-semibold text-zinc-100 mb-2">
        {#if downloadState === 'ready'}
          Ready to Install
        {:else if downloadState === 'downloading'}
          Downloading Update...
        {:else if downloadState === 'error'}
          Update Failed
        {:else}
          New Version Available
        {/if}
      </h3>

      <div class="text-sm text-zinc-400 mb-4 space-y-2">
        {#if downloadState === 'idle'}
          <p>
            A new version ({releaseInfo.tag_name}) is available for download.
          </p>
          <div class="bg-zinc-800/50 border border-zinc-700 p-3 mt-3 text-start">
            <p class="text-zinc-300 font-medium mb-3">Download Information</p>
            <p class="text-zinc-400 text-xs mb-2">
              File: <span class="text-zinc-300">{exeAsset.name}</span>
            </p>
            <p class="text-zinc-400 text-xs">
              Size: <span class="text-zinc-300">{formatFileSize(exeAsset.size)}</span>
            </p>
          </div>

        {:else if downloadState === 'downloading'}
          <div class="space-y-3">
            <p class="text-zinc-300">
              Downloading {releaseInfo.tag_name}...
            </p>
            <div class="w-full bg-zinc-800 h-2 overflow-hidden">
              <div 
                class="bg-blue-500 h-full transition-all duration-300"
                style="width: {downloadProgress}%"
              ></div>
            </div>
            <p class="text-zinc-400 text-xs">
              {formatFileSize(downloadedBytes)} / {formatFileSize(totalBytes)} ({downloadProgress.toFixed(1)}%)
            </p>
          </div>

        {:else if downloadState === 'ready'}
          <div class="space-y-3">
            <p class="text-zinc-300">
              Download complete! Click "Install Now" to update.
            </p>
            <p class="text-zinc-400 text-xs">
              The application will restart automatically.
            </p>
          </div>

        {:else if downloadState === 'error'}
          <div class="space-y-3">
            <p class="text-red-400">
              {errorMessage}
            </p>
            <p class="text-zinc-400 text-xs">
              You can try again or download manually from the releases page.
            </p>
          </div>
        {/if}
      </div>

      <div class="flex items-center justify-between gap-3">
        {#if downloadState === 'idle'}
          <button
            onclick={handleClose}
            class="cursor-pointer px-4 py-2 text-sm w-full bg-zinc-800 border border-zinc-700 text-zinc-300 hover:bg-zinc-700 transition-colors"
          >
            Later
          </button>
          <button
            onclick={downloadAndInstall}
            class="cursor-pointer px-4 py-2 text-sm w-full bg-blue-600 border border-blue-600 text-white hover:bg-blue-700 transition-colors"
          >
            Download & Install
          </button>

        {:else if downloadState === 'downloading'}
          <button
            disabled
            class="cursor-not-allowed px-4 py-2 text-sm w-full bg-zinc-800 border border-zinc-700 text-zinc-500"
          >
            Downloading...
          </button>

        {:else if downloadState === 'ready'}
          <button
            onclick={handleClose}
            class="cursor-pointer px-4 py-2 text-sm w-full bg-zinc-800 border border-zinc-700 text-zinc-300 hover:bg-zinc-700 transition-colors"
          >
            Later
          </button>
          <!-- svelte-ignore a11y_autofocus -->
          <button
            onclick={installNow}
            autofocus
            class="cursor-pointer px-4 py-2 text-sm w-full bg-green-600 border border-green-600 text-white hover:bg-green-700 transition-colors"
          >
            Install Now
          </button>

        {:else if downloadState === 'error'}
          <button
            onclick={handleManualDownload}
            class="cursor-pointer px-4 py-2 text-sm w-full bg-zinc-800 border border-zinc-700 text-zinc-300 hover:bg-zinc-700 transition-colors"
          >
            Download Manually
          </button>
          <button
            onclick={downloadAndInstall}
            class="cursor-pointer px-4 py-2 text-sm w-full bg-blue-600 border border-blue-600 text-white hover:bg-blue-700 transition-colors"
          >
            Retry
          </button>
        {/if}
      </div>
    </div>
  </div>
{/if}
