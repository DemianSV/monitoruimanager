<script setup>
import { onBeforeUnmount, ref, watch } from "vue";
import { Spinner } from "@/components/ui/spinner";

const props = defineProps({
	loading: {
		type: Boolean,
		default: false,
	},
});

const overlayVisible = ref(false);
let showTimer;

watch(
	() => props.loading,
	(loading) => {
		clearTimeout(showTimer);
		if (loading) {
			showTimer = setTimeout(() => {
				overlayVisible.value = true;
			}, 200);
		} else {
			overlayVisible.value = false;
		}
	},
);

onBeforeUnmount(() => clearTimeout(showTimer));
</script>

<template>
	<div :aria-busy="loading">
		<slot />
		<Teleport to="body">
			<Transition
				enter-active-class="transition-opacity duration-150"
				enter-from-class="opacity-0"
				leave-active-class="transition-opacity duration-150"
				leave-to-class="opacity-0">
				
                <div v-if="overlayVisible" class="fixed inset-0 z-100 grid place-items-center bg-black/50">
					<Spinner class="size-6 text-primary" />
				</div>
			</Transition>
		</Teleport>
	</div>
</template>