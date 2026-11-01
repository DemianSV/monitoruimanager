<script setup>
import { cn } from "@/lib/utils";

const props = defineProps({
  disabled: { type: Boolean, required: false, default: false },
  class: {
    type: [Boolean, null, String, Object, Array],
    required: false,
    skipCheck: true,
  },
});
</script>

<template>
  <div
    data-slot="card"
    :aria-disabled="props.disabled || undefined"
    :data-disabled="props.disabled || undefined"
    :class="
      cn(
        'bg-card text-card-foreground flex flex-col gap-6 rounded-xl border py-6 shadow-sm',
        props.disabled && 'pointer-events-none cursor-not-allowed opacity-50',
        props.class,
      )
    "
  >
    <slot />
  </div>
</template>
