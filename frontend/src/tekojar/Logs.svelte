<script lang="ts">
  import { ScrollArea } from "$lib/components/ui/scroll-area";
  import { selectedService } from "./store";
</script>

<ScrollArea class="flex-1 rounded border font-mono text-sm p-3">
  {#each $selectedService?.logs ?? [] as serviceLog}
    {#if serviceLog.is_error}
      <p class="text-destructive">{serviceLog.log}</p>
    {:else}
      <p>{serviceLog.log}</p>
    {/if}
  {/each}
</ScrollArea>
