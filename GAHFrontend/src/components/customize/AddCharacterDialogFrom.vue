<template>
    <el-dialog :model-value="modelValue" title="添加角色" width="720px"
        @update:model-value="(val) => emit('update:modelValue', val)">
        <el-form label-width="130px">
            <el-form-item label="角色名">
                <el-input v-model="form.character_name" placeholder="角色名" />
            </el-form-item>
            <el-form-item label="密码">
                <el-input v-model="form.password" type="password" placeholder="角色密码" show-password />
            </el-form-item>
            <el-form-item label="所属用户UUHash">
                <el-input v-model="form.user_uu_hash" placeholder="所属用户的 UUHash" />
            </el-form-item>
            <el-form-item label="MetaData">
                <MetaDataEditor v-model="form.meta_data" height="240px" />
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

// 表单字段对应 characterAdd.go 的请求体
const form = reactive({
    character_name: '',
    password: '',
    user_uu_hash: '',
    meta_data: '',
})
</script>
