<template>
    <el-dialog :model-value="modelValue" title="Add Character" width="720px"
        @update:model-value="(val) => emit('update:modelValue', val)">
        <el-form label-width="130px">
            <el-form-item label="Character Name">
                <el-input v-model="form.character_name" placeholder="Character Name" />
            </el-form-item>
            <el-form-item label="Password">
                <el-input v-model="form.password" type="password" placeholder="Character password" show-password />
            </el-form-item>
            <el-form-item label="Owner User UUHash">
                <el-input v-model="form.user_uu_hash" placeholder="Owner user UUHash" />
            </el-form-item>
            <el-form-item label="MetaData">
                <MetaDataEditor v-model="form.meta_data" height="240px" />
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

// Form fields map to the request body of characterAdd.go
const form = reactive({
    character_name: '',
    password: '',
    user_uu_hash: '',
    meta_data: '',
})
</script>
