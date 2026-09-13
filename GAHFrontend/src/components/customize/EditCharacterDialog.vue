<template>
    <el-dialog :model-value="modelValue" :title="`Edit Character: ${character?.CharacterName ?? ''}`" width="720px"
        @update:model-value="(val) => emit('update:modelValue', val)">
        <el-form label-width="130px">
            <el-form-item label="Character Name">
                <el-input v-model="form.character_name" placeholder="Leave empty to keep current character name" />
            </el-form-item>
            <el-form-item label="Password">
                <el-input v-model="form.password" type="password" placeholder="Leave empty to keep current password" show-password />
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
import { reactive, watch } from 'vue'
import MetaDataEditor from './MetaDataEditor.vue'

const props = defineProps({
    modelValue: {
        type: Boolean,
        default: false,
    },
    // Character row being edited, used to populate the form
    character: {
        type: Object,
        default: null,
    },
})

const emit = defineEmits(['update:modelValue', 'submit'])

// Form fields map to the request body of characterEdit.go
const form = reactive({
    user_uu_hash: '',
    character_uu_hash: '',
    character_name: '',
    password: '',
    meta_data: '',
})

// Populate the form when the dialog opens or the target row changes
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
