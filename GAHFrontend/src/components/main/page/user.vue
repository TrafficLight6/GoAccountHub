<template>
    <el-card>
        <template #header>
            Actions
        </template>
        <p><el-input v-model="searchForm.username" placeholder="Username"></el-input></p>
        <p><el-input v-model="searchForm.uu_hash" placeholder="User UUHash"></el-input></p>
        <template #footer>
            <div style="text-align: right">
                <el-button type="primary" @click="handleSearch">Search</el-button>
                <el-button type="warning" @click="handleReset">Reset</el-button>
                <el-button type="success" @click="handleAdd">Add User</el-button>
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
                <el-button type="danger" @click="handleBatchDelete">Delete Selected Users</el-button>
            </p>
        </el-card>
        <br>
    </template>
    <el-card>
        <template #header>
            User List ({{ userCount }} total)
        </template>
        <div style="height: 600px">
            <el-auto-resizer>
                <template #default="{ height, width }">
                    <el-table-v2 :columns="columns" :data="userList" :width="width" :height="height"
                        row-key="ID" :footer-height="noMore ? 32 : 0" fixed @end-reached="handleEndReached">
                        <template #footer>
                            <el-text type="success" v-if="noMore" size="large">All {{ userList.length }} loaded</el-text>
                        </template>
                    </el-table-v2>
                </template>
            </el-auto-resizer>
        </div>
    </el-card>

    <AddUserDialogFrom v-model="dialogVisible" @submit="handleAddSubmit" />
    <EditUserDialog v-model="editVisible" :user="editRow" @submit="handleEditSubmit" />
</template>
<script setup>
import { post, put, del } from '.././../../lib/request.js'
import rangeUser from '.././../../lib/rangeUser.js'
import AddUserDialogFrom from '../../customize/AddUserDialogFrom.vue'
import EditUserDialog from '../../customize/EditUserDialog.vue'
import { h, ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElCheckbox, ElButton, ElIcon, ElMessage, ElMessageBox } from 'element-plus'
import { DocumentCopy } from '@element-plus/icons-vue'

const router = useRouter()
const adminInfo = ref({})
const userList = ref([])
// Total user count (/user/count)
const userCount = ref(0)

const PAGE_SIZE = 100
const loading = ref(false)
// Whether the backend has no more data
const noMore = ref(false)

// Search condition form (filtered by the backend after submit)
const searchForm = reactive({
    username: '',
    uu_hash: '',
})

// Currently active search conditions, sent to the backend with every range request
const searchCondition = ref({})

// Add user dialog
const dialogVisible = ref(false)

// Edit user dialog
const editVisible = ref(false)
const editRow = ref(null)

const canGetUser = computed(() => {
    if (adminInfo.value.is_root) return true
    return Boolean(adminInfo.value.permission?.can_operate_user)
})

// Selection state: tracked by row ID
const selectedIds = ref([])
const isSelected = (id) => selectedIds.value.includes(id)
// Select-all check is based on the current data
const allSelected = computed(
    () => userList.value.length > 0 && userList.value.every((row) => isSelected(row.ID))
)

const toggleRow = (id) => {
    selectedIds.value = isSelected(id)
        ? selectedIds.value.filter((item) => item !== id)
        : [...selectedIds.value, id]
}

const toggleAll = () => {
    const shownIds = userList.value.map((item) => item.ID)
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
        username: searchForm.username.trim(),
        uu_hash: searchForm.uu_hash.trim(),
    }
    range(1)
}

// Reset: clear the conditions and re-fetch everything from the first page
const handleReset = () => {
    searchForm.username = ''
    searchForm.uu_hash = ''
    searchCondition.value = {}
    range(1)
}

// Open the add user dialog
const handleAdd = () => {
    dialogVisible.value = true
}

// Submit add: call /user/add, then refresh the list on success
const handleAddSubmit = async (form) => {
    try {
        await post('/user/add', {
            username: form.username.trim(),
            password: form.password,
            meta_data: form.meta_data,
        }, { autoRedirect401: false })
        ElMessage.success('User added')
        dialogVisible.value = false
        refresh()
        fetchCount()
    } catch {
        // Error messages are already shown by the request wrapper
    }
}

// Delete and sync the list (the backend accepts only one uu_hash at a time, so call it row by row)
const deleteUsers = async (rows) => {
    const results = await Promise.allSettled(
        rows.map((item) => del('/user/delete', { uu_hash: item.UUHash }, { autoRedirect401: false, showError: false }))
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
        await ElMessageBox.confirm(`Delete user "${row.Username}"?`, 'Delete Confirmation', {
            type: 'warning',
            confirmButtonText: 'Confirm',
            cancelButtonText: 'Cancel',
        })
    } catch {
        // Delete cancelled
        return
    }
    deleteUsers([row])
}

// Batch delete: delete all selected rows after confirmation
const handleBatchDelete = async () => {
    const rows = userList.value.filter((item) => selectedIds.value.includes(item.ID))
    if (rows.length === 0) return
    try {
        await ElMessageBox.confirm(`Delete ${rows.length} selected users?`, 'Delete Confirmation', {
            type: 'warning',
            confirmButtonText: 'Confirm',
            cancelButtonText: 'Cancel',
        })
    } catch {
        // Delete cancelled
        return
    }
    deleteUsers(rows)
}

// Open the edit user dialog
const handleEdit = (row) => {
    editRow.value = row
    editVisible.value = true
}

// Submit edit: call /user/edit, then refresh the list on success
const handleEditSubmit = async (form) => {
    try {
        await put('/user/edit', {
            uu_hash: form.uu_hash,
            username: form.username.trim(),
            password: form.password,
            meta_data: form.meta_data,
        }, { autoRedirect401: false })
        ElMessage.success('User updated')
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
    { key: 'Username', dataKey: 'Username', title: 'Username', width: 140 },
    {
        key: 'UUHash',
        title: 'UUHash',
        width: 600,
        cellRenderer: ({ rowData }) => h('div', { style: 'display: flex; align-items: center; gap: 8px;' }, [
            h('span', rowData.UUHash),
            h(ElButton, {
                size: 'small',
                title: 'Copy UUHash',
                onClick: () => handleCopy(rowData.UUHash),
            }, () => h(ElIcon, null, () => h(DocumentCopy))),
        ]),
    },
    { key: 'MetaData', title: 'MetaData', width: 300, cellRenderer: metaDataCell },
    {
        key: 'actions',
        title: 'Actions',
        width: 200,
        cellRenderer: ({ rowData }) => h('div', { style: 'display: flex; gap: 8px;' }, [
            h(ElButton, { type: 'danger', size: 'small', onClick: () => handleDelete(rowData) }, () => 'Delete User'),
            h(ElButton, { type: 'primary', size: 'small', onClick: () => handleEdit(rowData) }, () => 'Edit User'),
        ]),
    },
]

const range = async (start = 1) => {
    if (loading.value) return
    loading.value = true
    try {
        const list = await rangeUser(searchCondition.value, start, PAGE_SIZE)
        // When scrolling further after refresh(), the requested page may overlap with already loaded data, so dedupe by ID
        const map = new Map()
        for (const row of start === 1 ? list : [...userList.value, ...list]) {
            map.set(row.ID, row)
        }
        userList.value = [...map.values()].sort((a, b) => a.ID - b.ID)
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
        const length = Math.max(PAGE_SIZE, userList.value.length)
        const list = await rangeUser(searchCondition.value, 1, length)
        userList.value = list.sort((a, b) => a.ID - b.ID)
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
    range(Math.floor(userList.value.length / PAGE_SIZE) + 1)
}

// Total user count
const fetchCount = async () => {
    try {
        const res = await post('/user/count')
        userCount.value = res.data.user_count
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
    if (!canGetUser.value) {
        router.push('/')
        return
    }
    range()
    fetchCount()
})
</script>
