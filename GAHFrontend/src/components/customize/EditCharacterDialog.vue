<template>
    <el-dialog :model-value="modelValue" :title="`编辑角色：${character?.CharacterName ?? ''}`" width="720px"
        @update:model-value="(val) => emit('update:modelValue', val)">
        <el-form label-width="130px">
            <el-form-item label="角色名">
                <el-input v-model="form.character_name" placeholder="留空则不修改角色名" />
            </el-form-item>
            <el-form-item label="密码">
                <el-input v-model="form.password" type="password" placeholder="留空则不修改密码" show-password />
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
import { reactive, watch } from 'vue'
import MetaDataEditor from './MetaDataEditor.vue'

const props = defineProps({
    modelValue: {
        type: Boolean,
        default: false,
    },
    // 被编辑的角色行数据，用于回填表单
    character: {
        type: Object,
        default: null,
    },
})

const emit = defineEmits(['update:modelValue', 'submit'])

// 表单字段对应 characterEdit.go 的请求体
const form = reactive({
    user_uu_hash: '',
    character_uu_hash: '',
    character_name: '',
    password: '',
    meta_data: '',
})

// 打开弹窗或切换目标行时回填表单
watch(
    () => [props.modelValue, props.character],
    () => {
        if (!props.modelValue || !props.character) return
        form.user_uu_hash = props.character.UserUUHash ?? ''
        form.character_uu_hash = props.character.UUHash ?? ''
        form.character_name = props.character.CharacterName ?? ''
        form.password = ''
        form.meta_data = props.character.MetaData ?? ''
    },
    { immediate: true }
)
</script>
