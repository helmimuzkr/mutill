<script lang="ts">
  import { Badge } from "$lib/components/ui/badge";
  import { Button } from "$lib/components/ui/button";
  import {
    restartService,
    selectedService,
    startService,
    stopService,
  } from "./store";
</script>

<div class="flex justify-between items-center">
  <div class="flex items-center gap-2">
    <span class="font-bold">{$selectedService.name}</span>
    <Badge>{$selectedService.status}</Badge>
  </div>
  <div class="flex gap-2">
    <Button
      variant="default"
      size="sm"
      onclick={() => startService($selectedService.name)}>Start</Button
    >
    <Button
      variant="destructive"
      size="sm"
      onclick={() => stopService($selectedService.name)}>Stop</Button
    >
    <Button
      variant="outline"
      size="sm"
      onclick={() => restartService($selectedService.name)}>Restart</Button
    >
  </div>
</div>
