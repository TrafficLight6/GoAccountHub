<template>
    <el-card>
        <template #header>
            Actions
        </template>
        <p><el-input v-model="searchForm.character_name" placeholder="Character Name"></el-input></p>
        <p><el-input v-model="searchForm.user_uu_hash" placeholder="Owner User UUHash"></el-input></p>
        <p><el-input v-model="searchForm.uu_hash" placeholder="Character UUHash"></el-input></p>
        <template #footer>
            <div style="text-align: right">
                <el-button type="primary" @click="handleSearch">Search</el-button>
                <el-button type="warning" @click="handleReset">Reset</el-button>
                <el-button type="success" @click="handleAdd">Add Character</el-button>
            </div>
        </template>
    </el-card>
    <br>
    <template v-if="selectedIds.length > 0">
        <el-card>
            <template #header>
                Batch Actions
            </template>
            <el-text type="primary" size="large" style="text-align: center">Selected {{ selectedIds.length }} items</el-text>
            <p>
                <el-button type="primary" @click="handleCheckboxCancel">Clear Selection</el-button>
                <el-button type="danger" @click="handleBatchDelete">Delete Selected Characters</el-button>
            </p>
        </el-card>
        <br>
    </template>
    <el-card>
        <template #header>
            <div style="display: flex; align-items: center; justify-content: space-between">
                <span>Character List ({{ characterCount }} total)</span>
                <el-switch v-model="hideSameName" active-text="Hide Same-name Characters" />
            </div>
        </template>
        <div style="height: 600px">
            <el-auto-resizer>
                <template #default="{ height, width }">
                    <el-table-v2 :columns="columns" :data="visibleList" :width="width" :height="height"
                        row-key="ID" :footer-height="noMore ? 32 : 0" fixed @end-reached="handleEndReached">
                        <template #footer>
                            <el-text type="success" v-if="noMore" size="large">All {{ visibleList.length }} loaded</el-text>
                        </template>
                    </el-table-v2>
                </template>
            </el-auto-resizer>
        </div>
    </el-card>

    <AddCharacterDialogFrom v-model="dialogVisible" @submit="handleAddSubmit" />
    <EditCharacterDialog v-model="editVisible" :character="editRow" @submit="handleEditSubmit" />
</template>
<script setup>
import { post, put, del } from '.././../../lib/request.js'
import rangeCharacter from '.././../../lib/rangeCharacter.js'
import AddCharacterDialogFrom from '../../customize/AddCharacterDialogFrom.vue'
import EditCharacterDialog from '../../customize/EditCharacterDialog.vue'
import { h, ref, reactive, computed, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElCheckbox, ElButton, ElIcon, ElText, ElMessage, ElMessageBox } from 'element-plus'
import { DocumentCopy } from '@element-plus/icons-vue'

const router = useRouter()
const adminInfo = ref({})
const characterList = ref([])
// Total character count (/character/count)
const characterCount = ref(0)

const PAGE_SIZE = 100
const loading = ref(false)
// Whether the backend has no more data
const noMore = ref(false)

// Search condition form (filtered by the backend after submit)
const searchForm = reactive({
    character_name: '',
    user_uu_hash: '',
    uu_hash: '',
})

// Currently active search conditions, sent to the backend with every range request
const searchCondition = ref({})

// Add character dialog
const dialogVisible = ref(false)

// Edit character dialog
const editVisible = ref(false)
const editRow = ref(null)

const canGetCharacter = computed(() => {
    if (adminInfo.value.is_root) return true
    return Boolean(adminInfo.value.permission?.can_operate_character)
})

// Same-name character: the character uu_hash equals the owner user uu_hash; it cannot be deleted
// through the character API, so it is hidden by default
const hideSameName = ref(true)
const isSameNameCharacter = (row) => row.UUHash === row.UserUUHash
const visibleList = computed(() => (
    hideSameName.value
        ? characterList.value.filter((row) => !isSameNameCharacter(row))
        : characterList.value
))

// Selection state: tracked by row ID
const selectedIds = ref([])
const isSelected = (id) => selectedIds.value.includes(id)

// Clear the selection when the display scope changes so hidden rows are not left selected
watch(hideSameName, () => {
    selectedIds.value = []
})

// Select-all check is based on the currently displayed data
const allSelected = computed(
    () => visibleList.value.length > 0 && visibleList.value.every((row) => isSelected(row.ID))
)

const toggleRow = (id) => {
    selectedIds.value = isSelected(id)
        ? selectedIds.value.filter((item) => item !== id)
        : [...selectedIds.value, id]
}

const toggleAll = () => {
    const shownIds = visibleList.value.map((item) => item.ID)
    selectedIds.value = allSelected.value
        ? selectedIds.value.filter((id) => !shownIds.includes(id))
        : [...new Set([...selectedIds.value, ...shownIds])]
}

// Clear all selections
const handleCheckboxCancel = () => {
    selectedIds.value = []
}

// Search: hand the conditions to the backend and re-fetch from the first page
const handleSearch = () => {
    searchCondition.value = {
        character_name: searchForm.character_name.trim(),
        user_uu_hash: searchForm.user_uu_hash.trim(),
        uu_hash: searchForm.uu_hash.trim(),
    }
    range(1)
}

// Reset: clear the conditions and re-fetch everything from the first page
const handleReset = () => {
    searchForm.character_name = ''
    searchForm.user_uu_hash = ''
    searchForm.uu_hash = ''
    searchCondition.value = {}
    range(1)
}

// Open the add character dialog
const handleAdd = () => {
    dialogVisible.value = true
}

// Submit add: call /character/add, then refresh the list on success
const handleAddSubmit = async (form) => {
    try {
        // The character API returns 403 when the server has the multi-character switch turned off.
        // We show our own message here, so the shared error message is disabled
        await post('/character/add', {
            character_name: form.character_name.trim(),
            password: form.password,
            user_uu_hash: form.user_uu_hash.trim(),
            meta_data: form.meta_data,
        }, { autoRedirect401: false, showError: false })
        ElMessage.success('Character added')
        dialogVisible.value = false
        refresh()
        fetchCount()
    } catch (err) {
        if (err?.body?.code === 403) {
            ElMessage.warning('Multi-character is not enabled on the server; adding characters is not allowed')
        } else {
            ElMessage.error(err?.message || 'Add character failed')
        }
    }
}

// Delete and sync the list (the backend requires both the owner user and character uu_hash, so call it row by row)
const deleteCharacters = async (rows) => {
    const results = await Promise.allSettled(
        rows.map((item) => del('/character/delete', {
            user_uu_hash: item.UserUUHash,
            character_uu_hash: item.UUHash,
        }, { autoRedirect401: false, showError: false }))
    )
    const successIds = rows.filter((_, index) => results[index].status === 'fulfilled').map((item) => item.ID)
    const failed = results.filter((result) => result.status === 'rejected')
    // Uncheck the rows that were deleted successfully
    selectedIds.value = selectedIds.value.filter((id) => !successIds.includes(id))
    if (failed.length === 0) {
        ElMessage.success(`Deleted, ${successIds.length} items`)
    } else if (successIds.length === 0) {
        ElMessage.error(failed[0].reason?.message || 'Delete failed')
    } else {
        ElMessage.warning(`Succeeded ${successIds.length}, failed ${failed.length}: ${failed[0].reason?.message ?? 'Unknown reason'}`)
    }
    if (successIds.length > 0) {
        refresh()
        fetchCount()
    }
}

// Single delete: delete the row after confirmation
const handleDelete = async (row) => {
    try {
        await ElMessageBox.confirm(`Delete character "${row.CharacterName}"?`, 'Delete Confirmation', {
            type: 'warning',
            confirmButtonText: 'Confirm',
            cancelButtonText: 'Cancel',
        })
    } catch {
        // Delete cancelled
        return
    }
    deleteCharacters([row])
}

// Batch delete: delete all selected rows after confirmation
const handleBatchDelete = async () => {
    const rows = visibleList.value.filter((item) => selectedIds.value.includes(item.ID))
    if (rows.length === 0) return
    try {
        await ElMessageBox.confirm(`Delete ${rows.length} selected characters?`, 'Delete Confirmation', {
            type: 'warning',
            confirmButtonText: 'Confirm',
            cancelButtonText: 'Cancel',
        })
    } catch {
        // Delete cancelled
        return
    }
    deleteCharacters(rows)
}

// Open the edit character dialog
const handleEdit = (row) => {
    editRow.value = row
    editVisible.value = true
}

// Submit edit: call /character/edit, then refresh the list on success
const handleEditSubmit = async (form) => {
    try {
        await put('/character/edit', {
            user_uu_hash: form.user_uu_hash,
            character_uu_hash: form.character_uu_hash,
            character_name: form.character_name.trim(),
            password: form.password,
            meta_data: form.meta_data,
        }, { autoRedirect401: false })
        ElMessage.success('Character updated')
        editVisible.value = false
        refresh()
    } catch {
        // Error messages are already shown by the request wrapper
    }
}

// Copy text to the clipboard
const handleCopy = async (text) => {
    try {
        await navigator.clipboard.writeText(text)
        ElMessage.success('Copied')
    } catch {
        ElMessage.error('Copy failed')
    }
}

// Shared renderer for UUHash columns: el-text adds the ellipsis automatically when too long
// (it also provides a title when truncated), with a copy button on the right.
// The row itself must be width-constrained and allowed to shrink, otherwise the text is not
// truncated and the copy button gets clipped by the cell
const hashCell = (key) => ({ rowData }) => h('div', {
    style: 'display: flex; align-items: center; gap: 8px; width: 100%; min-width: 0;',
}, [
    h(ElText, { truncated: true, style: 'flex: 1; min-width: 0;' }, () => rowData[key]),
    h(ElButton, {
        size: 'small',
        title: 'Copy',
        onClick: () => handleCopy(rowData[key]),
    }, () => h(ElIcon, null, () => h(DocumentCopy))),
])

// MetaData can be long; truncate it in the cell and show the full text on hover
const metaDataCell = ({ rowData }) => h('span', {
    style: 'display: block; overflow: hidden; white-space: nowrap; text-overflow: ellipsis;',
    title: rowData.MetaData ?? '',
}, rowData.MetaData ?? '')

const columns = [
    {
        key: 'selection',
        width: 50,
        headerRenderer: () => h(ElCheckbox, {
            modelValue: allSelected.value,
            'onUpdate:modelValue': toggleAll,
            size: 'large',
        }),
        cellRenderer: ({ rowData }) => h(ElCheckbox, {
            modelValue: isSelected(rowData.ID),
            'onUpdate:modelValue': () => toggleRow(rowData.ID),
            size: 'large',
        }),
    },
    { key: 'ID', dataKey: 'ID', title: 'ID', width: 70 },
    { key: 'CharacterName', dataKey: 'CharacterName', title: 'Character Name', width: 160 },
    { key: 'UUHash', title: 'Character UUHash', width: 260, cellRenderer: hashCell('UUHash') },
    { key: 'UserUUHash', title: 'Owner User UUHash', width: 260, cellRenderer: hashCell('UserUUHash') },
    { key: 'MetaData', title: 'MetaData', width: 300, cellRenderer: metaDataCell },
    {
        key: 'actions',
        title: 'Actions',
        width: 250,
        cellRenderer: ({ rowData }) => h('div', { style: 'display: flex; gap: 8px;' }, [
            h(ElButton, { type: 'danger', size: 'small', onClick: () => handleDelete(rowData) }, () => 'Delete Character'),
            h(ElButton, { type: 'primary', size: 'small', onClick: () => handleEdit(rowData) }, () => 'Edit Character'),
        ]),
    },
]

const range = async (start = 1) => {
    if (loading.value) return
    loading.value = true
    try {
        const list = await rangeCharacter(searchCondition.value, start, PAGE_SIZE)
        // When scrolling further after refresh(), the requested page may overlap with already loaded data, so dedupe by ID
        const map = new Map()
        for (const row of start === 1 ? list : [...characterList.value, ...list]) {
            map.set(row.ID, row)
        }
        characterList.value = [...map.values()].sort((a, b) => a.ID - b.ID)
        noMore.value = list.length < PAGE_SIZE
    } catch {
        // 401 is already handled by the request wrapper, which redirects to the login page
    } finally {
        loading.value = false
    }
}

// Refresh after add/edit/delete: re-fetch based on the currently loaded row count.
// Fetching only the first page would drop already loaded rows, making the affected row appear to "disappear"
const refresh = async () => {
    if (loading.value) return
    loading.value = true
    try {
        const length = Math.max(PAGE_SIZE, characterList.value.length)
        const list = await rangeCharacter(searchCondition.value, 1, length)
        characterList.value = list.sort((a, b) => a.ID - b.ID)
        noMore.value = list.length < length
    } catch {
        // 401 is already handled by the request wrapper, which redirects to the login page
    } finally {
        loading.value = false
    }
}

// Reached the bottom: keep requesting the next 100 rows while the backend has more data
const handleEndReached = () => {
    if (loading.value || noMore.value) return
    range(Math.floor(characterList.value.length / PAGE_SIZE) + 1)
}

// Total character count
const fetchCount = async () => {
    try {
        const res = await post('/character/count')
        characterCount.value = res.data.character_count
    } catch {
        // 401 is already handled by the request wrapper, which redirects to the login page
    }
}

onMounted(async () => {
    try {
        const res = await post('/admin/info')
        adminInfo.value = res.data
    } catch {
        // 401 is already handled by the request wrapper, which redirects to the login page; skip further checks on network errors
        return
    }
    if (!canGetCharacter.value) {
        router.push('/')
        return
    }
    range()
    fetchCount()
})
</script>
