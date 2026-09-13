<template>
    <el-dialog :model-value="modelValue" title="Add User" width="720px"
        @update:model-value="(val) => emit('update:modelValue', val)">
        <el-form label-width="80px">
            <el-form-item label="Username">
                <el-input v-model="form.username" placeholder="Username" />
            </el-form-item>
            <el-form-item label="Password">
                <el-input v-model="form.password" type="password" placeholder="User password" show-password />
            </el-form-item>
            <el-form-item label="MetaData">
                <MetaDataEditor v-model="form.meta_data" height="260px" />
            </el-form-item>
        </el-form>
        <template #footer>
            <el-button @click="emit('update:modelValue', false)">Cancel</el-button>
            <el-button type="primary" @click="emit('submit', form)">Confirm</el-button>
        </template>
    </el-dialog>
</template>

<script setup>
import { reactive } from 'vue'
import MetaDataEditor from './MetaDataEditor.vue'

defineProps({
    modelValue: {
        type: Boolean,
        default: false,
    },
})

const emit = defineEmits(['update:modelValue', 'submit'])

// Form fields map to the request body of userAdd.go
const form = reactive({
    username: '',
    password: '',
    meta_data: '',
})
</script>
