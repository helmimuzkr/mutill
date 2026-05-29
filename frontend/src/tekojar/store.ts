import { derived, writable } from "svelte/store";
import type { Service } from "./type";

import { EventsOn } from "../../wailsjs/runtime/runtime";
import { GetAll, Get, Start, Stop } from "../../wailsjs/go/backend/TekojarApp"

export const services = writable<Service[]>([]);
export const selectedServiceName = writable<string | null>(null);

export const selectedService = derived([services, selectedServiceName], ([$services, $selectedServiceName]) => {
  return $services.find(s => s.name === $selectedServiceName) ?? null
})

const logListeners = new Map<string, () => void>();

export async function getAllServices() {
  try {
    const result = await GetAll();
    services.set(result);
  } catch (err) {
    console.log(err instanceof Error ? err.message : "Failed to fetch services");
  }
}

export async function getServices(name: string) {
  try {
    const result = await Get(name);
    return result
  } catch (err) {
    console.log(err instanceof Error ? err.message : "Failed to fetch services");
  }
}

export async function startService(name: string) {
  await Start(name);
  subscribeServiceLogs(name)
  await getServices(name).then(res => {
    services.update(curr =>
      curr.map(s => s.name === name ? res : s)
    )
  })
}

export async function stopService(name: string) {
  await Stop(name);
  unsubscribeServiceLogs(name)
  await getServices(name).then(res => {
    services.update(curr =>
      curr.map(s => s.name === name ? res : s)
    )
  })
}

export async function restartService(name: string) {
  await stopService(name)
  await startService(name)
}

// EventsOn registers a listener for the "service:log" event from Go. Returns an unlisten function to stop listening.
export function subscribeLogs(name: string) {
  const unlisten = EventsOn("service:log", (data: any) => {
    if (data.name === name) {
      services.update(current =>
        current.map(s => s.name === name
          ? { ...s, logs: [...s.logs, data.logView] } as Service
          : s
        )
      );
    }
  });
  return unlisten; // cleanup function
}

export function subscribeServiceLogs(name: string) {
  if (logListeners.has(name)) return; // already subscribed

  const unlisten = subscribeLogs(name);
  logListeners.set(name, unlisten);
}

export function unsubscribeServiceLogs(name: string) {
  const unlisten = logListeners.get(name);
  if (unlisten) {
    unlisten();
    logListeners.delete(name);
  }
}


