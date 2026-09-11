<template>
    <el-card>
        <template #header>
            操作
        </template>
        <p><el-input v-model="searchForm.username" placeholder="管理员用户名"></el-input></p>
        <p><el-input v-model="searchForm.uu_hash" placeholder="管理员UUHash"></el-input></p>
        <p>权限筛选：</p>
        <el-checkbox v-model="searchForm.permission.can_add_admin" label="是否可以添加管理员"></el-checkbox>
        <el-checkbox v-model="searchForm.permission.can_delete_admin" label="是否可以删除管理员"></el-checkbox>
        <el-checkbox v-model="searchForm.permission.can_edit_admin" label="是否可以修改管理员"></el-checkbox>
        <el-checkbox v-model="searchForm.permission.can_get_admin" label="是否可以获取管理员列表"></el-checkbox>
        <el-checkbox v-model="searchForm.permission.can_operate_user" label="是否可以操作用户"></el-checkbox>
        <el-checkbox v-model="searchForm.permission.can_operate_character" label="是否可以操作角色"></el-checkbox>
        <template #footer>
            <div style="text-align: right">
                <el-button type="primary" @click="handleSearch">查询</el-button>
                <el-button type="warning" @click="handleReset">重置</el-button>
                <el-button type="success" @click="handleAdd">添加管理员</el-button>
            </div>
        </template>
    </el-card>
    <br>
    <template v-if="selectedIds.length > 0">
        <el-card>
            <template #header>
                批量操作
            </template>
            <el-text type="primary" size="large" style="text-align: center">已选中 {{ selectedIds.length }} 条</el-text>
            <p>
                <el-button type="primary" @click="handleCheckboxCancel">取消选中</el-button>
                <el-button type="danger" @click="handleBatchDelete">删除选中管理员</el-button>
            </p>
        </el-card>
        <br>
    </template>
    <el-card>
        <template #header>
            管理员列表
        </template>
        <div style="height: 600px">
            <el-auto-resizer>
                <template #default="{ height, width }">
                    <el-table-v2 :columns="columns" :data="adminList" :width="width" :height="height"
                        row-key="ID" :footer-height="noMore ? 32 : 0" fixed @end-reached="handleEndReached">
                        <template #footer>
                            <el-text type="success" v-if="noMore" size="large">已加载全部 {{ adminList.length }} 条</el-text>
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
import { ElCheckbox, ElButton, ElMessage, ElMessageBox } from 'element-plus'

const router = useRouter()
const adminInfo = ref({})
const adminList = ref([])

const PAGE_SIZE = 100
const loading = ref(false)
// 后端是否已无更多数据
const noMore = ref(false)

// 筛选条件表单（提交后由后端筛选）
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
    },
})

// 当前生效的筛选条件，随每次 range 请求发给后端
const searchCondition = ref({})

// 添加管理员弹窗
const dialogVisible = ref(false)

// 编辑管理员弹窗
const editVisible = ref(false)
const editRow = ref(null)

const canGetAdmin = computed(() => {
    if (adminInfo.value.is_root) return true
    return Boolean(adminInfo.value.permission?.can_get_admin)
})

// 权限列统一渲染：命中为“是”，否则“否”
const permCell = (key) => ({ rowData }) => h('span', rowData.Permission?.[key] ? '是' : '否')

// 勾选状态：以行 ID 记录
const selectedIds = ref([])
const isSelected = (id) => selectedIds.value.includes(id)
// 全选判定基于当前数据
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

// 取消全部选中
const handleCheckboxCancel = () => {
    selectedIds.value = []
}

// 查询：把筛选条件交给后端，重新从第一页取
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
    }
    range(1)
}

// 重置：清空筛选条件并重新从第一页取全部
const handleReset = () => {
    searchForm.username = ''
    searchForm.uu_hash = ''
    Object.keys(searchForm.permission).forEach((key) => {
        searchForm.permission[key] = false
    })
    searchCondition.value = {}
    range(1)
}

// 打开添加管理员弹窗
const handleAdd = () => {
    dialogVisible.value = true
}

// 提交添加：调用 /admin/add，成功后刷新列表
const handleAddSubmit = async (form) => {
    try {
        await post('/admin/add', {
            username: form.username.trim(),
            password: form.password,
            permission: { ...form.permission },
        }, { autoRedirect401: false })
        ElMessage.success('添加管理员成功')
        dialogVisible.value = false
        range(1)
    } catch {
        // 失败提示已由 request 封装统一弹出
    }
}

// 执行删除并同步列表（后端每次只接收一个 uu_hash，故逐条调用）
const deleteAdmins = async (rows) => {
    const results = await Promise.allSettled(
        rows.map((item) => del('/admin/delete', { uu_hash: item.UUHash }, { autoRedirect401: false, showError: false }))
    )
    const successIds = rows.filter((_, index) => results[index].status === 'fulfilled').map((item) => item.ID)
    const failed = results.filter((result) => result.status === 'rejected')
    // 成功的行取消勾选
    selectedIds.value = selectedIds.value.filter((id) => !successIds.includes(id))
    if (failed.length === 0) {
        ElMessage.success(`删除成功，共 ${successIds.length} 条`)
    } else if (successIds.length === 0) {
        ElMessage.error(failed[0].reason?.message || '删除失败')
    } else {
        ElMessage.warning(`成功 ${successIds.length} 条，失败 ${failed.length} 条：${failed[0].reason?.message ?? '未知原因'}`)
    }
    if (successIds.length > 0) range(1)
}

// 单个删除：确认后删除该行
const handleDelete = async (row) => {
    try {
        await ElMessageBox.confirm(`确定删除管理员"${row.Username}"吗？`, '删除确认', {
            type: 'warning',
            confirmButtonText: '确定',
            cancelButtonText: '取消',
        })
    } catch {
        // 取消删除
        return
    }
    deleteAdmins([row])
}

// 批量删除：确认后删除所有选中行
const handleBatchDelete = async () => {
    const rows = adminList.value.filter((item) => selectedIds.value.includes(item.ID))
    if (rows.length === 0) return
    try {
        await ElMessageBox.confirm(`确定删除选中的 ${rows.length} 条管理员吗？`, '删除确认', {
            type: 'warning',
            confirmButtonText: '确定',
            cancelButtonText: '取消',
        })
    } catch {
        // 取消删除
        return
    }
    deleteAdmins(rows)
}

// 打开编辑管理员弹窗
const handleEdit = (row) => {
    editRow.value = row
    editVisible.value = true
}

// 提交编辑：调用 /admin/edit，成功后刷新列表
const handleEditSubmit = async (form) => {
    try {
        await put('/admin/edit', {
            uu_hash: form.uu_hash,
            admin_info: {
                password: form.admin_info.password,
                permission: { ...form.admin_info.permission },
            },
        }, { autoRedirect401: false })
        ElMessage.success('修改管理员成功')
        editVisible.value = false
        range(1)
    } catch {
        // 失败提示已由 request 封装统一弹出
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
    { key: 'Username', dataKey: 'Username', title: '用户名', width: 140 },
    { key: 'UUHash', dataKey: 'UUHash', title: 'UUHash', width: 550 },
    { key: 'can_add_admin', title: '添加管理员', width: 110, cellRenderer: permCell('can_add_admin') },
    { key: 'can_delete_admin', title: '删除管理员', width: 110, cellRenderer: permCell('can_delete_admin') },
    { key: 'can_edit_admin', title: '修改管理员', width: 110, cellRenderer: permCell('can_edit_admin') },
    { key: 'can_get_admin', title: '查看管理员', width: 110, cellRenderer: permCell('can_get_admin') },
    { key: 'can_operate_user', title: '操作用户', width: 100, cellRenderer: permCell('can_operate_user') },
    { key: 'can_operate_character', title: '操作角色', width: 100, cellRenderer: permCell('can_operate_character') },
    {
        key: 'actions',
        title: '操作',
        width: 200,
        cellRenderer: ({ rowData }) => h('div', { style: 'display: flex; gap: 8px;' }, [
            h(ElButton, { type: 'danger', size: 'small', onClick: () => handleDelete(rowData) }, () => '删除用户'),
            h(ElButton, { type: 'primary', size: 'small', onClick: () => handleEdit(rowData) }, () => '编辑用户'),
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
        // 401 已由 request 封装自动跳转登录页
    } finally {
        loading.value = false
    }
}

// 滚动到底部：后端还有数据时继续请求下一页 100 行
const handleEndReached = () => {
    if (loading.value || noMore.value) return
    range(Math.floor(adminList.value.length / PAGE_SIZE) + 1)
}

onMounted(async () => {
    try {
        const res = await post('/admin/info')
        adminInfo.value = res.data
    } catch {
        // 401 已由 request 封装自动跳转登录页，网络错误时不再继续鉴权
        return
    }
    if (!canGetAdmin.value) {
        router.push('/')
        return
    }
    range()
})
</script>
