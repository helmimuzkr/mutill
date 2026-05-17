<!-- ServiceManager.svelte -->
<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { Badge } from "$lib/components/ui/badge";
  import { ScrollArea } from "$lib/components/ui/scroll-area";

  let selected = $state("user-service");
</script>

<div class="flex h-full">
  <nav class="flex flex-col w-52 h-full border-r p-2 gap-1">
    <p class="text-xs text-muted-foreground px-2 py-1">Services</p>
    <button
      class="flex items-center gap-2 px-2 py-1 rounded hover:bg-accent text-sm"
      onclick={() => (selected = "user-service")}
    >
      <Badge variant="success">●</Badge>
      user-service
    </button>
  </nav>

  <div class="flex flex-col flex-1 p-4 gap-3">
    <div class="flex justify-between items-center">
      <div class="flex items-center gap-2">
        <span class="font-bold">{selected}</span>
        <Badge>Running</Badge>
      </div>
      <div class="flex gap-2">
        <Button variant="default" size="sm">Start</Button>
        <Button variant="destructive" size="sm">Stop</Button>
        <Button variant="outline" size="sm">Restart</Button>
      </div>
    </div>

    <ScrollArea class="flex-1 rounded border font-mono text-sm p-3">
      <p>[INFO] Started on port 8081</p>
      <p class="text-destructive">[ERROR] DB connection refused</p>
    </ScrollArea>
  </div>
</div>
