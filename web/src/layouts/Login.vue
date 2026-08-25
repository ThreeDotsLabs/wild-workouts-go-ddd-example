<template>
    <div class="ww-login-layout">
        <slot></slot>
    </div>
</template>

<script>
    export default {
        name: 'Login',
        components: {},
        beforeCreate: function() {
            document.body.className = 'login';
        },
    }
</script>

<style scoped>
    .ww-login-layout {
        min-height: 100vh;
    }
</style>
