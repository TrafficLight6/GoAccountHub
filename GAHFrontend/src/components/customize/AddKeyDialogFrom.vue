<template>
    <el-dialog :model-value="modelValue" title="Add Key" width="480px"
        @update:model-value="(val) => emit('update:modelValue', val)">
        <el-alert type="warning" :closable="false" show-icon
            title="The key user name cannot be modified once submitted. Please check it before submitting." />
        <el-form ref="keyFormRef" :model="form" :rules="rules" label-width="110px" style="margin-top: 12px">
            <el-form-item label="Key User Name" prop="key_user_name">
                <el-input v-model="form.key_user_name" placeholder="Printable ASCII characters only" />
            </el-form-item>
        </el-form>
        <template #footer>
            <el-button @click="emit('update:modelValue', false)">Cancel</el-button>
            <el-button type="primary" @click="handleSubmit">Confirm</el-button>
        </template>
    </el-dialog>
</template>

<script setup>
import { reactive, ref } from 'vue'

defineProps({
    modelValue: {
        type: Boolean,
        default: false,
    },
})

const emit = defineEmits(['update:modelValue', 'submit'])

// Printable ASCII only (0x20-0x7E), the same check as keyAdd.go
const printableAsciiPattern = /^[\x20-\x7E]+$/

const keyFormRef = ref(null)

// Form fields map to the request body of keyAdd.go
const form = reactive({
    key_user_name: '',
})

const rules = {
    key_user_name: [
        { required: true, message: 'Please enter key user name', trigger: 'blur' },
        { pattern: printableAsciiPattern, message: 'Key user name must be printable ASCII characters only', trigger: 'blur' },
    ],
}

const handleSubmit = async () => {
    if (!keyFormRef.value) return
    try {
        await keyFormRef.value.validate()
    } catch {
        return
    }
    emit('submit', form)
}
</script>
