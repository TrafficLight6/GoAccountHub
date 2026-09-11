<template>
    <el-dialog :model-value="modelValue" :title="`编辑用户：${user?.Username ?? ''}`" width="720px"
        @update:model-value="(val) => emit('update:modelValue', val)">
        <el-form label-width="80px">
            <el-form-item label="用户名">
                <el-input v-model="form.username" placeholder="留空则不修改用户名" />
            </el-form-item>
            <el-form-item label="密码">
                <el-input v-model="form.password" type="password" placeholder="留空则不修改密码" show-password />
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
import { reactive, watch } from 'vue'
import MetaDataEditor from './MetaDataEditor.vue'

const props = defineProps({
    modelValue: {
        type: Boolean,
        default: false,
    },
    // 被编辑的用户行数据，用于回填表单
    user: {
        type: Object,
        default: null,
    },
})

const emit = defineEmits(['update:modelValue', 'submit'])

// 表单字段对应 userEdit.go 的请求体
const form = reactive({
    uu_hash: '',
    username: '',
    password: '',
    meta_data: '',
})

// 打开弹窗或切换目标行时回填表单
watch(
    () => [props.modelValue, props.user],
    () => {
        if (!props.modelValue || !props.user) return
        form.uu_hash = props.user.UUHash ?? ''
        form.username = props.user.Username ?? ''
        form.password = ''
        form.meta_data = props.user.MetaData ?? ''
    },
    { immediate: true }
)
</script>
