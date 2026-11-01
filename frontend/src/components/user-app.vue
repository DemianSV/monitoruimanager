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
                    <BreadcrumbPage>{{ $t('message.users') }}</BreadcrumbPage>
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
                @click="putUserSelect(pages[currentPage])">
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
                <DialogTitle>{{ $t('user.message08') }}</DialogTitle>
            </DialogHeader>
            <div class="grid min-h-0 flex-1 content-start gap-4 overflow-y-auto p-6 lg:flex-none lg:p-0">
                <div class="grid gap-3">
                    <Label for="userID">User ID</Label>
                    <Input id="userID" v-model="userID" name="userID" disabled />
                </div>
                <div class="grid gap-3">
                    <Label for="userName">{{ $t('user.username') }}</Label>
                    <Input id="userName" v-model="userName" name="userName" autocomplete="off" disabled />
                </div>
                <div class="grid gap-3">
                    <Label for="newPassword">{{ $t('password.new') }}</Label>
                    <Input type="password" id="newPassword" v-model="newPassword" name="newPassword" />
                </div>
                <div class="grid gap-3">
                    <Label for="firstName">{{ $t('user.firstname') }}</Label>
                    <Input id="firstName" v-model="firstName" name="firstName" />
                </div>
                <div class="grid gap-3">
                    <Label for="lastName">{{ $t('user.lastname') }}</Label>
                    <Input id="lastName" v-model="lastName" name="lastName" />
                </div>
                <Label class="flex flex-col gap-3 items-start">
                    {{ $t('user.role') }}
                    <Select v-model="roleValue">
                        <SelectTrigger>
                            <SelectValue :placeholder="$t('user.role')" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectGroup>
                                <SelectLabel>{{ $t('user.role') }}</SelectLabel>
                                <SelectItem v-for="option in roleOptions" :key="option.value" :value="option">
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
                <Button @click="putUserUpdate()">{{ $t('message.save') }}</Button>
            </DialogFooter>
        </DialogContent>
    </Dialog>

    <Dialog v-model:open="showModal02">
        <DialogContent class="fixed! inset-0! top-0! left-0! z-50 flex! h-dvh! w-full! max-w-none! max-h-none! translate-x-0! translate-y-0! flex-col! gap-0! overflow-hidden! rounded-none! border-0! p-0! bg-card text-foreground! lg:inset-auto! lg:top-[50%]! lg:left-[50%]! lg:grid! lg:h-auto! lg:max-w-lg! lg:max-h-[80vh]! lg:translate-x-[-50%]! lg:translate-y-[-50%]! lg:gap-4! lg:overflow-visible! lg:rounded-lg! lg:border! lg:p-6!">
            <DialogHeader class="px-6 pt-6 pr-12 text-left lg:px-0 lg:pt-0 lg:pr-0">
                <DialogTitle>{{ $t('user.new') }}</DialogTitle>
            </DialogHeader>
            <div class="grid min-h-0 flex-1 content-start gap-4 overflow-y-auto p-6 lg:flex-none lg:p-0">
                <div class="grid gap-3">
                    <Label for="userName">{{ $t('user.username') }}</Label>
                    <Input id="userName" v-model="userName" name="userName" autocomplete="off" />
                </div>
                <div class="grid gap-3">
                    <Label for="newPassword">{{ $t('password.new') }}</Label>
                    <Input type="password" id="newPassword" v-model="newPassword" name="newPassword" />
                </div>
                <div class="grid gap-3">
                    <Label for="firstName">{{ $t('user.firstname') }}</Label>
                    <Input id="firstName" v-model="firstName" name="firstName" />
                </div>
                <div class="grid gap-3">
                    <Label for="lastName">{{ $t('user.lastname') }}</Label>
                    <Input id="lastName" v-model="lastName" name="lastName" />
                </div>
                <Label class="flex flex-col gap-3 items-start">
                    {{ $t('user.role') }}
                    <Select v-model="roleValue">
                        <SelectTrigger>
                            <SelectValue :placeholder="$t('user.role')" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectGroup>
                                <SelectLabel>{{ $t('user.role') }}</SelectLabel>
                                <SelectItem v-for="option in roleOptions" :key="option.value" :value="option">
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
                <Button @click="putUserCreate()">{{ $t('message.save') }}</Button>
            </DialogFooter>
        </DialogContent>
    </Dialog>

    <Dialog v-model:open="showModal03">
        <DialogContent class="fixed! inset-0! top-0! left-0! z-50 flex! h-dvh! w-full! max-w-none! max-h-none! translate-x-0! translate-y-0! flex-col! gap-0! overflow-hidden! rounded-none! border-0! p-0! bg-card text-foreground! lg:inset-auto! lg:top-[50%]! lg:left-[50%]! lg:grid! lg:h-auto! lg:max-w-lg! lg:max-h-[80vh]! lg:translate-x-[-50%]! lg:translate-y-[-50%]! lg:gap-4! lg:overflow-visible! lg:rounded-lg! lg:border! lg:p-6!">
            <DialogHeader class="px-6 pt-6 pr-12 text-left lg:px-0 lg:pt-0 lg:pr-0">
                    <DialogTitle>{{ $t('user.message07') }}</DialogTitle>
            </DialogHeader>
            <div class="flex-1 p-6 lg:flex-none lg:p-0">{{ $t('user.message09') }} {{ userIDRemove }}</div>
            <DialogFooter class="mt-auto flex-row px-6 pb-[max(1.5rem,env(safe-area-inset-bottom))] *:flex-1 lg:mt-0 lg:px-0 lg:pb-0 lg:*:flex-none">
                <DialogClose as-child>
                    <Button>{{ $t('message.cancel') }}</Button>
                </DialogClose>
                <Button @click="putUserRemove(rowIndexRemove)">{{ $t('message.delete') }}</Button>
            </DialogFooter>
        </DialogContent>
    </Dialog>

    </div>
</template>

<script>
    import axios from 'axios'
    import moment from 'moment'
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
                { text: this.$t('message.on'), value: 1 },
                { text: this.$t('message.off'), value: 0 },
            ];
            const roleOptions = [
                { text: this.$t('role.user'), value: "user" },
                { text: this.$t('role.admin'), value: "admin" },                
                { text: this.$t('role.superadmin'), value: "superadmin" },
            ];

            return {
                userRole: null,
                userBlock: true,
                rowIndexRemove: null,
                userNameRemove: "",
                noData: this.$t('message.nodata'),
                showModal01: false,
                showModal02: false,
                showModal03: false,
                tableData: [],
                userID: "",
                userIDRemove: "",
                userName: "",
                newPassword: "",
                firstName: "",
                lastName: "",
                statusOptions: statusOptions,
                statusValue: statusOptions[0],
                roleOptions: roleOptions,
                roleValue: roleOptions[0],
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
                    { field: "userid", headerName: this.$t('user.userid') },
                    { field: "username", headerName: this.$t('user.username') },
                    { field: "firstname", headerName: this.$t('user.firstname') },
                    { field: "lastname", headerName: this.$t('user.lastname') },
                    { field: "role", headerName: this.$t('user.role') },
                    { field: "status", headerName: this.$t('message.status') },
                    { field: "logintime", headerName: this.$t('user.logintime') },
                    { field: "ownerid", headerName: this.$t('user.ownerid') },
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
            rowClickEdit(value) {
                if (this.userRole == "admin") {
                    this.roleOptions = [
                        { text: this.$t('role.user'), value: "user" },
                    ];
                } else {
                    this.roleOptions = [
                        { text: this.$t('role.user'), value: "user" },
                        { text: this.$t('role.admin'), value: "admin" },                
                        { text: this.$t('role.superadmin'), value: "superadmin" },
                    ];
                }
                this.statusOptions.forEach(item => {
                    if (item.value == this.tableData[value].status) this.statusValue = item;
                });
                this.roleOptions.forEach(item => {
                    if (item.value == this.tableData[value].role) this.roleValue = item;
                });
                this.showModal01 = true;
                this.userID = this.tableData[value].userid;
                this.userName = this.tableData[value].username;
                this.newPassword = "";
                this.firstName = this.tableData[value].firstname;
                this.lastName = this.tableData[value].lastname;
                this.description = this.tableData[value].description;
            },
            buttonClickNew() {
                this.userName = "";
                this.newPassword = "";
                this.firstName = "";
                this.lastName = "";
                this.description = "";
                this.statusValue = this.statusOptions[0];
                this.showModal02 = true;
                if (this.userRole == "admin") {
                    this.roleOptions = [
                        { text: this.$t('role.user'), value: "user" },
                    ];
                } else {
                    this.roleOptions = [
                        { text: this.$t('role.user'), value: "user" },
                        { text: this.$t('role.admin'), value: "admin" },                
                        { text: this.$t('role.superadmin'), value: "superadmin" },
                    ];
                }
                this.roleValue = this.roleOptions[0];
            },
            selectRight() {
                this.currentPage++;
                this.putUserSelect(this.selectParam.next);
                this.pages[this.currentPage] = this.selectParam.next;
            },
            selectLeft() {
                if (this.currentPage > 0) {
                    this.currentPage--;
                }
                this.putUserSelect(this.pages[this.currentPage]);
            },
            async putUserSelect(current) {
                this.tableLoading = true;
                try {
                    const response = await axios.put(
                        "/api/v1/admin/user/select?" + Math.random(),
                        { current: current || "", pagesize: 100 },
                    );
                    const data = response.data;
                    const users = Array.isArray(data?.data) ? data.data : [];
                    this.tableData = users.map(item => ({
                        ...item,
                        logintime: moment(item.logintime).format("DD.MM.YYYY HH:mm:ss"),
                    }));
                    this.selectParam.current = data.current;
                    this.selectParam.next = data.next;
                    this.isDisabledRight = !data.next;
                    this.isDisabledLeft = !data.current;
                    this.componentKey += 1;
                } catch (error) {
                    console.error("Failed to load users:", error);
                    this.tableData = [];
                } finally {
                    this.tableLoading = false;
                }
            },
            async putUserUpdate() {
                if (this.newPassword.length < 8) {
                    toast(this.$t('password.message04'), { type: 'error',
                        style: {
                            background: 'var(--destructive)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--destructive)'
                        }
                    });
                    return;
                }

                let dataPut = {
                    username: this.userName,
                    password: this.newPassword,
                    firstname: this.firstName,
                    lastname: this.lastName,
                    status: parseInt(this.statusValue.value),
                    description: this.description,
                    role: this.roleValue.value,
                };
                try {
                    await axios.put("/api/v1/admin/user/update", dataPut);
                    toast(this.$t('message.message01'), {
                        type: 'success',
                        style: {
                            background: 'var(--success)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--success)'
                        }
                    });
                    this.putUserSelect();
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
                        console.error("Failed to update user:", error);
                    }
                }
            },
            rowClickRemove(value) {
                this.rowIndexRemove = value;
                this.userNameRemove = this.tableData[value].username;
                this.userIDRemove = this.tableData[value].userid;
                this.showModal03 = true;
            },
            async putUserRemove(value) {
                let dataPut = {
                    username: this.tableData[value].username,
                    userid: this.tableData[value].userid,
                    ownerid: this.tableData[value].ownerid,
                };
                try {
                    await axios.put("/api/v1/admin/user/remove", dataPut);
                    toast(this.$t('user.message01'), {
                        type: 'success',
                        style: {
                            background: 'var(--success)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--success)'
                        }
                    });
                    this.putUserSelect();
                    this.showModal03 = false;
                } catch (error) {
                    if (error.response?.status === 500) {
                        toast(this.$t('user.message02'), {
                            type: 'error',
                            style: {
                                background: 'var(--destructive)',
                                color: 'var(--primary-foreground)',
                                border: '1px solid var(--destructive)'
                            }
                        });
                        this.showModal03 = true;
                    } else {
                        console.error("Failed to remove user:", error);
                    }
                }
            },
            async putUserCreate() {
                if (this.userName == "") {
                    this.showModal02 = true;
                    toast(this.$t('user.message03'), {
                        type: 'warning',
                        style: {
                            background: 'var(--warning)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--warning)'
                        }
                    });
                    return;
                }
                if (this.newPassword == "") {
                    this.showModal02 = true;
                    toast(this.$t('user.message04'), {
                        type: 'warning',
                        style: {
                            background: 'var(--warning)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--warning)'
                        }
                    });
                    return;
                }
                let dataPut = {
                    username: this.userName,
                    password: this.newPassword,
                    firstname: this.firstName,
                    lastname: this.lastName,
                    status: parseInt(this.statusValue.value),
                    description: this.description,
                    role: this.roleValue.value,
                };
                try {
                    await axios.put("/api/v1/admin/user/create", dataPut);
                    toast(this.$t('user.message05'), {
                        type: 'success',
                        style: {
                            background: 'var(--success)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--success)'
                        }
                    });
                    this.putUserSelect();
                    this.showModal02 = false;
                } catch (error) {
                    if (error.response?.status === 500) {
                        this.$vaToast.init({ message: this.$t('user.message06'), color: 'danger' });
                        toast(this.$t('user.message06'), {
                            type: 'error',
                            style: {
                                background: 'var(--destructive)',
                                color: 'var(--primary-foreground)',
                                border: '1px solid var(--destructive)'
                            }
                        });
                        this.showModal02 = true;
                    } else {
                        console.error("Failed to create user:", error);
                    }
                }
            },
            async userInfo() {
                try {
                    const response = await axios.get("/api/v1/userinfo?" + Math.random());
                    this.userRole = response.data.role;
                    this.userBlock = this.userRole !== "superadmin" && this.userRole !== "admin";
                    if (this.userBlock) this.userData = this.$t('message.noaccess');
                } catch (error) {
                    console.error("Failed to load user info:", error);
                    return;
                }
                this.putUserSelect();
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
                        { field: "userid", headerName: this.$t('user.userid') },
                        { field: "username", headerName: this.$t('user.username') },
                        { field: "firstname", headerName: this.$t('user.firstname') },
                        { field: "lastname", headerName: this.$t('user.lastname') },
                        { field: "role", headerName: this.$t('user.role') },
                        { field: "status", headerName: this.$t('message.status') },
                        { field: "logintime", headerName: this.$t('user.logintime') },
                        { field: "ownerid", headerName: this.$t('user.ownerid') },
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
                        { text: this.$t('message.on'), value: 1 },
                        { text: this.$t('message.off'), value: 0 },
                    ];
                    this.roleOptions = [
                        { text: this.$t('role.user'), value: "user" },
                        { text: this.$t('role.admin'), value: "admin" },                
                        { text: this.$t('role.superadmin'), value: "superadmin" },
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