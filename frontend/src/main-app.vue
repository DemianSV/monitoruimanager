<template>
    <Dialog v-model:open="showPasswordDialog">
    <Toaster position="top-center" :offset="24" />
	<div class="flex! min-h-dvh min-w-0 max-w-[100vw] flex-col! overflow-x-hidden bg-background">
        <header class="fixed inset-x-0 top-[env(safe-area-inset-top,0px)] z-50 flex! h-14! w-full! max-w-[100vw] shrink-0! flex-row! items-center! justify-between! overflow-hidden! border-b bg-card p-2! shadow-sm lg:h-16!">
            <RouterLink to="/" aria-label="На главную" class="flex min-w-0 items-center">
                <div class="flex! min-w-0 items-center! gap-2 lg:gap-4">
                    <img :src="logo" alt="Polyxena UI Manager" class="block! size-10! object-contain lg:size-12!" />
                    <div class="hidden truncate whitespace-nowrap text-lg! font-bold text-foreground landscape:block landscape:text-xl! lg:block lg:text-2xl!">Πολυξένη&nbsp;UI&nbsp;Manager</div>
                </div>
            </RouterLink>

            <div class="flex! shrink-0! items-center! gap-3 lg:gap-4">
                <Button
                    type="button"
                    class="m-0! inline-flex! size-6! items-center! justify-center! border-0! bg-transparent! p-0! text-foreground transition-colors hover:text-primary"
                    :aria-label="currentLocale === 'ru' ? 'Переключить на английский язык' : 'Переключить на русский язык'"
                    :title="currentLocale === 'ru' ? 'English' : 'Русский'"
                    @click="toggleLocale">
                        <Languages class="size-5" />
                </Button>

                <Button
                    type="button"
                    class="m-0! inline-flex! size-6! items-center! justify-center! border-0! bg-transparent! p-0! text-foreground transition-colors hover:text-primary"
                    :aria-label="isDarkTheme ? 'Включить светлую тему' : 'Включить тёмную тему'"
                    :title="isDarkTheme ? 'Светлая тема' : 'Тёмная тема'"
                    @click="toggleTheme">
                        <Sun v-if="isDarkTheme" class="size-5" />
                        <Moon v-else class="size-5" />
                </Button>

                <Button
                    type="button"
                    class="m-0! inline-flex! size-6! items-center! justify-center! border-0! bg-transparent! p-0! text-foreground transition-colors hover:text-destructive"
                    aria-label="Выйти"
                    title="Выйти"
                    @click="exit">
                    <LogOut class="size-5" />
                </Button>

                <DialogTrigger as-child>
                    <Avatar :aria-label="user || 'Профиль пользователя'" :title="user" class="size-10! shrink-0! cursor-pointer lg:size-12!">
                        <AvatarFallback class="bg-primary text-sm text-primary-foreground">{{ userInitials }}</AvatarFallback>
                    </Avatar>
                </DialogTrigger>
            </div>
        </header>

        <main class="min-h-0 w-full max-w-full flex pl-2 pr-2 pb-2 pt-[calc(env(safe-area-inset-top,0px)+4.25rem)] lg:pl-6 lg:pr-6 lg:pb-6 lg:pt-[calc(env(safe-area-inset-top,0px)+5.5rem)]">
            <RouterView />
        </main>
    </div>

        <DialogContent class="fixed! inset-0! top-[env(safe-area-inset-top,0px)]! left-0! z-50 flex! h-[calc(100dvh-env(safe-area-inset-top,0px))]! w-full! max-w-none! max-h-none! translate-x-0! translate-y-0! flex-col! gap-0! overflow-hidden! rounded-none! border-0! p-0! bg-card text-foreground! lg:inset-auto! lg:top-[50%]! lg:left-[50%]! lg:grid! lg:h-auto! lg:max-w-lg! lg:max-h-[80vh]! lg:translate-x-[-50%]! lg:translate-y-[-50%]! lg:gap-4! lg:overflow-visible! lg:rounded-lg! lg:border! lg:p-6!">
            <form @submit.prevent="putPassword" class="grid min-h-0 flex-1 grid-rows-[auto_minmax(0,1fr)_auto] gap-0 overflow-hidden lg:flex-none lg:grid-rows-none lg:gap-6 lg:overflow-visible">
                <DialogHeader class="px-6 pt-6 pr-12 text-left lg:px-0 lg:pt-0 lg:pr-0">
                    <DialogTitle>{{ $t('message.settings') }}</DialogTitle>
                    <DialogDescription>
                        Для изменения пароля введите текущий пароль, новый пароль и подтвердите новый пароль.
                    </DialogDescription>
                </DialogHeader>
                <div class="grid min-h-0 content-start gap-3 overflow-y-auto p-6 lg:p-0">
                    <div class="grid gap-3">
                        <Label for="curPassword">{{ $t('password.current') }}</Label>
                        <Input id="curPassword" v-model="curPassword" name="curPassword" type="password" />
                    </div>
                    <div class="grid gap-3">
                        <Label for="newPassword">{{ $t('password.new') }}</Label>
                        <Input id="newPassword" v-model="newPassword" name="newPassword" type="password" />
                    </div>
                    <div class="grid gap-3">
                        <Label for="confirmPassword">{{ $t('password.again') }}</Label>
                        <Input id="confirmPassword" v-model="confirmPassword" name="confirmPassword" type="password" />
                    </div>
                </div>
                <DialogFooter class="mt-auto flex-row px-6 pb-[max(1.5rem,env(safe-area-inset-bottom))] *:flex-1 lg:mt-0 lg:px-0 lg:pb-0 lg:*:flex-none">
                    <DialogClose as-child>
                        <Button>{{ $t('message.cancel') }}</Button>
                    </DialogClose>
                    <Button type="submit" :disabled="isPasswordSubmitting">{{ $t('message.save') }}</Button>
                </DialogFooter>
            </form>
        </DialogContent>
    </Dialog>
</template>

<script>
import logo from '@/assets/logoPLXN02.svg';
import { Avatar, AvatarFallback } from '@/components/ui/avatar';
import { Languages, LogOut, Moon, Sun } from '@lucide/vue';
import {
    Dialog,
    DialogClose,
    DialogContent,
    DialogDescription,
    DialogFooter,
    DialogHeader,
    DialogTitle,
    DialogTrigger,
} from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Toaster } from '@/components/ui/sonner'
import { toast } from 'vue-sonner'
import { logout } from '@/auth.js';
import axios from 'axios';

import {
    themeQuartz,
    colorSchemeLight,
    colorSchemeDark,
    CellStyleModule,
    ClientSideRowModelModule,
    ClientSideRowModelApiModule,
    ColumnAutoSizeModule,
    QuickFilterModule,
    LocaleModule,
    ModuleRegistry,
} from 'ag-grid-community';
import { AG_GRID_LOCALE_RU } from './ag-grid-locale-ru.js';
import { provideGlobalGridOptions } from 'ag-grid-community';

ModuleRegistry.registerModules([
    CellStyleModule,
    ClientSideRowModelModule,
    ClientSideRowModelApiModule,
    ColumnAutoSizeModule,
    QuickFilterModule,
    LocaleModule,
]);

const colorAGLight = themeQuartz.withPart(colorSchemeLight)
    .withParams({
        backgroundColor: "rgba(255, 255, 255, 0.1)",
        headerBackgroundColor: "rgb(0, 0, 0, 0.03)",
    });

const colorAGDark = themeQuartz.withPart(colorSchemeDark)
    .withParams({
        backgroundColor: "rgba(0, 0, 0, 0)",
        headerBackgroundColor: "rgba(0, 0, 0, 0.2)",
    });

export default {
	name: 'main-app',
	components: {
		Avatar,
		AvatarFallback,
        Button,
        Dialog,
        DialogClose,
        DialogContent,
        DialogDescription,
        DialogFooter,
        DialogHeader,
        DialogTitle,
        DialogTrigger,
        Languages,
		Input,
		Label,
		LogOut,
		Moon,
		Sun,
		Toaster,
	},
	data() {
		return {
			logo,
			showPreloader: true,
            userBlock: true,
            user: '',
            showPasswordDialog: false,
            curPassword: '',
            newPassword: '',
            confirmPassword: '',
            isPasswordSubmitting: false,
		};
	},
    created() {
        this.applyTheme();
    },
	computed: {
		isDarkTheme() {
			return this.$store.state.mode;
		},
        currentLocale() {
            return this.$i18n.locale;
        },
        userInitials() {
            const nameParts = this.user.trim().split(/\s+/).filter(Boolean);
            return nameParts.length > 1
                ? `${nameParts[0][0]}${nameParts[1][0]}`.toUpperCase()
                : (nameParts[0]?.[0] || 'П').toUpperCase();
        },
	},
    watch: {
        showPasswordDialog(isOpen) {
            if (!isOpen) {
                this.resetPasswordFields();
            }
        },
    },
	mounted() {
		this.applyTheme();
	},
	methods: {
		async putPassword() {
            if (!this.curPassword || !this.newPassword || !this.confirmPassword) {
                toast(this.$t('password.message07'), { type: 'warning',
                    style: {
                        background: 'var(--warning)',
                        color: 'var(--primary-foreground)',
                        border: '1px solid var(--warning)'
                    }
                });
                return;
            }
            if (this.newPassword === this.curPassword) {
                toast(this.$t('password.message06'), { type: 'error',
                    style: {
                        background: 'var(--destructive)',
                        color: 'var(--primary-foreground)',
                        border: '1px solid var(--destructive)'
                    }
                });
                return;
            }
            if (this.newPassword !== this.confirmPassword) {
                toast(this.$t('password.message05'), { type: 'error',
                    style: {
                        background: 'var(--destructive)',
                        color: 'var(--primary-foreground)',
                        border: '1px solid var(--destructive)'
                    }
                });
                return;
            }
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
            if (!window.confirm(this.$t('password.confirm'))) {
                return;
            }

            this.isPasswordSubmitting = true;
            try {
                await axios.put('/api/v1/password', {
                    curpassword: this.curPassword,
                    newpassword: this.newPassword,
                });
                toast(this.$t('password.message01'), { type: 'success',
                    style: {
                        background: 'var(--success)',
                        color: 'var(--primary-foreground)',
                        border: '1px solid var(--success)'
                    }
                });
                this.showPasswordDialog = false;
            } catch (error) {
                const message = error.response?.status === 404
                    ? this.$t('password.message02')
                    : this.$t('password.message03');
                toast(message, { type: 'error' });
            } finally {
                this.isPasswordSubmitting = false;
            }
        },
        resetPasswordFields() {
            this.curPassword = '';
            this.newPassword = '';
            this.confirmPassword = '';
            this.isPasswordSubmitting = false;
        },
        async exit() {
            try {
                await logout();
                window.location.reload();
            } catch (error) {
                console.error('Failed to end the session', error);
                toast(this.$t('login.logoutFailed'), { type: 'error',
                    style: {
                        background: 'var(--destructive)',
                        color: 'var(--primary-foreground)',
                        border: '1px solid var(--destructive)'
                    }
                });
            }
        },
		async versionInfo() {
            try {
                const response = await axios.get("/api/v1/version?" + Math.random());
                this.version = response.data.Name;
            } catch (error) {
                console.error('Failed to load version information', error);
            }
        },
        async userInfo() {
            try {
                const response = await axios.get("/api/v1/userinfo?" + Math.random());
                const data = response.data;
                this.user = data.lname + " " + data.fname;
                this.userID = data.userid;
                this.userBlock = data.role !== "superadmin" && data.role !== "admin";
            } catch (error) {
                console.error('Failed to load user information', error);
            }
        },
		toggleTheme() {
			this.$store.commit('changemode', !this.$store.state.mode);
			window.localStorage.setItem('theme', this.$store.state.mode ? 'dark' : 'light');
			this.applyTheme();
		},
        toggleLocale() {
            const locale = this.currentLocale === 'ru' ? 'en' : 'ru';
            this.$i18n.locale = locale;
            this.$store.commit('changelocale', locale);
            window.localStorage.setItem('locale', locale);
            this.updateGridOptions();
        },
		applyTheme() {
			document.documentElement.classList.toggle('dark', this.isDarkTheme);
            this.updateGridOptions();
        },
        updateGridOptions() {
            provideGlobalGridOptions({
                ...(this.currentLocale === "ru" ? { localeText: AG_GRID_LOCALE_RU } : {}),
                theme: this.isDarkTheme ? colorAGDark : colorAGLight,
            }, 'shallow');
		},
        localeSelect() {
            let vm = this;
            this.localeOptions = [];
            if (this.$i18n.availableLocales.length > 0) {
                this.$i18n.availableLocales.forEach(function(item) {
                    vm.localeOptions.push(item);
                });

            }
            return vm.localeOptions;
        },
        localeUpdate() {
            this.$store.commit('changelocale', this.$i18n.locale);
            console.log("Change locale", this.$i18n.locale);
            this.updateGridOptions();
        },
	},
};
</script>