<template>
    <el-dialog :model-value="modelValue" :title="`Edit Admin: ${admin?.Username ?? ''}`" width="480px"
        @update:model-value="(val) => emit('update:modelValue', val)">
        <el-form label-width="80px">
            <el-form-item label="Password">
                <el-input v-model="form.admin_info.password" type="password" placeholder="Leave empty to keep current password" show-password />
            </el-form-item>
            <el-form-item label="Permissions">
                <el-checkbox v-model="form.admin_info.permission.can_add_admin">Can add admin</el-checkbox>
                <el-checkbox v-model="form.admin_info.permission.can_delete_admin">Can delete admin</el-checkbox>
                <el-checkbox v-model="form.admin_info.permission.can_edit_admin">Can edit admin</el-checkbox>
                <el-checkbox v-model="form.admin_info.permission.can_get_admin">Can get admin list</el-checkbox>
                <el-checkbox v-model="form.admin_info.permission.can_operate_user">Can operate users</el-checkbox>
                <el-checkbox v-model="form.admin_info.permission.can_operate_character">Can operate characters</el-checkbox>
                <el-checkbox v-model="form.admin_info.permission.can_operate_app_key">Can operate app key</el-checkbox>
            </el-form-item>
        </el-form>
        <template #footer>
            <el-button @click="emit('update:modelValue', false)">Cancel</el-button>
            <el-button type="primary" @click="emit('submit', form)">Confirm</el-button>
        </template>
    </el-dialog>
</template>

<script setup>
import { reactive, watch } from 'vue'

const props = defineProps({
    modelValue: {
        type: Boolean,
        default: false,
    },
    // Admin row being edited, used to populate the form
    admin: {
        type: Object,
        default: null,
    },
})

const emit = defineEmits(['update:modelValue', 'submit'])

// Form fields map to the request body of adminEdit.go
const form = reactive({
    uu_hash: '',
    admin_info: {
        password: '',
        permission: {
            can_add_admin: false,
            can_delete_admin: false,
            can_edit_admin: false,
            can_get_admin: false,
            can_operate_user: false,
            can_operate_character: false,
            can_operate_app_key: false,
        },
    },
})

// Populate the form when the dialog opens or the target row changes
watch(
    () => [props.modelValue, props.admin],
    () => {
        if (!props.modelValue || !props.admin) return
        form.uu_hash = props.admin.UUHash ?? ''
        form.admin_info.password = ''
        Object.assign(form.admin_info.permission, props.admin.Permission ?? {})
    },
    { immediate: true }
)
</script>
