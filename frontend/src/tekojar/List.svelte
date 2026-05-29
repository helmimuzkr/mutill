<script lang="ts">
  import { Button } from "$lib/components/ui/button";
  import { services, selectedServiceName } from "./store";
</script>

<nav class="flex flex-col w-52 h-full border-r p-2 gap-1">
  <p class="text-xs text-muted-foreground px-2 py-1">Services</p>
  {#each $services as service (service.name)}
    <Button
      class="w-full justify-start"
      size={"sm"}
      variant="ghost"
      onclick={() => selectedServiceName.set(service.name)}
    >
      {service.name}
    </Button>
  {/each}
</nav>
