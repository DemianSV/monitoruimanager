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
                    <BreadcrumbPage>{{ $t('message.dbstatus') }}</BreadcrumbPage>
                </BreadcrumbItem>
            </BreadcrumbList>
        </Breadcrumb>
    </div>

    <div class="flex! items-start! gap-3 w-full">
        <ag-grid-vue :key="componentKey" :rowData="tableData" :columnDefs="colDefs" :defaultColDef="defaultColDef" @grid-ready="onGridReady" domLayout="autoHeight" style="height: 100%; width: 100%" :loading="tableLoading" />
    </div>

    </div>
</template>

<script>
import axios from 'axios'
import { mapState } from 'vuex';
import { AgGridVue } from "ag-grid-vue3"; // Vue Data Grid Component

import {
    Breadcrumb,
    BreadcrumbItem,
    BreadcrumbLink,
    BreadcrumbList,
    BreadcrumbPage,
    BreadcrumbSeparator,
} from '@/components/ui/breadcrumb'

export default {
    components: {
        AgGridVue, // Add Vue Data Grid component
        Breadcrumb,
        BreadcrumbItem,
        BreadcrumbLink,
        BreadcrumbList,
        BreadcrumbPage,
        BreadcrumbSeparator,
    },
    data() {
        return {
            userRole: null,
            userBlock: true,
            adminBlock: true,
            tableData: [],
            tableLoading: false,
            gridApi: null,
            componentKey: 0,
            colDefs: [
                { field: "host", headerName: "Host" },
                { field: "dc", headerName: "DC" },
                { field: "rack", headerName: "Rack" },
                { field: "version", headerName: "Version" },
                { field: "stat", headerName: "Stat" },
            ],
            defaultColDef: {
                minWidth: 120,
                resizable: true,
            },
        }
    },
    methods: {
        async databaseStatus() {
            this.tableLoading = true;
            try {
                const response = await axios.get("/api/v1/admin/database/status?" + Math.random());
                const data = Array.isArray(response.data) ? response.data : [];
                this.tableData = data.map(item => ({
                    ...item,
                    dc: item.dc.toUpperCase(),
                    rack: item.rack.toUpperCase(),
                    stat: item.stat.toUpperCase(),
                }));
            } catch (error) {
                console.error("Failed to load database status:", error);
                this.tableData = [];
            } finally {
                this.tableLoading = false;
            }
        },
        async userInfo() {
            try {
                const response = await axios.get("/api/v1/userinfo?" + Math.random());
                const data = response.data;
                this.userRole = data.role;
                if (data.role == "superadmin") {
                    this.userBlock = false;
                    this.adminBlock = false;
                } else if (data.role == "admin") {
                    this.userBlock = false;
                    this.adminBlock = true;
                    this.databaseData = "Нет доступа";
                } else {
                    this.userBlock = true;
                    this.adminBlock = true;
                    this.databaseData = "Нет доступа";
                    this.userData = "Нет доступа";
                }
            } catch (error) {
                console.error("Failed to load user info:", error);
            }
            if (this.userRole == "superadmin" || this.userRole == "admin" || this.userRole == "user") {
                this.databaseStatus();
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
            this.colDefs = [
                    { field: "host", headerName: "Host" },
                    { field: "dc", headerName: "DC" },
                    { field: "rack", headerName: "Rack" },
                    { field: "version", headerName: "Version" },
                    { field: "stat", headerName: "Stat" },
                ];
            this.componentKey += 1;
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