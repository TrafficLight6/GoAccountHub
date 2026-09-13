<template>
    <el-dialog :model-value="modelValue" title="Key Created" width="600px" :close-on-click-modal="false"
        @update:model-value="(val) => emit('update:modelValue', val)">
        <el-alert type="warning" :closable="false" show-icon
            title="This key is shown only once. Please copy and store it now, it cannot be viewed again." />
        <el-form label-width="110px" style="margin-top: 12px">
            <el-form-item label="Key">
                <el-input :model-value="appKey" disabled>
                    <template #append>
                        <el-button @click="handleCopy">Copy</el-button>
                    </template>
                </el-input>
            </el-form-item>
        </el-form>
        <template #footer>
            <el-button type="primary" @click="emit('update:modelValue', false)">Confirm</el-button>
        </template>
    </el-dialog>
</template>

<script setup>
import { ElMessage } from 'element-plus'

const props = defineProps({
    modelValue: {
        type: Boolean,
        default: false,
    },
    // The key returned by /key/add, shown only once
    appKey: {
        type: String,
        default: '',
    },
})

const emit = defineEmits(['update:modelValue'])

// Copy text to the clipboard
const handleCopy = async () => {
    try {
        await navigator.clipboard.writeText(props.appKey)
        ElMessage.success('Copied')
    } catch {
        ElMessage.error('Copy failed')
    }
}
</script>
