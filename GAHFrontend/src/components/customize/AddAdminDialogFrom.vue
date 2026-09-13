<template>
    <el-dialog :model-value="modelValue" title="Add Admin" width="480px"
        @update:model-value="(val) => emit('update:modelValue', val)">
        <el-form label-width="80px">
            <el-form-item label="Username">
                <el-input v-model="form.username" placeholder="Admin username" />
            </el-form-item>
            <el-form-item label="Password">
                <el-input v-model="form.password" type="password" placeholder="Admin password" show-password />
            </el-form-item>
            <el-form-item label="Permissions">
                <el-checkbox v-model="form.permission.can_add_admin">Can add admin</el-checkbox>
                <el-checkbox v-model="form.permission.can_delete_admin">Can delete admin</el-checkbox>
                <el-checkbox v-model="form.permission.can_edit_admin">Can edit admin</el-checkbox>
                <el-checkbox v-model="form.permission.can_get_admin">Can get admin list</el-checkbox>
                <el-checkbox v-model="form.permission.can_operate_user">Can operate users</el-checkbox>
                <el-checkbox v-model="form.permission.can_operate_character">Can operate characters</el-checkbox>
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

defineProps({
    modelValue: {
        type: Boolean,
        default: false,
    },
})

const emit = defineEmits(['update:modelValue', 'submit'])

// Form fields map to the request body of adminAdd.go
const form = reactive({
    username: '',
    password: '',
    permission: {
        can_add_admin: false,
        can_delete_admin: false,
        can_edit_admin: false,
        can_get_admin: false,
        can_operate_user: false,
        can_operate_character: false,
    },
})
</script>
