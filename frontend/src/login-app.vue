<template>
    <Toaster
        position="top-center"
        :offset="24"/>

    <div class="bg-background fixed inset-0 z-999 flex items-center justify-center" v-if="showPreloader">
        <div class="bg-primary flex-none! h-auto w-max! min-w-0 box-border p-4! rounded-lg font-sans text-[15px] text-primary-foreground text-center">
            {{ $t('login.message05') }}
        </div>
    </div>
    <div class="flex flex-col h-screen min-h-dvh w-screen items-center landscape:flex-row lg:flex-row bg-background lg:bg-[url('assets/bg01.svg')] lg:bg-center lg:bg-no-repeat lg:bg-cover">
        <div class="flex float w-full h-1/2 items-end justify-center p-6! landscape:float-left landscape:h-full landscape:w-1/2 landscape:items-center lg:justify-end lg:float-left lg:h-full lg:w-1/2 lg:items-center">
            <div class="flex w-full h-auto justify-center lg:h-full lg:justify-end!">
                <img :src="logo01" class="h-auto max-w-75 landscape:max-w-82.5 lg:max-w-112.5!">
            </div>
        </div>
        <div class="flex w-full h-1/2 items-start justify-center p-6! landscape:h-screen landscape:float-right landscape:items-center landscape:justify-start! landscape:w-1/2 lg:h-screen lg:float-right lg:items-center lg:justify-start! lg:w-1/2">
            <div class="justify-center lg:w-112.5 lg:max-w-112.5 lg:justify-start">
                <div class="text-center font-bold text-[28px]! p-1! lg:flex lg:justify-start lg:text-left lg:text-[32px]!">Πολυξένη&nbsp;UI&nbsp;Manager</div>
                <div class="text-center p-1! lg:flex lg:justify-start lg:text-left">Polyxena - Opens Source Monitoring System</div>

                <div class="flex p-1!"><Input id="username" name="username" v-model="username" v-on:keydown.enter="login()" :placeholder="$t('login.username01')" :label="$t('login.username02')" autocomplete="username" /></div>
                <div class="flex p-1!"><Input id="password" name="password" type="password" v-model="password" v-on:keydown.enter="login()" :placeholder="$t('login.password01')" :label="$t('login.password02')" autocomplete="current-password" /></div>
                <div class="flex p-1!"><Button class="w-full" v-on:click="login()" v-bind:disabled="isButtonDisabled"> {{ $t('login.button') }} </Button></div>
            </div>
        </div>
    </div>

</template>

<script>
import logo01 from '@/assets/logoPLXN01.svg';
import { login } from '@/auth.js';

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Toaster } from '@/components/ui/sonner'
import { toast } from 'vue-sonner'
import 'vue-sonner/style.css'

export default {
    name: 'login-app',
    components: {
        Button,
        Input,
        Toaster,
    },
    data() {
        return {
            placeholder01: this.$t('login.username02'),
            placeholder02: this.$t('login.password02'),
            error: this.$t('login.message01'),
            errorColor: "",
            username: "",
            password: "",
            isDisabled: false,
            showPreloader: true,
            isBlock: true,
            isError1: false,
            isError2: false,
            isButtonDisabled: false,
            logo01
        }
    },
    methods: {
        async login() {
            if (this.isButtonDisabled) {
                return;
            }
            this.isButtonDisabled = true;
            try {
                await login(this.username, this.password);

                toast(this.$t('login.message03'), {
                    type: 'success',
                    style: {
                        background: 'var(--success)',
                        color: 'var(--primary-foreground)',
                        border: '1px solid var(--success)'
                    }
                });
            
                setTimeout(() => window.location.reload(), 1000);

            } catch (error) {
                toast(this.$t('login.message04'), {
                    type: 'error',
                    style: {
                        background: 'var(--destructive)',
                        color: 'var(--primary-foreground)',
                        border: '1px solid var(--destructive)'
                    }
                });
                this.isButtonDisabled = false;
            }
        }
    }
};
</script>