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
                    <BreadcrumbPage>{{ $t('message.rawdata') }}</BreadcrumbPage>
                </BreadcrumbItem>
            </BreadcrumbList>
        </Breadcrumb>
    </div>

    <div class="grid w-full items-stretch grid-cols-1 gap-2 landscape:grid-cols-2 lg:auto-rows-fr lg:grid-cols-2!">
        <div class="grid gap-2">
            <div class="grid gap-2">
                <Label class="flex flex-col gap-3 items-start">
                    {{ $t('rawdata.filter01') }}
                    <Select v-model="spaceValue" v-on:update:model-value="spaceSelectUpdateValue">
                        <SelectTrigger :aria-label="$t('rawdata.filter01')">
                            <SelectValue :placeholder="$t('rawdata.filter01')" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectGroup>
                                <SelectLabel>{{ $t('rawdata.filter01') }}</SelectLabel>
                                <SelectItem v-for="option in spaceOptions" :key="option.value" :value="option">
                                    {{ option.text }}
                                </SelectItem>
                            </SelectGroup>
                        </SelectContent>
                    </Select>
                </Label>
            </div>
        </div>
        <div class="grid gap-2">
            <div class="grid gap-2">
                <Label class="flex flex-col gap-3 items-start">
                    {{ $t('rawdata.filter02') }}
                    <Select v-model="metricValue" v-on:update:model-value="metricSelectUpdateValue">
                        <SelectTrigger :aria-label="$t('rawdata.filter02')">
                            <SelectValue :placeholder="$t('rawdata.filter02')" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectGroup>
                                <SelectLabel>{{ $t('rawdata.filter02') }}</SelectLabel>
                                <SelectItem v-for="option in metricOptions" :key="option.value" :value="option">
                                    {{ option.text }}
                                </SelectItem>
                            </SelectGroup>
                        </SelectContent>
                    </Select>
                </Label>
            </div>
        </div>
    </div>
    <div class="grid w-full items-stretch grid-cols-1 gap-2">
        <Label class="flex flex-col items-start gap-3">
            {{ $t('rawdata.message03') }}
            <Input v-model="input" v-on:update:model-value="filterUpdate" name="input" />
        </Label>
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
        <div class="grid grid-cols-1 gap-1" align="right">
            <Button
                type="button"
                @click="rawDataSelect(pages[currentPage])">
                <RefreshCw class="size-5" />
            </Button>
        </div>
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
                spaceValue: {},
                spaceOptions: [],
                metricOptions: [],
                metricValue: {},
                gridApi: null,
                tableData: [],
                tableLoading: false,
                input: "",
                wrapperSize: "800px",
                componentKey: 0,
                colDefs: [
                    { field: "spacedesc", headerName: this.$t('rawdata.table01') },
                    { field: "object", headerName: this.$t('rawdata.table02') },
                    { field: "metric", headerName: this.$t('rawdata.table03') },
                    { field: "value", headerName: this.$t('rawdata.table04') },
                    { field: "labels", headerName: "Labels" },
                    { field: "eventtime", headerName: this.$t('rawdata.table05') },
                ],
                defaultColDef: {
                    minWidth: 120,
                    resizable: true,
                },

                // Pagination
                selectParam: {},
                isDisabledLeft: true,
                isDisabledRight: true,
                currentPage: 0,
                pages: [""],
            }
        },
        methods: {
            selectRight() {
                this.currentPage++;
                this.rawDataSelect(this.selectParam.next);
                this.pages[this.currentPage] = this.selectParam.next;
            },
            selectLeft() {
                if (this.currentPage > 0) {
                    this.currentPage--;
                }
                this.rawDataSelect(this.pages[this.currentPage]);
            },
            async rawDataSelect(current = "") {
                if (!this.spaceValue.value || !this.metricValue.value) {
                    return;
                }

                this.tableLoading = true;
                const dataPut = {
                    current: current || "",
                    pagesize: 1000,
                    spaceid: this.spaceValue.value,
                    metric: this.metricValue.value,
                };

                try {
                    const response = await axios.put("/api/v1/admin/rawdata/select?" + Math.random(), dataPut);
                    const data = response.data;
                    const dataNormal = Array.isArray(data.data) ? data.data : [];
                    dataNormal.sort((a, b) => b.eventtime - a.eventtime);
                    dataNormal.forEach((item) => {
                        item.createtime = moment(item.createtime).format("DD.MM.YYYY HH:mm:ss");
                        item.eventtime = moment(item.eventtime).format("DD.MM.YYYY HH:mm:ss");
                        item.labels = item.labels != null ? JSON.stringify(item.labels) : "";
                    });
                    this.tableData = dataNormal;
                    this.selectParam.current = data.current;
                    this.selectParam.next = data.next;
                    this.isDisabledRight = !data.next;
                    this.isDisabledLeft = !data.current;
                } catch (error) {
                    this.tableData = [];
                } finally {
                    this.tableLoading = false;
                }
            },
            async spaceSelect() {
                try {
                    const response = await axios.put("/api/v1/admin/space/select?" + Math.random(), { pagesize: 0 });
                    this.spaceOptions = Array.isArray(response.data.data)
                        ? response.data.data.map(item => ({ text: item.description, value: item.id }))
                        : [];
                } catch (error) {
                    this.spaceOptions = [];
                }
            },
            async spaceSelectUpdateValue() {
                this.metricOptions = [];
                this.metricValue = {};
                if (this.spaceValue.value) {
                    try {
                        const response = await axios.put("/api/v1/admin/rawdata/metric/select?" + Math.random(), {
                            spaceid: this.spaceValue.value,
                        });
                        this.metricOptions = Array.isArray(response.data)
                            ? response.data.map(item => ({ text: item.metric, value: item.metric }))
                            : [];
                    } catch (error) {
                        this.metricOptions = [];
                    }
                    this.tableData = [];
                }
            },
            pageUpdateValue() {
                this.rawDataSelect();
            },
            metricSelectUpdateValue() {
                this.input = "";
                if (this.metricValue.value != "") {
                    this.selectParam = {};
                    this.isDisabledLeft = true;
                    this.isDisabledRight = true;
                    this.currentPage = 0;
                    this.pages = [""];
                    this.rawDataSelect();
                }
            },
            filterUpdate() {
                this.gridApi.setGridOption("quickFilterText", this.input);
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
            this.spaceSelect();
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
                    { field: "spacedesc", headerName: this.$t('rawdata.table01') },
                    { field: "object", headerName: this.$t('rawdata.table02') },
                    { field: "metric", headerName: this.$t('rawdata.table03') },
                    { field: "value", headerName: this.$t('rawdata.table04') },
                    { field: "labels", headerName: "Labels" },
                    { field: "eventtime", headerName: this.$t('rawdata.table05') },
                ]
                this.componentKey += 1;
            },
            mode(newValue, oldValue) {
                this.componentKey += 1;
            },
            tableData(newValue) {
                if (newValue.length > 0 && this.gridApi) {
                    setTimeout(() => {
                        // this.gridApi.autoSizeAllColumns();
                        this.gridApi.sizeColumnsToFit();
                    }, 50);
                }
            },
        }
    }
</script>

<style></style>