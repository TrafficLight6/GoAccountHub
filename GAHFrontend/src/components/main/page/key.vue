<template>
    <el-card>
        <template #header>
            Actions
        </template>
        <p><el-input v-model="searchForm.key_user" placeholder="Key User Name"></el-input></p>
        <template #footer>
            <div style="text-align: right">
                <el-button type="primary" @click="handleSearch">Search</el-button>
                <el-button type="warning" @click="handleReset">Reset</el-button>
                <el-button type="success" @click="handleAdd">Add Key</el-button>
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
                <el-button type="danger" @click="handleBatchDelete">Delete Selected Keys</el-button>
            </p>
        </el-card>
        <br>
    </template>
    <el-card>
        <template #header>
            Key List
        </template>
        <div style="height: 600px">
            <el-auto-resizer>
                <template #default="{ height, width }">
                    <el-table-v2 :columns="columns" :data="keyList" :width="width" :height="height"
                        row-key="ID" :footer-height="noMore ? 32 : 0" fixed @end-reached="handleEndReached">
                        <template #footer>
                            <el-text type="success" v-if="noMore" size="large">All {{ keyList.length }} loaded</el-text>
                        </template>
                    </el-table-v2>
                </template>
            </el-auto-resizer>
        </div>
    </el-card>

    <AddKeyDialogFrom v-model="addVisible" @submit="handleAddSubmit" />
    <ShowKeyDialog v-model="showVisible" :app-key="newKey" />
</template>
<script setup>
import { post, del } from '../../../lib/request.js'
import rangeKey from '../../../lib/rangeKey.js'
import AddKeyDialogFrom from '../../customize/AddKeyDialogFrom.vue'
import ShowKeyDialog from '../../customize/ShowKeyDialog.vue'
import { h, ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElCheckbox, ElButton, ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()
const adminInfo = ref({})
const keyList = ref([])

const PAGE_SIZE = 100
const loading = ref(false)
// Whether the backend has no more data
const noMore = ref(false)

// Search condition form (filtered by the backend after submit)
const searchForm = reactive({
    key_user: '',
})

// Currently active search conditions, sent to the backend with every range request
const searchCondition = ref({})

// Add key dialog
const addVisible = ref(false)

// The key returned by /key/add, shown only once
const showVisible = ref(false)
const newKey = ref('')

// Permission check: root is always allowed; regular admins read the permission field; returns false while data is not loaded
const canOperateAppKey = computed(() => {
    if (adminInfo.value.is_root) return true
    return Boolean(adminInfo.value.permission?.can_operate_app_key)
})

// Selection state: tracked by row ID
const selectedIds = ref([])
const isSelected = (id) => selectedIds.value.includes(id)
// Select-all check is based on the current data
const allSelected = computed(
    () => keyList.value.length > 0 && keyList.value.every((row) => isSelected(row.ID))
)

const toggleRow = (id) => {
    selectedIds.value = isSelected(id)
        ? selectedIds.value.filter((item) => item !== id)
        : [...selectedIds.value, id]
}

const toggleAll = () => {
    const shownIds = keyList.value.map((item) => item.ID)
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
        key_user: searchForm.key_user.trim(),
    }
    range(1)
}

// Reset: clear the conditions and re-fetch everything from the first page
const handleReset = () => {
    searchForm.key_user = ''
    searchCondition.value = {}
    range(1)
}

// Open the add key dialog
const handleAdd = () => {
    addVisible.value = true
}

// Submit add: call /key/add, then show the new key once and refresh the list
const handleAddSubmit = async (form) => {
    try {
        const res = await post('/key/add', {
            key_user_name: form.key_user_name.trim(),
        }, { autoRedirect401: false })
        newKey.value = res.key
        addVisible.value = false
        showVisible.value = true
        refresh()
    } catch {
        // Error messages are already shown by the request wrapper
    }
}

// Delete and sync the list (the backend accepts only one key_user_name at a time, so call it row by row)
const deleteKeys = async (rows) => {
    const results = await Promise.allSettled(
        rows.map((item) => del('/key/delete', { key_user_name: item.KeyUser }, { autoRedirect401: false, showError: false }))
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
    if (successIds.length > 0) refresh()
}

// Single delete: delete the row after confirmation
const handleDelete = async (row) => {
    try {
        await ElMessageBox.confirm(`Delete key "${row.KeyUser}"?`, 'Delete Confirmation', {
            type: 'warning',
            confirmButtonText: 'Confirm',
            cancelButtonText: 'Cancel',
        })
    } catch {
        // Delete cancelled
        return
    }
    deleteKeys([row])
}

// Batch delete: delete all selected rows after confirmation
const handleBatchDelete = async () => {
    const rows = keyList.value.filter((item) => selectedIds.value.includes(item.ID))
    if (rows.length === 0) return
    try {
        await ElMessageBox.confirm(`Delete ${rows.length} selected keys?`, 'Delete Confirmation', {
            type: 'warning',
            confirmButtonText: 'Confirm',
            cancelButtonText: 'Cancel',
        })
    } catch {
        // Delete cancelled
        return
    }
    deleteKeys(rows)
}

// Key creation time
const createdAtCell = ({ rowData }) =>
    h('span', rowData.CreatedAt ? new Date(rowData.CreatedAt).toLocaleString() : '')

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
    { key: 'KeyUser', dataKey: 'KeyUser', title: 'Key User Name', width: 240 },
    { key: 'Key', dataKey: 'Key', title: 'Key', width: 300 },
    { key: 'CreatedAt', title: 'Created At', width: 200, cellRenderer: createdAtCell },
    {
        key: 'actions',
        title: 'Actions',
        width: 150,
        cellRenderer: ({ rowData }) => h('div', { style: 'display: flex; gap: 8px;' }, [
            h(ElButton, { type: 'danger', size: 'small', onClick: () => handleDelete(rowData) }, () => 'Delete Key'),
        ]),
    },
]

const range = async (start = 1) => {
    if (loading.value) return
    loading.value = true
    try {
        const list = await rangeKey(searchCondition.value, start, PAGE_SIZE)
        // When scrolling further after refresh(), the requested page may overlap with already loaded data, so dedupe by ID
        const map = new Map()
        for (const row of start === 1 ? list : [...keyList.value, ...list]) {
            map.set(row.ID, row)
        }
        keyList.value = [...map.values()].sort((a, b) => a.ID - b.ID)
        noMore.value = list.length < PAGE_SIZE
    } catch {
        // 401 is already handled by the request wrapper, which redirects to the login page
    } finally {
        loading.value = false
    }
}

// Refresh after add/delete: re-fetch based on the currently loaded row count.
// Fetching only the first page would drop already loaded rows
const refresh = async () => {
    if (loading.value) return
    loading.value = true
    try {
        const length = Math.max(PAGE_SIZE, keyList.value.length)
        const list = await rangeKey(searchCondition.value, 1, length)
        keyList.value = list.sort((a, b) => a.ID - b.ID)
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
    range(Math.floor(keyList.value.length / PAGE_SIZE) + 1)
}

onMounted(async () => {
    try {
        const res = await post('/admin/info')
        adminInfo.value = res.data
    } catch {
        // 401 is already handled by the request wrapper, which redirects to the login page; skip further checks on network errors
        return
    }
    if (!canOperateAppKey.value) {
        router.push('/main/home')
        return
    }
    range()
})
</script>
