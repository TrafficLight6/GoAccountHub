<template>
    <el-dialog :model-value="modelValue" title="添加管理员" width="480px"
        @update:model-value="(val) => emit('update:modelValue', val)">
        <el-form label-width="80px">
            <el-form-item label="用户名">
                <el-input v-model="form.username" placeholder="管理员用户名" />
            </el-form-item>
            <el-form-item label="密码">
                <el-input v-model="form.password" type="password" placeholder="管理员密码" show-password />
            </el-form-item>
            <el-form-item label="权限">
                <el-checkbox v-model="form.permission.can_add_admin">是否可以添加管理员</el-checkbox>
                <el-checkbox v-model="form.permission.can_delete_admin">是否可以删除管理员</el-checkbox>
                <el-checkbox v-model="form.permission.can_edit_admin">是否可以修改管理员</el-checkbox>
                <el-checkbox v-model="form.permission.can_get_admin">是否可以获取管理员列表</el-checkbox>
                <el-checkbox v-model="form.permission.can_operate_user">是否可以操作用户</el-checkbox>
                <el-checkbox v-model="form.permission.can_operate_character">是否可以操作角色</el-checkbox>
            </el-form-item>
        </el-form>
        <template #footer>
            <el-button @click="emit('update:modelValue', false)">取消</el-button>
            <el-button type="primary" @click="emit('submit', form)">确定</el-button>
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

// 表单字段对应 adminAdd.go 的请求体
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
