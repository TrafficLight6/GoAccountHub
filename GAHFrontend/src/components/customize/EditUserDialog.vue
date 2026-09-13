<template>
    <el-dialog :model-value="modelValue" :title="`Edit User: ${user?.Username ?? ''}`" width="720px"
        @update:model-value="(val) => emit('update:modelValue', val)">
        <el-form label-width="80px">
            <el-form-item label="Username">
                <el-input v-model="form.username" placeholder="Leave empty to keep current username" />
            </el-form-item>
            <el-form-item label="Password">
                <el-input v-model="form.password" type="password" placeholder="Leave empty to keep current password" show-password />
            </el-form-item>
            <el-form-item label="MetaData">
                <MetaDataEditor v-model="form.meta_data" height="260px" />
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
import MetaDataEditor from './MetaDataEditor.vue'

const props = defineProps({
    modelValue: {
        type: Boolean,
        default: false,
    },
    // User row being edited, used to populate the form
    user: {
        type: Object,
        default: null,
    },
})

const emit = defineEmits(['update:modelValue', 'submit'])

// Form fields map to the request body of userEdit.go
const form = reactive({
    uu_hash: '',
    username: '',
    password: '',
    meta_data: '',
})

// Populate the form when the dialog opens or the target row changes
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
