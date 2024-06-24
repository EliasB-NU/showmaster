<template>
    <div class="backdrop" @click.self="closePopUp">
        <div class="popup">
            <h1>{{ title }}</h1>
            <p>{{ text }}</p>
        </div>
    </div>
</template>

<script>
export default {
    props: ['title', 'text'],
    methods: {
        closePopUp(event) {
            event.preventDefault()
            this.$emit('close')
        }
    }
}
</script>

<style scoped>
    .popup {
        width: 400px;
        padding: 20px;
        margin: 100px auto;
        background: white;
        border-radius: 10px;
    }
    .popup h1 {
        color: black;
        border-bottom: 1px solid #ddd;
        display: inline-block;
        padding-bottom: 10px;
    }
    .backdrop {
        top: 0;
        position: fixed;
        background: rgba(0, 0, 0, 0.5);
        width: 100%;
        height: 100%;
    }
</style>