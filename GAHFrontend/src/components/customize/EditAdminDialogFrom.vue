<template>
    <el-dialog :model-value="modelValue" :title="`编辑管理员：${admin?.Username ?? ''}`" width="480px"
        @update:model-value="(val) => emit('update:modelValue', val)">
        <el-form label-width="80px">
            <el-form-item label="密码">
                <el-input v-model="form.admin_info.password" type="password" placeholder="留空则不修改密码" show-password />
            </el-form-item>
            <el-form-item label="权限">
                <el-checkbox v-model="form.admin_info.permission.can_add_admin">是否可以添加管理员</el-checkbox>
                <el-checkbox v-model="form.admin_info.permission.can_delete_admin">是否可以删除管理员</el-checkbox>
                <el-checkbox v-model="form.admin_info.permission.can_edit_admin">是否可以修改管理员</el-checkbox>
                <el-checkbox v-model="form.admin_info.permission.can_get_admin">是否可以获取管理员列表</el-checkbox>
                <el-checkbox v-model="form.admin_info.permission.can_operate_user">是否可以操作用户</el-checkbox>
                <el-checkbox v-model="form.admin_info.permission.can_operate_character">是否可以操作角色</el-checkbox>
            </el-form-item>
        </el-form>
        <template #footer>
            <el-button @click="emit('update:modelValue', false)">取消</el-button>
            <el-button type="primary" @click="emit('submit', form)">确定</el-button>
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
    // 被编辑的管理员行数据，用于回填表单
    admin: {
        type: Object,
        default: null,
    },
})

const emit = defineEmits(['update:modelValue', 'submit'])

// 表单字段对应 adminEdit.go 的请求体
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
        },
    },
})

// 打开弹窗或切换目标行时回填表单
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
