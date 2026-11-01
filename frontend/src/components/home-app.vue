<template>
<div class="grid w-full gap-2">
<div class="grid w-full items-stretch grid-cols-1 gap-2 landscape:grid-cols-3 lg:auto-rows-fr lg:grid-cols-5!">
    <Card class="h-full w-full">
            <CardHeader>
                <CardTitle>{{ $t('message.time') }}</CardTitle>
            </CardHeader>
        <CardContent class="text-2xl! font-bold!">
            {{ currentDT }}
        </CardContent>
    </Card>

    <RouterLink
        to="/user"
        class="block h-full w-full"
            :aria-disabled="userBlock || undefined"
            :tabindex="userBlock ? -1 : undefined"
            :class="userBlock ? 'pointer-events-none cursor-not-allowed opacity-50' : undefined"
            @click="userBlock && $event.preventDefault()">
        
        <Card :disabled="userBlock" class="h-full w-full">
            <CardHeader>
                <CardTitle>{{ $t('message.users') }}</CardTitle>
            </CardHeader>
            <CardContent>
                <ShieldUser class="size-12" />
            </CardContent>
        </Card>
    </RouterLink>

    <RouterLink
        to="/space"
        class="block w-full">
        
            <Card class="h-full w-full">
            <CardHeader>
                <CardTitle>{{ $t('message.spaces') }}</CardTitle>
            </CardHeader>
            <CardContent>
                <CirclePile class="size-12" />
            </CardContent>
        </Card>
    </RouterLink>

    <RouterLink
        to="/module"
        class="block w-full">
        
            <Card class="h-full w-full">
            <CardHeader>
                <CardTitle>{{ $t('message.modules') }}</CardTitle>
            </CardHeader>
            <CardContent>
                <Server class="size-12" />
            </CardContent>
        </Card>
    </RouterLink>

    <RouterLink
        to="/task"
        class="block w-full">
        
            <Card class="h-full w-full">
            <CardHeader>
                <CardTitle>{{ $t('message.tasks') }}</CardTitle>
            </CardHeader>
            <CardContent>
                <CircleCheckBig class="size-12" />
            </CardContent>
        </Card>
    </RouterLink>

    <RouterLink
        to="/dashboard"
        class="block w-full">
        
            <Card class="h-full w-full">
            <CardHeader>
                <CardTitle>{{ $t('message.dashboards') }}</CardTitle>
            </CardHeader>
            <CardContent>
                <Gauge class="size-12" />
            </CardContent>
        </Card>
    </RouterLink>

    <RouterLink
        to="/notification"
        class="block h-full w-full"
            :aria-disabled="userBlock || undefined"
            :tabindex="userBlock ? -1 : undefined"
            :class="userBlock ? 'pointer-events-none cursor-not-allowed opacity-50' : undefined"
            @click="userBlock && $event.preventDefault()">
        
            <Card :disabled="userBlock" class="h-full w-full">
            <CardHeader>
                <CardTitle>{{ $t('message.notifications') }}</CardTitle>
            </CardHeader>
            <CardContent>
                <BellRing class="size-12" />
            </CardContent>
        </Card>
    </RouterLink>

    <RouterLink
        to="/database"
        class="block w-full">
        
            <Card :class="['h-full w-full', databaseCardClass]">
            <CardHeader>
                <CardTitle>{{ $t('message.statusdb') }}</CardTitle>
            </CardHeader>
            <CardContent class="text-2xl! font-bold!">
                {{ databaseData }}
            </CardContent>
        </Card>
    </RouterLink>

    <Card :class="['h-full w-full', datamanagerCardClass]">
        <CardHeader>
            <CardTitle>{{ $t('message.statusdm') }}</CardTitle>
        </CardHeader>
        <CardContent class="text-2xl! font-bold!">
            {{ datamanagerData }}
        </CardContent>
    </Card>

    <Card :class="['h-full w-full', uimanagerCardClass]">
        <CardHeader>
            <CardTitle>{{ $t('message.statusuim') }}</CardTitle>
        </CardHeader>
        <CardContent class="text-2xl! font-bold!">
            {{ uimanagerData }}
        </CardContent>
    </Card>
</div>

<div class="grid w-full items-stretch grid-cols-1 gap-2 landscape:grid-cols-1 lg:grid-cols-1!">
    <RouterLink
        to="/problem"
        class="block w-full">
        
        <Card class="h-full w-full">
            <CardHeader>
                <CardTitle>{{ $t('message.problems') }}</CardTitle>
            </CardHeader>
            <CardContent>
                <ag-grid-vue :key="componentKey" :rowData="tableData" :columnDefs="colDefs" :defaultColDef="defaultColDef" @grid-ready="onGridReady" style="height: 200px; width: 100%"></ag-grid-vue>
            </CardContent>
        </Card>
    </RouterLink>
</div>

<div class="grid w-full items-stretch grid-cols-1 landscape:grid-cols-1 gap-2 lg:auto-rows-fr lg:grid-cols-2!">
    <RouterLink
        to="/rawdata"
        class="block w-full">
        
        <Card class="h-full w-full">
            <CardHeader>
                <CardTitle>{{ $t('message.metricsr') }}</CardTitle>
            </CardHeader>
            <CardContent>
                <Bar :options="chartOptions" :data="chartData01" chart-id="barChart01" :key="componentKey"></Bar>
            </CardContent>
        </Card>
    </RouterLink>

    <RouterLink
        to="/rawtext"
        class="block w-full">
        
        <Card class="h-full w-full">
            <CardHeader>
                <CardTitle>{{ $t('message.logsr') }}</CardTitle>
            </CardHeader>
            <CardContent>
                <Bar :options="chartOptions" :data="chartData02" chart-id="barChart02" :key="componentKey"></Bar>
            </CardContent>
        </Card>
    </RouterLink>
</div>
</div>
</template>

<script>
    import axios from 'axios'
    import moment from 'moment'
    import { mapState } from 'vuex';
    import { AgGridVue } from "ag-grid-vue3"; // Vue Data Grid Component
    import { Bar } from 'vue-chartjs'
    import { Chart as ChartJS, Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale } from 'chart.js'

    import {
        Card,
        CardAction,
        CardContent,
        CardDescription,
        CardFooter,
        CardHeader,
        CardTitle,
    } from '@/components/ui/card';

    import { ShieldUser, CirclePile, Server, CircleCheckBig, Gauge, BellRing } from '@lucide/vue';
    
    ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale)

    export default {
        name: 'home-app',
        uimanagerIntervalID: null,
        datamanagerIntervalID: null,
        databaseIntervalID: null,
        createDTID: null,
        chartDataIntervalID: null,
        problemSelectID: null,
        components: {
            Card,
            CardAction,
            CardContent,
            CardDescription,
            CardFooter,
            CardHeader,
            CardTitle,
            AgGridVue, // Add Vue Data Grid component
            Bar,
            ShieldUser,
            CirclePile,
            Server,
            CircleCheckBig,
            Gauge,
            BellRing,
        },
        data() {
            return {
                currentDT: this.$t('message.nodata'),
                tab: "Home",
                userRole: null,
                userBlock: true,
                adminBlock: true,
                databaseData: this.$t('message.nodata'),
                databaseColor: "",
                uimanagerData: this.$t('message.nodata'),
                uimanagerColor: "",
                datamanagerData: this.$t('message.nodata'),
                datamanagerColor: "",
                tableData: [],
                modulid: "",
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
                chartData01: {
                    labels: [],
                    datasets: [
                        {
                            label: '',
                            backgroundColor: '',
                            data: [],
                            barPercentage: 1
                        }
                    ]
                },
                chartData02: {
                    labels: [],
                    datasets: [
                        {
                            label: '',
                            backgroundColor: '',
                            data: [],
                            barPercentage: 1,
                        }
                    ]
                },
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
            chartDataAPI(dataSRC) {
                let vm = this;
                let dataPut = {
                    datasrc: dataSRC,
                };
                axios.put("/api/v1/admin/chartdata?" + Math.random(), dataPut)
                    .then(function (response) {
                        const data = vm.normalizeResponseData(response.data);
                        if (data && Array.isArray(data.labels)) {
                            data.labels.forEach(function(item, i) {
                                data.labels[i] = moment(item).format("HH:mm");
                            });
                            if (dataSRC == "raw_data") {
                                vm.chartData01 = data;
                            }
                            if (dataSRC == "raw_text") {
                                vm.chartData02 = data;
                            }
                        }
                    });
            },
            async init() {
                let vm = this;
                try {
                    const response = await axios.get("/api/v1/userinfo?" + Math.random());
                    const data = response.data;
                            vm.userRole = data.role;
                            if (data.role == "superadmin") {
                                vm.userBlock = false
                                vm.adminBlock = false
                            } else if (data.role == "admin") {
                                vm.userBlock = false
                                vm.adminBlock = true
                                vm.databaseData = vm.$t('message.noaccess')
                            }
                            else {
                                vm.userBlock = true
                                vm.adminBlock = true
                                vm.databaseData = vm.$t('message.noaccess')
                                vm.userData = vm.$t('message.noaccess')
                            }
                            
                } catch(err) {
                    console.log("UIManager Error!");
                }

                vm.chartDataAPI("raw_data");
                vm.chartDataAPI("raw_text");
                vm.databaseStatus();
                vm.uimanagerStatus();
                vm.datamanagerStatus();
                vm.problemSelect();

                vm.createDTID = window.setInterval(vm.createDT, 1000);
                vm.uimanagerIntervalID = window.setInterval(vm.uimanagerStatus, 10000);
                vm.datamanagerIntervalID = window.setInterval(vm.datamanagerStatus, 10000);
                vm.databaseIntervalID = window.setInterval(vm.databaseStatus, 10000);
                vm.chartDataIntervalID = window.setInterval(vm.chartDataUpdate, 10000);
                vm.problemSelectID = window.setInterval(vm.problemSelect, 60000);
            },
            chartDataUpdate() {
                this.chartDataAPI("raw_data");
                this.chartDataAPI("raw_text");
            },
            databaseStatus() {
                let vm = this;
                axios.get("/api/v1/admin/database/status?" + Math.random())
                    .then(function (response) {
                        const data = vm.normalizeResponseList(response.data);
                        if (data.length === 0) {
                            vm.databaseData = vm.$t('message.notavailable');
                            vm.databaseColor = "danger";
                            return;
                        }
                        let status01 = 0;
                        let status02 = 0;
                        data.forEach(function(item) {
                            if (item.up) {
                                status01++;
                            } else {
                                status02++;
                            }
                        });
                        vm.databaseData = status01 + " " + vm.$t('message.of') + " " + (status01 + status02);
                        if (status02 > 0) {
                            vm.databaseColor = "danger";
                        } else {
                            vm.databaseColor = "success";
                        }
                    })
                    .catch(function(error) {
                        console.error('Failed to load database status', error);
                        vm.databaseData = vm.$t('message.notavailable');
                        vm.databaseColor = "danger";
                    });
            },
            uimanagerStatus() {
                let vm = this;
                axios.get("/api/v1/version?" + Math.random())
                    .then(function(response) {
                        const data = vm.normalizeResponseData(response.data);
                        if (!data || typeof data.Name !== 'string' || data.Name.length === 0) {
                            throw new Error('Invalid UIManager status response');
                        }
                            vm.uimanagerColor = "success";
                            vm.uimanagerData = vm.$t('message.available');
                    })
                    .catch(function(error) {
                        console.error('Failed to load UIManager status', error);
                        if (error.response?.status === 401) {
                            vm.uimanagerColor = "warning";
                            vm.uimanagerData = vm.$t('message.noaccess');
                        } else {
                            vm.uimanagerColor = "danger";
                            vm.uimanagerData = vm.$t('message.notavailable');
                        }
                    });
            },
            datamanagerStatus() {
                let vm = this;
                axios.get("/api/v1/admin/datamanager/status?" + Math.random())
                    .then(function (response) {
                        const data = vm.normalizeResponseList(response.data);
                        if (data.length === 0) {
                            vm.datamanagerData = vm.$t('message.notavailable');
                            vm.datamanagerColor = "danger";
                            return;
                        }
                        let status01 = 0;
                        let status02 = 0;
                        data.forEach(function(item) {
                            if (item.up) {
                                status01++;
                            } else {
                                status02++;
                            }
                        });
                        vm.datamanagerData = status01 + " " + vm.$t('message.of') + " " + (status01 + status02);
                        if (status02 > 0 && status01 > 0) {
                            vm.datamanagerColor = "warning";
                        } else if (status02 > 0 && status01 == 0) {
                            vm.datamanagerColor = "danger";
                        } else if (status02 == 0 && status01 == 0) {
                            vm.datamanagerColor = "danger";
                        } else {
                            vm.datamanagerColor = "success";
                        }
                    })
                    .catch(function(error) {
                        console.error('Failed to load DataManager status', error);
                        vm.datamanagerData = vm.$t('message.notavailable');
                        vm.datamanagerColor = "danger";
                    });
            },
            cellClassStatus(params) {
                if (params.value === "NORMAL") {
                    return "problem-green";
                } else if (params.value === "WARNING") {
                    return "problem-yellow";
                } else if (params.value === "CRITICAL") {
                    return "problem-red";
                }
            },
            statusCardClass(status) {
                return {
                    'border-[var(--success)]! bg-[var(--success)]! text-primary-foreground!': status === 'success',
                    'border-[var(--warning)]! bg-[var(--warning)]! text-foreground!': status === 'warning',
                    'border-[var(--destructive)]! bg-[var(--destructive)]! text-primary-foreground!': status === 'danger',
                };
            },
            normalizeResponseData(payload) {
                if (typeof payload === 'string') {
                    try {
                        return JSON.parse(payload);
                    } catch {
                        return null;
                    }
                }
                return payload?.data ?? payload;
            },
            normalizeResponseList(payload) {
                const data = this.normalizeResponseData(payload);
                return Array.isArray(data) ? data : [];
            },
            problemSelect() {
                let vm = this;
                vm.tableLoading = true;
                axios.get("/api/v1/admin/problem/select?" + Math.random())
                    .then(function (response) {
                        const data = vm.normalizeResponseList(response.data);
                        let dataNormal = [];
                        data.forEach(function(item, i) {
                            dataNormal[i] = item;
                            dataNormal[i].eventtime = moment(item.eventtime).format("DD.MM.YYYY HH:mm:ss");
                            dataNormal[i].eventtimestart = moment(item.eventtimestart).format("DD.MM.YYYY HH:mm:ss");
                        });
                        vm.tableData = dataNormal;
                        vm.tableLoading = false;
                    })
                    .catch(function() {
                        vm.tableLoading = false;
                    });
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
            createDT() {
                this.currentDT = moment().format('DD.MM.YYYY hh:mm:ss');
            },
        },
        unmounted() {
            clearInterval(this.createDTID);
            clearInterval(this.uimanagerIntervalID);
            clearInterval(this.datamanagerIntervalID);
            clearInterval(this.databaseIntervalID);
            clearInterval(this.chartDataIntervalID);
            clearInterval(this.problemSelectID);
        },
        mounted() {
            window.addEventListener('resize', this.handleResize);
            this.init();
        },
        beforeUnmount() {
            window.removeEventListener('resize', this.handleResize);
        },
        computed: {
            databaseCardClass() {
                return this.statusCardClass(this.databaseColor);
            },
            datamanagerCardClass() {
                return this.statusCardClass(this.datamanagerColor);
            },
            uimanagerCardClass() {
                return this.statusCardClass(this.uimanagerColor);
            },
            getMode() {
                return this.$store.state.mode;
            },
            ...mapState(['locale', 'mode']),

            chartOptions() {
                const themeRoot = this.mode
                    ? document.querySelector('.dark') || document.documentElement
                    : document.documentElement;

                const getColor = (variableName) => {
                    if (typeof window === 'undefined') return '#ccc';

                    return getComputedStyle(themeRoot)
                        .getPropertyValue(variableName)
                        .trim();
                };

                return {
                    responsive: true,
                    maintainAspectRatio: false,
                    color: getColor('--primary'),
                    borderColor: getColor('--primary'),
                    scales: {
                        x: {
                            border: {
                                width: 0,
                            },
                            grid: {
                                color: getColor('--border'),
                                display: true,
                            },
                            ticks: {
                                color: getColor('--primary'),
                            },
                        },
                        y: {
                            border: {
                                width: 0,
                            },
                            grid: {
                                color: getColor('--border'),
                            },
                            ticks: {
                                color: getColor('--primary'),
                            },
                        },
                    },
                    plugins: {
                        legend: {
                            display: false
                        }
                    },
                };
            },
        },
        watch: {
            mode(newValue, oldValue) {
                this.componentKey += 1;
            },
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

<style>
* {
    --ag-row-height: 30px;
    --ag-header-height: 48px;
    --ag-list-item-height: 24px;
}

#barChart01, #barChart02 {
    max-height: 200px;
}
.problem-red {
    background-color: var(--destructive) !important;
    color: white;
    font-weight: bold;
}
.problem-yellow {
    background-color: var(--warning) !important;
    color: black;
    font-weight: bold;
}
.problem-green {
    background-color: var(--success) !important;
    color: black;
    font-weight: bold;
}
</style>