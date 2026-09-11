<template>
    <el-dialog :model-value="modelValue" title="添加用户" width="720px"
        @update:model-value="(val) => emit('update:modelValue', val)">
        <el-form label-width="80px">
            <el-form-item label="用户名">
                <el-input v-model="form.username" placeholder="用户名" />
            </el-form-item>
            <el-form-item label="密码">
                <el-input v-model="form.password" type="password" placeholder="用户密码" show-password />
            </el-form-item>
            <el-form-item label="MetaData">
                <MetaDataEditor v-model="form.meta_data" height="260px" />
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
import MetaDataEditor from './MetaDataEditor.vue'

defineProps({
    modelValue: {
        type: Boolean,
        default: false,
    },
})

const emit = defineEmits(['update:modelValue', 'submit'])

// 表单字段对应 userAdd.go 的请求体
const form = reactive({
    username: '',
    password: '',
    meta_data: '',
})
</script>
