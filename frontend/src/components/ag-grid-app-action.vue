<template>
    <div class="flex items-center justify-start h-full w-full">
        <Button variant="ghost" v-bind:disabled="params.userBlock" @click="rowClickEdit(params.node.rowIndex)" class="h-auto w-auto p-1 hover:bg-transparent text-muted-foreground hover:text-foreground">
            <SquarePen class="size-5" />
        </Button>
        <Button variant="ghost" v-bind:disabled="params.userBlock" @click="rowClickRemove(params.node.rowIndex)" class="h-auto w-auto p-1 hover:bg-transparent text-muted-foreground hover:text-foreground">
            <Trash class="size-5" />
        </Button>
    </div>
</template>

<script>
    import { Button } from '@/components/ui/button'
    import { SquarePen, Trash } from '@lucide/vue';
    export default {
        components: {
            Button,

            SquarePen,
            Trash,
        },
        props: {
            params: {
                type: Object,
                required: true,
            },
        },
        methods: {
            rowClickEdit(value) {
                this.params.rowClickEdit(value);
            },
            rowClickRemove(value) {
                this.params.rowClickRemove(value);
            },
        }
    };
</script>