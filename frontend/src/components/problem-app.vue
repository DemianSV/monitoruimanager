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
                    <BreadcrumbPage>{{ $t('message.problems') }}</BreadcrumbPage>
                </BreadcrumbItem>
            </BreadcrumbList>
        </Breadcrumb>
    </div>

    <div class="flex! items-start! gap-2 w-full">
        <ag-grid-vue :key="componentKey" :rowData="tableData" :columnDefs="colDefs" :defaultColDef="defaultColDef" domLayout="autoHeight" @grid-ready="onGridReady" style="height: 100%; width: 100%" :loading="tableLoading" />
    </div>
    </div>
</template>

<script>
    import axios from 'axios'
    import moment from 'moment'
    import { mapState } from 'vuex';
    import { AgGridVue } from "ag-grid-vue3"; // Vue Data Grid Component

    import { Button } from '@/components/ui/button'
    import { Input } from '@/components/ui/input'
    import { Label } from '@/components/ui/label'
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

    import { ChevronsLeft, ChevronsRight, RefreshCw } from '@lucide/vue';

    export default {
        components: {
            AgGridVue, // Add Vue Data Grid component
            Breadcrumb,
            BreadcrumbItem,
            BreadcrumbLink,
            BreadcrumbList,
            BreadcrumbPage,
            BreadcrumbSeparator,
            Button,
            Input,
            Label,
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
        },
        data() {
            return {
                userRole: null,
                userBlock: true,
                tableData: [],
                moduleid: "",
                moduledesc: "",
                spaceid: "",
                spacedesc: "",
                metric: "",
                eventtime: "",
                eventtimestart: "",
                status: "",
                value: "",
                tableLoading: false,
                gridApi: null,
                componentKey: 0,
                colDefs: [
                    { field: "spacedesc", headerName: this.$t('problem.table01') },
                    { field: "moduledesc", headerName: this.$t('problem.table02') },
                    { field: "metric", headerName: this.$t('problem.table03') },
                    { field: "eventtime", headerName: this.$t('problem.table04') },
                    { field: "eventtimestart", headerName: this.$t('problem.table05') },
                    { field: "value", headerName: this.$t('problem.table06') },
                    {
                        field: "status", 
                        headerName: this.$t('message.status'),
                        cellClass: this.cellClassStatus,
                    },
                ],
                defaultColDef: {
                    minWidth: 120,
                    resizable: true,
                },
            }
        },
        methods: {
            cellClassStatus(params) {
                if (params.value === "NORMAL") {
                    return "problem-green";
                } else if (params.value === "WARNING") {
                    return "problem-yellow";
                } else if (params.value === "CRITICAL") {
                    return "problem-red";
                }
            },
            async problemSelect() {
                this.tableLoading = true;
                try {
                    const response = await axios.get("/api/v1/admin/problem/select?" + Math.random());
                    const problems = Array.isArray(response.data) ? response.data : [];
                    this.tableData = problems.map(item => ({
                        ...item,
                        eventtime: moment(item.eventtime).format("DD.MM.YYYY HH:mm:ss"),
                        eventtimestart: moment(item.eventtimestart).format("DD.MM.YYYY HH:mm:ss"),
                    }));
                } catch (error) {
                    console.error("Failed to load problems:", error);
                    this.tableData = [];
                } finally {
                    this.tableLoading = false;
                }
            },
            async userInfo() {
                try {
                    const response = await axios.get("/api/v1/userinfo?" + Math.random());
                    this.userRole = response.data.role;
                    this.userBlock = false;
                    await this.problemSelect();
                } catch (error) {
                    console.error("Failed to load user info:", error);
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
                    { field: "spacedesc", headerName: this.$t('problem.table01') },
                    { field: "moduledesc", headerName: this.$t('problem.table02') },
                    { field: "metric", headerName: this.$t('problem.table03') },
                    { field: "eventtime", headerName: this.$t('problem.table04') },
                    { field: "eventtimestart", headerName: this.$t('problem.table05') },
                    { field: "value", headerName: this.$t('problem.table06') },
                    {
                        field: "status", 
                        headerName: this.$t('message.status'),
                        cellClass: this.cellClassStatus,
                    },
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
        },
    }
</script>

<style>
.problem-red {
    background-color: var(--va-danger);
    color: white;
    font-weight: bold;
}
.problem-yellow {
    background-color: var(--va-warning);
    color: black;
    font-weight: bold;
}
.problem-green {
    background-color: var(--va-success);
    color: black;
    font-weight: bold;
}
</style>