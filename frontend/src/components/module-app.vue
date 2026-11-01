<template>
    <div class="grid w-full gap-2">
    <div class="flex! items-center! justify-start w-full">
        <Breadcrumb>
            <BreadcrumbList>
                <BreadcrumbItem>
                    <BreadcrumbLink as-child>
                        <RouterLink to="/">
                            {{ $t('message.home') }}
                        </RouterLink>
                    </BreadcrumbLink>
                </BreadcrumbItem>
                <BreadcrumbSeparator />
                <BreadcrumbItem>
                    <BreadcrumbPage>{{ $t('message.modules') }}</BreadcrumbPage>
                </BreadcrumbItem>
            </BreadcrumbList>
        </Breadcrumb>
    </div>

    <div class="flex! items-center! justify-end gap-3 w-full">
        <div class="grid grid-cols-3 gap-1" align="left">
            <Button
                type="button"
                @click="selectLeft()"
                :disabled="isDisabledLeft || undefined">
                <ChevronsLeft class="size-5" />
            </Button>
            <Button type="button">{{ currentPage + 1 }}</Button>
            <Button
                type="button"
                @click="selectRight()"
                :disabled="isDisabledRight || undefined">
                <ChevronsRight class="size-5" />
            </Button>
        </div>
        <div class="grid grid-cols-2 gap-1" align="right">
            <Button
                type="button"
                @click="putModuleSelect(pages[currentPage])">
                <RefreshCw class="size-5" />
            </Button>
            <Button
                type="button"
                @click="buttonClickNew"
                :disabled="userBlock || undefined">
                <Plus class="size-5" />
            </Button>
        </div>
    </div>

    <div class="flex! items-start! gap-3 w-full">
        <ag-grid-vue :key="componentKey" :rowData="tableData" :columnDefs="colDefs" :defaultColDef="defaultColDef" domLayout="autoHeight" :components="components" @grid-ready="onGridReady" style="height: 100%; width: 100%" :enableCellTextSelection="true" :ensureDomOrder="true" :loading="tableLoading" />
    </div>

    <Dialog v-model:open="showModal01">
        <DialogContent class="fixed! inset-0! top-0! left-0! z-50 flex! h-dvh! w-full! max-w-none! max-h-none! translate-x-0! translate-y-0! flex-col! gap-0! overflow-hidden! rounded-none! border-0! p-0! bg-card text-foreground! lg:inset-auto! lg:top-[50%]! lg:left-[50%]! lg:grid! lg:h-auto! lg:max-w-lg! lg:max-h-[80vh]! lg:translate-x-[-50%]! lg:translate-y-[-50%]! lg:gap-4! lg:overflow-visible! lg:rounded-lg! lg:border! lg:p-6!">
            <DialogHeader class="px-6 pt-6 pr-12 text-left lg:px-0 lg:pt-0 lg:pr-0">
                <DialogTitle>{{ $t('space.message01') }}</DialogTitle>
            </DialogHeader>
            <div class="grid min-h-0 flex-1 content-start gap-4 overflow-y-auto p-6 lg:flex-none lg:p-0">
                <div class="grid gap-3">
                    <Label for="id">{{ $t('module.id') }}</Label>
                    <Input id="id" v-model="id" name="id" disabled />
                </div>
                <Label class="flex flex-col gap-3 items-start">
                    {{ $t('module.message11') }}
                    <Select v-model="typeValue" v-on:update:model-value="typeUpdate()">
                        <SelectTrigger>
                            <SelectValue :placeholder="$t('module.message11')" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectGroup>
                                <SelectLabel>{{ $t('module.message11') }}</SelectLabel>
                                <SelectItem v-for="option in typeOptions" :key="option.value" :value="option">
                                    {{ option.text }}
                                </SelectItem>
                            </SelectGroup>
                        </SelectContent>
                    </Select>
                </Label>
                <div class="grid gap-3">
                    <Label for="address">{{ $t('message.address') }}</Label>
                    <Input id="address" v-model="address" name="address" :disabled="isDisabled" autocomplete="off" />
                </div>
                <Label class="flex flex-col gap-3 items-start">
                    {{ $t('message.status') }}
                    <Select v-model="statusValue">
                        <SelectTrigger>
                            <SelectValue :placeholder="$t('message.status')" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectGroup>
                                <SelectLabel>{{ $t('message.status') }}</SelectLabel>
                                <SelectItem v-for="option in statusOptions" :key="option.value" :value="option">
                                    {{ option.text }}
                                </SelectItem>
                            </SelectGroup>
                        </SelectContent>
                    </Select>
                </Label>
                <div class="grid gap-3">
                    <Label for="description">{{ $t('message.description') }}</Label>
                    <Input id="description" v-model="description" name="description" />
                </div>
            </div>
            <DialogFooter class="mt-auto flex-row px-6 pb-[max(1.5rem,env(safe-area-inset-bottom))] *:flex-1 lg:mt-0 lg:px-0 lg:pb-0 lg:*:flex-none">
                <DialogClose as-child>
                    <Button>{{ $t('message.cancel') }}</Button>
                </DialogClose>
                <Button @click="putModuleUpdate()">{{ $t('message.save') }}</Button>
            </DialogFooter>
        </DialogContent>
    </Dialog>

    <Dialog v-model:open="showModal02">
        <DialogContent class="fixed! inset-0! top-0! left-0! z-50 flex! h-dvh! w-full! max-w-none! max-h-none! translate-x-0! translate-y-0! flex-col! gap-0! overflow-hidden! rounded-none! border-0! p-0! bg-card text-foreground! lg:inset-auto! lg:top-[50%]! lg:left-[50%]! lg:grid! lg:h-auto! lg:max-w-lg! lg:max-h-[80vh]! lg:translate-x-[-50%]! lg:translate-y-[-50%]! lg:gap-4! lg:overflow-visible! lg:rounded-lg! lg:border! lg:p-6!">
                <DialogHeader class="px-6 pt-6 pr-12 text-left lg:px-0 lg:pt-0 lg:pr-0">
                    <DialogTitle>{{ $t('module.new') }}</DialogTitle>
                </DialogHeader>
                <div class="grid min-h-0 flex-1 content-start gap-4 overflow-y-auto p-6 lg:flex-none lg:p-0">
                    <Label class="flex flex-col gap-3 items-start">
                        {{ $t('module.message11') }}
                        <Select v-model="typeValue" v-on:update:model-value="typeUpdate()">
                            <SelectTrigger>
                                <SelectValue :placeholder="$t('module.message11')" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectGroup>
                                    <SelectLabel>{{ $t('module.message11') }}</SelectLabel>
                                    <SelectItem v-for="option in typeOptions" :key="option.value" :value="option">
                                        {{ option.text }}
                                    </SelectItem>
                                </SelectGroup>
                            </SelectContent>
                        </Select>
                    </Label>
                    <div class="grid gap-3">
                        <Label for="address">{{ $t('message.address') }}</Label>
                        <Input id="address" v-model="address" name="address" :disabled="isDisabled" autocomplete="off" />
                    </div>
                    <Label class="flex flex-col gap-3 items-start">
                        {{ $t('message.status') }}
                        <Select v-model="statusValue">
                            <SelectTrigger>
                                <SelectValue :placeholder="$t('message.status')" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectGroup>
                                    <SelectLabel>{{ $t('message.status') }}</SelectLabel>
                                    <SelectItem v-for="option in statusOptions" :key="option.value" :value="option">
                                        {{ option.text }}
                                    </SelectItem>
                                </SelectGroup>
                            </SelectContent>
                        </Select>
                    </Label>
                    <div class="grid gap-3">
                        <Label for="description">{{ $t('message.description') }}</Label>
                        <Input id="description" v-model="description" name="description" />
                    </div>
                </div>
                <DialogFooter class="mt-auto flex-row px-6 pb-[max(1.5rem,env(safe-area-inset-bottom))] *:flex-1 lg:mt-0 lg:px-0 lg:pb-0 lg:*:flex-none">
                    <DialogClose as-child>
                        <Button>{{ $t('message.cancel') }}</Button>
                    </DialogClose>
                    <Button @click="putModuleCreate()">{{ $t('message.save') }}</Button>
                </DialogFooter>
        </DialogContent>
    </Dialog>

    <Dialog v-model:open="showModal03">
        <DialogContent class="fixed! inset-0! top-0! left-0! z-50 flex! h-dvh! w-full! max-w-none! max-h-none! translate-x-0! translate-y-0! flex-col! gap-0! overflow-hidden! rounded-none! border-0! p-0! bg-card text-foreground! lg:inset-auto! lg:top-[50%]! lg:left-[50%]! lg:grid! lg:h-auto! lg:max-w-lg! lg:max-h-[80vh]! lg:translate-x-[-50%]! lg:translate-y-[-50%]! lg:gap-4! lg:overflow-visible! lg:rounded-lg! lg:border! lg:p-6!">
            <DialogHeader class="px-6 pt-6 pr-12 text-left lg:px-0 lg:pt-0 lg:pr-0">
                <DialogTitle>{{ $t('module.message12') }}</DialogTitle>
            </DialogHeader>
            <div class="flex-1 p-6 lg:flex-none lg:p-0">{{ $t('module.message08') }} {{ registrationIDRemove }}</div>
            <DialogFooter class="mt-auto flex-row px-6 pb-[max(1.5rem,env(safe-area-inset-bottom))] *:flex-1 lg:mt-0 lg:px-0 lg:pb-0 lg:*:flex-none">
                <DialogClose as-child>
                    <Button>{{ $t('message.cancel') }}</Button>
                </DialogClose>
                <Button @click="putModuleRemove(rowIndexRemove)">{{ $t('message.delete') }}</Button>
            </DialogFooter>
        </DialogContent>
    </Dialog>

    </div>
</template>

<script>
    import axios from 'axios'
    import { mapState } from 'vuex';
    import { AgGridVue } from "ag-grid-vue3"; // Vue Data Grid Component
    import ActionRenderer from './ag-grid-app-action.vue'

    import { Button } from '@/components/ui/button'
    import { Input } from '@/components/ui/input'
    import { Label } from '@/components/ui/label'
    import { Toaster } from '@/components/ui/sonner'
    import { toast } from 'vue-sonner'
    import {
        Dialog,
        DialogClose,
        DialogContent,
        DialogDescription,
        DialogFooter,
        DialogHeader,
        DialogTitle,
    } from '@/components/ui/dialog'
    import {
        Breadcrumb,
        BreadcrumbItem,
        BreadcrumbLink,
        BreadcrumbList,
        BreadcrumbPage,
        BreadcrumbSeparator,
    } from '@/components/ui/breadcrumb'
    import {
        Select,
        SelectContent,
        SelectGroup,
        SelectItem,
        SelectLabel,
        SelectTrigger,
        SelectValue,
    } from '@/components/ui/select'

    import { ChevronsLeft, ChevronsRight, RefreshCw, Plus } from '@lucide/vue';

    export default {
        components: {
            AgGridVue,
            Breadcrumb,
            BreadcrumbItem,
            BreadcrumbLink,
            BreadcrumbList,
            BreadcrumbPage,
            BreadcrumbSeparator,
            Button,
            Dialog,
            DialogClose,
            DialogContent,
            DialogDescription,
            DialogFooter,
            DialogHeader,
            DialogTitle,
            Input,
            Label,
            Toaster,
            Select,
            SelectContent,
            SelectGroup,
            SelectItem,
            SelectLabel,
            SelectTrigger,
            SelectValue,
            ChevronsLeft,
            ChevronsRight,
            RefreshCw,
            Plus,
        },
        data() {
            const statusOptions = [
                {text: this.$t('message.on'), value: 1},
                {text: this.$t('message.off'), value: 0},
            ];
            const typeOptions = [
                {text: this.$t('message.notdefined'), value: ""},
                {text: "zbxagent", value: "zbxagent"},
                {text: "prometheus", value: "prometheus"},
                {text: "uimanager", value: "uimanager"},
                {text: "datamanager", value: "datamanager"},
            ];

            return {
                isDisabled: true,
                userRole: null,
                userBlock: true,
                rowIndexRemove: null,
                registrationIDRemove: "",
                noData: this.$t('message.nodata'),
                showModal01: false,
                showModal02: false,
                showModal03: false,
                tableData: [],
                id: "",
                statusOptions: statusOptions,
                statusValue: statusOptions[0],
                typeOptions: typeOptions,
                typeValue: typeOptions[0],
                address: "",
                description: "",
                tableLoading: false,
                gridApi: null,
                componentKey: 0,

                // Pagination
                selectParam: {},
                isDisabledLeft: true,
                isDisabledRight: true,
                currentPage: 0,
                pages: [""],

                // Column Definitions: Defines the columns to be displayed.
                colDefs: [
                    { field: "id", headerName: this.$t('module.id') },
                    { field: "type", headerName: this.$t('message.type') },
                    { field: "address", headerName: this.$t('message.address') },
                    { field: "status", headerName: this.$t('message.status') },
                    { field: "description", headerName: this.$t('message.description') },
                    { 
                        field: "action",
                        headerName: this.$t('message.action'),
                        cellRenderer: 'actionRenderer',
                        cellRendererParams: {
                            rowClickEdit: this.rowClickEdit,
                            rowClickRemove: this.rowClickRemove,
                            userBlock: this.userBlock,
                        },
                    },
                ],
                defaultColDef: {
                    minWidth: 120,
                    resizable: true,
                },
                components: {
                    actionRenderer: ActionRenderer,
                },
            }
        },
        methods: {
            typeUpdate() {
                if (this.typeValue.value == "prometheus") {
                    this.isDisabled = false;
                } else {
                    this.isDisabled = true;
                }
            },
            rowClickEdit(value) {
                this.statusOptions.forEach(item => {
                    if (item.value == this.tableData[value].status) this.statusValue = item;
                });
                this.typeOptions.forEach(item => {
                    if (item.value == this.tableData[value].type) this.typeValue = item;
                });
                this.showModal01 = true;
                this.id = this.tableData[value].id;
                this.type = this.tableData[value].type;
                this.address = this.tableData[value].address;
                this.description = this.tableData[value].description;

                if (this.typeValue.value == "prometheus") {
                    this.isDisabled = false;
                } else {
                    this.isDisabled = true;
                }
            },
            buttonClickNew() {
                this.id = "";
                this.description = "";
                this.address = "";
                this.statusValue = this.statusOptions[0];
                this.typeValue = this.typeOptions[0];
                this.showModal02 = true;

                if (this.typeValue.value == "prometheus") {
                    this.isDisabled = false;
                } else {
                    this.isDisabled = true;
                }
            },
            selectRight() {
                this.currentPage++;
                this.putModuleSelect(this.selectParam.next);
                this.pages[this.currentPage] = this.selectParam.next;
            },
            selectLeft() {
                if (this.currentPage > 0) {
                    this.currentPage--;
                }
                this.putModuleSelect(this.pages[this.currentPage]);
            },
            async putModuleSelect(current) {
                this.tableLoading = true;
                try {
                    const response = await axios.put(
                        "/api/v1/admin/module/select?" + Math.random(),
                        { current: current || "", pagesize: 100 },
                    );
                    const data = response.data;
                    this.tableData = data.data ?? [];
                    this.selectParam.current = data.current;
                    this.selectParam.next = data.next;
                    this.isDisabledRight = !data.next;
                    this.isDisabledLeft = !data.current;
                } catch (error) {
                    console.error("Failed to load registrations:", error);
                } finally {
                    this.tableLoading = false;
                }
            },
            async putModuleUpdate() {
                if (this.typeValue.value == "prometheus" && (this.address == "" || this.address == undefined)) {
                    this.showModal01 = true;
                    toast(this.$t('module.message14'), {
                        type: 'error',
                        style: {
                            background: 'var(--destructive)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--destructive)'
                        }
                    });
                    return;
                }

                let dataPut = {
                    id: this.id,
                    type: this.typeValue.value,
                    address: this.address,
                    status: parseInt(this.statusValue.value),
                    description: this.description,
                };
                try {
                    await axios.put("/api/v1/admin/module/update", dataPut);
                    toast(this.$t('message.message01'), {
                        type: 'success',
                        style: {
                            background: 'var(--success)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--success)'
                        }
                    });
                    this.putModuleSelect();
                    this.showModal01 = false;
                } catch (error) {
                    if (error.response?.status === 500) {
                        toast(this.$t('message.message02'), {
                            type: 'error',
                            style: {
                                background: 'var(--destructive)',
                                color: 'var(--primary-foreground)',
                                border: '1px solid var(--destructive)'
                            }
                        });
                        this.showModal01 = true;
                    } else {
                        console.error("Failed to update registration:", error);
                    }
                }
            },
            rowClickRemove(value) {
                this.showModal03 = true;
                this.rowIndexRemove = value;
                this.registrationIDRemove = this.tableData[value].id;
            },
            async putModuleRemove(value) {
                let dataPut = {
                    id: this.tableData[value].id
                };
                try {
                    await axios.put("/api/v1/admin/module/remove", dataPut);
                    toast(this.$t('module.message02'), {
                        type: 'success',
                        style: {
                            background: 'var(--success)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--success)'
                        }
                    });
                    this.putModuleSelect();
                    this.showModal03 = false;
                } catch (error) {
                    if (error.response?.status === 500) {
                        toast(this.$t('module.message03'), {
                            type: 'error',
                            style: {
                                background: 'var(--destructive)',
                                color: 'var(--primary-foreground)',
                                border: '1px solid var(--destructive)'
                            }
                        });
                        this.showModal03 = true;
                    } else {
                        console.error("Failed to remove registration:", error);
                    }
                }
            },
            async putModuleCreate() {
                
                if (this.description == "" || this.description == undefined) {
                    this.showModal02 = true;
                    toast(this.$t('module.message04'), {
                        type: 'error',
                        style: {
                            background: 'var(--destructive)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--destructive)'
                        }
                    });
                    return;
                }
                if (this.description.length > 60) {
                    this.showModal02 = true;
                    toast(this.$t('module.message05'), {
                        type: 'error',
                        style: {
                            background: 'var(--destructive)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--destructive)'
                        }
                    });
                    return;
                }
                if (this.typeValue.value == "prometheus" && (this.address == "" || this.address == undefined)) {
                    this.showModal02 = true;
                    toast(this.$t('module.message14'), {
                        type: 'error',
                        style: {
                            background: 'var(--destructive)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--destructive)'
                        }
                    });
                    return;
                }

                let dataPut = {
                    type: this.typeValue.value,
                    address: this.address,
                    status: parseInt(this.statusValue.value),
                    description: this.description,
                };
                try {
                    await axios.put("/api/v1/admin/module/create", dataPut);
                    toast(this.$t('module.message06'), {
                        type: 'success',
                        style: {
                            background: 'var(--success)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--success)'
                        }
                    });
                    this.putModuleSelect();
                    this.showModal02 = false;
                } catch (error) {
                    if (error.response?.status === 500) {
                        toast(this.$t('module.message07'), {
                            type: 'error',
                            style: {
                                background: 'var(--destructive)',
                                color: 'var(--primary-foreground)',
                                border: '1px solid var(--destructive)'
                            }
                        });
                        this.showModal02 = true;
                    } else {
                        console.error("Failed to create registration:", error);
                    }
                }
            },
            async userInfo() {
                try {
                    const response = await axios.get("/api/v1/userinfo?" + Math.random());
                    this.userRole = response.data.role;
                    this.userBlock = this.userRole !== "superadmin" && this.userRole !== "admin";
                } catch (error) {
                    console.error("Failed to load user info:", error);
                    return;
                }

                this.colDefs = [
                    { field: "id", headerName: this.$t('module.id') },
                    { field: "type", headerName: this.$t('message.type') },
                    { field: "address", headerName: this.$t('message.address') },
                    { field: "status", headerName: this.$t('message.status') },
                    { field: "description", headerName: this.$t('message.description') },
                    { 
                        field: "action",
                        headerName: this.$t('message.action'),
                        cellRenderer: 'actionRenderer',
                        cellRendererParams: {
                            rowClickEdit: this.rowClickEdit,
                            rowClickRemove: this.rowClickRemove,
                            userBlock: this.userBlock,
                        },
                    },
                ];

                this.putModuleSelect();
            },
            onGridReady(params) {
                this.gridApi = params.api;
                this.gridApi.sizeColumnsToFit();
            },
            handleResize() {
                if (this.gridApi) {
                    this.gridApi.sizeColumnsToFit();
                }
            },
        },
        mounted() {
            window.addEventListener('resize', this.handleResize);
            this.userInfo();
        },
        beforeUnmount() {
            window.removeEventListener('resize', this.handleResize);
        },
        computed: {
            ...mapState(['locale', 'mode'])
        },
        watch: {
            locale(newValue, oldValue) {
                if (newValue != oldValue && oldValue != undefined) {
                    this.showModal01 = false;
                    this.showModal02 = false;
                    this.colDefs = [
                        { field: "id", headerName: this.$t('module.id') },
                        { field: "type", headerName: this.$t('message.type') },
                        { field: "address", headerName: this.$t('message.address') },
                        { field: "status", headerName: this.$t('message.status') },
                        { field: "description", headerName: this.$t('message.description') },
                        { 
                            field: "action",
                            headerName: this.$t('message.action'),
                            cellRenderer: 'actionRenderer',
                            cellRendererParams: {
                                rowClickEdit: this.rowClickEdit,
                                rowClickRemove: this.rowClickRemove,
                                userBlock: this.userBlock,
                            },
                        },
                    ];
                    this.statusOptions = [
                        {text: this.$t('message.on'), value: 1},
                        {text: this.$t('message.off'), value: 0},
                    ];
                    this.typeOptions = [
                        {text: this.$t('message.notdefined'), value: ""},
                        {text: "zbxagent", value: "zbxagent"},
                        {text: "prometheus", value: "prometheus"},
                        {text: "uimanager", value: "uimanager"},
                    ];
                    this.componentKey += 1;
                }
            },
            mode(newValue, oldValue) {
                this.componentKey += 1;
            },
            tableData(newValue) {
                if (newValue.length > 0 && this.gridApi) {
                    setTimeout(() => {
                        this.gridApi.sizeColumnsToFit();
                    }, 50);
                }
            },
        }
    }
</script>