<template>
    <Toaster position="top-center" :offset="24"/>

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
                    <BreadcrumbPage>{{ $t('message.notifications') }}</BreadcrumbPage>
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
                @click="putNotificationSelect(pages[currentPage])">
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
                <DialogTitle>{{ $t('notilfication.message31') }}</DialogTitle>
            </DialogHeader>
            <div class="grid min-h-0 flex-1 content-start gap-4 overflow-y-auto p-6 lg:flex-none lg:p-0">
                <div class="grid gap-3">
                    <Label for="id">ID</Label>
                    <Input id="id" v-model="id" name="id" disabled />
                </div>
                <div class="grid gap-3">
                    <Label for="emailListEditName">{{ $t('notilfication.message02') }}</Label>
                    <Input id="emailListEditName" v-model="emailListEditName" name="emailListEditName" />
                </div>
                <div class="grid gap-3">
                    <Label for="emailListEditDescription">{{ $t('notilfication.message04') }}</Label>
                    <Input id="emailListEditDescription" v-model="emailListEditDescription" name="emailListEditDescription" />
                </div>
                <Label class="flex flex-col gap-3 items-start">
                    {{ $t('message.status') }}
                    <Select v-model="emailListEditStatusValue">
                        <SelectTrigger>
                            <SelectValue :placeholder="$t('notilfication.message11')" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectGroup>
                                <SelectLabel>{{ $t('message.status') }}</SelectLabel>
                                <SelectItem v-for="option in emailListStatusOptions" :key="option.value" :value="option">
                                    {{ option.text }}
                                </SelectItem>
                            </SelectGroup>
                        </SelectContent>
                    </Select>
                </Label>
                <Separator class="my-4" />

                <Label for="emailNew" class="px-1">
                    {{ $t('notilfication.message12') }}
                </Label>
                <Button id="emailNew" type="button" @click="buttonClickDataNew()">
                    <Plus class="size-5" />
                </Button>
                <div class="grid gap-4 w-full" v-for="(emailListData, index) in tableDataRow.data" :key="index">
                    <div class="grid gap-3">
                        <Label :for="'email' + index">{{ $t('notilfication.message13') }}</Label>
                        <Input :id="'email' + index" v-model="email[index]" :name="'email' + index" />
                    </div>
                    <div class="grid gap-3">
                        <Label :for="'description' + index">{{ $t('notilfication.message23') }}</Label>
                        <Input :id="'description' + index" v-model="description[index]" :name="'description' + index" />
                    </div>
                    <Label class="flex flex-col gap-3 items-start">
                        {{ $t('message.status') }}
                        <Select v-model="status[index]">
                            <SelectTrigger>
                                <SelectValue :placeholder="$t('notilfication.message11')" />
                            </SelectTrigger>
                            <SelectContent>
                                <SelectGroup>
                                    <SelectLabel>{{ $t('message.status') }}</SelectLabel>
                                    <SelectItem v-for="option in emailListStatusOptions" :key="option.value" :value="option">
                                        {{ option.text }}
                                    </SelectItem>
                                </SelectGroup>
                            </SelectContent>
                        </Select>
                    </Label>
                    <Button :id="'emailNew' + index" type="button" @click="buttonClickDataRemove(index)">
                        <Trash class="size-5" />
                    </Button>
                    <Separator class="my-4" />
                </div>
            </div>
            <DialogFooter class="mt-auto flex-row px-6 pb-[max(1.5rem,env(safe-area-inset-bottom))] *:flex-1 lg:mt-0 lg:px-0 lg:pb-0 lg:*:flex-none">
                <DialogClose as-child>
                    <Button>{{ $t('message.cancel') }}</Button>
                </DialogClose>
                <Button @click="putEMailListUpdate()">{{ $t('message.save') }}</Button>
            </DialogFooter>
        </DialogContent>
    </Dialog>

    <Dialog v-model:open="showModal02">
        <DialogContent class="fixed! inset-0! top-0! left-0! z-50 flex! h-dvh! w-full! max-w-none! max-h-none! translate-x-0! translate-y-0! flex-col! gap-0! overflow-hidden! rounded-none! border-0! p-0! bg-card text-foreground! lg:inset-auto! lg:top-[50%]! lg:left-[50%]! lg:grid! lg:h-auto! lg:max-w-lg! lg:max-h-[80vh]! lg:translate-x-[-50%]! lg:translate-y-[-50%]! lg:gap-4! lg:overflow-visible! lg:rounded-lg! lg:border! lg:p-6!">
            <DialogHeader class="px-6 pt-6 pr-12 text-left lg:px-0 lg:pt-0 lg:pr-0">
                <DialogTitle>{{ $t('notilfication.message26') }}</DialogTitle>
            </DialogHeader>
            <div class="grid min-h-0 flex-1 content-start gap-4 overflow-y-auto p-6 lg:flex-none lg:p-0">
                <div class="grid gap-3">
                    <Label for="emailListName">{{ $t('notilfication.message02') }}</Label>
                    <Input id="emailListName" v-model="emailListName" name="emailListName" />
                </div>
                <div class="grid gap-3">
                    <Label for="emailListDescription">{{ $t('notilfication.message04') }}</Label>
                    <Input id="emailListDescription" v-model="emailListDescription" name="emailListDescription" />
                </div>
                <Label class="flex flex-col gap-3 items-start">
                    {{ $t('message.status') }}
                    <Select v-model="emailListStatusValue">
                        <SelectTrigger>
                            <SelectValue :placeholder="$t('notilfication.message11')" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectGroup>
                                <SelectLabel>{{ $t('message.status') }}</SelectLabel>
                                <SelectItem v-for="option in emailListStatusOptions" :key="option.value" :value="option">
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
                <Button @click="putNotificationCreate()">{{ $t('message.save') }}</Button>
            </DialogFooter>
        </DialogContent>
    </Dialog>

    <Dialog v-model:open="showModal03">
        <DialogContent class="fixed! inset-0! top-0! left-0! z-50 flex! h-dvh! w-full! max-w-none! max-h-none! translate-x-0! translate-y-0! flex-col! gap-0! overflow-hidden! rounded-none! border-0! p-0! bg-card text-foreground! lg:inset-auto! lg:top-[50%]! lg:left-[50%]! lg:grid! lg:h-auto! lg:max-w-lg! lg:max-h-[80vh]! lg:translate-x-[-50%]! lg:translate-y-[-50%]! lg:gap-4! lg:overflow-visible! lg:rounded-lg! lg:border! lg:p-6!">
            <DialogHeader class="px-6 pt-6 pr-12 text-left lg:px-0 lg:pt-0 lg:pr-0">
                <DialogTitle>{{ $t('notilfication.message27') }}</DialogTitle>
            </DialogHeader>
            <div class="flex-1 p-6 lg:flex-none lg:p-0">{{ $t('notilfication.message28') }} {{ emailListIDRemove }}</div>
            <DialogFooter class="mt-auto flex-row px-6 pb-[max(1.5rem,env(safe-area-inset-bottom))] *:flex-1 lg:mt-0 lg:px-0 lg:pb-0 lg:*:flex-none">
                <DialogClose as-child>
                    <Button>{{ $t('message.cancel') }}</Button>
                </DialogClose>
                <Button @click="putNotificationRemove(rowIndexRemove)">{{ $t('message.delete') }}</Button>
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
    import { Separator } from '@/components/ui/separator'
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

    import { ChevronsLeft, ChevronsRight, RefreshCw, Plus, Trash } from '@lucide/vue';

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
            Separator,
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
            Trash,
        },
        data() {
            const emailListStatusOptions = [
                {text: this.$t('message.on'), value: 1},
                {text: this.$t('message.off'), value: 0},
            ];

            return {
                userRole: null,
                userBlock: true,
                rowIndexRemove: null,
                emailListIDRemove: "",
                showModal01: false,
                showModal02: false,
                showModal03: false,
                tableData: [],
                tableDataRow: {},
                id: "",
                emailListName: "",
                emailListDescription: "",
                emailListEditName: "",
                emailListStatusValue: emailListStatusOptions[0],
                emailListStatusOptions: emailListStatusOptions,
                emailListEditDescription: "",
                emailListEditStatusValue: emailListStatusOptions[0],
                emailListStatusOptions: emailListStatusOptions,
                email: [],
                description: [],
                status: [],
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
                    { field: "name", headerName: this.$t('notilfication.message02') },
                    { field: "description", headerName: this.$t('notilfication.message04') },
                    { field: "status", headerName: this.$t('message.status') },
                    { field: "owner", headerName: this.$t('message.owner') },
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
            rowClickEdit(index) {
                const row = this.tableData[index];
                if (!row) return;

                this.tableDataRow = {
                    ...row,
                    data: Array.isArray(row.data) ? row.data.map(item => ({ ...item })) : [],
                };
                this.email = this.tableDataRow.data.map(item => item.email ?? "");
                this.description = this.tableDataRow.data.map(item => item.description ?? "");
                this.status = this.tableDataRow.data.map(item =>
                    this.emailListStatusOptions.find(option => option.value == item.status)
                        ?? this.emailListStatusOptions[0]
                );
                this.rowIndexEdit = index;
                this.id = this.tableDataRow.id;
                this.emailListEditName = this.tableDataRow.name;
                this.emailListEditDescription = this.tableDataRow.description;
                this.emailListEditStatusValue = this.emailListStatusOptions.find(
                    item => item.value == this.tableDataRow.status
                ) ?? this.emailListStatusOptions[0];
                this.showModal01 = true;
            },
            buttonClickNew() {
                this.id = "";
                this.emailListName = "";
                this.emailListDescription = "";
                this.emailListStatusValue = this.emailListStatusOptions[0];
                this.showModal02 = true;
            },
            buttonClickDataNew() {
                if (!Array.isArray(this.tableDataRow.data)) {
                    this.tableDataRow.data = [];
                }
                this.tableDataRow.data.push({
                    id: "",
                    email: "",
                    description: "",
                    status: 1,
                });
                this.email = [...this.email, ""];
                this.description = [...this.description, ""];
                this.status = [...this.status, this.emailListStatusOptions[0]];
            },
            buttonClickDataRemove(index) {
                if (!Array.isArray(this.tableDataRow.data)) return;
                this.tableDataRow.data.splice(index, 1);
                this.email.splice(index, 1);
                this.description.splice(index, 1);
                this.status.splice(index, 1);
            },
            putNotificationSelect(current) { 
                let dataPut = {};
                if (current != null && current != undefined && current != "") {
                    dataPut.current = current; 
                } else { 
                    dataPut.current = ''; 
                } 
                dataPut.pagesize = 100;
    
                this.tableLoading = true;
    
                axios.put('/api/v1/admin/notification/select?' + Math.random(), dataPut) 
                    .then((response) => {
                        const data = response.data;
                        console.log(data);
                        if (data && (Array.isArray(data.data) || data.data === null)) {
                            const rows = Array.isArray(data.data) ? data.data : [];

                            rows.sort(function (a, b) {
                                if (a.name > b.name) {
                                    return 1;
                                }
                                if (a.name < b.name) {
                                    return -1;
                                }
                                return 0;
                            });

                            this.tableData = rows; 
                            this.selectParam.current = data.current; 
                            this.selectParam.next = data.next;
                            
                            if (data.next != null && data.next != undefined && data.next != "" ) { 
                                this.isDisabledRight = false; 
                            } else { 
                                this.isDisabledRight = true; 
                            } 
                        
                            if (data.current != null && data.current != undefined && data.current != "") { 
                                this.isDisabledLeft = false; 
                            } else { 
                                this.isDisabledLeft = true; 
                            }

                            this.tableLoading = false;
                        } else { 
                            console.error('Invalid response format:', response.data); 
                            this.tableLoading = false;
                        } 
                    }) 
                    .catch(() => { 
                        this.tableLoading = false; 
                    }); 
            },
            async putEMailListUpdate() {
                const dataPut = {
                    id: this.id,
                    name: this.emailListEditName,
                    description: this.emailListEditDescription,
                    status: parseInt(this.emailListEditStatusValue.value),
                    data: (Array.isArray(this.tableDataRow.data) ? this.tableDataRow.data : []).map((item, index) => ({
                        id: item.id ?? "",
                        email: this.email[index] ?? "",
                        description: this.description[index] ?? "",
                        status: Number(this.status[index]?.value ?? item.status ?? 0),
                    })),
                };

                try {
                    await axios.put("/api/v1/admin/notification/edit", dataPut);
                    toast(this.$t('message.message01'), {
                        type: 'success',
                        style: {
                            background: 'var(--success)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--success)'
                        }
                    });
                    this.putNotificationSelect();
                    this.showModal01 = false;
                } catch (error) {
                    toast(this.$t('message.message02'), {
                        type: 'error',
                        style: {
                            background: 'var(--destructive)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--destructive)'
                        }
                    });
                    this.showModal01 = true;
                }
            },
            rowClickRemove(index) {
                this.rowIndexRemove = index;
                this.emailListIDRemove = this.tableData[index].id;
                this.showModal03 = true;
            },
            async putNotificationRemove(index) {
                const dataPut = {
                    id: this.tableData[index].id,
                };
                try {
                    await axios.put("/api/v1/admin/notification/remove", dataPut);
                    toast(this.$t('notilfication.message29'), {
                        type: 'success',
                        style: {
                            background: 'var(--success)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--success)'
                        }
                    });
                    this.putNotificationSelect();
                    this.showModal03 = false;
                } catch (error) {
                    toast(this.$t('notilfication.message30'), {
                        type: 'error',
                        style: {
                            background: 'var(--destructive)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--destructive)'
                        }
                    });
                    this.showModal03 = true;
                }
            },
            selectRight() {
                this.currentPage++;
                this.putNotificationSelect(this.selectParam.next);
                this.pages[this.currentPage] = this.selectParam.next;
            },
            selectLeft() {
                if (this.currentPage > 0) {
                    this.currentPage--;
                }
                this.putNotificationSelect(this.pages[this.currentPage]);
            },
            async putNotificationCreate() {
                if (this.emailListName == "") {
                    this.showModal02 = true;
                    toast(this.$t('notilfication.message14'), {
                        type: 'error',
                        style: {
                            background: 'var(--destructive)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--destructive)'
                        }
                    });
                    return;
                }

                const dataPut = {
                    name: this.emailListName,
                    description: this.emailListDescription,
                    status: 1
                };
                try {
                    await axios.put("/api/v1/admin/notification/create", dataPut);
                    toast(this.$t('notilfication.message15'), {
                        type: 'success',
                        style: {
                            background: 'var(--success)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--success)'
                        }
                    });
                    this.putNotificationSelect();
                    this.showModal02 = false;
                } catch (error) {
                    toast(this.$t('notilfication.message16'), {
                        type: 'error',
                        style: {
                            background: 'var(--destructive)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--destructive)'
                        }
                    });
                    this.showModal02 = true;
                }
            },
            async userInfo() {
                try {
                    const response = await axios.get("/api/v1/userinfo?" + Math.random());
                    const data = response.data;
                            this.userRole = data.role;
                            if (data.role == "superadmin" || data.role == "admin") {
                                this.userBlock = false
                            } 
                            else {
                                this.userBlock = true
                            }
                            
                } catch(err) {
                    console.log("UIManager Error!");
                }

                this.colDefs = [
                    { field: "name", headerName: this.$t('notilfication.message02') },
                    { field: "description", headerName: this.$t('notilfication.message04') },
                    { field: "status", headerName: this.$t('message.status') },
                    { field: "owner", headerName: this.$t('message.owner') },
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

                this.putNotificationSelect();
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
                        { field: "name", headerName: this.$t('notilfication.message02') },
                        { field: "description", headerName: this.$t('notilfication.message04') },
                        { field: "status", headerName: this.$t('message.status') },
                        { field: "owner", headerName: this.$t('message.owner') },
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