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
                    <BreadcrumbPage>{{ $t('message.dashboards') }}</BreadcrumbPage>
                </BreadcrumbItem>
            </BreadcrumbList>
        </Breadcrumb>
    </div>

    <div class="flex! items-center! justify-end gap-3 w-full">
        <div class="grid grid-cols-1 gap-1" align="right">
            <Button
                type="button"
                @click="buttonClickNew">
                <Plus class="size-5" />
            </Button>
        </div>
    </div>

    <Dialog v-model:open="showModalCreate">
        <DialogContent class="fixed! inset-0! top-0! left-0! z-50 flex! h-dvh! w-full! max-w-none! max-h-none! translate-x-0! translate-y-0! flex-col! gap-0! overflow-hidden! rounded-none! border-0! p-0! bg-card text-foreground! lg:inset-auto! lg:top-[50%]! lg:left-[50%]! lg:grid! lg:h-auto! lg:max-w-lg! lg:max-h-[90vh]! lg:translate-x-[-50%]! lg:translate-y-[-50%]! lg:gap-4! lg:overflow-visible! lg:rounded-lg! lg:border! lg:p-6!">
            <DialogHeader class="px-6 pt-6 pr-12 text-left lg:px-0 lg:pt-0 lg:pr-0">
                <DialogTitle>{{ $t('dashboard.message07') }}</DialogTitle>
            </DialogHeader>
            <div class="grid min-h-0 flex-1 content-start gap-4 w-full overflow-y-auto p-6 lg:flex-none lg:max-w-fit lg:p-0">
                <div class="grid gap-4 w-full">
                    <Label for="dashboardName">{{ $t('dashboard.message02') }}</Label>
                    <Input id="dashboardName" v-model="dashboardName" name="dashboardName" />
                </div>

                <div class="grid grid-cols-1 gap-4 w-full lg:grid-cols-2">
                    <div class="flex flex-col gap-3 w-full">
                        <Label for="dashboardStartDate" class="px-1">
                            {{ $t('dashboard.message03') }}
                        </Label>
                        <Popover v-model:open="showCalendar01">
                            <PopoverTrigger as-child>
                                <Button
                                    id="dashboardStartDate"
                                    variant="outline"
                                    class="w-full justify-between font-normal bg-popover">
                                        {{ dashboardStartDate ? formatDate(dashboardStartDate) : "Выберите дату" }}
                                        <ChevronDownIcon />
                                </Button>
                            </PopoverTrigger>
                            <PopoverContent class="w-full overflow-hidden p-0" align="start">
                                <Calendar
                                    :model-value="dashboardStartDate"
                                    @update:model-value="(value) => {
                                        if (value) {
                                            dashboardStartDate = value
                                            showCalendar01 = false
                                        }
                                    }"
                                    :locale="locale"
                                />
                            </PopoverContent>
                        </Popover>
                    </div>
                    <div class="flex flex-col gap-3 w-full">
                        <Label for="dashboardStartTime" class="px-1">
                            {{ $t('dashboard.message05') }}
                        </Label>
                        <Input
                            v-model="dashboardStartTime"
                            id="dashboardStartTime"
                            type="time"
                            step="1"
                            class="appearance-none [&::-webkit-calendar-picker-indicator]:hidden [&::-webkit-calendar-picker-indicator]:appearance-none"
                        />
                    </div>
                </div>

                <div class="grid grid-cols-1 gap-4 w-full lg:grid-cols-2">
                    <div class="flex flex-col gap-3 w-full">
                        <Label for="dashboardStopDate" class="px-1">
                            {{ $t('dashboard.message04') }}
                        </Label>
                        <Popover v-model:open="showCalendar02">
                            <PopoverTrigger as-child>
                                <Button
                                    id="dashboardStopDate"
                                    variant="outline"
                                    class="w-full justify-between font-normal bg-popover">
                                        {{ dashboardStopDate ? formatDate(dashboardStopDate) : "Выберите дату" }}
                                        <ChevronDownIcon />
                                </Button>
                            </PopoverTrigger>
                            <PopoverContent class="w-full overflow-hidden p-0" align="start">
                                <Calendar
                                    :model-value="dashboardStopDate"
                                    @update:model-value="(value) => {
                                        if (value) {
                                            dashboardStopDate = value
                                            showCalendar02 = false
                                        }
                                    }"
                                    :locale="locale"
                                />
                            </PopoverContent>
                        </Popover>
                    </div>
                    <div class="flex flex-col gap-3 w-full">
                        <Label for="dashboardStopTime" class="px-1">
                            {{ $t('dashboard.message06') }}
                        </Label>
                        <Input
                            v-model="dashboardStopTime"
                            id="dashboardStopTime"
                            type="time"
                            step="1"
                            class="appearance-none [&::-webkit-calendar-picker-indicator]:hidden [&::-webkit-calendar-picker-indicator]:appearance-none"
                        />
                    </div>
                </div>

                <Label for="dashboardInterval" class="px-1">
                    {{ $t('dashboard.message37') }}
                </Label>
                <div class="grid grid-cols-3 gap-1 w-full items-stretch lg:grid-cols-6">
                    <Button type="button" @click="buttonClickInterval(3600)">{{ $t('dashboard.message38') }}</Button>
                    <Button type="button" @click="buttonClickInterval(10800)">{{ $t('dashboard.message39') }}</Button>
                    <Button type="button" @click="buttonClickInterval(21600)">{{ $t('dashboard.message40') }}</Button>
                    <Button type="button" @click="buttonClickInterval(43200)">{{ $t('dashboard.message41') }}</Button>
                    <Button type="button" @click="buttonClickInterval(86400)">{{ $t('dashboard.message42') }}</Button>
                    <Button type="button" @click="buttonClickInterval(259200)">{{ $t('dashboard.message43') }}</Button>
                </div>
                <div class="grid gap-4 w-full">
                    <Label for="dashboardInterval">{{ $t('dashboard.message45') }}</Label>
                    <Input id="dashboardInterval" v-model="dashboardInterval" name="dashboardInterval" />
                </div>
            </div>
            <DialogFooter class="mt-auto flex-row px-6 pb-[max(1.5rem,env(safe-area-inset-bottom))] *:flex-1 lg:mt-0 lg:px-0 lg:pb-0 lg:*:flex-none">
                <DialogClose as-child>
                    <Button>{{ $t('message.cancel') }}</Button>
                </DialogClose>
                <Button @click="putDashboardCreate()">{{ $t('message.save') }}</Button> 
            </DialogFooter>
        </DialogContent>
    </Dialog>

    <InnerLoading :loading="innerLoading">
        <div class="grid w-full items-stretch grid-cols-1 gap-2 landscape:grid-cols-2 lg:auto-rows-fr lg:grid-cols-2!">
            <Dialog v-if="activeDashboardIndex !== null" v-model:open="showModalRemove[activeDashboardIndex]">
                    <DialogContent class="fixed! inset-0! top-0! left-0! z-50 flex! h-dvh! w-full! max-w-none! max-h-none! translate-x-0! translate-y-0! flex-col! gap-0! overflow-hidden! rounded-none! border-0! p-0! bg-card text-foreground! lg:inset-auto! lg:top-[50%]! lg:left-[50%]! lg:grid! lg:h-auto! lg:max-w-lg! lg:max-h-[80vh]! lg:translate-x-[-50%]! lg:translate-y-[-50%]! lg:gap-4! lg:overflow-visible! lg:rounded-lg! lg:border! lg:p-6!">
                        <DialogHeader class="px-6 pt-6 pr-12 text-left lg:px-0 lg:pt-0 lg:pr-0">
                            <DialogTitle>{{ $t('dashboard.message08') }}</DialogTitle>
                        </DialogHeader>
                        <div class="flex-1 p-6 lg:flex-none lg:p-0">{{ $t('dashboard.message10') }} {{ activeDashboardIndex + 1 }} ({{ dashboards[activeDashboardIndex]?.name }})</div>
                        <DialogFooter class="mt-auto flex-row px-6 pb-[max(1.5rem,env(safe-area-inset-bottom))] *:flex-1 lg:mt-0 lg:px-0 lg:pb-0 lg:*:flex-none">
                            <DialogClose as-child>
                                <Button>{{ $t('message.cancel') }}</Button>
                            </DialogClose>
                            <Button @click="putDashboardRemove(activeDashboardIndex)">{{ $t('message.delete') }}</Button>
                        </DialogFooter>
                    </DialogContent>
            </Dialog>

            <Dialog v-if="activeDashboardIndex !== null" v-model:open="showModalEdit[activeDashboardIndex]">
                    <DialogContent class="fixed! inset-0! top-0! left-0! z-50 flex! h-dvh! w-full! max-w-none! max-h-none! translate-x-0! translate-y-0! flex-col! gap-0! overflow-hidden! rounded-none! border-0! p-0! bg-card text-foreground! lg:inset-auto! lg:top-[50%]! lg:left-[50%]! lg:grid! lg:h-auto! lg:max-w-lg! lg:max-h-[90vh]! lg:translate-x-[-50%]! lg:translate-y-[-50%]! lg:gap-4! lg:overflow-visible! lg:rounded-lg! lg:border! lg:p-6!">
                        <DialogHeader class="px-6 pt-6 pr-12 text-left lg:px-0 lg:pt-0 lg:pr-0">
                            <DialogTitle>{{ $t('dashboard.message09') }}</DialogTitle>
                        </DialogHeader>
                        <div class="grid min-h-0 flex-1 content-start gap-4 w-full overflow-y-auto p-6 lg:flex-none lg:max-w-fit lg:p-0">
                            <div class="grid gap-4 w-full">
                                <Label for="dashboardEditName">{{ $t('dashboard.message02') }}</Label>
                                <Input id="dashboardEditName" v-model="dashboardEditName" name="dashboardEditName" />
                            </div>

                            <div class="grid grid-cols-1 gap-4 w-full lg:grid-cols-2">
                                <div class="flex flex-col gap-3 w-full">
                                    <Label for="dashboardEditStartDate" class="px-1">
                                        {{ $t('dashboard.message03') }}
                                    </Label>
                                    <Popover v-model:open="showCalendar01">
                                        <PopoverTrigger as-child>
                                            <Button
                                                id="dashboardEditStartDate"
                                                variant="outline"
                                                class="w-full justify-between font-normal bg-popover">
                                                    {{ dashboardEditStartDate ? formatDate(dashboardEditStartDate) : "Выберите дату" }}
                                                    <ChevronDownIcon />
                                            </Button>
                                        </PopoverTrigger>
                                        <PopoverContent class="w-full overflow-hidden p-0" align="start">
                                            <Calendar
                                                :model-value="dashboardEditStartDate"
                                                @update:model-value="(value) => {
                                                    if (value) {
                                                        dashboardEditStartDate = value
                                                        showCalendar01 = false
                                                    }
                                                }"
                                                :locale="locale"
                                            />
                                        </PopoverContent>
                                    </Popover>
                                </div>
                                <div class="flex flex-col gap-3 w-full">
                                    <Label for="dashboardEditStartTime" class="px-1">
                                        {{ $t('dashboard.message05') }}
                                    </Label>
                                    <Input
                                        v-model="dashboardEditStartTime"
                                        id="dashboardEditStartTime"
                                        type="time"
                                        step="1"
                                        class="appearance-none [&::-webkit-calendar-picker-indicator]:hidden [&::-webkit-calendar-picker-indicator]:appearance-none"
                                    />
                                </div>
                            </div>

                            <div class="grid grid-cols-1 gap-4 w-full lg:grid-cols-2">
                                <div class="flex flex-col gap-3 w-full">
                                    <Label for="dashboardEditStopDate" class="px-1">
                                        {{ $t('dashboard.message04') }}
                                    </Label>
                                    <Popover v-model:open="showCalendar02">
                                        <PopoverTrigger as-child>
                                            <Button
                                                id="dashboardEditStopDate"
                                                variant="outline"
                                                class="w-full justify-between font-normal bg-popover">
                                                    {{ dashboardEditStopDate ? formatDate(dashboardEditStopDate) : "Выберите дату" }}
                                                    <ChevronDownIcon />
                                            </Button>
                                        </PopoverTrigger>
                                        <PopoverContent class="w-full overflow-hidden p-0" align="start">
                                            <Calendar
                                                :model-value="dashboardEditStopDate"
                                                @update:model-value="(value) => {
                                                    if (value) {
                                                        dashboardEditStopDate = value
                                                        showCalendar02 = false
                                                    }
                                                }"
                                                :locale="locale"
                                            />
                                        </PopoverContent>
                                    </Popover>
                                </div>
                                <div class="flex flex-col gap-3 w-full">
                                    <Label for="dashboardEditStopTime" class="px-1">
                                        {{ $t('dashboard.message06') }}
                                    </Label>
                                    <Input
                                        v-model="dashboardEditStopTime"
                                        id="dashboardEditStopTime"
                                        type="time"
                                        step="1"
                                        class="appearance-none [&::-webkit-calendar-picker-indicator]:hidden [&::-webkit-calendar-picker-indicator]:appearance-none"
                                    />
                                </div>
                            </div>

                            <Label for="dashboardEditInterval" class="px-1">
                                {{ $t('dashboard.message37') }}
                            </Label>
                            <div class="grid grid-cols-3 gap-1 w-full items-stretch lg:grid-cols-6">
                                <Button type="button" @click="buttonClickEditInterval(3600)">{{ $t('dashboard.message38') }}</Button>
                                <Button type="button" @click="buttonClickEditInterval(10800)">{{ $t('dashboard.message39') }}</Button>
                                <Button type="button" @click="buttonClickEditInterval(21600)">{{ $t('dashboard.message40') }}</Button>
                                <Button type="button" @click="buttonClickEditInterval(43200)">{{ $t('dashboard.message41') }}</Button>
                                <Button type="button" @click="buttonClickEditInterval(86400)">{{ $t('dashboard.message42') }}</Button>
                                <Button type="button" @click="buttonClickEditInterval(259200)">{{ $t('dashboard.message43') }}</Button>
                            </div>
                            <div class="grid gap-4 w-full">
                                <Label for="dashboardEditInterval">{{ $t('dashboard.message45') }}</Label>
                                <Input id="dashboardEditInterval" v-model="dashboardEditInterval" name="dashboardEditInterval" />
                            </div>

                            <Separator class="my-4" />

                            <Label for="chartNew" class="px-1">
                                {{ $t('dashboard.message11') }}
                            </Label>
                            <Button id="chartNew" type="button" @click="buttonClickChartNew(activeDashboardIndex)">
                                <Plus class="size-5" />
                            </Button>

                            <div class="grid gap-4 w-full" v-for="(chartData, index02) in dashboardsEdit[activeDashboardIndex]?.chartdata ?? []" :key="index02">
                                <Separator class="my-4" />
                                {{ $t('dashboard.message46') }} {{ index02 + 1 }}

                                <Label class="flex flex-col gap-3 items-start">
                                    {{ $t('dashboard.message29') }}
                                    <Select :model-value="spaceValue[activeDashboardIndex][index02]" @update:model-value="(value) => onSpaceSelectChange(activeDashboardIndex, index02, value)">
                                        <SelectTrigger>
                                            <SelectValue :placeholder="$t('dashboard.message30')" />
                                        </SelectTrigger>
                                        <SelectContent>
                                            <SelectGroup>
                                                <SelectLabel>{{ $t('dashboard.message29') }}</SelectLabel>
                                                <SelectItem v-for="option in spaceOptions" :key="option.value" :value="option.value">
                                                    {{ option.text }}
                                                </SelectItem>
                                            </SelectGroup>
                                        </SelectContent>
                                    </Select>
                                </Label>

                                <Label class="flex flex-col gap-3 items-start">
                                    {{ $t('dashboard.message31') }}
                                    <Select v-model="metricValue[activeDashboardIndex][index02]">
                                        <SelectTrigger>
                                            <SelectValue :placeholder="$t('dashboard.message32')" />
                                        </SelectTrigger>
                                        <SelectContent>
                                            <SelectGroup>
                                                <SelectLabel>{{ $t('dashboard.message31') }}</SelectLabel>
                                                <SelectItem v-for="option in metricOptions[activeDashboardIndex][index02]" :value="option.value">{{ option.text }}</SelectItem>
                                            </SelectGroup>
                                        </SelectContent>
                                    </Select>
                                </Label>

                                <Label class="flex flex-col gap-3 items-start">
                                    {{ $t('dashboard.message33') }}
                                    <Select v-model="graphtypeValue[activeDashboardIndex][index02]">
                                        <SelectTrigger>
                                            <SelectValue :placeholder="$t('dashboard.message33')" />
                                        </SelectTrigger>
                                        <SelectContent>
                                            <SelectGroup>
                                                <SelectLabel>{{ $t('dashboard.message33') }}</SelectLabel>
                                                <SelectItem v-for="option in graphtypeOptions" :key="option.value" :value="option.value">
                                                    {{ option.text }}
                                                </SelectItem>
                                            </SelectGroup>
                                        </SelectContent>
                                    </Select>
                                </Label>

                                <Label :for="'graphColor'+index02" class="px-1">
                                    {{ $t('dashboard.message34') }}
                                </Label>
                                <div class="grid grid-cols-5 gap-2 w-full">
                                    <Button type="button" title="--chart-1" aria-label="--chart-1" class="bg-chart-1" @click="buttonClickColor(activeDashboardIndex, index02, '--chart-1')"><Palette class="size-5" /></Button>
                                    <Button type="button" title="--chart-2" aria-label="--chart-2" class="bg-chart-2" @click="buttonClickColor(activeDashboardIndex, index02, '--chart-2')"><Palette class="size-5" /></Button>
                                    <Button type="button" title="--chart-3" aria-label="--chart-3" class="bg-chart-3" @click="buttonClickColor(activeDashboardIndex, index02, '--chart-3')"><Palette class="size-5" /></Button>
                                    <Button type="button" title="--chart-4" aria-label="--chart-4" class="bg-chart-4" @click="buttonClickColor(activeDashboardIndex, index02, '--chart-4')"><Palette class="size-5" /></Button>
                                    <Button type="button" title="--chart-5" aria-label="--chart-5" class="bg-chart-5" @click="buttonClickColor(activeDashboardIndex, index02, '--chart-5')"><Palette class="size-5" /></Button>
                                </div>
                                <div class="grid gap-4 w-full">
                                    <Input :id="'graphColor'+index02" v-model="dashboardsEdit[activeDashboardIndex].chartdata[index02].graphcolor" name="graphColor" />
                                </div>

                                <Label class="flex flex-col gap-3 items-start">
                                    {{ $t('dashboard.message35') }}
                                    <Select v-model="groupfuncValue[activeDashboardIndex][index02]">
                                        <SelectTrigger>
                                            <SelectValue :placeholder="$t('dashboard.message35')" />
                                        </SelectTrigger>
                                        <SelectContent>
                                            <SelectGroup>
                                                <SelectLabel>{{ $t('dashboard.message35') }}</SelectLabel>
                                                <SelectItem v-for="option in groupfuncOptions" :key="option.value" :value="option.value">
                                                    {{ option.text }}
                                                </SelectItem>
                                            </SelectGroup>
                                        </SelectContent>
                                    </Select>
                                </Label>

                                <div class="grid gap-4 w-full">
                                    <Label :for="'graphOrder'+index02">{{ $t('dashboard.message36') }}</Label>
                                    <Input :id="'graphOrder'+index02" v-model="dashboardsEdit[activeDashboardIndex].chartdata[index02].graphorder" name="graphOrder" />
                                </div>

                                <Button type="button" @click="buttonClickChartRemove(activeDashboardIndex, index02)">
                                    <Trash class="size-5" />
                                </Button>
                            </div>
                            <Separator class="my-4" />
                        </div>
                        <DialogFooter class="mt-auto flex-row px-6 pb-[max(1.5rem,env(safe-area-inset-bottom))] *:flex-1 lg:mt-0 lg:px-0 lg:pb-0 lg:*:flex-none">
                            <DialogClose as-child>
                                <Button>{{ $t('message.cancel') }}</Button>
                            </DialogClose>
                            <Button @click="putDashboardUpdate(activeDashboardIndex)">{{ $t('message.save') }}</Button>
                        </DialogFooter>
                    </DialogContent>
            </Dialog>

            <div v-for="(dashboard, index01) in dashboards" :key="index01">
                <Card class="h-full w-full">
                    <CardHeader>
                        <CardTitle>
                            <div class="grid grid-cols-2 items-center! justify-between! gap-3 w-full">
                                <div class="flex gap-1 justify-start">
                                    {{ index01 + 1 }}
                                    {{ dashboardTitle[index01] }}
                                </div>
                                <div class="flex gap-1 justify-end">
                                    <Button
                                        type="button"
                                        @click="putDashboardUpdateForm(index01)">
                                        <SquarePen class="size-5" />
                                    </Button>
                                    <Button
                                        type="button"
                                        @click="putDashboardRemoveConfirm(index01)">
                                        <Trash class="size-5" />
                                    </Button>
                                </div>
                            </div>
                        </CardTitle>
                    </CardHeader>
                    <CardContent>
                        <Bar v-if="barLoading" :options="chartOptions" :data="chartDataArr[index01]" chart-id="index" height="300vh" :key="componentKey"></Bar>
                    </CardContent>
                </Card>
            </div>
        </div>
    </InnerLoading>
    </div>
</template>

<script>
    import axios from 'axios'
    import moment from 'moment'
    // import { CalendarDate } from '@internationalized/date'
    import { mapState } from 'vuex'
    import { Bar } from 'vue-chartjs'
    import { Chart as ChartJS, Title, Tooltip, Legend, BarElement, LineElement, CategoryScale, LinearScale, LineController, PointElement } from 'chart.js'

    import { InnerLoading } from '@/components/ui/innerloading'
    import { Button } from '@/components/ui/button'
    import { Input } from '@/components/ui/input'
    import { Label } from '@/components/ui/label'
    import { Toaster } from '@/components/ui/sonner'
    import { toast } from 'vue-sonner'
    import {
        Breadcrumb,
        BreadcrumbItem,
        BreadcrumbLink,
        BreadcrumbList,
        BreadcrumbPage,
        BreadcrumbSeparator,
    } from '@/components/ui/breadcrumb'
    import {
        Card,
        CardAction,
        CardContent,
        CardDescription,
        CardFooter,
        CardHeader,
        CardTitle,
    } from '@/components/ui/card'
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
        Select,
        SelectContent,
        SelectGroup,
        SelectItem,
        SelectLabel,
        SelectTrigger,
        SelectValue,
    } from '@/components/ui/select'
    import { Calendar } from '@/components/ui/calendar'
    import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
    import { Separator } from '@/components/ui/separator'

    import { Plus, SquarePen, Trash, ChevronDownIcon, Palette } from '@lucide/vue';

    ChartJS.register(Title, Tooltip, Legend, BarElement, LineElement, CategoryScale, LinearScale, LineController, PointElement)

    export default {
        components: {
            Bar,
            Breadcrumb,
            BreadcrumbItem,
            BreadcrumbLink,
            BreadcrumbList,
            BreadcrumbPage,
            BreadcrumbSeparator,
            Card,
            CardAction,
            CardContent,
            CardDescription,
            CardFooter,
            CardHeader,
            CardTitle,
            Button,
            Input,
            Label,
            Dialog,
            DialogClose,
            DialogContent,
            DialogDescription,
            DialogFooter,
            DialogHeader,
            DialogTitle,
            Calendar,
            Popover,
            PopoverContent,
            PopoverTrigger,
            Separator,
            Select,
            SelectContent,
            SelectGroup,
            SelectItem,
            SelectLabel,
            SelectTrigger,
            SelectValue,
            Toaster,
            InnerLoading,

            Plus,
            SquarePen,
            Trash,
            ChevronDownIcon,
            Palette,
        },
        chartDataIntervalID: null,
        data() {
            const graphtypeOptions = [
                {text: "Bar", value: "bar"},
                {text: "Line", value: "line"},
            ];
            const groupfuncOptions = [
                {text: "Avg", value: "avg"},
            ];

            return {
                spaceOptions: [],
                spaceValue: [],
                metricOptions: [],
                metricValue: [],
                graphtypeOptions: graphtypeOptions,
                graphtypeValue: [],
                groupfuncOptions: groupfuncOptions,
                groupfuncValue: [],

                dashboardTitle: [],
                dashboardName: "",
                dashboardStartDate: undefined,
                dashboardStartTime: undefined,
                dashboardStopDate: undefined,
                dashboardStopTime: undefined,
                dashboardInterval: undefined,

                dashboardEditName: "",
                dashboardEditStartDate: undefined,
                dashboardEditStartTime: undefined,
                dashboardEditStopDate: undefined,
                dashboardEditStopTime: undefined,
                dashboardEditInterval: undefined,

                /* Переменные состояния модальных окон */
                showModalCreate: false,
                showModalRemove: [],
                showModalEdit: [],
                activeDashboardIndex: null,

                showCalendar01: false,
                showCalendar02: false,

                stripeColorChart: [],

                dashboards: [],
                dashboardsEdit: [],
                innerLoading:true,
                barLoading: false,
                chartDataArr: [],
                componentKey: 0,
            }
        },
        methods: {
            truncateString(str, maxLength = 30) {
                if (str.length > maxLength) {
                    return str.substring(0, maxLength) + '...';
                }
                return str;
            },
            toCalendarDate(value) {
                const date = value instanceof Date ? value : new Date(value);
                return new CalendarDate(date.getFullYear(), date.getMonth() + 1, date.getDate());
            },
            toNativeDate(value) {
                if (value instanceof Date) {
                    return value;
                }

                return new Date(value.year, value.month - 1, value.day);
            },
            getTimeParts(value) {
                if (value instanceof Date) {
                    return {
                        hours: value.getHours(),
                        minutes: value.getMinutes(),
                        seconds: value.getSeconds(),
                    };
                }

                const [hours, minutes, seconds = 0] = value.split(':').map(Number);
                return { hours, minutes, seconds };
            },
            formatTimeInput(value) {
                const timeParts = this.getTimeParts(value);
                return [timeParts.hours, timeParts.minutes, timeParts.seconds]
                    .map((part) => String(part).padStart(2, '0'))
                    .join(':');
            },
            buttonClickColor(index01, index02, variableName) {
                const cssColor = getComputedStyle(document.documentElement)
                    .getPropertyValue(variableName)
                    .trim();
                if (!cssColor || !CSS.supports('color', cssColor)) return;

                const context = document.createElement('canvas').getContext('2d', {
                    willReadFrequently: true,
                });
                if (!context) return;

                context.fillStyle = cssColor;
                context.fillRect(0, 0, 1, 1);
                const [red, green, blue] = context.getImageData(0, 0, 1, 1).data;
                const toHex = (value) => value.toString(16).padStart(2, '0');

                this.dashboardsEdit[index01].chartdata[index02].graphcolor =
                    `#${toHex(red)}${toHex(green)}${toHex(blue)}`;
            },
            buttonClickInterval(sec) {
                this.dashboardInterval = sec;
                this.dashboardStartDate = undefined;
                this.dashboardStartTime = undefined;
            },
            buttonClickEditInterval(sec) {
                this.dashboardEditInterval = sec;
                this.dashboardEditStartDate = undefined;
                this.dashboardEditStartTime = undefined;
            },
            buttonClickNew() {
                this.dashboardName = "";
                this.dashboardStartDate = undefined;
                this.dashboardStopDate = undefined;
                this.dashboardStartTime = undefined;
                this.dashboardStopTime = undefined;
                this.dashboardInterval = undefined;

                this.showModalCreate = true;
            },
            buttonClickChartNew(index) {
                let chartSize = 0;
                if (this.dashboardsEdit[index].chartdata == null) {
                    this.dashboardsEdit[index].chartdata = [];
                    chartSize = this.dashboardsEdit[index].chartdata.push(
                        {
                            id: "",
                            metric: "",
                            spaceid: "",
                            graphtype: "",
                            graphcolor: "#ade2ffbe",
                            groupfunc: "",
                            graphorder: 0
                        }
                    );

                    let chart01 = [];
                    chart01[0] = [];
                    this.metricOptions[index] = chart01;
                    let chart02 = [];
                    chart02[0] = "";
                    this.metricValue[index] = chart02;
                    let chart03 = [];
                    chart03[0] = "";
                    this.spaceValue[index] = chart03;
                    let chart04 = [];
                    chart04[0] = "";
                    this.graphtypeValue[index] = chart04;
                    let chart05 = [];
                    chart05[0] = "";
                    this.groupfuncValue[index] = chart05;
                } else {
                    chartSize = this.dashboardsEdit[index].chartdata.push(
                        {
                            id: "",
                            metric: "",
                            spaceid: "",
                            graphtype: "",
                            graphcolor: "#ade2ffbe",
                            groupfunc: "",
                            graphorder: 0
                        }
                    );

                    this.metricOptions[index][chartSize - 1] = [];
                    this.metricValue[index][chartSize -1] = "";
                    this.spaceValue[index][chartSize -1] = "";
                    this.graphtypeValue[index][chartSize -1] = "";
                    this.groupfuncValue[index][chartSize -1] = "";
                }

                this.stripeColorChart[chartSize - 1] = "warning";
            },
            buttonClickChartRemove(index01, index02) {
                if (!this.dashboardsEdit?.[index01]?.chartdata) return;

                this.dashboardsEdit[index01].chartdata.splice(index02, 1);

                if (this.dashboardsEdit[index01].chartdata == null) return;

                const syncSelectionState = () => {
                    this.dashboardsEdit[index01].chartdata.forEach((item, i) => {
                        this.spaceOptions.forEach((space) => {
                            if (space.value == this.dashboardsEdit[index01].chartdata[i].spaceid) {
                                this.spaceValue[index01] = [...(this.spaceValue[index01] ?? [])];
                                this.spaceValue[index01][i] = space.value;
                                this.$nextTick(() => this.spaceSelectUpdateValue(index01, i));
                                return;
                            }
                        });

                        if (Array.isArray(this.metricOptions[index01]) && Array.isArray(this.metricOptions[index01][i])) {
                            this.metricOptions[index01][i].forEach((metricOption) => {
                                if (metricOption.value == this.dashboardsEdit[index01].chartdata[i].metric) {
                                    this.metricValue[index01] = [...(this.metricValue[index01] ?? [])];
                                    this.metricValue[index01][i] = metricOption.value;
                                    return;
                                }
                            });
                        }

                        this.graphtypeOptions.forEach((option) => {
                            if (option.value == this.dashboardsEdit[index01].chartdata[i].graphtype) {
                                this.graphtypeValue[index01] = [...(this.graphtypeValue[index01] ?? [])];
                                this.graphtypeValue[index01][i] = option.value;
                                return;
                            }
                        });

                        this.groupfuncOptions.forEach((option) => {
                            if (option.value == this.dashboardsEdit[index01].chartdata[i].groupfunc) {
                                this.groupfuncValue[index01] = [...(this.groupfuncValue[index01] ?? [])];
                                this.groupfuncValue[index01][i] = option.value;
                                return;
                            }
                        });
                    });
                };

                this.$nextTick(syncSelectionState);
            },
            putDashboardRemoveConfirm(index) {
                this.activeDashboardIndex = index;
                this.showModalRemove[index] = true;
            },
            putDashboardUpdateForm(index) {
                this.activeDashboardIndex = index;
                this.dashboardsEdit = JSON.parse(JSON.stringify(this.dashboards));
                this.dashboardEditName = this.dashboardsEdit[index].name;

                if (this.dashboardsEdit[index].starttime == 0 || this.dashboardsEdit[index].starttime == undefined) {
                    this.dashboardEditStartDate = undefined;
                    this.dashboardEditStartTime = undefined;
                } else {
                    const startDate = new Date(this.dashboardsEdit[index].starttime);
                    this.dashboardEditStartDate = this.toCalendarDate(startDate);
                    this.dashboardEditStartTime = this.formatTimeInput(startDate);
                }
                if (this.dashboardsEdit[index].stoptime == 0 || this.dashboardsEdit[index].stoptime == undefined) {
                    this.dashboardEditStopDate = undefined;
                    this.dashboardEditStopTime = undefined;
                } else {
                    const stopDate = new Date(this.dashboardsEdit[index].stoptime);
                    this.dashboardEditStopDate = this.toCalendarDate(stopDate);
                    this.dashboardEditStopTime = this.formatTimeInput(stopDate);
                }

                if (this.dashboardsEdit[index].interval == 0) {
                    this.dashboardEditInterval = undefined;
                } else {
                    this.dashboardEditInterval = this.dashboardsEdit[index].interval;
                }

                /* Инициализация metricOptions и metricValue для формы */
                if (this.dashboardsEdit[index].chartdata != null) {
                    const chart01 = [];
                    const chart02 = [];
                    const chart03 = [];
                    const chart04 = [];
                    const chart05 = [];
                    this.dashboardsEdit[index].chartdata.forEach((item, i) => {
                        chart01[i] = [];
                        chart02[i] = "";
                        chart03[i] = "";
                        chart04[i] = "";
                        chart05[i] = "";
                    });
                    this.metricOptions[index] = chart01;
                    this.metricValue[index] = chart02;
                    this.spaceValue[index] = chart03;
                    this.graphtypeValue[index] = chart04;
                    this.groupfuncValue[index] = chart05;
                }

                if (this.dashboards[index].chartdata != null) {
                    const syncEditFormSelection = () => {
                        this.dashboards[index].chartdata.forEach((item, i) => {
                            this.stripeColorChart[i] = "success";

                            this.spaceOptions.forEach((spaceOption) => {
                                if (spaceOption.value == this.dashboardsEdit[index].chartdata[i].spaceid) {
                                    this.spaceValue[index] = [...(this.spaceValue[index] ?? [])];
                                    this.spaceValue[index][i] = spaceOption.value;
                                    return;
                                }
                            });

                            if (Array.isArray(this.metricOptions[index]) && Array.isArray(this.metricOptions[index][i])) {
                                this.metricOptions[index][i].forEach((metricOption) => {
                                    if (metricOption.value == this.dashboardsEdit[index].chartdata[i].metric) {
                                        this.metricValue[index] = [...(this.metricValue[index] ?? [])];
                                        this.metricValue[index][i] = metricOption.value;
                                        return;
                                    }
                                });
                            }

                            this.graphtypeOptions.forEach((option) => {
                                if (option.value == this.dashboardsEdit[index].chartdata[i].graphtype) {
                                    this.graphtypeValue[index] = [...(this.graphtypeValue[index] ?? [])];
                                    this.graphtypeValue[index][i] = option.value;
                                    return;
                                }
                            });

                            this.groupfuncOptions.forEach((option) => {
                                if (option.value == this.dashboardsEdit[index].chartdata[i].groupfunc) {
                                    this.groupfuncValue[index] = [...(this.groupfuncValue[index] ?? [])];
                                    this.groupfuncValue[index][i] = option.value;
                                    return;
                                }
                            });

                            if (this.dashboardsEdit[index].chartdata[i].spaceid) {
                                this.$nextTick(() => this.spaceSelectUpdateValue(index, i));
                            }

                            if (this.dashboardsEdit[index].chartdata[i].graphcolor == "") {
                                this.dashboardsEdit[index].chartdata[i].graphcolor = "#ade2ffbe";
                            }
                        });
                    };

                    this.$nextTick(syncEditFormSelection);
                }

                this.showModalEdit[index] = true;
            },
            async putDashboardUpdate(index) {
                let dashboardStartDateTime = parseInt(0);
                let dashboardStopDateTime = parseInt(0);
                let dashboardIntervalTime = parseInt(0);

                if (this.dashboardEditStartDate != undefined && this.dashboardEditStartTime == undefined) {
                    this.showModalEdit[index] = true;
                    toast(this.$t('dashboard.message12'), {
                        type: 'warning',
                        style: {
                            background: 'var(--warning)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--warning)'
                        }
                    });

                    return;
                }
                if (this.dashboardEditStartDate == undefined && this.dashboardEditStartTime != undefined) {
                    this.showModalEdit[index] = true;
                    toast(this.$t('dashboard.message13'), {
                        type: 'warning',
                        style: {
                            background: 'var(--warning)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--warning)'
                        }
                    });
                    return;
                }
                if (this.dashboardEditStopDate != undefined && this.dashboardEditStopTime == undefined) {
                    this.showModalEdit[index] = true;
                    toast(this.$t('dashboard.message14'), {
                        type: 'warning',
                        style: {
                            background: 'var(--warning)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--warning)'
                        }
                    });
                    return;
                }
                if (this.dashboardEditStopDate == undefined && this.dashboardEditStopTime != undefined) {
                    this.showModalEdit[index] = true;
                    toast(this.$t('dashboard.message15'), {
                        type: 'warning',
                        style: {
                            background: 'var(--warning)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--warning)'
                        }
                    });
                    return;
                }

                if (this.dashboardEditStartDate != undefined && this.dashboardEditStartTime != undefined) {
                    let dashboardEditStartDate = this.toNativeDate(this.dashboardEditStartDate);
                    let timeParts = this.getTimeParts(this.dashboardEditStartTime);
                    let dashboardStartJSDateTime = new Date(dashboardEditStartDate.getFullYear(), dashboardEditStartDate.getMonth(), dashboardEditStartDate.getDate(), timeParts.hours, timeParts.minutes, timeParts.seconds);
                    dashboardStartDateTime = parseInt((dashboardStartJSDateTime.getTime()));
                }
                if (this.dashboardEditStopDate != undefined && this.dashboardEditStopTime != undefined) {
                    let dashboardEditStopDate = this.toNativeDate(this.dashboardEditStopDate);
                    let timeParts = this.getTimeParts(this.dashboardEditStopTime);
                    let dashboardStopJSDateTime = new Date(dashboardEditStopDate.getFullYear(), dashboardEditStopDate.getMonth(), dashboardEditStopDate.getDate(), timeParts.hours, timeParts.minutes, timeParts.seconds);
                    dashboardStopDateTime = parseInt((dashboardStopJSDateTime.getTime()));
                }

                if (dashboardStartDateTime > 0 && dashboardStopDateTime > 0) {
                    if (dashboardStopDateTime - dashboardStartDateTime < 0) {
                        this.showModalEdit[index] = true;
                        toast(this.$t('dashboard.message16'), {
                            type: 'warning',
                            style: {
                                background: 'var(--warning)',
                                color: 'var(--primary-foreground)',
                                border: '1px solid var(--warning)'
                            }
                        });
                        return;
                    }
                    if (dashboardStopDateTime - dashboardStartDateTime > 3 * 31 * 24 * 60 * 60 * 1000) {
                        this.showModalEdit[index] = true;
                        toast(this.$t('dashboard.message17'), {
                            type: 'warning',
                            style: {
                                background: 'var(--warning)',
                                color: 'var(--primary-foreground)',
                                border: '1px solid var(--warning)'
                            }
                        });
                        return;
                    }
                }
                if (this.dashboardEditInterval == undefined || this.dashboardEditInterval == 0 || this.dashboardEditInterval == "" || isNaN(parseInt(this.dashboardEditInterval))) {
                    dashboardIntervalTime = parseInt(0)
                } else {
                    if (this.dashboardEditInterval < 0 || /^\d+$/.test(this.dashboardEditInterval) == false) {
                        this.showModalEdit[index] = true;
                        toast(this.$t('dashboard.message47'), {
                            type: 'warning',
                            style: {
                                background: 'var(--warning)',
                                color: 'var(--primary-foreground)',
                                border: '1px solid var(--warning)'
                            }
                        });
                        return;
                    } else {
                        dashboardIntervalTime = parseInt(this.dashboardEditInterval);
                    }
                }

                let errorFlag = 0;
                if (this.dashboardsEdit[index].chartdata != undefined) {
                    this.dashboardsEdit[index].chartdata.forEach((item, i) => {
                        this.dashboardsEdit[index].chartdata[i].graphorder = parseInt(this.dashboardsEdit[index].chartdata[i].graphorder);
                        this.dashboardsEdit[index].chartdata[i].spaceid = this.spaceValue[index][i];
                        this.dashboardsEdit[index].chartdata[i].metric = this.metricValue[index][i];
                        this.dashboardsEdit[index].chartdata[i].graphtype = this.graphtypeValue[index][i];
                        this.dashboardsEdit[index].chartdata[i].groupfunc = this.groupfuncValue[index][i];
                    
                        if (this.dashboardsEdit[index].chartdata[i].spaceid == undefined) {
                            this.showModalEdit[index] = true;
                            toast(this.$t('dashboard.message18') + ' ' + (i + 1), {
                                type: 'warning',
                                style: {
                                    background: 'var(--warning)',
                                    color: 'var(--primary-foreground)',
                                    border: '1px solid var(--warning)'
                                }
                            });
                            errorFlag = 1;
                            return;
                        }
                        if (this.dashboardsEdit[index].chartdata[i].metric == undefined) {
                            this.showModalEdit[index] = true;
                            toast(this.$t('dashboard.message19') + ' ' + (i + 1), {
                                type: 'warning',
                                style: {
                                    background: 'var(--warning)',
                                    color: 'var(--primary-foreground)',
                                    border: '1px solid var(--warning)'
                                }
                            });
                            errorFlag = 1;
                            return;
                        }
                        if (this.dashboardsEdit[index].chartdata[i].graphtype == undefined) {
                            this.showModalEdit[index] = true;
                            toast(this.$t('dashboard.message20') + ' ' + (i + 1), {
                                type: 'warning',
                                style: {
                                    background: 'var(--warning)',
                                    color: 'var(--primary-foreground)',
                                    border: '1px solid var(--warning)'
                                }
                            });
                            errorFlag = 1;
                            return;
                        }
                        if (this.dashboardsEdit[index].chartdata[i].groupfunc == undefined) {
                            this.showModalEdit[index] = true;
                            toast(this.$t('dashboard.message21') + ' ' + (i + 1), {
                                type: 'warning',
                                style: {
                                    background: 'var(--warning)',
                                    color: 'var(--primary-foreground)',
                                    border: '1px solid var(--warning)'
                                }
                            });
                            errorFlag = 1;
                            return;
                        }
                    });
                }

                if (errorFlag == 1) {
                    return;
                }

                let dataPut = {
                    id: this.dashboardsEdit[index].id,
                    name: this.dashboardEditName,
                    starttime: parseInt(dashboardStartDateTime),
                    stoptime: parseInt(dashboardStopDateTime),
                    interval: parseInt(dashboardIntervalTime),
                    status: 1,
                    chartdata: this.dashboardsEdit[index].chartdata
                };

                try {
                    await axios.put("/api/v1/admin/dashboard/edit", dataPut);
                    toast(this.$t('dashboard.message22'), {
                        type: 'success',
                        style: {
                            background: 'var(--success)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--success)'
                        }
                    });
                    this.showModalEdit[index] = false;
                    this.innerLoading = true;
                    this.barLoading = false;
                    await this.getDashboard();
                } catch (error) {
                    toast(this.$t('dashboard.message23'), {
                        type: 'error',
                        style: {
                            background: 'var(--destructive)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--destructive)'
                        }
                    });
                    this.showModalEdit[index] = true;
                }
            },
            async putDashboardRemove(index) {
                try {
                    await axios.put("/api/v1/admin/dashboard/remove", this.dashboards[index]);
                    toast(this.$t('dashboard.message24'), {
                        type: 'success',
                        style: {
                            background: 'var(--success)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--success)'
                        }
                    });
                    this.showModalRemove[index] = false;
                    this.innerLoading = true;
                    this.barLoading = false;
                    await this.getDashboard();
                } catch (error) {
                    toast(this.$t('dashboard.message25'), {
                        type: 'error',
                        style: {
                            background: 'var(--destructive)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--destructive)'
                        }
                    });
                    this.showModalRemove[index] = true;
                }
            },
            async putDashboardCreate() {
                let dashboardStartDateTime = parseInt(0);
                let dashboardStopDateTime = parseInt(0);
                let dashboardIntervalTime = parseInt(0);

                if (this.dashboardName == "") {
                    this.showModalCreate = true;
                    toast(this.$t('dashboard.message26'), {
                        type: 'error',
                        style: {
                            background: 'var(--destructive)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--destructive)'
                        }
                    });
                    return;
                }
                if (this.dashboardStartDate != undefined && this.dashboardStartTime == undefined) {
                    this.showModalCreate = true;
                    toast(this.$t('dashboard.message12'), {
                        type: 'error',
                        style: {
                            background: 'var(--destructive)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--destructive)'
                        }
                    });
                    return;
                }
                if (this.dashboardStartDate == undefined && this.dashboardStartTime != undefined) {
                    this.showModalCreate = true;
                    toast(this.$t('dashboard.message13'), {
                        type: 'error',
                        style: {
                            background: 'var(--destructive)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--destructive)'
                        }
                    });
                    return;
                }
                if (this.dashboardStopDate != undefined && this.dashboardStopTime == undefined) {
                    this.showModalCreate = true;
                    toast(this.$t('dashboard.message14'), {
                        type: 'error',
                        style: {
                            background: 'var(--destructive)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--destructive)'
                        }
                    });
                    return;
                }
                if (this.dashboardStopDate == undefined && this.dashboardStopTime != undefined) {
                    this.showModalCreate = true;
                    toast(this.$t('dashboard.message15'), {
                        type: 'error',
                        style: {
                            background: 'var(--destructive)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--destructive)'
                        }
                    });
                    return;
                }
                
                if (this.dashboardStartDate != undefined && this.dashboardStartTime != undefined) {
                    let dashboardStartDate = this.toNativeDate(this.dashboardStartDate);
                    let timeParts = this.getTimeParts(this.dashboardStartTime);
                    let dashboardStartJSDateTime = new Date(dashboardStartDate.getFullYear(), dashboardStartDate.getMonth(), dashboardStartDate.getDate(), timeParts.hours, timeParts.minutes, timeParts.seconds);
                    dashboardStartDateTime = parseInt((dashboardStartJSDateTime.getTime()));
                }
                if (this.dashboardStopDate != undefined && this.dashboardStopTime != undefined) {
                    let dashboardStopDate = this.toNativeDate(this.dashboardStopDate);
                    let timeParts = this.getTimeParts(this.dashboardStopTime);
                    let dashboardStopJSDateTime = new Date(dashboardStopDate.getFullYear(), dashboardStopDate.getMonth(), dashboardStopDate.getDate(), timeParts.hours, timeParts.minutes, timeParts.seconds);
                    dashboardStopDateTime = parseInt((dashboardStopJSDateTime.getTime()));
                }

                if (dashboardStartDateTime > 0 && dashboardStopDateTime > 0) {
                    if (dashboardStopDateTime - dashboardStartDateTime < 0) {
                        this.showModalCreate = true;
                        toast(this.$t('dashboard.message16'), {
                            type: 'warning',
                            style: {
                                background: 'var(--warning)',
                                color: 'var(--primary-foreground)',
                                border: '1px solid var(--warning)'
                            }
                        });
                        return;
                    }
                    if (dashboardStopDateTime - dashboardStartDateTime > 3 * 31 * 24 * 60 * 60 * 1000) {
                        this.showModalCreate = true;
                        toast(this.$t('dashboard.message47'), {
                            type: 'warning',
                            style: {
                                background: 'var(--warning)',
                                color: 'var(--primary-foreground)',
                                border: '1px solid var(--warning)'
                            }
                        });
                        return;
                    }
                }

                if (this.dashboardInterval == undefined || this.dashboardInterval == 0 || this.dashboardInterval == "" || isNaN(parseInt(this.dashboardInterval))) {
                    dashboardIntervalTime = parseInt(0)
                } else {
                    if (this.dashboardInterval < 0 || /^\d+$/.test(this.dashboardInterval) == false) {
                        this.showModalCreate = true;
                        toast(this.$t('dashboard.message47'), {
                            type: 'warning',
                            style: {
                                background: 'var(--warning)',
                                color: 'var(--primary-foreground)',
                                border: '1px solid var(--warning)'
                            }
                        });
                        return;
                    } else {
                        dashboardIntervalTime = parseInt(this.dashboardInterval);
                    }
                }

                let dataPut = {
                    name: this.dashboardName,
                    starttime: parseInt(dashboardStartDateTime),
                    stoptime: parseInt(dashboardStopDateTime),
                    interval: parseInt(dashboardIntervalTime),
                    status: 1,
                };

                try {
                    await axios.put("/api/v1/admin/dashboard/create", dataPut);
                    toast(this.$t('dashboard.message27'), {
                        type: 'success',
                        style: {
                            background: 'var(--success)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--success)'
                        }
                    });
                    this.showModalCreate = false;
                    this.innerLoading = true;
                    this.barLoading = false;
                    await this.getDashboard();
                } catch (error) {
                    toast(this.$t('dashboard.message28'), {
                        type: 'error',
                        style: {
                            background: 'var(--destructive)',
                            color: 'var(--primary-foreground)',
                            border: '1px solid var(--destructive)'
                        }
                    });
                    this.showModalCreate = true;
                }
            },
            async getDashboard() {
                try {
                    const response = await axios.get("/api/v1/admin/dashboard/select?" + Math.random());
                    const data = response.data;
                    if (data == undefined || data == null || data.length == 0) {
                        this.dashboards = [];
                        this.innerLoading = false;
                        return;
                    }
                    data.sort((a, b) => {
                        if (a.name > b.name) return 1;
                        if (a.name < b.name) return -1;
                        return 0;
                    });
                    this.dashboards = data;
                    data.forEach((item, i) => {
                        this.showModalRemove[i] = false;
                        this.showModalEdit[i] = false;

                        let dateStartTitle = "-30m";
                        let dateStopTitle = "Now";
                        if (this.dashboards[i].starttime > 0) {
                            dateStartTitle = moment(this.dashboards[i].starttime).format("DD.MM.YYYY HH:mm");
                        }
                        if (this.dashboards[i].stoptime > 0) {
                            dateStopTitle = moment(this.dashboards[i].stoptime).format("DD.MM.YYYY HH:mm");
                        }

                        let dateTitle = "";
                        if (this.dashboards[i].stoptime == 0) {
                            dateTitle = " (Now)";
                        } else if (this.dashboards[i].starttime > 0) {
                            dateTitle = " (" + dateStartTitle + " - " + dateStopTitle + ")";
                        }
                        this.dashboardTitle[i] = this.dashboards[i].name + dateTitle;
                    });
                    await this.putChartData();
                } catch (error) {
                    console.error("Failed to load dashboards", error);
                    this.innerLoading = false;
                }
            },
            async putChartData() {
                if (this.dashboards == null || this.dashboards.length == 0) {
                    this.innerLoading = false;
                    return;
                }
                try {
                    const response = await axios.put("/api/v1/admin/dashboard/data?" + Math.random(), this.dashboards);
                    const data = response.data;
                    data.forEach((item01, i01) => {
                        if (item01.labels != null) {
                            item01.labels.forEach((item02, i02) => {
                                if (item01.stoptime - item01.starttime <= (180 * 1000)) {
                                    data[i01].labels[i02] = moment(item02).format("HH:mm:ss");
                                } else {
                                    data[i01].labels[i02] = moment(item02).format("DD.MM.YYYY HH:mm");
                                }
                            });
                        }
                        if (item01.datasets == null) {
                            data[i01].datasets = [];
                        }
                    });
                    this.chartDataArr = data;
                    this.innerLoading = false;
                    this.barLoading = true;
                } catch (error) {
                    console.error("Failed to load dashboard chart data", error);
                    this.innerLoading = false;
                }
            },
            chartDataUpdate() {
                this.putChartData();
            },
            async spaceSelect() {
                let dataPut = {
                    pagesize: 0,
                };

                try {
                    const response = await axios.put("/api/v1/admin/space/select?" + Math.random(), dataPut);
                    const data = response.data;
                    if (data.data == null) {
                        this.spaceOptions = [];
                    } else {
                        this.spaceOptions = data.data.map((item) => ({
                            text: item.description,
                            value: item.id,
                        }));
                    }
                } catch (error) {
                    console.error("Failed to load spaces", error);
                }
            },
            onSpaceSelectChange(index01, index02, value) {
                const nextValue = value ?? "";
                this.spaceValue[index01] = [...(this.spaceValue[index01] ?? [])];
                this.spaceValue[index01][index02] = nextValue;

                const schedule = () => this.spaceSelectUpdateValue(index01, index02);
                if (typeof requestAnimationFrame === 'function') {
                    requestAnimationFrame(schedule);
                } else {
                    setTimeout(schedule, 0);
                }
            },
            async spaceSelectUpdateValue(index01, index02) {
                const selectedSpaceValue = this.spaceValue?.[index01]?.[index02] ?? "";
                const metricOptionsList = this.metricOptions?.[index01] ?? [];
                const metricValueList = this.metricValue?.[index01] ?? [];

                if (!selectedSpaceValue) {
                    this.metricOptions[index01] = [...metricOptionsList];
                    this.metricOptions[index01][index02] = [];
                    this.metricValue[index01] = [...metricValueList];
                    this.metricValue[index01][index02] = "";
                    if (this.dashboardsEdit?.[index01]?.chartdata?.[index02]) {
                        this.dashboardsEdit[index01].chartdata[index02].metric = "";
                    }
                    return;
                }

                const dataPut = {
                    spaceid: selectedSpaceValue,
                };

                try {
                    this.metricOptions[index01] = [...metricOptionsList];
                    this.metricOptions[index01][index02] = [];
                    this.metricValue[index01] = [...metricValueList];
                    this.metricValue[index01][index02] = "";

                    const response = await axios.put("/api/v1/admin/rawdata/metric/select?" + Math.random(), dataPut);
                    const data = response.data || [];
                    const metricOptions = data.map((item) => ({
                        text: item.metric,
                        value: item.metric,
                    }));
                    const currentMetric = this.dashboardsEdit?.[index01]?.chartdata?.[index02]?.metric ?? "";
                    const matchedMetric = metricOptions.find((item) => item.value == currentMetric)?.value || metricOptions[0]?.value || "";

                    this.metricOptions[index01] = [...(this.metricOptions?.[index01] ?? [])];
                    this.metricOptions[index01][index02] = metricOptions;
                    this.metricValue[index01] = [...(this.metricValue?.[index01] ?? [])];
                    this.metricValue[index01][index02] = matchedMetric;

                    if (this.dashboardsEdit?.[index01]?.chartdata?.[index02]) {
                        this.dashboardsEdit[index01].chartdata[index02].metric = matchedMetric;
                    }
                } catch (error) {
                    console.error("Failed to load metrics", error);
                    this.metricOptions[index01] = [...metricOptionsList];
                    this.metricOptions[index01][index02] = [];
                    this.metricValue[index01] = [...metricValueList];
                    this.metricValue[index01][index02] = "";
                }
            },
            formatDate(value) {
                if (!value) return "";

                const date = value instanceof Date
                    ? value
                    : new Date(value.year, value.month - 1, value.day);

                return new Intl.DateTimeFormat(this.locale, {
                    day: "2-digit",
                    month: "2-digit",
                    year: "numeric",
                }).format(date);
            },
        },  
        created() {
            this.getDashboard();
            this.spaceSelect();
            this.chartDataIntervalID = window.setInterval(this.chartDataUpdate, 10000);
        },
        unmounted() {
            clearInterval(this.chartDataIntervalID);
        },
        computed: {
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
            locale(newValue, oldValue) {
                this.componentKey += 1;
            },
            mode(newValue, oldValue) {
                this.componentKey += 1;
            }
        }
    }
</script>