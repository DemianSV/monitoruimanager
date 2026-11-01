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
                    <BreadcrumbPage>{{ $t('message.tasks') }}</BreadcrumbPage>
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
                @click="putTaskSelect(pages[currentPage])">
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
                <DialogTitle>{{ $t('task.message21') }}</DialogTitle>
            </DialogHeader>
            <div class="grid min-h-0 flex-1 content-start gap-4 overflow-y-auto p-6 lg:flex-none lg:p-0">
            <Label class="flex flex-col gap-3 items-start">
                {{ $t('message.module') }}
                <Select v-model="moduleValue" disabled>
                    <SelectTrigger>
                        <SelectValue :placeholder="$t('module.message02')" />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectGroup>
                            <SelectLabel>{{ $t('message.module') }}</SelectLabel>
                            <SelectItem v-for="option in moduleOptions" :key="option.value" :value="option">
                                {{ option.text }}
                            </SelectItem>
                        </SelectGroup>
                    </SelectContent>
                </Select>
            </Label>
            <Label class="flex flex-col gap-3 items-start">
                {{ $t('message.space') }}
                <Select v-model="spaceValue" disabled>
                    <SelectTrigger>
                        <SelectValue :placeholder="$t('module.message01')" />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectGroup>
                            <SelectLabel>{{ $t('message.space') }}</SelectLabel>
                            <SelectItem v-for="option in spaceOptions" :key="option.value" :value="option">
                                {{ option.text }}
                            </SelectItem>
                        </SelectGroup>
                    </SelectContent>
                </Select>
            </Label>
            <Label class="flex flex-col gap-3 items-start">
                {{ $t('task.message04') }}
                <Select v-model="objectValue">
                    <SelectTrigger>
                        <SelectValue :placeholder="$t('module.message03')" />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectGroup>
                            <SelectLabel>{{ $t('task.message04') }}</SelectLabel>
                            <SelectItem v-for="option in objectOptions" :key="option.value" :value="option">
                                {{ option.text }}
                            </SelectItem>
                        </SelectGroup>
                    </SelectContent>
                </Select>
            </Label>
            <Label class="flex flex-col gap-3 items-start">
                {{ $t('message.dm') }}
                <Select v-model="dmValue" disabled>
                    <SelectTrigger>
                        <SelectValue :placeholder="$t('task.message24')" />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectGroup>
                            <SelectLabel>{{ $t('message.dm') }}</SelectLabel>
                            <SelectItem v-for="option in dmOptions" :key="option.value" :value="option">
                                {{ option.text }}
                            </SelectItem>
                        </SelectGroup>
                    </SelectContent>
                </Select>
            </Label>
            <div class="grid gap-3">
                <Label for="metric">{{ $t('task.message06') }}</Label>
                <Input id="metric" v-model="metric" name="metric" disabled />
            </div>
            <Label class="flex flex-col gap-3 items-start">
                {{ $t('task.message07') }}
                <Select v-model="datatypeValue">
                    <SelectTrigger>
                        <SelectValue :placeholder="$t('task.message08')" />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectGroup>
                            <SelectLabel>{{ $t('task.message07') }}</SelectLabel>
                            <SelectItem v-for="option in datatypeOptions" :key="option.value" :value="option">
                                {{ option.text }}
                            </SelectItem>
                        </SelectGroup>
                    </SelectContent>
                </Select>
            </Label>
            <div class="grid gap-3">
                <Label for="interval">{{ $t('task.message10') }}</Label>
                <Input id="interval" v-model="interval" name="interval" />
            </div>
            <div class="grid gap-3">
                <Label for="critical">{{ $t('task.message12') }}</Label>
                <Input id="critical" v-model="critical" name="critical" />
            </div>
            <div class="grid gap-3">
                <Label for="warning">{{ $t('task.message14') }}</Label>
                <Input id="warning" v-model="warning" name="warning" />
            </div>
            <Label class="flex flex-col gap-3 items-start">
                {{ $t('task.message22') }}
                <Select v-model="emailListValue" :disabled="emailListOptions.length === 0">
                    <SelectTrigger>
                        <SelectValue :placeholder="$t('task.message23')" />
                    </SelectTrigger>
                    <SelectContent>
                        <SelectGroup>
                            <SelectLabel>{{ $t('task.message22') }}</SelectLabel>
                            <SelectItem v-for="option in emailListOptions" :key="option.value" :value="option">
                                {{ option.text }}
                            </SelectItem>
                        </SelectGroup>
                    </SelectContent>
                </Select>
            </Label>
            <Label class="flex flex-col gap-3 items-start">
                {{ $t('message.status') }}
                <Select v-model="statusValue">
                    <SelectTrigger>
                        <SelectValue :placeholder="$t('task.message15')" />
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
            </div>
            <DialogFooter class="mt-auto flex-row px-6 pb-[max(1.5rem,env(safe-area-inset-bottom))] *:flex-1 lg:mt-0 lg:px-0 lg:pb-0 lg:*:flex-none">
                <DialogClose as-child>
                    <Button>{{ $t('message.cancel') }}</Button>
                </DialogClose>
                <Button @click="putTaskUpdate()">{{ $t('message.save') }}</Button>
            </DialogFooter>
        </DialogContent>
    </Dialog>

    <Dialog v-model:open="showModal02">
        <DialogContent class="fixed! inset-0! top-0! left-0! z-50 flex! h-dvh! w-full! max-w-none! max-h-none! translate-x-0! translate-y-0! flex-col! gap-0! overflow-hidden! rounded-none! border-0! p-0! bg-card text-foreground! lg:inset-auto! lg:top-[50%]! lg:left-[50%]! lg:grid! lg:h-auto! lg:max-w-lg! lg:max-h-[80vh]! lg:translate-x-[-50%]! lg:translate-y-[-50%]! lg:gap-4! lg:overflow-visible! lg:rounded-lg! lg:border! lg:p-6!">
            <DialogHeader class="px-6 pt-6 pr-12 text-left lg:px-0 lg:pt-0 lg:pr-0">
                <DialogTitle>{{ $t('task.message21') }}</DialogTitle>
            </DialogHeader>
            <div class="grid min-h-0 flex-1 content-start gap-4 overflow-y-auto p-6 lg:flex-none lg:p-0">
                <Label class="flex flex-col gap-3 items-start">
                    {{ $t('message.module') }}
                    <Select v-model="moduleValue">
                        <SelectTrigger>
                            <SelectValue :placeholder="$t('module.message02')" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectGroup>
                                <SelectLabel>{{ $t('message.module') }}</SelectLabel>
                                <SelectItem v-for="option in moduleOptions" :key="option.value" :value="option">
                                    {{ option.text }}
                                </SelectItem>
                            </SelectGroup>
                        </SelectContent>
                    </Select>
                </Label>
                <Label class="flex flex-col gap-3 items-start">
                    {{ $t('message.space') }}
                    <Select v-model="spaceValue">
                        <SelectTrigger>
                            <SelectValue :placeholder="$t('module.message01')" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectGroup>
                                <SelectLabel>{{ $t('message.space') }}</SelectLabel>
                                <SelectItem v-for="option in spaceOptions" :key="option.value" :value="option">
                                    {{ option.text }}
                                </SelectItem>
                            </SelectGroup>
                        </SelectContent>
                    </Select>
                </Label>
                <Label class="flex flex-col gap-3 items-start">
                    {{ $t('task.message04') }}
                    <Select v-model="objectValue">
                        <SelectTrigger>
                            <SelectValue :placeholder="$t('module.message03')" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectGroup>
                                <SelectLabel>{{ $t('task.message04') }}</SelectLabel>
                                <SelectItem v-for="option in objectOptions" :key="option.value" :value="option">
                                    {{ option.text }}
                                </SelectItem>
                            </SelectGroup>
                        </SelectContent>
                    </Select>
                </Label>
                <Label class="flex flex-col gap-3 items-start">
                    {{ $t('message.dm') }}
                    <Select v-model="dmValue">
                        <SelectTrigger>
                            <SelectValue :placeholder="$t('task.message24')" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectGroup>
                                <SelectLabel>{{ $t('message.dm') }}</SelectLabel>
                                <SelectItem v-for="option in dmOptions" :key="option.value" :value="option">
                                    {{ option.text }}
                                </SelectItem>
                            </SelectGroup>
                        </SelectContent>
                    </Select>
                </Label>
                <div class="grid gap-3">
                    <Label for="metric">{{ $t('task.message06') }}</Label>
                    <Input id="metric" v-model="metric" name="metric" />
                </div>
                <Label class="flex flex-col gap-3 items-start">
                    {{ $t('task.message07') }}
                    <Select v-model="datatypeValue">
                        <SelectTrigger>
                            <SelectValue :placeholder="$t('task.message08')" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectGroup>
                                <SelectLabel>{{ $t('task.message07') }}</SelectLabel>
                                <SelectItem v-for="option in datatypeOptions" :key="option.value" :value="option">
                                    {{ option.text }}
                                </SelectItem>
                            </SelectGroup>
                        </SelectContent>
                    </Select>
                </Label>
                <div class="grid gap-3">
                    <Label for="interval">{{ $t('task.message10') }}</Label>
                    <Input id="interval" v-model="interval" name="interval" />
                </div>
                <div class="grid gap-3">
                    <Label for="critical">{{ $t('task.message12') }}</Label>
                    <Input id="critical" v-model="critical" name="critical" />
                </div>
                <div class="grid gap-3">
                    <Label for="warning">{{ $t('task.message14') }}</Label>
                    <Input id="warning" v-model="warning" name="warning" />
                </div>
                <Label class="flex flex-col gap-3 items-start">
                    {{ $t('task.message22') }}
                    <Select v-model="emailListValue" :disabled="emailListOptions.length === 0">
                        <SelectTrigger>
                            <SelectValue :placeholder="$t('task.message23')" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectGroup>
                                <SelectLabel>{{ $t('task.message22') }}</SelectLabel>
                                <SelectItem v-for="option in emailListOptions" :key="option.value" :value="option">
                                    {{ option.text }}
                                </SelectItem>
                            </SelectGroup>
                        </SelectContent>
                    </Select>
                </Label>
                <Label class="flex flex-col gap-3 items-start">
                    {{ $t('message.status') }}
                    <Select v-model="statusValue">
                        <SelectTrigger>
                            <SelectValue :placeholder="$t('task.message15')" />
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
            </div>
            <DialogFooter class="mt-auto flex-row px-6 pb-[max(1.5rem,env(safe-area-inset-bottom))] *:flex-1 lg:mt-0 lg:px-0 lg:pb-0 lg:*:flex-none">
                <DialogClose as-child>
                    <Button>{{ $t('message.cancel') }}</Button>
                </DialogClose>
                <Button @click="putTaskCreate()">{{ $t('message.save') }}</Button>
            </DialogFooter>
        </DialogContent>
    </Dialog>

    <Dialog v-model:open="showModal03">
        <DialogContent class="fixed! inset-0! top-0! left-0! z-50 flex! h-dvh! w-full! max-w-none! max-h-none! translate-x-0! translate-y-0! flex-col! gap-0! overflow-hidden! rounded-none! border-0! p-0! bg-card text-foreground! lg:inset-auto! lg:top-[50%]! lg:left-[50%]! lg:grid! lg:h-auto! lg:max-w-lg! lg:max-h-[80vh]! lg:translate-x-[-50%]! lg:translate-y-[-50%]! lg:gap-4! lg:overflow-visible! lg:rounded-lg! lg:border! lg:p-6!">
            <DialogHeader class="px-6 pt-6 pr-12 text-left lg:px-0 lg:pt-0 lg:pr-0">
                <DialogTitle>{{ $t('task.message16') }}</DialogTitle>
            </DialogHeader>
            <DialogFooter class="mt-auto flex-row px-6 pb-[max(1.5rem,env(safe-area-inset-bottom))] *:flex-1 lg:mt-0 lg:px-0 lg:pb-0 lg:*:flex-none">
                <DialogClose as-child>
                    <Button>{{ $t('message.cancel') }}</Button>
                </DialogClose>
                <Button @click="putTaskRemove(rowIndexRemove)">{{ $t('message.delete') }}</Button>
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
            const datatypeOptions = [
                    {text: this.$t('message.number'), value: "float"},
                    {text: this.$t('message.text'), value: "text"},
                ];
            const objectOptions = [
                    {text: "zbxagent", value: "zbxagent"},
                    {text: "prometheus", value: "prometheus"},
                    {text: "zbxagentpassive", value: "zbxagentpassive"},
                ];

            return {
                isDisabled: false,
                isDisabledDM: true,
                spaceOptions: [],
                moduleOptions: [],
                dmOptions: [],
                emailListOptions: [],
                emailListValue: {},
                spaceValue: {},
                moduleValue: {},
                dmValue: {},
                userRole: null,
                userBlock: true,
                rowIndexRemove: null,
                taskModuleIDRemove: "",
                taskSpaceIDRemove: "",
                taskObjectRemove: "",
                taskMetricRemove: "",
                noData: "Нет данных",
                showModal01: false,
                showModal02: false,
                showModal03: false,
                tableData: [],
                id: "",
                moduleid: "",
                spaceid: "",
                metric: "",
                critical: "",
                warning: "",
                interval: "",
                datatype: "",
                moduledesc: "",
                spacedesc: "",
                datatypeOptions: datatypeOptions,
                datatypeValue: datatypeOptions[0],
                statusOptions: statusOptions,
                statusValue: statusOptions[0],
                objectOptions: objectOptions,
                objectValue: objectOptions[0],
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
                    { field: "moduledesc", headerName: this.$t('message.module') },
                    { field: "spacedesc", headerName: this.$t('message.space') },
                    { field: "object", headerName: this.$t('task.object') },
                    { field: "metric", headerName: this.$t('task.metric') },
                    { field: "interval", headerName: this.$t('task.interval') },
                    { field: "status", headerName: this.$t('message.status') },
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
            objectUpdate() {
                if (this.objectValue.value == "prometheus") {
                    this.isDisabled = true;
                    this.isDisabledDM = false;
                } else {
                    this.isDisabled = false;
                    this.isDisabledDM = true;
                    this.dmValue = {};
                }
            },
            rowClickEdit(value) {
                this.emailListValue = "";

                this.moduleOptions.forEach(item => {
                    if (item.value == this.tableData[value].moduleid) this.moduleValue = item;
                });
                this.dmOptions.forEach(item => {
                    if (item.value == this.tableData[value].dmid) this.dmValue = item;
                });
                this.spaceOptions.forEach(item => {
                    if (item.value == this.tableData[value].spaceid) this.spaceValue = item;
                });
                this.statusOptions.forEach(item => {
                    if (item.value == this.tableData[value].status) this.statusValue = item;
                });
                this.datatypeOptions.forEach(item => {
                    if (item.value == this.tableData[value].datatype) this.datatypeValue = item;
                });
                this.objectOptions.forEach(item => {
                    if (item.value == this.tableData[value].object) this.objectValue = item;
                });
                this.emailListOptions.forEach(item => {
                    if (item.value == this.tableData[value].emaillistid) this.emailListValue = item;
                });

                this.showModal01 = true;
                this.metric = this.tableData[value].metric;
                this.critical = this.tableData[value].critical;
                this.warning = this.tableData[value].warning;
                this.interval = this.tableData[value].interval;

                if (this.objectValue.value == "prometheus") {
                    this.isDisabled = true;
                    this.isDisabledDM = false;
                } else {
                    this.isDisabled = false;
                    this.isDisabledDM = true;
                    this.dmValue = {};
                }
            },
            buttonClickNew() {
                this.moduleValue = "";
                this.dmValue = "";
                this.spaceValue = "";
                this.metric = "";
                this.critical = "";
                this.warning = "";
                this.interval = "";
                this.statusValue = this.statusOptions[0];
                this.datatypeValue = this.datatypeOptions[0];
                this.objectValue = this.objectOptions[0];
                this.showModal02 = true;
                this.emailListValue = "";

                if (this.objectValue.value == "prometheus") {
                    this.isDisabled = true;
                    this.isDisabledDM = false;
                } else {
                    this.isDisabled = false;
                    this.isDisabledDM = true;
                    this.dmValue = {};
                }
            },
            selectRight() {
                this.currentPage++;
                this.putTaskSelect(this.selectParam.next);
                this.pages[this.currentPage] = this.selectParam.next;
            },
            selectLeft() {
                if (this.currentPage > 0) {
                    this.currentPage--;
                }
                this.putTaskSelect(this.pages[this.currentPage]);
            },
            async putTaskSelect(current) {
                this.tableLoading = true;
                try {
                    const response = await axios.put(
                        "/api/v1/admin/task/select?" + Math.random(),
                        { current: current || "", pagesize: 100 },
                    );
                    const data = response.data;
                    this.tableData = data.data ?? [];
                    this.selectParam.current = data.current;
                    this.selectParam.next = data.next;
                    this.isDisabledRight = !data.next;
                    this.isDisabledLeft = !data.current;
                } catch (error) {
                    console.error("Failed to load tasks:", error);
                } finally {
                    this.tableLoading = false;
                }
            },
            async putTaskUpdate() {
                let dataPut = {
                    moduleid: this.moduleValue.value,
                    dmid: this.dmValue.value,
                    spaceid: this.spaceValue.value,
                    object: this.objectValue.value,
                    metric: this.metric,
                    status: parseInt(this.statusValue.value),
                    critical: this.critical,
                    warning: this.warning,
                    interval: parseInt(this.interval),
                    datatype: this.datatypeValue.value,
                    emaillistid: this.emailListValue.value,
                };

                try {
                    await axios.put("/api/v1/admin/task/update", dataPut);
                    toast(this.$t('message.message01'), {
                        type: 'success',
                        style: {
                            background: 'var(--success)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--success)'
                        }
                    });
                    this.putTaskSelect();
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
                        console.error("Failed to update task:", error);
                    }
                }
            },
            rowClickRemove(value) {
                this.showModal03 = true;
                this.rowIndexRemove = value;
                this.taskModulIDRemove = this.tableData[value].moduleid;
                this.taskSpaceIDRemove = this.tableData[value].spaceid;
                this.taskObjectRemove = this.tableData[value].object;
                this.taskMetricRemove = this.tableData[value].metric;
            },
            async putTaskRemove(value) {
                let dataPut = {
                    moduleid: this.tableData[value].moduleid,
                    spaceid: this.tableData[value].spaceid,
                    object: this.tableData[value].object,
                    metric: this.tableData[value].metric,
                };
                try {
                    await axios.put("/api/v1/admin/task/remove", dataPut);
                    toast(this.$t('task.message17'), {
                        type: 'success',
                        style: {
                            background: 'var(--success)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--success)'
                        }
                    });
                    this.putTaskSelect();
                    this.showModal03 = false;
                } catch (error) {
                    if (error.response?.status === 500) {
                        toast(this.$t('task.message18'), {
                            type: 'error',
                            style: {
                                background: 'var(--destructive)',
                                color: 'var(--primary-foreground)',
                                border: '1px solid var(--destructive)'
                            }
                        });
                        this.showModal03 = true;
                    } else {
                        console.error("Failed to remove task:", error);
                    }
                }
            },
            putTaskCreate() {
                if (this.status == "") {
                    this.status = 1;
                }
                if (this.moduleValue.value == "" || this.moduleValue.value == undefined) {
                    this.showModal02 = true;
                    toast(this.$t('task.message02'), {
                        type: 'error',
                        style: {
                            background: 'var(--destructive)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--destructive)'
                        }
                    });
                    return;
                }
                if (this.spaceValue.value == "" || this.spaceValue.value == undefined) {
                    this.showModal02 = true;
                    toast(this.$t('task.message01'), {
                        type: 'error',
                        style: {
                            background: 'var(--destructive)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--destructive)'
                        }
                    });
                    return;
                }
                if (this.objectValue.value == "" || this.objectValue.value == undefined) {
                    this.showModal02 = true;
                    toast(this.$t('task.message03'), {
                        type: 'error',
                        style: {
                            background: 'var(--destructive)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--destructive)'
                        }
                    });
                    return;
                }
                if (this.metric == "") {
                    this.showModal02 = true;
                    toast(this.$t('task.message05'), {
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
                    moduleid: this.moduleValue.value,
                    dmid: this.dmValue.value,
                    spaceid: this.spaceValue.value,
                    object: this.objectValue.value,
                    metric: this.metric,
                    status: parseInt(this.statusValue.value),
                    critical: this.critical,
                    warning: this.warning,
                    interval: parseInt(this.interval),
                    datatype: this.datatypeValue.value,
                    emaillistid: this.emailListValue.value,
                };
                axios.put("/api/v1/admin/task/create", dataPut)
                    .then(() => {
                        toast(this.$t('task.message19'), {
                            type: 'success',
                            style: {
                                background: 'var(--success)',
                                color: 'var(--primary-foreground)',
                                border: '1px solid var(--success)'
                            }
                        });
                        this.putTaskSelect();
                        this.showModal02 = false;
                    })
                    .catch(error => {
                        if (error.response?.status === 500) {
                            toast(this.$t('task.message20'), {
                                type: 'error',
                                style: {
                                    background: 'var(--destructive)',
                                    color: 'var(--primary-foreground)',
                                    border: '1px solid var(--destructive)'
                                }
                            });
                            this.showModal02 = true;
                        } else {
                            console.error("Failed to create task:", error);
                        }
                    });
            },
            async userInfo() {
                const response = await axios.get("/api/v1/userinfo?" + Math.random());
                const data = response.data;
                this.userRole = data.role;
                this.userBlock = data.role !== "superadmin" && data.role !== "admin";

                this.colDefs = [
                    { field: "moduledesc", headerName: this.$t('message.module') },
                    { field: "spacedesc", headerName: this.$t('message.space') },
                    { field: "object", headerName: this.$t('task.object') },
                    { field: "metric", headerName: this.$t('task.metric') },
                    { field: "interval", headerName: this.$t('task.interval') },
                    { field: "status", headerName: this.$t('message.status') },
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

                this.putTaskSelect();
                this.spaceSelect();
                this.moduleSelect();
                this.notificationListSelect();
            },
            async spaceSelect() {
                try {
                    const response = await axios.put(
                        "/api/v1/admin/space/select?" + Math.random(),
                        { pagesize: 0 },
                    );
                    const spaces = Array.isArray(response.data?.data) ? response.data.data : [];
                    this.spaceOptions = spaces.map(item => ({ text: item.description, value: item.id }));
                } catch (error) {
                    console.error("Failed to load spaces:", error);
                }
            },
            async moduleSelect() {
                try {
                    const response = await axios.put(
                        "/api/v1/admin/module/select?" + Math.random(),
                        { pagesize: 0 },
                    );
                    const modules = Array.isArray(response.data?.data) ? response.data.data : [];
                    this.moduleOptions = modules.map(item => ({ text: item.description, value: item.id }));
                    this.dmOptions = modules
                        .filter(item => item.type == "datamanager")
                        .map(item => ({ text: item.description, value: item.id }));
                } catch (error) {
                    console.error("Failed to load modules:", error);
                }
            },
            async notificationListSelect() {
                try {
                    const response = await axios.put(
                        "/api/v1/admin/notification/select?" + Math.random(),
                        { pagesize: 0 },
                    );
                    const notificationLists = Array.isArray(response.data?.data) ? response.data.data : [];
                    this.emailListOptions = notificationLists
                        .filter(item => item.status == 1)
                        .map(item => ({ text: item.name, value: item.id }));
                } catch (error) {
                    console.error("Failed to load notification lists:", error);
                }
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
                        { field: "moduledesc", headerName: this.$t('message.module') },
                        { field: "spacedesc", headerName: this.$t('message.space') },
                        { field: "object", headerName: this.$t('task.object') },
                        { field: "metric", headerName: this.$t('task.metric') },
                        { field: "interval", headerName: this.$t('task.interval') },
                        { field: "status", headerName: this.$t('message.status') },
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
                    this.datatypeOptions = [
                        {text: this.$t('message.number'), value: "float"},
                        {text: this.$t('message.text'), value: "text"},
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