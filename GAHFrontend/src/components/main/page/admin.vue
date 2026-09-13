<template>
    <el-card>
        <template #header>
            Actions
        </template>
        <p><el-input v-model="searchForm.username" placeholder="Admin Username"></el-input></p>
        <p><el-input v-model="searchForm.uu_hash" placeholder="Admin UUHash"></el-input></p>
        <p>Permission Filter:</p>
        <el-checkbox v-model="searchForm.permission.can_add_admin" label="Can add admin"></el-checkbox>
        <el-checkbox v-model="searchForm.permission.can_delete_admin" label="Can delete admin"></el-checkbox>
        <el-checkbox v-model="searchForm.permission.can_edit_admin" label="Can edit admin"></el-checkbox>
        <el-checkbox v-model="searchForm.permission.can_get_admin" label="Can get admin list"></el-checkbox>
        <el-checkbox v-model="searchForm.permission.can_operate_user" label="Can operate users"></el-checkbox>
        <el-checkbox v-model="searchForm.permission.can_operate_character" label="Can operate characters"></el-checkbox>
        <el-checkbox v-model="searchForm.permission.can_operate_app_key" label="Can operate app key"></el-checkbox>
        <template #footer>
            <div style="text-align: right">
                <el-button type="primary" @click="handleSearch">Search</el-button>
                <el-button type="warning" @click="handleReset">Reset</el-button>
                <el-button type="success" @click="handleAdd">Add Admin</el-button>
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
                <el-button type="danger" @click="handleBatchDelete">Delete Selected Admins</el-button>
            </p>
        </el-card>
        <br>
    </template>
    <el-card>
        <template #header>
            Admin List
        </template>
        <div style="height: 600px">
            <el-auto-resizer>
                <template #default="{ height, width }">
                    <el-table-v2 :columns="columns" :data="adminList" :width="width" :height="height"
                        row-key="ID" :footer-height="noMore ? 32 : 0" fixed @end-reached="handleEndReached">
                        <template #footer>
                            <el-text type="success" v-if="noMore" size="large">All {{ adminList.length }} loaded</el-text>
                        </template>
                    </el-table-v2>
                </template>
            </el-auto-resizer>
        </div>
    </el-card>

    <AddAdminDialogFrom v-model="dialogVisible" @submit="handleAddSubmit" />
    <EditAdminDialogFrom v-model="editVisible" :admin="editRow" @submit="handleEditSubmit" />
</template>
<script setup>
import { post, put, del } from '.././../../lib/request.js'
import rangeAdmin from '.././../../lib/rangeAdmin.js'
import AddAdminDialogFrom from '../../customize/AddAdminDialogFrom.vue'
import EditAdminDialogFrom from '../../customize/EditAdminDialogFrom.vue'
import { h, ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElCheckbox, ElButton, ElIcon, ElMessage, ElMessageBox } from 'element-plus'
import { DocumentCopy } from '@element-plus/icons-vue'

const router = useRouter()
const adminInfo = ref({})
const adminList = ref([])

const PAGE_SIZE = 100
const loading = ref(false)
// Whether the backend has no more data
const noMore = ref(false)

// Filter form (filtered by the backend after submit)
const searchForm = reactive({
    username: '',
    uu_hash: '',
    permission: {
        can_add_admin: false,
        can_delete_admin: false,
        can_edit_admin: false,
        can_get_admin: false,
        can_operate_user: false,
        can_operate_character: false,
        can_operate_app_key: false,
    },
})

// Currently active filter, sent to the backend with each range request
const searchCondition = ref({})

// Add admin dialog
const dialogVisible = ref(false)

// Edit admin dialog
const editVisible = ref(false)
const editRow = ref(null)

const canGetAdmin = computed(() => {
    if (adminInfo.value.is_root) return true
    return Boolean(adminInfo.value.permission?.can_get_admin)
})

// Unified permission cell rendering: "Yes" when granted, otherwise "No"
const permCell = (key) => ({ rowData }) => h('span', rowData.Permission?.[key] ? 'Yes' : 'No')

// Selection state: recorded by row ID
const selectedIds = ref([])
const isSelected = (id) => selectedIds.value.includes(id)
// Select-all is based on the current data
const allSelected = computed(
    () => adminList.value.length > 0 && adminList.value.every((row) => isSelected(row.ID))
)

const toggleRow = (id) => {
    selectedIds.value = isSelected(id)
        ? selectedIds.value.filter((item) => item !== id)
        : [...selectedIds.value, id]
}

const toggleAll = () => {
    const shownIds = adminList.value.map((item) => item.ID)
    selectedIds.value = allSelected.value
        ? selectedIds.value.filter((id) => !shownIds.includes(id))
        : [...new Set([...selectedIds.value, ...shownIds])]
}

// Clear all selections
const handleCheckboxCancel = () => {
    selectedIds.value = []
}

// Search: send the filter to the backend and refetch from the first page
const handleSearch = () => {
    searchCondition.value = {
        username: searchForm.username.trim(),
        uu_hash: searchForm.uu_hash.trim(),
        can_add_admin: searchForm.permission.can_add_admin,
        can_delete_admin: searchForm.permission.can_delete_admin,
        can_edit_admin: searchForm.permission.can_edit_admin,
        can_get_admin: searchForm.permission.can_get_admin,
        can_operate_user: searchForm.permission.can_operate_user,
        can_operate_character: searchForm.permission.can_operate_character,
        can_operate_app_key: searchForm.permission.can_operate_app_key,
    }
    range(1)
}

// Reset: clear the filter and refetch everything from the first page
const handleReset = () => {
    searchForm.username = ''
    searchForm.uu_hash = ''
    Object.keys(searchForm.permission).forEach((key) => {
        searchForm.permission[key] = false
    })
    searchCondition.value = {}
    range(1)
}

// Open the add admin dialog
const handleAdd = () => {
    dialogVisible.value = true
}

// Submit add: call /admin/add and refresh the list on success
const handleAddSubmit = async (form) => {
    try {
        await post('/admin/add', {
            username: form.username.trim(),
            password: form.password,
            permission: { ...form.permission },
        }, { autoRedirect401: false })
        ElMessage.success('Admin added successfully')
        dialogVisible.value = false
        range(1)
    } catch {
        // Failure messages are already shown by the request wrapper
    }
}

// Perform deletion and sync the list (the backend accepts one uu_hash per request, so call one by one)
const deleteAdmins = async (rows) => {
    const results = await Promise.allSettled(
        rows.map((item) => del('/admin/delete', { uu_hash: item.UUHash }, { autoRedirect401: false, showError: false }))
    )
    const successIds = rows.filter((_, index) => results[index].status === 'fulfilled').map((item) => item.ID)
    const failed = results.filter((result) => result.status === 'rejected')
    // Deselect successfully deleted rows
    selectedIds.value = selectedIds.value.filter((id) => !successIds.includes(id))
    if (failed.length === 0) {
        ElMessage.success(`Deleted successfully, ${successIds.length} in total`)
    } else if (successIds.length === 0) {
        ElMessage.error(failed[0].reason?.message || 'Delete failed')
    } else {
        ElMessage.warning(`Succeeded: ${successIds.length}, failed: ${failed.length}: ${failed[0].reason?.message ?? 'Unknown reason'}`)
    }
    if (successIds.length > 0) range(1)
}

// Single delete: confirm, then delete the row
const handleDelete = async (row) => {
    try {
        await ElMessageBox.confirm(`Are you sure you want to delete admin "${row.Username}"?`, 'Delete Confirmation', {
            type: 'warning',
            confirmButtonText: 'Confirm',
            cancelButtonText: 'Cancel',
        })
    } catch {
        // Deletion canceled
        return
    }
    deleteAdmins([row])
}

// Batch delete: confirm, then delete all selected rows
const handleBatchDelete = async () => {
    const rows = adminList.value.filter((item) => selectedIds.value.includes(item.ID))
    if (rows.length === 0) return
    try {
        await ElMessageBox.confirm(`Are you sure you want to delete the ${rows.length} selected admins?`, 'Delete Confirmation', {
            type: 'warning',
            confirmButtonText: 'Confirm',
            cancelButtonText: 'Cancel',
        })
    } catch {
        // Deletion canceled
        return
    }
    deleteAdmins(rows)
}

// Open the edit admin dialog
const handleEdit = (row) => {
    editRow.value = row
    editVisible.value = true
}

// Submit edit: call /admin/edit and refresh the list on success
const handleEditSubmit = async (form) => {
    try {
        await put('/admin/edit', {
            uu_hash: form.uu_hash,
            admin_info: {
                password: form.admin_info.password,
                permission: { ...form.admin_info.permission },
            },
        }, { autoRedirect401: false })
        ElMessage.success('Admin updated successfully')
        editVisible.value = false
        range(1)
    } catch {
        // Failure messages are already shown by the request wrapper
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
        width: 300,
        cellRenderer: ({ rowData }) => h('div', { style: 'display: flex; align-items: center; gap: 8px; min-width: 0;' }, [
            h('span', { style: 'flex: 1; min-width: 0; overflow: hidden; text-overflow: ellipsis; white-space: nowrap;' }, rowData.UUHash),
            h(ElButton, {
                size: 'small',
                title: 'Copy UUHash',
                onClick: () => handleCopy(rowData.UUHash),
            }, () => h(ElIcon, null, () => h(DocumentCopy))),
        ]),
    },
    { key: 'can_add_admin', title: 'Add Admin', width: 110, cellRenderer: permCell('can_add_admin') },
    { key: 'can_delete_admin', title: 'Delete Admin', width: 110, cellRenderer: permCell('can_delete_admin') },
    { key: 'can_edit_admin', title: 'Edit Admin', width: 110, cellRenderer: permCell('can_edit_admin') },
    { key: 'can_get_admin', title: 'View Admin', width: 110, cellRenderer: permCell('can_get_admin') },
    { key: 'can_operate_user', title: 'Operate Users', width: 100, cellRenderer: permCell('can_operate_user') },
    { key: 'can_operate_character', title: 'Operate Characters', width: 100, cellRenderer: permCell('can_operate_character') },
    { key: 'can_operate_app_key', title: 'Operate App Key', width: 110, cellRenderer: permCell('can_operate_app_key') },
    {
        key: 'actions',
        title: 'Actions',
        width: 200,
        cellRenderer: ({ rowData }) => h('div', { style: 'display: flex; gap: 8px;' }, [
            h(ElButton, { type: 'danger', size: 'small', onClick: () => handleDelete(rowData) }, () => 'Delete Admin'),
            h(ElButton, { type: 'primary', size: 'small', onClick: () => handleEdit(rowData) }, () => 'Edit Admin'),
        ]),
    },
]

const range = async (start = 1) => {
    if (loading.value) return
    loading.value = true
    try {
        const list = await rangeAdmin(searchCondition.value, start, PAGE_SIZE)
        const merged = start === 1 ? list : [...adminList.value, ...list]
        adminList.value = merged.sort((a, b) => a.ID - b.ID)
        noMore.value = list.length < PAGE_SIZE
    } catch {
        // 401 already redirects to the login page via the request wrapper
    } finally {
        loading.value = false
    }
}

// Scroll to bottom: keep requesting the next 100 rows while the backend has more data
const handleEndReached = () => {
    if (loading.value || noMore.value) return
    range(Math.floor(adminList.value.length / PAGE_SIZE) + 1)
}

onMounted(async () => {
    try {
        const res = await post('/admin/info')
        adminInfo.value = res.data
    } catch {
        // 401 already redirects to the login page via the request wrapper; stop auth check on network errors
        return
    }
    if (!canGetAdmin.value) {
        router.push('/')
        return
    }
    range()
})
</script>