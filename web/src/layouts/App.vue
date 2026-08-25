<template>
    <div class="ww-app">
        <WebMenu/>
        <div class="ww-container ww-app__content">
            <slot></slot>
            <WebFooter/>
        </div>
    </div>
</template>

<script>
    import WebFooter from "../components/WebFooter.vue";
    import WebMenu from "../components/WebMenu.vue";
    import {Auth} from "../repositories/auth";

    export default {
        name: 'Main',
        components: {
            WebFooter,
            WebMenu,
        },
        beforeCreate: function () {
            document.body.className = 'app';

            if (!Auth.isLoggedIn()) {
                this.$router.push({name: 'login'});
            }
        },
    }
</script>

<style scoped>
    .ww-app {
        min-height: 100vh;
        display: flex;
        flex-direction: column;
    }

    .ww-app__content {
        display: flex;
        flex-direction: column;
        flex-grow: 1;
        width: 100%;
    }
</style>
