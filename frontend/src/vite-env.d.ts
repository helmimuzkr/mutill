/// <reference types="svelte" />
/// <reference types="vite/client" />

// handle .svelte file imports. Without it, TS doesn't know what a .svelte module exports.
declare module "*.svelte" {
  import type { Component } from "svelte";
  const component: Component;
  export default component;
}
